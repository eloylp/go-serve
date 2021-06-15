package server

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"go.eloylp.dev/kit/http/middleware"

	"go.eloylp.dev/go-serve/config"
	"go.eloylp.dev/go-serve/metrics"
)

func router(cfg *config.Settings, logger *logrus.Logger, docRoot string, info Info) http.Handler {
	r := httprouter.New()
	var userMiddlewares []middleware.Middleware
	if cfg.MetricsEnabled {
		userMiddlewares = append(userMiddlewares, configureMetrics(cfg)...)
	}
	if cfg.MetricsEnabled && cfg.MetricsListenAddr == "" {
		r.Handler(http.MethodGet, cfg.MetricsPath, promhttp.Handler())
		logger.Infof("configuring metrics at %s endpoint", cfg.MetricsPath)
	}
	userMiddlewares = append(userMiddlewares,
		middleware.RequestLogger(logger),
		middleware.ServerHeader(fmt.Sprintf("go-serve %s", Version)),
	)
	if len(cfg.ReadAuthorizations) > 0 {
		logger.Info("configuring read authorizations in server")
		authReadCfg := readAuthConfig(cfg)
		userMiddlewares = append(userMiddlewares, middleware.AuthChecker(authReadCfg))
	}
	if len(cfg.WriteAuthorizations) > 0 {
		logger.Info("configuring write authorizations in server")
		authWriteCfg := writeAuthConfig(cfg)
		userMiddlewares = append(userMiddlewares, middleware.AuthChecker(authWriteCfg))
	}
	r.Handler(http.MethodGet, "/status", StatusHandler(info))
	if cfg.DownloadEndpoint != "" {
		r.Handler(http.MethodGet, cfg.DownloadEndpoint, middleware.For(DownloadTARGZHandler(logger, cfg.DocRoot), userMiddlewares...))
		logger.Infof("configuring downloads at %s endpoint", cfg.DownloadEndpoint)
	}
	if cfg.UploadEndpoint != "" {
		r.Handler(http.MethodPost, cfg.UploadEndpoint, middleware.For(UploadTARGZHandler(logger, cfg.DocRoot), userMiddlewares...))
		logger.Infof("configuring uploads at %s endpoint", cfg.UploadEndpoint)
	}
	fileHandler := http.FileServer(http.Dir(docRoot))
	r.GET(cfg.Prefix+"/*filepath", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		r.URL.Path = p.ByName("filepath")
		middleware.For(fileHandler, userMiddlewares...).ServeHTTP(w, r)
	})
	return r
}

func writeAuthConfig(cfg *config.Settings) *middleware.AuthConfig {
	return middleware.NewAuthConfig().
		WithAuth(middleware.Authorization(cfg.WriteAuthorizations)).
		WithMethod(http.MethodPost).
		WithPathRegex(fmt.Sprintf("^%s$", cfg.UploadEndpoint))
}

func readAuthConfig(cfg *config.Settings) *middleware.AuthConfig {
	return middleware.NewAuthConfig().
		WithAuth(middleware.Authorization(cfg.ReadAuthorizations)).
		WithMethod(http.MethodGet).
		WithPathRegex(".*")
}

func configureMetrics(cfg *config.Settings) []middleware.Middleware {
	metrics.Initialize(cfg)
	mapper := configureEndpointMapper(cfg)
	var metricsMiddlewares []middleware.Middleware
	durationObserver := middleware.RequestDurationObserver(
		"",
		prometheus.DefaultRegisterer,
		cfg.MetricsRequestDurationBuckets,
		mapper,
	)
	metricsMiddlewares = append(metricsMiddlewares, durationObserver)
	responseSizeObserver := middleware.ResponseSizeObserver(
		"",
		prometheus.DefaultRegisterer,
		cfg.MetricsSizeBuckets,
		mapper,
	)
	return append(metricsMiddlewares, responseSizeObserver)
}

func configureEndpointMapper(cfg *config.Settings) *endpointMapper {
	em := newEndpointMapper()
	em.Declare(cfg.Prefix, cfg.Prefix)
	if cfg.UploadEndpoint != "" {
		em.Declare(cfg.UploadEndpoint, cfg.UploadEndpoint)
	}
	if cfg.DownloadEndpoint != "" {
		em.Declare(cfg.DownloadEndpoint, cfg.DownloadEndpoint)
	}
	return em
}
