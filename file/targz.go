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

	pathFileInfo, err := os.Stat(path)
	if err != nil {
		return 0, fmt.Errorf("CreateTARGZ(): %w", err) //nolint:golint
	}
	if !pathFileInfo.IsDir() {
		b, err := tarFromFile(path, tarStream)
		if err != nil {
			return 0, fmt.Errorf("CreateTARGZ(): %w", err) //nolint:golint
		}
		return b, nil
	}
	bytes, err := tarFromDir(path, tarStream)
	if err != nil {
		return 0, fmt.Errorf("CreateTARGZ(): %w", err) //nolint:golint
	}
	return bytes, nil
}

func tarFromDir(path string, tarStream *tar.Writer) (int64, error) {
	var totalBytes int64
	err := filepath.Walk(path, func(currentPath string, fileInfo fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		header, err := tar.FileInfoHeader(fileInfo, "")
		if err != nil {
			return err
		}
		header.Name, err = RelativePath(path, currentPath)
		if err != nil {
			return err
		}
		if err := tarStream.WriteHeader(header); err != nil {
			return err
		}
		if !fileInfo.IsDir() {
			bytesWritten, err := appendToWriter(tarStream, currentPath)
			if err != nil {
				return err
			}
			totalBytes += bytesWritten
		}
		return nil
	})
	if err != nil {
		return 0, err //nolint:golint
	}
	return totalBytes, nil
}

func tarFromFile(path string, tarStream *tar.Writer) (int64, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	header, err := tar.FileInfoHeader(fileInfo, "")
	if err != nil {
		return 0, fmt.Errorf("CreateTARGZ(): %w", err) //nolint:golint
	}
	header.Name = filepath.Base(path)
	if err := tarStream.WriteHeader(header); err != nil {
		return 0, fmt.Errorf("CreateTARGZ(): %w", err) //nolint:golint
	}
	b, err := appendToWriter(tarStream, path)
	if err != nil {
		return 0, fmt.Errorf("CreateTARGZ(): %w", err) //nolint:golint
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
	uncompressedStream, err := gzip.NewReader(stream)
	if err != nil {
		return 0, fmt.Errorf("ExtractTARGZ(): failed reading compressed gzip: %w " + err.Error())
	}
	var writtenBytes int64
	tarReader := tar.NewReader(uncompressedStream)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, fmt.Errorf("ExtractTARGZ(): failed reading next part of tar: %w", err)
		}
		extractionPath := filepath.Join(path, header.Name) //nolint:gosec
		// Start processing types
		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(extractionPath, 0755); err != nil {
				return 0, fmt.Errorf("ExtractTARGZ(): failed creating dir %s part of tar: %w", path, err)
			}
		case tar.TypeReg:
			dir := filepath.Dir(extractionPath)
			if err := os.MkdirAll(dir, 0755); err != nil {
				return 0, fmt.Errorf("ExtractTARGZ(): failed creating dir %s part of tar: %w ", dir, err)
			}
			outFile, err := os.Create(extractionPath)
			if err != nil {
				return 0, fmt.Errorf("ExtractTARGZ(): failed creating file part %s of tar: %w", path, err)
			}
			fileBytes, err := io.Copy(outFile, tarReader) // nolinter: gosec (must be controlled by read/write timeouts)
			if err != nil {
				return 0, fmt.Errorf("ExtractTARGZ(): failed copying data of file %s part of tar: %v", path, err)
			}
			writtenBytes += fileBytes
			if err := outFile.Close(); err != nil {
				return writtenBytes, fmt.Errorf("ExtractTARGZ(): failed closing file %s part of tar: %v", path, err)
			}
		default:
			return 0, fmt.Errorf("ExtractTARGZ(): unknown part of tar: type: %v in %s", header.Typeflag, header.Name)
		}
	}
	return writtenBytes, nil
}
