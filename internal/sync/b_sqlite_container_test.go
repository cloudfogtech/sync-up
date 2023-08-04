package sync

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func GetMockSqliteContainerBackupper() Backupper {
	return NewSqliteContainer(SqliteContainerBackupperOptions{
		ContainerName: "sqlite",
		DbFilePath:    "/data/data.db",
	})
}

func TestSqliteContainerBackupper_Type(t *testing.T) {
	Convey("TestSqliteContainerBackupper_Type", t, func() {
		backupper := GetMockSqliteContainerBackupper()
		So(backupper.Type(), ShouldEqual, "SqliteContainerBackupper")
	})
}

func TestSqliteContainerBackupper_ContainerName(t *testing.T) {
	Convey("TestSqliteContainerBackupper_ContainerName", t, func() {
		backupper := GetMockSqliteContainerBackupper()
		So(backupper.ContainerName(), ShouldEqual, "sqlite")
	})
}

func TestSqliteContainerBackupper_BackupCommand(t *testing.T) {
	Convey("TestSqliteContainerBackupper_BackupCommand", t, func() {
		backupper := GetMockSqliteContainerBackupper()
		var expected []string
		actual := backupper.BackupCommand()
		So(len(actual), ShouldEqual, len(expected))
		for i, s := range actual {
			So(s, ShouldEqual, expected[i])
		}
	})
}

func TestSqliteContainerBackupper_CheckResultError(t *testing.T) {
	Convey("TestSqliteContainerBackupper_CheckResultError", t, func() {
		backupper := GetMockSqliteContainerBackupper()
		So(backupper.CheckResultError("test"), ShouldBeFalse)
	})
}

func TestSqliteContainerBackupper_ContainerDumpFilePath(t *testing.T) {
	Convey("TestSqliteContainerBackupper_ContainerDumpFilePath", t, func() {
		backupper := GetMockSqliteContainerBackupper()
		So(backupper.ContainerDumpFilePath(), ShouldEqual, "/data/data.db")
	})
}

func TestSqliteContainerBackupper_ContainerDumpFileAutoRemove(t *testing.T) {
	Convey("TestSqliteContainerBackupper_ContainerDumpFileAutoRemove", t, func() {
		backupper := GetMockSqliteContainerBackupper()
		So(backupper.ContainerDumpFileAutoRemove(), ShouldBeFalse)
	})
}

func TestSqliteContainerBackupper_CheckVersionCommands(t *testing.T) {
	Convey("TestSqliteContainerBackupper_CheckVersionCommands", t, func() {
		backupper := GetMockSqliteContainerBackupper()
		expected := []string{
			"file",
			"/data/data.db",
		}
		actual := backupper.CheckVersionCommands()
		So(len(actual), ShouldEqual, len(expected))
		for i, s := range actual {
			So(s, ShouldEqual, expected[i])
		}
	})
}
