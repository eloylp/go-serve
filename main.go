package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

func main() {

	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	var docRoot string
	var listenAddr string

	flag.StringVar(&docRoot, "d", currentDir, "Defines the document root docRoot")
	flag.StringVar(&listenAddr, "l", "0.0.0.0:8080", "Defines the listen address")
	flag.Parse()

	fileHandler := http.FileServer(http.Dir(docRoot))
	if err := http.ListenAndServe(listenAddr, fileHandler); err != http.ErrServerClosed && err != nil {
		log.Fatal(err)
	}
}
