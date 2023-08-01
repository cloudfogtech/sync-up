package sync

import (
	"strings"
)

type RedisRDBContainerBackupper struct {
	BackupperOptions
	Password    string
	RDBFilePath string
}

type RedisRDBContainerBackupperOptions struct {
	ContainerName string
	Password      string
	DumpFilePath  string
}

func NewRedisContainerRDBBackupper(options RedisRDBContainerBackupperOptions) Backupper {
	return &RedisRDBContainerBackupper{
		BackupperOptions{
			ContainerName: options.ContainerName,
		},
		options.Password,
		options.DumpFilePath,
	}
}

func (r *RedisRDBContainerBackupper) Type() string {
	return "RedisRDBContainerBackupper"
}

func (r *RedisRDBContainerBackupper) ContainerName() string {
	return r.BackupperOptions.ContainerName
}

func (r *RedisRDBContainerBackupper) BackupCommand() []string {
	commands := []string{"redis-cli"}
	if r.Password != "" {
		commands = append(commands, "-a")
		commands = append(commands, r.Password)
	}
	commands = append(commands, "save")
	return commands
}

func (r *RedisRDBContainerBackupper) CheckResultError(result string) bool {
	return !strings.Contains(result, "OK")
}

func (r *RedisRDBContainerBackupper) ContainerDumpFilePath() string {
	return r.RDBFilePath
}

func (r *RedisRDBContainerBackupper) ContainerDumpFileAutoRemove() bool {
	return false
}

func (r *RedisRDBContainerBackupper) CheckVersionCommands() []string {
	return []string{
		"redis-cli",
		"-v",
	}
}
