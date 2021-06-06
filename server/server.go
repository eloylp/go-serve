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

	"go.eloylp.dev/go-serve/config"
)

type Server struct {
	identity                     string
	servingRoot                  string
	internalHTTPServer           *http.Server
	alternativeMetricsHTTPServer *http.Server
	logger                       *logrus.Logger
	cfg                          *config.Settings
	wg                           *sync.WaitGroup
	ctx                          context.Context
	cancl                        context.CancelFunc
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
	handler := router(cfg, logger, docRoot, Info{
		Name:      Name,
		Version:   Version,
		Build:     Build,
		BuildTime: BuildTime,
	})
	s := &http.Server{
		Addr:         cfg.ListenAddr,
		Handler:      handler,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}
	ctx, cancl := context.WithCancel(context.Background())
	server := &Server{
		identity:           identity,
		internalHTTPServer: s,
		logger:             logger,
		cfg:                cfg,
		wg:                 &sync.WaitGroup{},
		servingRoot:        docRoot,
		ctx:                ctx,
		cancl:              cancl,
	}
	return server, nil
}

func (s *Server) ListenAndServe() error {
	s.wg.Add(1)
	s.logger.Info(s.identity)
	s.logger.Infof("starting to serve %s at %s ...", s.servingRoot, s.cfg.ListenAddr)
	if s.cfg.MetricsEnabled && s.cfg.MetricsListenAddr != "" {
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
	s.logger.Infof("starting to serve metrics at %s ...", s.cfg.MetricsListenAddr)
	h := promhttp.Handler()
	mux := http.NewServeMux()
	mux.Handle(s.cfg.MetricsPath, h)
	s.alternativeMetricsHTTPServer = &http.Server{
		Handler: mux,
		Addr:    s.cfg.MetricsListenAddr,
	}
	go s.awaitShutdownSignalFor(s.alternativeMetricsHTTPServer)
	go func() {
		if err := s.alternativeMetricsHTTPServer.ListenAndServe(); err != http.ErrServerClosed {
			s.logger.WithError(err).Error("go-serve: metrics: server error")
		}
	}()
}

func (s *Server) awaitShutdownSignalFor(instance *http.Server) {
	defer s.wg.Done()
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGHUP, syscall.SIGTERM)
	select {
	case <-s.ctx.Done():
		return
	case <-signals:
		break
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := instance.Shutdown(ctx); err != nil {
		s.logger.Error("await shutdown: " + err.Error())
		return
	}
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.cancl()
	s.logger.Info("started gracefully shutdown of server ...")
	if err := s.internalHTTPServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("go-serve: shutdown: %w", err)
	}
	if s.cfg.MetricsEnabled && s.cfg.MetricsListenAddr != "" {
		if err := s.alternativeMetricsHTTPServer.Shutdown(ctx); err != nil {
			return fmt.Errorf("go-serve: metrics: shutdown: %w", err)
		}
	}
	s.logger.Info("server is now shutdown !")
	return nil
}
