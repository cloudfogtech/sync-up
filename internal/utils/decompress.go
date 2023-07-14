package utils

import (
	"archive/tar"
	"compress/gzip"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"strings"
)

func Decompress(source, target string) error {
	sourceFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer func(sourceFile *os.File) {
		err := sourceFile.Close()
		if err != nil {
			log.Errorf("sourceFile close error: %v", err)
		}
	}(sourceFile)
	gzipReader, err := gzip.NewReader(sourceFile)
	if err != nil {
		log.Errorf("gzip new reader error: %v", err)
		return err
	}
	defer func(gzipReader *gzip.Reader) {
		err := gzipReader.Close()
		if err != nil {
			log.Errorf("gzipReader close error: %v", err)
		}
	}(gzipReader)
	tarReader := tar.NewReader(gzipReader)
	for {
		hdr, err := tarReader.Next()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				log.Errorf("tarReader read next error: %v", err)
				return err
			}
		}
		filename := target + hdr.Name
		file, err := createFile(filename)
		if err != nil {
			return err
		}
		_, err = io.Copy(file, tarReader)
		if err != nil {
			log.Errorf("data copy error: %v", err)
			return err
		}
	}
	return nil
}

func createFile(filePath string) (*os.File, error) {
	err := os.MkdirAll(string([]rune(filePath)[0:strings.LastIndex(filePath, "/")]), 0755)
	if err != nil {
		return nil, err
	}
	return os.Create(filePath)
}
