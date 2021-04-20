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

func CreateTARGZ(writer io.Writer, path string) (int64, error) {
	gzipReader := gzip.NewWriter(writer)
	defer gzipReader.Close()
	tarReader := tar.NewWriter(gzipReader)
	defer tarReader.Close()

	pathInfo, err := os.Stat(path)
	if err != nil {
		return 0, fmt.Errorf("CreateTARGZ(): %w", err) //nolint:golint
	}
	if !pathInfo.IsDir() {
		b, err := tarFromFile(path, tarReader)
		if err != nil {
			return 0, fmt.Errorf("CreateTARGZ(): %w", err) //nolint:golint
		}
		return b, nil
	}
	b, err := tarFromDir(path, tarReader)
	if err != nil {
		return 0, fmt.Errorf("CreateTARGZ(): %w", err) //nolint:golint
	}
	return b, nil
}

func tarFromDir(path string, tarWriter *tar.Writer) (int64, error) {
	var totalFileBytes int64
	err := filepath.Walk(path, func(currentPath string, info fs.FileInfo, err error) error {
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
		if err := tarWriter.WriteHeader(header); err != nil {
			return err
		}
		if !info.IsDir() {
			b, err := appendToWriter(tarWriter, currentPath)
			if err != nil {
				return err
			}
			totalFileBytes += b
		}
		return nil
	})
	if err != nil {
		return 0, err //nolint:golint
	}
	return totalFileBytes, nil
}

func tarFromFile(path string, tarStream *tar.Writer) (int64, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	header, err := tar.FileInfoHeader(fileInfo, "")
	if err != nil {
		return 0, err
	}
	header.Name = filepath.Base(path)
	if err := tarStream.WriteHeader(header); err != nil {
		return 0, err
	}
	b, err := appendToWriter(tarStream, path)
	if err != nil {
		return 0, err
	}
	return b, nil
}

func appendToWriter(w io.Writer, path string) (int64, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer file.Close()
	b, err := io.Copy(w, file)
	if err != nil {
		return 0, err
	}
	return b, nil
}

func ExtractTARGZ(stream io.Reader, path string) (int64, error) {
	gzipReader, err := gzip.NewReader(stream)
	if err != nil {
		return 0, fmt.Errorf("ExtractTARGZ(): failed reading compressed gzip: %w " + err.Error())
	}
	tarReader := tar.NewReader(gzipReader)
	var totalFileBytes int64
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, fmt.Errorf("ExtractTARGZ(): failed reading next part of tar: %w", err) //nolint:golint
		}
		extractionPath := filepath.Join(path, header.Name) //nolint:gosec
		// Start processing types
		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(extractionPath, 0755); err != nil {
				return 0, fmt.Errorf("ExtractTARGZ(): failed creating dir %s part of tar: %w", path, err) //nolint:golint
			}
		case tar.TypeReg:
			dir := filepath.Dir(extractionPath)
			if err := os.MkdirAll(dir, 0755); err != nil {
				return 0, fmt.Errorf("ExtractTARGZ(): failed creating dir %s part of tar: %w ", dir, err) //nolint:golint
			}
			outFile, err := os.Create(extractionPath)
			if err != nil {
				return 0, fmt.Errorf("ExtractTARGZ(): failed creating file part %s of tar: %w", path, err) //nolint:golint
			}
			b, err := io.Copy(outFile, tarReader) // nolinter: gosec (must be controlled by read/write timeouts)
			if err != nil {
				return 0, fmt.Errorf("ExtractTARGZ(): failed copying data of file %s part of tar: %v", path, err) //nolint:golint
			} //nolint:golint
			totalFileBytes += b
			if err := outFile.Close(); err != nil {
				return totalFileBytes, fmt.Errorf("ExtractTARGZ(): failed closing file %s part of tar: %v", path, err) //nolint:golint
			} //nolint:golint
		default:
			return 0, fmt.Errorf("ExtractTARGZ(): unknown part of tar: type: %v in %s", header.Typeflag, header.Name) //nolint:golint
		}
	}
	return totalFileBytes, nil
}
