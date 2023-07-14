package sync

import (
	"errors"
	"fmt"
	"github.com/catfishlty/sync-up/internal/docker"
	"github.com/catfishlty/sync-up/internal/utils"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
)

type Syncer struct {
	id                   string
	d                    *docker.Docker
	r                    *RClone
	backupper            Backupper
	localDir             string
	rCloneDir            string
	rCloneBandwidthLimit string
	rCloneRemoteName     string
	rCloneRemotePath     string
	cron                 string
}

type Info struct {
	Id       string
	Cron     string
	Backuper BackupperInfo
	RClone   RCloneInfo
}

type BackupperInfo struct {
	Type                  string
	ContainerName         string
	ContainerDumpFilePath string
	Cmd                   string
}

type RCloneInfo struct {
	ContainerName string
	Bandwidth     string
	RemoteName    string
	RemotePath    string
}

type SyncerOptions struct {
	Id        string
	RClone    *RClone
	Backupper Backupper
	LocalDir  string
	Cron      string
}

type RCloneOptions struct {
	Dir            string
	BandwidthLimit string
	RemoteName     string
	RemotePath     string
}

func NewSyncer(d *docker.Docker, syncerOptions SyncerOptions, rCloneOptions RCloneOptions) *Syncer {
	return &Syncer{
		id:                   syncerOptions.Id,
		d:                    d,
		r:                    syncerOptions.RClone,
		backupper:            syncerOptions.Backupper,
		localDir:             syncerOptions.LocalDir,
		rCloneDir:            rCloneOptions.Dir,
		rCloneBandwidthLimit: rCloneOptions.BandwidthLimit,
		rCloneRemoteName:     rCloneOptions.RemoteName,
		rCloneRemotePath:     rCloneOptions.RemotePath,
		cron:                 syncerOptions.Cron,
	}
}

func (s *Syncer) Info() Info {
	return Info{
		Id:   s.id,
		Cron: s.cron,
		Backuper: BackupperInfo{
			Type:                  s.backupper.Type(),
			ContainerName:         s.backupper.ContainerName(),
			ContainerDumpFilePath: s.backupper.ContainerDumpFilePath(),
			Cmd:                   strings.Join(s.backupper.BackupCommand(), " "),
		},
		RClone: RCloneInfo{
			ContainerName: s.r.containerName,
			Bandwidth:     s.rCloneBandwidthLimit,
			RemotePath:    s.rCloneRemotePath,
			RemoteName:    s.rCloneRemoteName,
		},
	}
}

func (s *Syncer) backup() (string, error) {
	container, err := s.d.GetContainerByName(s.backupper.ContainerName())
	if err != nil {
		return "", err
	}
	result, err := s.d.RunCommandSync(container, s.backupper.BackupCommand())
	if err != nil {
		return "", err
	}
	if s.backupper.CheckResultError(result) {
		log.Errorf("[%s - Backup] error: %v", s.backupper.Type(), result)
		return "", errors.New(result)
	}
	containerDumpFilePath := s.backupper.ContainerDumpFilePath()
	log.Infof("[Debug] containerDumpFilePath=%s", containerDumpFilePath)
	dumpFileName := filepath.Base(containerDumpFilePath)
	tempFilePath := fmt.Sprintf("%s/%s", s.localDir, dumpFileName)
	err = s.d.GetFileFromContainer(container, containerDumpFilePath, tempFilePath)
	if err != nil {
		return "", err
	}
	if s.backupper.ContainerDumpFileAutoRemove() {
		_, err = s.d.RunCommandSync(container, []string{
			"rm",
			"-f",
			containerDumpFilePath,
		})
		if err != nil {
			log.Errorf("[%s - Backup] remove container file error,'%s' : %v", s.backupper.Type(), containerDumpFilePath, err)
			return "", err
		}
	}
	compressFilePath, err := utils.CompressToDir(tempFilePath, s.localDir)
	if err != nil {
		log.Errorf("[%s - Backup] compress file error,'%s' -> '%s': %v", s.backupper.Type(), tempFilePath, s.localDir, err)
		return "", err
	}
	err = os.Remove(tempFilePath)
	if err != nil {
		log.Errorf("[%s - Backup] remove sync temp file error: %v", s.backupper.Type(), err)
		return "", err
	}
	return compressFilePath, nil
}

func (s *Syncer) rClone(filePath, rCloneBaseDir, rCloneBandwidthLimit, rCloneRemoteName, rCloneRemotePath string) error {
	rCloneFilePath, err := s.r.PushFile(RClonePushOptions{
		FilePath: filePath,
		TempDir:  rCloneBaseDir,
	})
	if err != nil {
		return err
	}
	err = s.r.Sync(RCloneSyncOptions{
		FilePath:       rCloneFilePath,
		Progress:       true,
		BandwidthLimit: rCloneBandwidthLimit,
		RemoteName:     rCloneRemoteName,
		RemotePath:     rCloneRemotePath,
	})
	return err
}

func (s *Syncer) Sync() error {
	dumpFilePath, err := s.backup()
	if err != nil {
		return err
	}
	return s.rClone(dumpFilePath, s.rCloneDir, s.rCloneBandwidthLimit, s.rCloneRemoteName, s.rCloneRemotePath)
}

func (s *Syncer) CheckVersion() (string, error) {
	container, err := s.d.GetContainerByName(s.backupper.ContainerName())
	if err != nil {
		return "", err
	}
	return s.d.RunCommandSync(container, s.backupper.CheckVersionCommands())
}
