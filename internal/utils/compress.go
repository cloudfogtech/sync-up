package utils

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
	"time"
)

func CompressToDir(sourcePath, targetDir string) (string, error) {
	filename := filepath.Base(sourcePath)
	timePattern := time.Now().Format("20060102T150405Z07")
	targetPath := filepath.Join(targetDir, fmt.Sprintf("%s.%s.tar.gz", filename, timePattern))
	return targetPath, CompressToFile(sourcePath, targetPath)
}

func CompressToFile(sourcePath, targetPath string) error {
	file, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	fileWriter, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	defer func(writer *os.File) {
		err := writer.Close()
		if err != nil {
			log.Errorf("file writer close error, %v", err)
		}
	}(fileWriter)
	gzipWriter := gzip.NewWriter(fileWriter)
	gzipWriter.Name = filepath.Base(sourcePath)
	defer func(writer *gzip.Writer) {
		err := writer.Close()
		if err != nil {
			log.Errorf("gzip writer close error, %v", err)
		}
	}(gzipWriter)
	tarWriter := tar.NewWriter(gzipWriter)
	defer func(tarWriter *tar.Writer) {
		err := tarWriter.Close()
		if err != nil {
			log.Errorf("tar writer close error, %v", err)
		}
	}(tarWriter)
	err = compress(file, "", tarWriter)
	if err != nil {
		log.Errorf("compress file error: %v", err)
		return err
	}
	return nil
}

func compress(currentFile *os.File, prefix string, tarWriter *tar.Writer) error {
	info, err := currentFile.Stat()
	if err != nil {
		log.Errorf("get currentFile info error: %v", err)
		return err
	}
	if info.IsDir() {
		err = compressDir(info, currentFile, prefix, tarWriter)
		if err != nil {
			log.Errorf("compress dir %s %s error", prefix, currentFile.Name())
			return err
		}
	} else {
		err = compressFile(info, currentFile, prefix, tarWriter)
		if err != nil {
			log.Errorf("compress file %s %s error", prefix, currentFile.Name())
			return err
		}
	}
	return nil
}

func compressDir(info os.FileInfo, currentFile *os.File, prefix string, tarWriter *tar.Writer) error {
	prefix = fmt.Sprintf("%s/%s", prefix, info.Name())
	fileInfos, err := currentFile.Readdir(-1)
	if err != nil {
		log.Errorf("read dir -1 error: %v", err)
		return err
	}
	for _, fileInfo := range fileInfos {
		file, err := os.Open(fmt.Sprintf("%s/%s", currentFile.Name(), fileInfo.Name()))
		if err != nil {
			log.Errorf("open sub file error: %v", err)
			return err
		}
		err = compress(file, prefix, tarWriter)
		if err != nil {
			log.Errorf("compress %s %s error", prefix, file.Name())
			return err
		}
	}
	return nil
}

func compressFile(info os.FileInfo, currentFile *os.File, prefix string, tarWriter *tar.Writer) error {
	header, err := tar.FileInfoHeader(info, "")
	if err != nil {
		log.Errorf("get tar file header error: %v", err)
		return err
	}
	header.Name = fmt.Sprintf("%s/%s", prefix, header.Name)
	err = tarWriter.WriteHeader(header)
	if err != nil {
		log.Errorf("write tar file header error: %v", err)
		return err
	}
	_, err = io.Copy(tarWriter, currentFile)
	if err != nil {
		log.Errorf("write tar file error: %v", err)
		return err
	}
	err = currentFile.Close()
	if err != nil {
		log.Errorf("tar file close error: %v", err)
		return err
	}
	return nil
}
