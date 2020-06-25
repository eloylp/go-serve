package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
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

func main() {
	logger := logging.NewConsoleLogger()
	s, err := config.FromArgs(os.Args)
	if errors.Is(err, flag.ErrHelp) {
		return
	}
	if err != nil {
		log.Fatal(err)
	}

	serverIdentity := fmt.Sprintf("%s %s %s %s", Name, Version, Build, BuildTime)
	fmt.Println(serverIdentity)
	log.Println(fmt.Sprintf("Starting to serve %s at %s ...", s.DocRoot, s.ListenAddr))
	fileHandler := http.FileServer(http.Dir(s.DocRoot))

	middlewares := []www.Middleware{
		handler.ServerHeader(Version),
		handler.RequestLogger(logger),
	}
	if s.AuthFile != "" {
		middlewares = append(middlewares, handler.AuthChecker(serverIdentity, s.AuthFile))
	}
	m := http.NewServeMux()
	m.Handle(s.Prefix, http.StripPrefix(s.Prefix, www.Apply(fileHandler, middlewares...)))
	server := &http.Server{
		Addr:         s.ListenAddr,
		Handler:      m,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}
	www.Shutdown(server, 20*time.Second) //nolint:gomnd
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}
	log.Println("server shutdown gracefully")
}
