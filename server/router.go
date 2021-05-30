package server

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"go.eloylp.dev/kit/http/middleware"

	"github.com/eloylp/go-serve/config"
)

func router(cfg *config.Settings, logger *logrus.Logger, docRoot string, info Info) http.Handler {
	r := httprouter.New()
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
	var userMiddlewares []middleware.Middleware
	if cfg.MetricsEnabled && cfg.MetricsListenAddr == "" {
		observer := middleware.RequestDurationObserver("goserve", cfg.PrometheusRegistry, cfg.MetricsRequestDurationBuckets)
		userMiddlewares = append(userMiddlewares, observer)
		r.Handler(http.MethodGet, cfg.MetricsPath, promhttp.HandlerFor(cfg.PrometheusRegistry, promhttp.HandlerOpts{}))
	}
	userMiddlewares = append(userMiddlewares,
		middleware.RequestLogger(logger),
		middleware.ServerHeader(fmt.Sprintf("go-serve %s", Version)),
		middleware.AuthChecker(authReadCfg),
		middleware.AuthChecker(authWriteCfg),
	)
	r.Handler(http.MethodGet, "/status", StatusHandler(info))
	if cfg.DownloadEndpoint != "" {
		r.Handler(http.MethodGet, cfg.DownloadEndpoint, middleware.InFrontOf(DownloadTARGZHandler(logger, cfg.DocRoot), userMiddlewares...))
	}
	if cfg.UploadEndpoint != "" {
		r.Handler(http.MethodPost, cfg.UploadEndpoint, middleware.InFrontOf(UploadTARGZHandler(logger, cfg.DocRoot), userMiddlewares...))
	}
	fileHandler := http.FileServer(http.Dir(docRoot))
	r.GET(cfg.Prefix+"/*filepath", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		r.URL.Path = p.ByName("filepath")
		middleware.InFrontOf(fileHandler, userMiddlewares...).ServeHTTP(w, r)
	})
	return r
}
