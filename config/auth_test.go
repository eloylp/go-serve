package config_test

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/eloylp/go-serve/config"
)

func TestAuthorization_Decode(t *testing.T) {
	htpasswdContent := `
user1:$2y$10$4N6UIL11veX3dDP3n5TEquYrYVPSxF/ZAya3eqXXLTbRqDPDYlMr2
user2:$2y$10$TO.aZyNrGPGWuI2m55TsNe6XoOT7kh70idr6fMMLOaHSUK5guuEXi

`
	b64EnvVar := base64.StdEncoding.EncodeToString([]byte(htpasswdContent))

	a := config.Authorization{}
	err := a.Decode(b64EnvVar)
	assert.NoError(t, err)
	assert.Equal(t, config.Authorization{
		"user1": "$2y$10$4N6UIL11veX3dDP3n5TEquYrYVPSxF/ZAya3eqXXLTbRqDPDYlMr2",
		"user2": "$2y$10$TO.aZyNrGPGWuI2m55TsNe6XoOT7kh70idr6fMMLOaHSUK5guuEXi",
	}, a)
}
