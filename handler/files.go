package handler

import (
	"fmt"
	"path/filepath"
	"strings"
)

func checkPath(docRoot, path string) error {
	absRoot, err := filepath.Abs(docRoot)
	if err != nil {
		return err
	}
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	if !strings.HasPrefix(absPath, absRoot) {
		return fmt.Errorf("the path you provided %s is not a suitable one", path)
	}
	return nil
}

func headerName(root, requiredPath string) (string, error) {
	rel, err := filepath.Rel(root, requiredPath)
	if err != nil {
		return "", err
	}
	result := filepath.ToSlash(rel)
	return result, err
}
