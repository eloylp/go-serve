package packer

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

func WriteTARGZ(writer io.Writer, path string) (totalBytes int64, err error) {
	gzipStream := gzip.NewWriter(writer)
	defer gzipStream.Close()
	tarStream := tar.NewWriter(gzipStream)
	defer tarStream.Close()
	var bytesWritten int64

	err = filepath.Walk(path, func(currentPath string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}
		header.Name, err = pathDiff(path, currentPath)
		if err != nil {
			return err
		}
		if err = tarStream.WriteHeader(header); err != nil {
			return err
		}
		if !info.IsDir() {
			file, err := os.Open(currentPath)
			if err != nil {
				return err
			}
			defer file.Close()
			bytesWritten, err = io.Copy(tarStream, file)
			if err != nil {
				return err
			}
			totalBytes += bytesWritten
		}
		return nil
	})
	if err != nil {
		return 0, fmt.Errorf("WriteTARGZ(): %w", err)
	}
	return bytesWritten, nil
}

func pathDiff(root, requiredPath string) (string, error) {
	rel, err := filepath.Rel(root, requiredPath)
	if err != nil {
		return "", err
	}
	result := filepath.ToSlash(rel)
	return result, nil
}
