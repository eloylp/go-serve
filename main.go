package main

import (
	"context"
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

	docRoot, prefix, listenAddr, err := config.FromArgs(os.Args)
	if err != nil {
		log.Fatal(err)
	}
	token := config.FromEnvAuthToken()

	fmt.Println(fmt.Sprintf("go-serve %s", version))
	log.Println(fmt.Sprintf("Starting to serve %s at %s ...", docRoot, listenAddr))
	fileHandler := http.FileServer(http.Dir(docRoot))

	middlewares := []www.Middleware{
		handler.ServerHeader(version),
		handler.RequestLogger(logger),
		handler.AuthChecker(token),
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
