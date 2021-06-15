package server

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"go.eloylp.dev/kit/http/middleware"

	"go.eloylp.dev/go-serve/config"
	"go.eloylp.dev/go-serve/metrics"
)

func router(cfg *config.Settings, logger *logrus.Logger, docRoot string, info Info) http.Handler {
	r := httprouter.New()
	var authReadCfg *middleware.AuthConfig
	if len(cfg.ReadAuthorizations) > 0 {
		authReadCfg = middleware.NewAuthConfig().
			WithAuth(middleware.Authorization(cfg.ReadAuthorizations)).
			WithMethod(http.MethodGet).
			WithPathRegex(".*")
		logger.Info("configuring read authorizations in server")
	}
	var authWriteCfg *middleware.AuthConfig
	if len(cfg.WriteAuthorizations) > 0 {
		authWriteCfg = middleware.NewAuthConfig().
			WithAuth(middleware.Authorization(cfg.WriteAuthorizations)).
			WithMethod(http.MethodPost).
			WithPathRegex(fmt.Sprintf("^%s$", cfg.UploadEndpoint))
		logger.Info("configuring write authorizations in server")
	}
	var userMiddlewares []middleware.Middleware
	if cfg.MetricsEnabled {
		metrics.Initialize(cfg)
		mapper := configureEndpointMapper(cfg)
		durationObserver := middleware.RequestDurationObserver(
			"",
			prometheus.DefaultRegisterer,
			cfg.MetricsRequestDurationBuckets,
			mapper,
		)
		userMiddlewares = append(userMiddlewares, durationObserver)
		responseSizeObserver := middleware.ResponseSizeObserver(
			"",
			prometheus.DefaultRegisterer,
			cfg.MetricsSizeBuckets,
			mapper,
		)
		userMiddlewares = append(userMiddlewares, responseSizeObserver)
	}
	if cfg.MetricsListenAddr == "" {
		r.Handler(http.MethodGet, cfg.MetricsPath, promhttp.Handler())
		logger.Infof("configuring metrics at %s endpoint", cfg.MetricsPath)
	}
	userMiddlewares = append(userMiddlewares,
		middleware.RequestLogger(logger),
		middleware.ServerHeader(fmt.Sprintf("go-serve %s", Version)),
		middleware.AuthChecker(authReadCfg),
		middleware.AuthChecker(authWriteCfg),
	)
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
