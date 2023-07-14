package sync

import (
	"github.com/catfishlty/sync-up/internal/common"
	"github.com/catfishlty/sync-up/internal/env"
)

type Parser func(m *env.Manager, id string, t string) Backupper

type ParserManager struct {
	manager *env.Manager
	parser  map[string]Parser
}

func NewParserManager(m *env.Manager) *ParserManager {
	return &ParserManager{
		manager: m,
		parser: map[string]Parser{
			"redis-rdb-c": newRedisRDBContainer,
			"postgres-c":  newPostgresContainer,
		},
	}
}

func (bpm *ParserManager) NewBackupper(t, id string) Backupper {
	containerName := bpm.manager.GetEnvWithNilCheck(env.ContainerTpl, id)
	return bpm.parser[t](bpm.manager, id, containerName)
}

func (bpm *ParserManager) CheckType(t string) bool {
	return bpm.parser[t] != nil
}

func newRedisRDBContainer(m *env.Manager, id, containerName string) Backupper {
	return NewRedisContainerRDBBackupper(RedisRDBContainerBackupperOptions{
		ContainerName: containerName,
		Password:      m.GetEnv(env.ServicePasswordTpl, id),
		DumpFilePath:  m.GetEnvWithDefault(env.RedisRDBFilePathTpl, id, common.DefaultRedisRDBDumpPath),
	})
}

func newPostgresContainer(m *env.Manager, id, containerName string) Backupper {
	return NewPostgresContainer(PostgresContainerBackupperOptions{
		ContainerName: containerName,
		Username:      m.GetEnvWithNilCheck(env.ServiceUsernameTpl, id),
	})
}
