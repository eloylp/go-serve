package main

import (
	"log"

	server "github.com/eloylp/go-serve"
	"github.com/eloylp/go-serve/config"
)

func main() {
	settings, err := config.FromEnv()
	if err != nil {
		log.Fatal(err)
	}
	s := server.New(settings)
	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
