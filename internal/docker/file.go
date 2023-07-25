package docker

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/cloudfogtech/sync-up/internal/common"
	"github.com/cloudfogtech/sync-up/internal/utils"
	"github.com/docker/docker/api/types"
	log "github.com/sirupsen/logrus"
)

func (d *Docker) GetFileFromContainer(container types.Container, src, dest string) error {
	reader, stat, err := d.cli.CopyFromContainer(context.Background(), container.ID, src)
	if err != nil {
		log.Errorf("CopyFromContainer error, container=%s, src='%s', %v", container.ID, src, err)
		return err
	}
	defer reader.Close()
	log.Debugf("CopyFromContainer, container=%s, src='%s'", container.ID, src)
	log.Infof("Container Path Stat, name=%s, size=%d", stat.Name, stat.Size)
	i := strings.LastIndex(dest, "/")
	if i != -1 {
		err = os.MkdirAll(dest[:i], 0o755)
		if err != nil {
			log.Errorf("Create Dest File Path,destDir=%s, dest=%s, %v", dest[:i], dest, err)
			return err
		}
	}
	f, err := os.OpenFile(dest, os.O_RDWR|os.O_CREATE, 0o755)
	if err != nil {
		log.Errorf("Open Dest File error, dest=%s, %v", dest, err)
		return err
	}
	defer f.Close()
	br := bufio.NewReader(reader)
	bw := bufio.NewWriter(f)
	_, err = io.Copy(bw, br)
	if err != nil {
		log.Errorf("Src to Dest File Data error, src=%s, dest=%s, %v", src, dest, err)
		return err
	}
	if dest, err = filepath.Abs(dest); err != nil {
		log.Errorf("can't find abs path for '%s'", dest)
	}
	log.Infof("File from Src to Dest, src=%s:%s, dest=local:%s", strings.Join(container.Names, ","), src, dest)
	return nil
}

func (d *Docker) SendFileToContainer(container types.Container, filePath, dest string) error {
	filePath, err := filepath.Abs(filePath)
	if err != nil {
		log.Errorf("can't get absolute path of '%s'", filePath)
	}
	err = os.MkdirAll(common.TempDirPath, os.ModeDir)
	if err != nil {
		log.Errorf("create temp file dir error: '%s', %v", common.TempDirPath, err)
		return err
	}
	archiveTempFilePath := fmt.Sprintf("%s/temp.tar.gz", common.TempDirPath)
	err = utils.CompressToFile(filePath, archiveTempFilePath)
	if err != nil {
		log.Errorf("create compress temp file error: '%s', %v", fmt.Sprintf("%s/temp.tar.gz", common.TempDirPath), err)
		return err
	}
	log.Debugf("compress file='%s' to temp compress file='%s'", filePath, archiveTempFilePath)
	file, err := os.Open(archiveTempFilePath)
	if err != nil {
		log.Errorf("open compress temp file error: '%s', %v", fmt.Sprintf("%s/temp.tar.gz", common.TempDirPath), err)
		return err
	}
	_, err = d.RunCommandSync(container, []string{
		"mkdir",
		"-p",
		dest,
	})
	if err != nil {
		log.Errorf("Crete Dest Dir in Container error, container=%s, dest='%s', %v", strings.Join(container.Names, ","), dest, err)
		return err
	}
	err = d.cli.CopyToContainer(context.Background(), container.ID, dest, file, types.CopyToContainerOptions{})
	if err != nil {
		log.Errorf("CopyToContainer error, container=%s, dest='%s', %v", strings.Join(container.Names, ","), dest, err)
		return err
	}
	log.Infof("File from src='%s->%s' to dest='%s:%s'", filePath, archiveTempFilePath, container.ID, dest)
	return nil
}
