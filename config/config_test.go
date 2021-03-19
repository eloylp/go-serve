package config_test

import (
	"flag"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/eloylp/go-serve/config"
)

// nolint: GoLinterLocal
func TestFromArgs(t *testing.T) {
	samples := []struct {
		context          string
		args             []string
		expectedSettings config.Settings
		err              error
	}{
		{
			"Can be called with a full argument list",
			[]string{"go-serve", "-d", "/root", "-p", "/prefix", "-l", "127.0.0.1:8080", "-a", "./.htpasswd"},
			config.Settings{
				"/root",
				"/prefix",
				"127.0.0.1:8080",
				"./.htpasswd",
				time.Second,
			},
			nil,
		},
		{
			"Can be called with any params falling back in defaults",
			[]string{"go-serve"},
			config.Settings{"^(.*)/go-serve/config$",
				"/",
				"0.0.0.0:8080",
				"",
				time.Second,
			},
			nil,
		},
		{
			"Invoke help with -h must cause errHelp",
			[]string{"go-serve", "-h"},
			config.Settings{
				"^(.*)/go-serve/config$",
				"/",
				"0.0.0.0:8080",
				"",
				time.Second,
			},
			flag.ErrHelp,
		}, {
			"Invoke help with -help must cause errHelp",
			[]string{"go-serve", "-help"},
			config.Settings{"^(.*)/go-serve/config$",
				"/",
				"0.0.0.0:8080",
				"",
				time.Second,
			},
			flag.ErrHelp,
		},
	}
	for _, s := range samples {
		t.Run(s.context, func(t *testing.T) {
			flag.CommandLine = flag.NewFlagSet(s.context, flag.ContinueOnError)
			c, err := config.FromArgs(s.args)
			if err != nil {
				assert.Equal(t, s.err, err, "Error is not expected")
				return
			}
			assert.Regexp(t, s.expectedSettings.DocRoot, c.DocRoot, "Not expected doc root")
			assert.Equal(t, s.expectedSettings.Prefix, c.Prefix, "Not expected prefix")
			assert.Equal(t, s.expectedSettings.ListenAddr, c.ListenAddr, "Not expected listen addr")
			assert.Equal(t, s.expectedSettings.AuthFile, c.AuthFile, "Not expected auth file")
		})
	}
}
