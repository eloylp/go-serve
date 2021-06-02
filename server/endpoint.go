package server

import (
	iradix "github.com/hashicorp/go-immutable-radix"
)

type endpointMapper struct {
	t *iradix.Tree
}

func newEndpointMapper() *endpointMapper {
	return &endpointMapper{
		t: iradix.New(),
	}
}
func (e *endpointMapper) Declare(endpoint, name string) {
	e.t, _, _ = e.t.Insert([]byte(endpoint), name)
}
func (e *endpointMapper) Map(url string) string {
	_, name, ok := e.t.Root().LongestPrefix([]byte(url))
	if !ok {
		return ""
	}
	return name.(string)
}
