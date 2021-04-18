package file

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

func CreateTARGZ(writer io.Writer, path string) (totalBytes int64, err error) {
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
		header.Name, err = RelativePath(path, currentPath)
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
		return 0, fmt.Errorf("CreateTARGZ(): %w", err)
	}
	return bytesWritten, nil
}

func ExtractTARGZ(stream io.Reader, path string) (int64, error) {
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
		extractionPath := filepath.Join(path, header.Name)
		// Start processing types
		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(extractionPath, 0755); err != nil {
				return 0, fmt.Errorf("failed creating dir %s part of tar: %w", path, err)
			}
		case tar.TypeReg:
			dir := filepath.Dir(extractionPath)
			if err := os.MkdirAll(dir, 0755); err != nil {
				return 0, fmt.Errorf("failed creating dir %s part of tar: %w ", dir, err)
			}
			outFile, err := os.Create(extractionPath)
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
