package sync

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func GetMockRedisRDBContainerBackupper() Backupper {
	return NewRedisContainerRDBBackupper(RedisRDBContainerBackupperOptions{
		ContainerName: "redis",
		Password:      "",
		DumpFilePath:  "/data/dump.rdb",
	})
}

func GetMockRedisRDBContainerBackupperWithPass() Backupper {
	return NewRedisContainerRDBBackupper(RedisRDBContainerBackupperOptions{
		ContainerName: "redis",
		Password:      "123456",
		DumpFilePath:  "/data/dump.rdb",
	})
}

func Test_RedisRDBContainerBackupper_Type(t *testing.T) {
	Convey("TestRedisRDBContainerBackupper_Type", t, func() {
		backupper := GetMockRedisRDBContainerBackupper()
		So(backupper.Type(), ShouldEqual, "RedisRDBContainerBackupper")
	})
}

func Test_RedisRDBContainerBackupper_ContainerName(t *testing.T) {
	Convey("TestRedisRDBContainerBackupper_ContainerName", t, func() {
		backupper := GetMockRedisRDBContainerBackupper()
		So(backupper.ContainerName(), ShouldEqual, "redis")
	})
}
func Test_RedisRDBContainerBackupper_BackupCommand(t *testing.T) {
	Convey("Test_RedisRDBContainerBackupper_BackupCommand", t, func() {
		Convey("no pass", func() {
			backupper := GetMockRedisRDBContainerBackupper()
			expected := []string{
				"redis-cli",
				"save",
			}
			actual := backupper.BackupCommand()
			So(len(actual), ShouldEqual, len(expected))
			for i, s := range actual {
				So(s, ShouldEqual, expected[i])
			}
		})
		Convey("with pass", func() {
			backupper := GetMockRedisRDBContainerBackupperWithPass()
			expected := []string{
				"redis-cli",
				"-a",
				"123456",
				"save",
			}
			actual := backupper.BackupCommand()
			So(len(actual), ShouldEqual, len(expected))
			for i, s := range actual {
				So(s, ShouldEqual, expected[i])
			}
		})
	})
}
func Test_RedisRDBContainerBackupper_CheckResultError(t *testing.T) {
	Convey("Test_RedisRDBContainerBackupper_CheckResultError", t, func() {
		Convey("True", func() {
			backupper := GetMockRedisRDBContainerBackupper()
			So(backupper.CheckResultError("OK"), ShouldBeFalse)

		})
		Convey("False", func() {
			backupper := GetMockRedisRDBContainerBackupper()
			So(backupper.CheckResultError("ERROR"), ShouldBeTrue)
		})
	})
}

func Test_RedisRDBContainerBackupper_ContainerDumpFilePath(t *testing.T) {
	Convey("Test_RedisRDBContainerBackupper_ContainerDumpFilePath", t, func() {
		backupper := GetMockRedisRDBContainerBackupper()
		So(backupper.ContainerDumpFilePath(), ShouldEqual, "/data/dump.rdb")
	})
}

func Test_RedisRDBContainerBackupper_ContainerDumpFileAutoRemove(t *testing.T) {
	Convey("Test_RedisRDBContainerBackupper_ContainerDumpFileAutoRemove", t, func() {
		backupper := GetMockRedisRDBContainerBackupper()
		So(backupper.ContainerDumpFileAutoRemove(), ShouldBeFalse)
	})
}

func Test_RedisRDBContainerBackupper_CheckVersionCommands(t *testing.T) {
	Convey("Test_RedisRDBContainerBackupper_CheckVersionCommands", t, func() {
		backupper := GetMockRedisRDBContainerBackupper()
		expected := []string{
			"redis-cli",
			"-v",
		}
		actual := backupper.CheckVersionCommands()
		So(len(actual), ShouldEqual, len(expected))
		for i, s := range actual {
			So(s, ShouldEqual, expected[i])
		}
	})
}
