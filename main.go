package main

import (
	"flag"
	"fmt"
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
	http.Handle(prefix, http.StripPrefix(prefix, requestLogger(versionHeader(fileHandler))))
	if err := http.ListenAndServe(listenAddr, nil); err != http.ErrServerClosed && err != nil {
		log.Fatal(err)
	}
}

func versionHeader(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Server", fmt.Sprintf("go-serve %s", version))
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func requestLogger(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		go log.Printf("%s %s from client %s", r.Method, r.RequestURI, r.RemoteAddr)
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
