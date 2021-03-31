package server_test

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
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

// Copy will copy directories and files recursively from one root to another.
// This function will abort the test if any operation files. This function will
// not do any type on clean up on fail. This is recommended to use in conjunction
// to testing.TempDir() .
func Copy(t *testing.T, source, dest string) {
	err := filepath.Walk(source, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		sourceRel, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}
		if info.IsDir() {
			err = os.MkdirAll(filepath.Join(dest, sourceRel), 0777)
			if err != nil {
				return err
			}
			return nil
		}
		fileFrom, err := os.Open(path)
		if err != nil {
			return err
		}
		defer fileFrom.Close()
		fileTo, err := os.Create(filepath.Join(dest, sourceRel))
		if err != nil {
			return err
		}
		defer fileTo.Close()
		if _, err := io.Copy(fileTo, fileFrom); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
