package main

import (
	"errors"
	"flag"
	"log"
	"os"

	"github.com/eloylp/go-serve/config"
)

var (
	Name      string
	Version   string
	Build     string
	BuildTime string
)

func main() {
	s, err := config.FromArgs(os.Args)
	if errors.Is(err, flag.ErrHelp) {
		return
	}
	if err != nil {
		log.Fatal(err)
	}
	server := NewServer(&s)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
