package handler

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func ProcessTARGZStream(stream io.Reader, root, deployPath string) (int64, error) {
	uncompressedStream, err := gzip.NewReader(stream)
	if err != nil {
		return 0, fmt.Errorf("failed reading compressed gzip: %w " + err.Error())
	}
	var writtenBytes int64
	tarReader := tar.NewReader(uncompressedStream)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, fmt.Errorf("failed reading next part of tar: %w", err)
		}
		// Check that path does not go outside the document root
		path := filepath.Join(root, deployPath, header.Name) // nolinter: gosec
		if err := checkPath(root, path); err != nil {
			return 0, fmt.Errorf("incorrect deploy path: %w", err)
		}
		// Start processing types
		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(path, 0755); err != nil {
				return 0, fmt.Errorf("failed creating dir %s part of tar: %w", path, err)
			}
		case tar.TypeReg:
			dir := filepath.Dir(path)
			if err := os.MkdirAll(dir, 0755); err != nil {
				return 0, fmt.Errorf("failed creating dir %s part of tar: %w ", dir, err)
			}
			outFile, err := os.Create(path)
			if err != nil {
				return 0, fmt.Errorf("failed creating file part %s of tar: %w", path, err)
			}
			fileBytes, err := io.Copy(outFile, tarReader) // nolinter: gosec (controlled by read/write timeouts)
			if err != nil {
				return 0, fmt.Errorf("failed copying data of file %s part of tar: %v", path, err)
			}
			writtenBytes += fileBytes
			_ = outFile.Close()
		default:
			return 0, fmt.Errorf("unknown part of tar: type: %v in %s", header.Typeflag, header.Name)
		}
	}
	return writtenBytes, nil
}

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
		header.Name, err = headerName(path, currentPath)
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

func headerName(root, requiredPath string) (string, error) {
	rel, err := filepath.Rel(root, requiredPath)
	if err != nil {
		return "", err
	}
	result := filepath.ToSlash(rel)
	return result, err
}
