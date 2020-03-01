package main

import (
	"fmt"
	"github.com/eloylp/go-serve/config"
	"github.com/eloylp/go-serve/handler"
	"github.com/eloylp/go-serve/logging"
	"github.com/eloylp/go-serve/www"
	"log"
	"net/http"
	"os"
)

var version string

func main() {

	logger := logging.NewConsoleLogger()

	docRoot, prefix, listenAddr, err := config.FromArgs(os.Args)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(fmt.Sprintf("go-serve %s", version))
	log.Println(fmt.Sprintf("Starting to serve %s at %s ...", docRoot, listenAddr))

	fileHandler := http.FileServer(http.Dir(docRoot))
	http.Handle(prefix, http.StripPrefix(prefix, www.Apply(fileHandler, handler.ServerHeader(version), handler.RequestLogger(logger))))
	if err := http.ListenAndServe(listenAddr, nil); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
