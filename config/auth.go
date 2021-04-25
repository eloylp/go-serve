package config

import (
	"encoding/base64"
	"strings"
)

type Authorization map[string]string

func (a Authorization) Decode(value string) error {
	decodedValue, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return err
	}
	for _, line := range strings.Split(string(decodedValue), "\n") {
		if line == "" {
			continue
		}
		parts := strings.Split(line, ":")
		a[parts[0]] = parts[1]
	}
	return nil
}
