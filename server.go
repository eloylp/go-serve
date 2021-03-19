package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/eloylp/go-serve/config"
	"github.com/eloylp/go-serve/handler"
	"github.com/eloylp/go-serve/logging"
	"github.com/eloylp/go-serve/www"
)

var (
	Name      string
	Version   string
	Build     string
	BuildTime string
)

type Server struct {
	internalHTTPServer *http.Server
	logger             logging.Logger
	cfg                config.Settings
	wg                 sync.WaitGroup
}

func New(cfg *config.Settings) *Server {
	serverIdentity := fmt.Sprintf("%s %s %s %s", Name, Version, Build, BuildTime)
	fmt.Println(serverIdentity)
	log.Println(fmt.Sprintf("Starting to serve %s at %s ...", cfg.DocRoot, cfg.ListenAddr))
	fileHandler := http.FileServer(http.Dir(cfg.DocRoot))
	logger := logging.NewConsoleLogger()

	middlewares := []www.Middleware{
		handler.ServerHeader(Version),
		handler.RequestLogger(logger),
	}
	if cfg.AuthFile != "" {
		middlewares = append(middlewares, handler.AuthChecker(serverIdentity, cfg.AuthFile))
	}
	m := http.NewServeMux()
	m.Handle(cfg.Prefix, http.StripPrefix(cfg.Prefix, www.Apply(fileHandler, middlewares...)))
	s := &http.Server{
		Addr:         cfg.ListenAddr,
		Handler:      m,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}
	server := &Server{
		internalHTTPServer: s,
		logger:             logger,
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
	log.Println("started gracefully shutdown of server ...")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		s.logger.Error("await shutdown: " + err.Error())
		return
	}
	log.Println("server shutdown gracefully")
}

func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.internalHTTPServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("shutdown: %w", err)
	}
	return nil
}
