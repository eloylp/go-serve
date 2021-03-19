package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
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
	internalHTTPServer *http.Server
	logger             *logrus.Logger
	cfg                *config.Settings
	wg                 *sync.WaitGroup
}

func New(cfg *config.Settings) *Server {
	logger := loggerFrom(cfg.Logger)
	serverIdentity := fmt.Sprintf("%s %s %s %s", Name, Version, Build, BuildTime)
	logger.Info(serverIdentity)
	logger.Infof("Starting to serve %s at %s ...", cfg.DocRoot, cfg.ListenAddr)
	m := router(cfg, logger, serverIdentity)
	s := &http.Server{
		Addr:         cfg.ListenAddr,
		Handler:      m,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}
	server := &Server{
		internalHTTPServer: s,
		logger:             logger,
		cfg:                cfg,
		wg:                 &sync.WaitGroup{},
	}
	return server
}

func (s *Server) ListenAndServe() error {
	s.wg.Add(1)
	go s.awaitShutdownSignal()
	if err := s.internalHTTPServer.ListenAndServe(); err != http.ErrServerClosed {
		return err
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
		return fmt.Errorf("shutdown: %w", err)
	}
	s.logger.Info("server is now shutdown !")
	return nil
}
