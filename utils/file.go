package utils

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	httpFileMaxBytes = 100 * (1 << 20) // 50 MB
)

func WriteFile(path string, content string) error {
	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}

func CreateTarFromZip(zipReader *zip.Reader) ([]byte, error) {
	var tarBuffer bytes.Buffer
	err := writeTarArchive(&tarBuffer, zipReader)
	if err != nil {
		return nil, err
	}
	return tarBuffer.Bytes(), nil
}

func writeTarArchive(w io.Writer, zipReader *zip.Reader) error {
	tarWriter := tar.NewWriter(w)
	defer tarWriter.Close()

	for _, file := range zipReader.File {
		err := processFileInZipArchive(file, tarWriter)
		if err != nil {
			return err
		}
	}
	return nil
}

func processFileInZipArchive(file *zip.File, tarWriter *tar.Writer) error {
	fileReader, err := file.Open()
	if err != nil {
		return err
	}
	defer fileReader.Close()

	err = tarWriter.WriteHeader(&tar.Header{
		Name:    file.Name,
		Size:    file.FileInfo().Size(),
		Mode:    int64(file.Mode()),
		ModTime: file.Modified,
		// Note: Zip archives do not store ownership information.
		Uid: 1000,
		Gid: 1000,
	})
	if err != nil {
		return err
	}

	n, err := io.CopyN(tarWriter, fileReader, httpFileMaxBytes)
	log.Println(file.Name, n, err)
	if errors.Is(err, io.EOF) {
		err = nil
	}
	return err
}

func CreateZipFromTar(tarReader *tar.Reader) ([]byte, error) {
	var zipBuffer bytes.Buffer
	err := WriteZipArchive(&zipBuffer, tarReader)
	if err != nil {
		return nil, err
	}
	return zipBuffer.Bytes(), nil
}

func WriteZipArchive(w io.Writer, tarReader *tar.Reader) error {
	zipWriter := zip.NewWriter(w)
	defer zipWriter.Close()

	for {
		tarHeader, err := tarReader.Next()
		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			return err
		}

		zipHeader, err := zip.FileInfoHeader(tarHeader.FileInfo())
		if err != nil {
			return err
		}
		zipHeader.Name = tarHeader.Name
		// Some versions of unzip do not check the mode on a file entry and
		// simply assume that entries with a trailing path separator (/) are
		// directories, and that everything else is a file. Give them a hint.
		if tarHeader.FileInfo().IsDir() && !strings.HasSuffix(tarHeader.Name, "/") {
			zipHeader.Name += "/"
		}

		zipEntry, err := zipWriter.CreateHeader(zipHeader)
		if err != nil {
			return err
		}

		_, err = io.CopyN(zipEntry, tarReader, httpFileMaxBytes)
		if errors.Is(err, io.EOF) {
			err = nil
		}
		if err != nil {
			return err
		}
	}
	return nil // don't need to flush as we call `writer.Close()`
}
