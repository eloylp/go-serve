package main

import (
	"errors"
	"flag"
	"log"
	"os"

	"github.com/eloylp/go-serve"
	"github.com/eloylp/go-serve/config"
)

func main() {
	settings, err := config.FromArgs(os.Args)
	if errors.Is(err, flag.ErrHelp) {
		return
	}
	if err != nil {
		log.Fatal(err)
	}
	s := server.New(&settings)
	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
