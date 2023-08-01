package sync

import (
	"fmt"
	"path/filepath"

	"github.com/cloudfogtech/sync-up/internal/docker"
	log "github.com/sirupsen/logrus"
)

type RCloneSyncOptions struct {
	FilePath       string
	Progress       bool
	BandwidthLimit string
	RemoteName     string
	RemotePath     string
}

type RClonePushOptions struct {
	TempDir  string
	FilePath string
}

type RClone struct {
	d             *docker.Docker
	containerName string
}

func NewRClone(d *docker.Docker, containerName string) *RClone {
	return &RClone{
		d,
		containerName,
	}
}

func (r *RClone) PushFile(options RClonePushOptions) (string, error) {
	container, err := r.d.GetContainerByName(r.containerName)
	if err != nil {
		return "", err
	}
	err = r.d.SendFileToContainer(container, options.FilePath, options.TempDir)
	if err != nil {
		return "", err
	}
	fileName := filepath.Base(options.FilePath)
	return fmt.Sprintf("%s/%s", options.TempDir, fileName), nil
}

func (r *RClone) Sync(options RCloneSyncOptions) error {
	container, err := r.d.GetContainerByName(r.containerName)
	if err != nil {
		return err
	}
	commands := []string{
		"rclone",
		"sync",
		options.FilePath,
		fmt.Sprintf("%s:%s", options.RemoteName, options.RemotePath),
	}
	if options.Progress {
		commands = append(commands, "-P")
	}
	if options.BandwidthLimit != "" {
		// TODO validate BandwidthLimit
		commands = append(commands, "--bwlimit", options.BandwidthLimit)
	}
	result, err := r.d.RunCommandSync(container, commands)
	if err != nil {
		log.Errorf("RClone sync error: %v", err)
		return err
	}
	log.Infof("RClone Sync Local:'%s' to Remote:'%s:%s' success: %s", options.FilePath, options.RemoteName, options.RemotePath, result)
	_, err = r.d.RunCommandSync(container, []string{
		"rm",
		"-f",
		options.FilePath,
	})
	if err != nil {
		log.Errorf("RClone sync file delete error: %v", err)
		return err
	}
	log.Debugf("RClone Sync delete local file:'%s' success", options.FilePath)

	return nil
}
