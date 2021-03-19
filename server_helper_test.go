package main

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func BodyFrom(t *testing.T, url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	return d
}

func md5From(data []byte) string {
	md5.New()
	return fmt.Sprintf("%x", md5.Sum(data))
}
