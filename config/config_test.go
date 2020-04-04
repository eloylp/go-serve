package config_test

import (
	"flag"
	"github.com/eloylp/go-serve/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFromArgs(t *testing.T) {

	type sample struct {
		context                               string
		args                                  []string
		docRoot, prefix, listenAddr, authFile string
		err                                   error
	}

	samples := []sample{
		{
			"Can be called with a full argument list",
			[]string{"-d", "/root", "-p", "/prefix", "-l", "127.0.0.1:8080", "-a", "./.htpasswd"},
			"/root",
			"/prefix",
			"127.0.0.1:8080",
			"./.htpasswd",
			nil,
		},
		{
			"Can be called with any params falling back in defaults",
			[]string{},
			"^(.*)/go-serve/config$",
			"/",
			"0.0.0.0:8080",
			"",
			nil,
		},
		{
			"Invoke help with -h must cause errHelp",
			[]string{"-h"},
			"^(.*)/go-serve/config$",
			"/",
			"0.0.0.0:8080",
			"",
			flag.ErrHelp,
		}, {
			"Invoke help with -help must cause errHelp",
			[]string{"-help"},
			"^(.*)/go-serve/config$",
			"/",
			"0.0.0.0:8080",
			"",
			flag.ErrHelp,
		},
	}

	for _, s := range samples {
		t.Run(s.context, func(t *testing.T) {

			flag.CommandLine = flag.NewFlagSet(s.context, flag.ContinueOnError)
			docRoot, prefix, listenAddr, authFile, err := config.FromArgs(s.args)
			if err != nil {
				assert.Equal(t, s.err, err, "Error is not expected")
				return
			}
			assert.Regexp(t, s.docRoot, docRoot, "Not expected doc root")
			assert.Equal(t, s.prefix, prefix, "Not expected prefix")
			assert.Equal(t, s.listenAddr, listenAddr, "Not expected listen addr")
			assert.Equal(t, s.authFile, authFile, "Not expected auth file")
		})
	}
}
