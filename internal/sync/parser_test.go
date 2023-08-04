package sync

import (
	. "github.com/agiledragon/gomonkey/v2"
	"github.com/cloudfogtech/sync-up/internal/env"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestNewParserManager(t *testing.T) {
	Convey("TestNewParserManager", t, func() {
		envManager := &env.Manager{}
		parserManager := NewParserManager(envManager)
		So(parserManager, ShouldNotBeNil)
		So(parserManager.manager, ShouldNotBeNil)
		So(parserManager.parser, ShouldNotBeNil)
		So(len(parserManager.parser), ShouldEqual, 3)
	})
}

func TestParserManager_NewBackupper(t *testing.T) {
	envManager := &env.Manager{}
	parserManager := NewParserManager(envManager)
	Convey("TestParserManager_NewBackupper", t, func() {
		Convey("error", func() {
			patches := ApplyMethod(envManager, "GetEnvWithNilCheck", func(manager *env.Manager, tpl, id string) string {
				return "test"
			})
			defer patches.Reset()
			backupper, err := parserManager.NewBackupper("test", "test")
			So(backupper, ShouldBeNil)
			So(err, ShouldBeError)
		})
		Convey("success", func() {
			patches := ApplyMethod(envManager, "GetEnvWithNilCheck", func(manager *env.Manager, tpl, id string) string {
				return "test"
			})
			patches.ApplyFunc(newSqliteContainer, func(m *env.Manager, id, containerName string) Backupper {
				return &SqliteContainerBackupper{}
			})
			defer patches.Reset()
			backupper, err := parserManager.NewBackupper("sqlite-c", "test")
			So(backupper, ShouldNotBeNil)
			So(err, ShouldBeNil)
		})
	})
}

func TestParserManager_newRedisRDBContainer(t *testing.T) {
	Convey("TestParserManager_newRedisRDBContainer", t, func() {
		envManager := &env.Manager{}
		patches := ApplyFunc(NewRedisContainerRDBBackupper, func(options RedisRDBContainerBackupperOptions) Backupper {
			return &RedisRDBContainerBackupper{}
		})
		patches.ApplyMethod(envManager, "GetEnv", func(_ *env.Manager, _, _ string) string {
			return "test"
		})
		patches.ApplyMethod(envManager, "GetEnvWithDefault", func(_ *env.Manager, _, _, _ string) string {
			return "test"
		})
		defer patches.Reset()
		backupper := newRedisRDBContainer(envManager, "test", "test")
		So(backupper, ShouldNotBeNil)
	})
}

func TestParserManager_newPostgresContainer(t *testing.T) {
	Convey("TestParserManager_newPostgresContainer", t, func() {
		envManager := &env.Manager{}
		patches := ApplyFunc(NewPostgresContainer, func(options PostgresContainerBackupperOptions) Backupper {
			return &PostgresContainerBackupper{}
		})
		patches.ApplyMethod(envManager, "GetEnvWithNilCheck", func(_ *env.Manager, _, _ string) string {
			return "test"
		})
		defer patches.Reset()
		backupper := newPostgresContainer(envManager, "test", "test")
		So(backupper, ShouldNotBeNil)
	})
}

func TestParserManager_newSqliteContainer(t *testing.T) {
	Convey("TestParserManager_newSqliteContainer", t, func() {
		envManager := &env.Manager{}
		patches := ApplyFunc(NewSqliteContainer, func(options SqliteContainerBackupperOptions) Backupper {
			return &SqliteContainerBackupper{}
		})
		patches.ApplyMethod(envManager, "GetEnvWithNilCheck", func(_ *env.Manager, _, _ string) string {
			return "test"
		})
		defer patches.Reset()
		backupper := newSqliteContainer(envManager, "test", "test")
		So(backupper, ShouldNotBeNil)
	})
}
