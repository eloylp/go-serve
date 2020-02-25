package main

import (
	"flag"
	"fmt"
	"github.com/eloylp/go-serve/handler"
	"github.com/eloylp/go-serve/logging"
	"github.com/eloylp/go-serve/www"
	"log"
	"net/http"
	"os"
)

var version string

func main() {

	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	var docRoot string
	var prefix string
	var listenAddr string

	flag.StringVar(&docRoot, "d", currentDir, "Defines the document root")
	flag.StringVar(&prefix, "p", "/", "Defines prefix to use for serve files")
	flag.StringVar(&listenAddr, "l", "0.0.0.0:8080", "Defines the listen address")
	flag.Parse()

	fmt.Println(fmt.Sprintf("go-serve %s", version))
	log.Println(fmt.Sprintf("Starting to serve %s at %s ...", docRoot, listenAddr))

	fileHandler := http.FileServer(http.Dir(docRoot))
	logger := logging.NewConsoleLogger()
	http.Handle(prefix, http.StripPrefix(prefix, www.Apply(fileHandler, handler.ServerHeader(version), handler.RequestLogger(logger))))
	if err := http.ListenAndServe(listenAddr, nil); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
