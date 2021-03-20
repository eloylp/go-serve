package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/eloylp/go-serve/config"
)

var (
	Name      string
	Version   string
	Build     string
	BuildTime string
)

type Server struct {
	identity           string
	servingRoot        string
	internalHTTPServer *http.Server
	logger             *logrus.Logger
	cfg                *config.Settings
	wg                 *sync.WaitGroup
	lock               *sync.RWMutex
}

func New(cfg *config.Settings) (*Server, error) {
	logger, err := loggerFrom(cfg.Logger)
	if err != nil {
		return nil, fmt.Errorf("go-serve: %w", err)
	}
	identity := fmt.Sprintf("%s %s %s %s", Name, Version, Build, BuildTime)
	docRoot, err := filepath.Abs(cfg.DocRoot)
	if err != nil {
		return nil, fmt.Errorf("go-serve: %w", err)
	}
	m := router(cfg, logger, docRoot, identity)
	s := &http.Server{
		Addr:         cfg.ListenAddr,
		Handler:      m,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}
	server := &Server{
		identity:           identity,
		internalHTTPServer: s,
		logger:             logger,
		cfg:                cfg,
		wg:                 &sync.WaitGroup{},
		servingRoot:        docRoot,
		lock:               &sync.RWMutex{},
	}
	return server, nil
}

func (s *Server) ListenAndServe() error {
	s.wg.Add(1)
	s.logger.Info(s.identity)
	s.logger.Infof("starting to serve %s at %s ...", s.servingRoot, s.cfg.ListenAddr)
	go s.awaitShutdownSignal()
	if err := s.internalHTTPServer.ListenAndServe(); err != http.ErrServerClosed {
		return fmt.Errorf("go-serve: %w", err)
	}
	s.wg.Wait()
	return nil
}

func (s *Server) awaitShutdownSignal() {
	defer s.wg.Done()
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGHUP, syscall.SIGTERM)
	<-signals
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		s.logger.Error("await shutdown: " + err.Error())
		return
	}
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("started gracefully shutdown of server ...")
	if err := s.internalHTTPServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("go-serve: shutdown: %w", err)
	}
	s.logger.Info("server is now shutdown !")
	return nil
}
