package main

import (
	"log"

	"github.com/eloylp/go-serve/config"
	"github.com/eloylp/go-serve/server"
)

func main() {
	settings, err := config.FromEnv()
	if err != nil {
		log.Fatal(err)
	}
	s, err := server.New(settings)
	if err != nil {
		log.Fatal(err)
	}
	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
