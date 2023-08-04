package sync

import (
	"errors"
	"fmt"

	"github.com/cloudfogtech/sync-up/internal/common"
	"github.com/cloudfogtech/sync-up/internal/env"
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
			"sqlite-c":    newSqliteContainer,
		},
	}
}

func (bpm *ParserManager) NewBackupper(t, id string) (Backupper, error) {
	containerName := bpm.manager.GetEnvWithNilCheck(env.ContainerTpl, id)
	handler := bpm.parser[t]
	if handler == nil {
		return nil, errors.New(fmt.Sprintf("%s - %s create function is not found", t, id))
	}
	return handler(bpm.manager, id, containerName), nil
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

func newSqliteContainer(m *env.Manager, id, containerName string) Backupper {
	return NewSqliteContainer(SqliteContainerBackupperOptions{
		ContainerName: containerName,
		DbFilePath:    m.GetEnvWithNilCheck(env.SQLiteFilePathTpl, id),
	})
}
