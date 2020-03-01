// Package config is a local use package for grouping
// app config related functions
package config

import (
	"flag"
	"os"
)

// FromArgs will receive and argument list as parameter, returning the proper
// variables with all the values.
func FromArgs(args []string) (docRoot, prefix, listenAddr string, err error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return
	}
	flag.StringVar(&docRoot, "d", currentDir, "Defines the document root")
	flag.StringVar(&prefix, "p", "/", "Defines prefix to use for serve files")
	flag.StringVar(&listenAddr, "l", "0.0.0.0:8080", "Defines the listen address")

	if err = flag.CommandLine.Parse(args); err != nil {
		return
	}
	return docRoot, prefix, listenAddr, nil
}

// FromEnvAuthToken encapsulates the gathering of the imposed
// auth token
func FromEnvAuthToken() string {
	return os.Getenv("GO_SERVE_AUTH_TOKEN")
}
