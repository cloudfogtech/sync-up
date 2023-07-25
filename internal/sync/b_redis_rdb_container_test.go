package sync

import (
	"testing"

	"github.com/cloudfogtech/sync-up/internal/docker"
	log "github.com/sirupsen/logrus"
)

func TestRedisRDB_Backup(t *testing.T) {
	var redisRDB Backupper
	d, err := docker.NewDockerCli()
	if err != nil {
		panic(err)
	}
	redisRDB = NewRedisContainerRDBBackupper(RedisRDBContainerBackupperOptions{
		ContainerName: "juicefs-redis",
		Password:      "E8sJCHto2VNYZzin",
		DumpFilePath:  "/data/dump.rdb",
	})
	if err != nil {
		panic(err)
	}
	syncer := NewSyncer(d, SyncerOptions{Id: "test", Backupper: redisRDB}, RCloneOptions{})
	path, err := syncer.backup()
	if err != nil {
		panic(err)
	}
	log.Infof("result=%s", path)
}

func TestRedisRDB_CheckVersion(t *testing.T) {
	var redisRDB Backupper
	d, err := docker.NewDockerCli()
	if err != nil {
		panic(err)
	}
	redisRDB = NewRedisContainerRDBBackupper(RedisRDBContainerBackupperOptions{
		ContainerName: "juicefs-redis",
		Password:      "E8sJCHto2VNYZzin",
		DumpFilePath:  "/data/dump.rdb",
	})
	if err != nil {
		panic(err)
	}
	syncer := NewSyncer(d, SyncerOptions{Id: "test", Backupper: redisRDB}, RCloneOptions{})
	path, err := syncer.CheckVersion()
	if err != nil {
		panic(err)
	}
	log.Infof("result=%s", path)
}
