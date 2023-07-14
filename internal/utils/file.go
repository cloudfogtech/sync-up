package utils

import (
	"io"
	"os"

	log "github.com/sirupsen/logrus"
)

func MoveFile(src, dest string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		log.Errorf("open src file: %s error: %v", src, err)
		return err
	}
	destFile, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o755)
	if err != nil {
		log.Errorf("open dest file: %s error: %v", dest, err)
		return err
	}
	length, err := io.Copy(destFile, srcFile)
	if err != nil {
		log.Errorf("file copy error, src=%s, dest=%s error: %v", src, dest, err)
		return err
	}
	log.Infof("file copy from src=%s to dest=%s size: %d", src, dest, length)
	err = os.Remove(src)
	if err != nil {
		log.Errorf("remove src=%s file error, %v", src, err)
		return err
	}
	return nil
}
