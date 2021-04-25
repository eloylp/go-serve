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

	"github.com/prometheus/client_golang/prometheus/promhttp"
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
	identity                     string
	servingRoot                  string
	internalHTTPServer           *http.Server
	alternativeMetricsHTTPServer *http.Server
	logger                       *logrus.Logger
	cfg                          *config.Settings
	wg                           *sync.WaitGroup
}

func New(cfg *config.Settings) (*Server, error) {
	logger, err := logger(cfg.Logger)
	if err != nil {
		return nil, fmt.Errorf("go-serve: %w", err)
	}
	identity := fmt.Sprintf("%s %s %s %s", Name, Version, Build, BuildTime)
	docRoot, err := filepath.Abs(cfg.DocRoot)
	if err != nil {
		return nil, fmt.Errorf("go-serve: %w", err)
	}
	handler := router(cfg, logger, docRoot)
	s := &http.Server{
		Addr:         cfg.ListenAddr,
		Handler:      handler,
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
	}
	return server, nil
}

func (s *Server) ListenAndServe() error {
	s.wg.Add(1)
	s.logger.Info(s.identity)
	s.logger.Infof("starting to serve %s at %s ...", s.servingRoot, s.cfg.ListenAddr)
	if s.cfg.MetricsEnabled && s.cfg.MetricsAlternativeListenAddr != "" {
		s.startAlternateMetricsServer()
	}
	go s.awaitShutdownSignalFor(s.internalHTTPServer)
	if err := s.internalHTTPServer.ListenAndServe(); err != http.ErrServerClosed {
		return fmt.Errorf("go-serve: %w", err)
	}
	s.wg.Wait()
	return nil
}

func (s *Server) startAlternateMetricsServer() {
	s.wg.Add(1)
	s.logger.Infof("starting to serve metrics at %s ...", s.cfg.MetricsAlternativeListenAddr)
	h := promhttp.HandlerFor(s.cfg.PrometheusRegistry, promhttp.HandlerOpts{})
	mux := http.NewServeMux()
	mux.Handle(s.cfg.MetricsPath, h)
	s.alternativeMetricsHTTPServer = &http.Server{
		Handler: mux,
		Addr:    s.cfg.MetricsAlternativeListenAddr,
	}
	go s.awaitShutdownSignalFor(s.alternativeMetricsHTTPServer)
	go func() {
		if err := s.alternativeMetricsHTTPServer.ListenAndServe(); err != http.ErrServerClosed {
			s.logger.WithError(err).Error("go-serve: metrics: server error")
		}
	}()
}

// TOdo, this function is leaking resources. Pass context and
// react on its cancellation.
func (s *Server) awaitShutdownSignalFor(instance *http.Server) {
	defer s.wg.Done()
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGHUP, syscall.SIGTERM)
	<-signals
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := instance.Shutdown(ctx); err != nil {
		s.logger.Error("await shutdown: " + err.Error())
		return
	}
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("started gracefully shutdown of server ...")
	if err := s.internalHTTPServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("go-serve: shutdown: %w", err)
	}
	if s.cfg.MetricsEnabled && s.cfg.MetricsAlternativeListenAddr != "" {
		if err := s.alternativeMetricsHTTPServer.Shutdown(ctx); err != nil {
			return fmt.Errorf("go-serve: metrics: shutdown: %w", err)
		}
	}
	s.logger.Info("server is now shutdown !")
	return nil
}
