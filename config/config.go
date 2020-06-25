// Package config is a local use package for grouping
// app config related functions
package config

import (
	"flag"
	"os"
)

type Settings struct {
	DocRoot, Prefix, ListenAddr, AuthFile string
}

// FromArgs will receive and argument list as parameter, including
// the program name and returning the proper variables with all the values.
func FromArgs(args []string) (s Settings, err error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return
	}
	flag.StringVar(&s.DocRoot, "d", currentDir, "Defines the document root")
	flag.StringVar(&s.Prefix, "p", "/", "Defines prefix to use for serve files")
	flag.StringVar(&s.ListenAddr, "l", "0.0.0.0:8080", "Defines the listen address")
	flag.StringVar(&s.AuthFile, "a", "", "Defines the .htpasswd file path for auth")
	argsFiltered := args[1:] // exclude program name
	if err = flag.CommandLine.Parse(argsFiltered); err != nil {
		return
	}
	return
}
