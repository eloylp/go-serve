package server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/eloylp/kit/http/middleware"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"

	"github.com/eloylp/go-serve/config"
)

func router(cfg *config.Settings, logger *logrus.Logger, docRoot string, info Info) http.Handler {
	r := mux.NewRouter()
	var authReadCfg *middleware.AuthConfig
	if len(cfg.ReadAuthorizations) > 0 {
		authReadCfg = middleware.NewAuthConfig().
			WithAuth(middleware.Authorization(cfg.ReadAuthorizations)).
			WithMethod(http.MethodGet).
			WithPathRegex(".*")
	}
	var authWriteCfg *middleware.AuthConfig
	if len(cfg.WriteAuthorizations) > 0 {
		authWriteCfg = middleware.NewAuthConfig().
			WithAuth(middleware.Authorization(cfg.WriteAuthorizations)).
			WithMethod(http.MethodPost).
			WithPathRegex(fmt.Sprintf("^%s$", cfg.UploadEndpoint))
	}
	r.Use(mux.MiddlewareFunc(middleware.RequestLogger(logger)))
	if cfg.MetricsEnabled && cfg.MetricsAlternativeListenAddr == "" {
		r.Use(mux.MiddlewareFunc(middleware.RequestDurationObserver("goserve", cfg.PrometheusRegistry, cfg.MetricsRequestDurationBuckets)))
		r.Methods(http.MethodGet).
			Path(cfg.MetricsPath).
			Handler(promhttp.HandlerFor(cfg.PrometheusRegistry, promhttp.HandlerOpts{}))
	}
	r.Methods(http.MethodGet).Path("/status").Handler(StatusHandler(info))
	r.Use(
		mux.MiddlewareFunc(middleware.ServerHeader(fmt.Sprintf("go-serve %s", Version))),
		mux.MiddlewareFunc(middleware.AuthChecker(authReadCfg)),
		mux.MiddlewareFunc(middleware.AuthChecker(authWriteCfg)),
	)
	if cfg.DownloadEndpoint != "" {
		r.Methods(http.MethodGet).
			Path(cfg.DownloadEndpoint).
			Handler(DownloadTARGZHandler(logger, cfg.DocRoot)).
			Headers("Accept", "application/tar+gzip")
	}
	if cfg.UploadEndpoint != "" {
		r.Methods(http.MethodPost).
			Path(cfg.UploadEndpoint).
			Handler(UploadTARGZHandler(logger, cfg.DocRoot)).
			Headers("Content-Type", "application/tar+gzip")
	}
	fileHandler := http.FileServer(http.Dir(docRoot))
	r.Methods(http.MethodGet).PathPrefix(cfg.Prefix).Handler(http.StripPrefix(cfg.Prefix, fileHandler))
	debugRouter(r, logger)
	return r
}

func debugRouter(r *mux.Router, logger *logrus.Logger) {
	err := r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, err := route.GetPathTemplate()
		if err == nil {
			logger.Debugf("router: registering ROUTE: %s", pathTemplate)
		}
		pathRegexp, err := route.GetPathRegexp()
		if err == nil {
			logger.Debugf("router: path regexp: %s", pathRegexp)
		}
		queriesTemplates, err := route.GetQueriesTemplates()
		if err == nil {
			logger.Debugf("router: queries templates: %s", strings.Join(queriesTemplates, ","))
		}
		queriesRegexps, err := route.GetQueriesRegexp()
		if err == nil {
			logger.Debugf("router: Queries regexps: %s", strings.Join(queriesRegexps, ","))
		}
		methods, err := route.GetMethods()
		if err == nil {
			logger.Debugf("router: methods: %s", strings.Join(methods, ","))
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
}
