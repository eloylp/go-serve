package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/eloylp/go-serve/config"
	"github.com/eloylp/go-serve/handler"
	"github.com/eloylp/go-serve/logging"
	"github.com/eloylp/go-serve/www"
	"log"
	"net/http"
	"os"
	"time"
)

var version string

func main() {

	logger := logging.NewConsoleLogger()

	docRoot, prefix, listenAddr, authFile, err := config.FromArgs(os.Args[1:])
	if errors.Is(err, flag.ErrHelp) {
		return
	}
	if err != nil {
		log.Fatal(err)
	}

	serverIdentity := fmt.Sprintf("go-serve %s", version)
	fmt.Println(serverIdentity)
	log.Println(fmt.Sprintf("Starting to serve %s at %s ...", docRoot, listenAddr))
	fileHandler := http.FileServer(http.Dir(docRoot))

	middlewares := []www.Middleware{
		handler.ServerHeader(version),
		handler.RequestLogger(logger),
	}
	if authFile != "" {
		middlewares = append(middlewares, handler.AuthChecker(serverIdentity, authFile))
	}
	m := http.NewServeMux()
	m.Handle(prefix, http.StripPrefix(prefix, www.Apply(fileHandler, middlewares...)))
	s := &http.Server{
		Addr:    listenAddr,
		Handler: m,
	}
	ctx, _ := context.WithTimeout(context.Background(), 20*time.Second)
	www.Shutdown(ctx, s)
	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}
	log.Println("server shutdown gracefully")
}
