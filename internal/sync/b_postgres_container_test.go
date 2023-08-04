package sync

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func GetPostgresContainerBackupper() Backupper {
	return NewPostgresContainer(PostgresContainerBackupperOptions{
		ContainerName: "postgres",
		Username:      "test",
	})
}

func GetPostgresContainerBackupperWithoutUsername() Backupper {
	return NewPostgresContainer(PostgresContainerBackupperOptions{
		ContainerName: "postgres",
		Username:      "",
	})
}

func TestPostgresContainerBackupper_Type(t *testing.T) {
	Convey("TestPostgresContainerBackupper_Type", t, func() {
		backupper := GetPostgresContainerBackupper()
		So(backupper.Type(), ShouldEqual, "PostgresContainerBackupper")
	})
}

func TestPostgresContainerBackupper_ContainerName(t *testing.T) {
	Convey("TestPostgresContainerBackupper_ContainerName", t, func() {
		backupper := GetPostgresContainerBackupper()
		So(backupper.ContainerName(), ShouldEqual, "postgres")
	})
}

func TestPostgresContainerBackupper_BackupCommand(t *testing.T) {
	Convey("TestPostgresContainerBackupper_BackupCommand", t, func() {
		Convey("no username", func() {
			backupper := GetPostgresContainerBackupperWithoutUsername()
			expected := []string{
				"pg_dumpall",
				"-f",
				PostgresDumpFilePath,
			}
			actual := backupper.BackupCommand()
			So(len(actual), ShouldEqual, len(expected))
			for i, s := range actual {
				So(s, ShouldEqual, expected[i])
			}
		})
		Convey("with username", func() {
			backupper := GetPostgresContainerBackupper()
			expected := []string{
				"pg_dumpall",
				"-U",
				"test",
				"-f",
				PostgresDumpFilePath,
			}
			actual := backupper.BackupCommand()
			So(len(actual), ShouldEqual, len(expected))
			for i, s := range actual {
				So(s, ShouldEqual, expected[i])
			}
		})
	})
}

func TestPostgresContainerBackupper_CheckResultError(t *testing.T) {
	Convey("TestPostgresContainerBackupper_CheckResultError", t, func() {
		backupper := GetPostgresContainerBackupper()
		So(backupper.CheckResultError("OK"), ShouldBeFalse)
	})
}

func TestPostgresContainerBackupper_ContainerDumpFilePath(t *testing.T) {
	Convey("TestPostgresContainerBackupper_ContainerDumpFilePath", t, func() {
		backupper := GetPostgresContainerBackupper()
		So(backupper.ContainerDumpFilePath(), ShouldEqual, PostgresDumpFilePath)
	})
}

func TestPostgresContainerBackupper_ContainerDumpFileAutoRemove(t *testing.T) {
	Convey("TestPostgresContainerBackupper_ContainerDumpFileAutoRemove", t, func() {
		backupper := GetPostgresContainerBackupper()
		So(backupper.ContainerDumpFileAutoRemove(), ShouldBeTrue)
	})
}

func TestPostgresContainerBackupper_CheckVersionCommands(t *testing.T) {
	Convey("TestPostgresContainerBackupper_CheckVersionCommands", t, func() {
		backupper := GetPostgresContainerBackupper()
		expected := []string{
			"pg_dumpall",
			"-V",
		}
		actual := backupper.CheckVersionCommands()
		So(len(actual), ShouldEqual, len(expected))
		for i, s := range actual {
			So(s, ShouldEqual, expected[i])
		}
	})
}
