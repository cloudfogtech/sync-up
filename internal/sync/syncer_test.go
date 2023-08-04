package sync

import (
	"errors"
	. "github.com/agiledragon/gomonkey/v2"
	"github.com/cloudfogtech/sync-up/internal/docker"
	"github.com/cloudfogtech/sync-up/internal/utils"
	"github.com/docker/docker/api/types"
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"path/filepath"
	"testing"
)

func TestNewSyncer(t *testing.T) {
	Convey("TestNewSyncer", t, func() {
		s := NewSyncer(&docker.Docker{}, SyncerOptions{}, RCloneOptions{})
		So(s, ShouldNotBeNil)
	})
}

func TestSyncer_Info(t *testing.T) {
	Convey("TestSyncer_Info", t, func() {
		s := NewSyncer(&docker.Docker{}, SyncerOptions{
			RClone:    &RClone{},
			Backupper: &SqliteContainerBackupper{},
		}, RCloneOptions{})
		So(s.Info(), ShouldNotBeNil)
	})
}

func TestSyncer_backup(t *testing.T) {
	d := &docker.Docker{}
	b := &SqliteContainerBackupper{}
	s := NewSyncer(d, SyncerOptions{
		RClone:    &RClone{},
		Backupper: b,
	}, RCloneOptions{})
	Convey("TestSyncer_Info", t, func() {
		Convey("GetContainerByName error", func() {
			patches := ApplyMethod(d, "GetContainerByName", func(_ *docker.Docker, name string) (types.Container, error) {
				return types.Container{}, errors.New("GetContainerByName")
			})
			defer patches.Reset()
			_, err := s.backup()
			So(err, ShouldBeError)
			So(err.Error(), ShouldEqual, "GetContainerByName")
		})
		Convey("RunCommandSync error1", func() {
			patches := ApplyMethod(d, "GetContainerByName", func(_ *docker.Docker, name string) (types.Container, error) {
				return types.Container{}, nil
			})
			patches.ApplyMethodSeq(d, "RunCommandSync", []OutputCell{
				{Values: Params{"", errors.New("RunCommandSync")}},
			})
			defer patches.Reset()
			_, err := s.backup()
			So(err, ShouldBeError)
			So(err.Error(), ShouldEqual, "RunCommandSync")
		})
		Convey("CheckResultError error", func() {
			patches := ApplyMethod(d, "GetContainerByName", func(_ *docker.Docker, name string) (types.Container, error) {
				return types.Container{}, nil
			})
			patches.ApplyMethodSeq(d, "RunCommandSync", []OutputCell{
				{Values: Params{"CheckResultError", nil}},
			})
			patches.ApplyMethodSeq(b, "CheckResultError", []OutputCell{
				{Values: Params{true}},
			})
			defer patches.Reset()
			_, err := s.backup()
			So(err, ShouldBeError)
			So(err.Error(), ShouldEqual, "CheckResultError")
		})
		Convey("GetFileFromContainer error", func() {
			patches := ApplyMethod(d, "GetContainerByName", func(_ *docker.Docker, name string) (types.Container, error) {
				return types.Container{}, nil
			})
			patches.ApplyMethodSeq(d, "RunCommandSync", []OutputCell{
				{Values: Params{"OK", nil}},
			})
			patches.ApplyMethodSeq(b, "CheckResultError", []OutputCell{
				{Values: Params{false}},
			})
			patches.ApplyMethodSeq(b, "ContainerDumpFilePath", []OutputCell{
				{Values: Params{"test"}},
			})
			patches.ApplyFunc(filepath.Base, func(_ string) string {
				return "test"
			})
			patches.ApplyMethodSeq(d, "GetFileFromContainer", []OutputCell{
				{Values: Params{errors.New("GetFileFromContainer")}},
			})
			defer patches.Reset()
			_, err := s.backup()
			So(err, ShouldBeError)
			So(err.Error(), ShouldEqual, "GetFileFromContainer")
		})
		Convey("CompressToDir error", func() {
			patches := ApplyMethod(d, "GetContainerByName", func(_ *docker.Docker, name string) (types.Container, error) {
				return types.Container{}, nil
			})
			patches.ApplyMethodSeq(d, "RunCommandSync", []OutputCell{
				{Values: Params{"OK", nil}},
			})
			patches.ApplyMethodSeq(b, "CheckResultError", []OutputCell{
				{Values: Params{false}},
			})
			patches.ApplyMethodSeq(b, "ContainerDumpFilePath", []OutputCell{
				{Values: Params{"test"}},
			})
			patches.ApplyFunc(filepath.Base, func(_ string) string {
				return "test"
			})
			patches.ApplyMethodSeq(d, "GetFileFromContainer", []OutputCell{
				{Values: Params{nil}},
			})
			patches.ApplyMethodSeq(b, "ContainerDumpFileAutoRemove", []OutputCell{
				{Values: Params{false}},
			})
			patches.ApplyFuncSeq(utils.CompressToDir, []OutputCell{
				{Values: Params{"", errors.New("CompressToDir")}},
			})
			defer patches.Reset()
			_, err := s.backup()
			So(err, ShouldBeError)
			So(err.Error(), ShouldEqual, "CompressToDir")
		})
		Convey("os.Remove error", func() {
			patches := ApplyMethod(d, "GetContainerByName", func(_ *docker.Docker, name string) (types.Container, error) {
				return types.Container{}, nil
			})
			patches.ApplyMethodSeq(d, "RunCommandSync", []OutputCell{
				{Values: Params{"OK", nil}},
			})
			patches.ApplyMethodSeq(b, "CheckResultError", []OutputCell{
				{Values: Params{false}},
			})
			patches.ApplyMethodSeq(b, "ContainerDumpFilePath", []OutputCell{
				{Values: Params{"test"}},
			})
			patches.ApplyFunc(filepath.Base, func(_ string) string {
				return "test"
			})
			patches.ApplyMethodSeq(d, "GetFileFromContainer", []OutputCell{
				{Values: Params{nil}},
			})
			patches.ApplyMethodSeq(b, "ContainerDumpFileAutoRemove", []OutputCell{
				{Values: Params{false}},
			})
			patches.ApplyFuncSeq(utils.CompressToDir, []OutputCell{
				{Values: Params{"", nil}},
			})
			patches.ApplyFuncSeq(os.Remove, []OutputCell{
				{Values: Params{errors.New("os.Remove")}},
			})
			defer patches.Reset()
			_, err := s.backup()
			So(err, ShouldBeError)
			So(err.Error(), ShouldEqual, "os.Remove")
		})
		Convey("success condition 1", func() {
			patches := ApplyMethod(d, "GetContainerByName", func(_ *docker.Docker, name string) (types.Container, error) {
				return types.Container{}, nil
			})
			patches.ApplyMethodSeq(d, "RunCommandSync", []OutputCell{
				{Values: Params{"OK", nil}},
			})
			patches.ApplyMethodSeq(b, "CheckResultError", []OutputCell{
				{Values: Params{false}},
			})
			patches.ApplyMethodSeq(b, "ContainerDumpFilePath", []OutputCell{
				{Values: Params{"test"}},
			})
			patches.ApplyFunc(filepath.Base, func(_ string) string {
				return "test"
			})
			patches.ApplyMethodSeq(d, "GetFileFromContainer", []OutputCell{
				{Values: Params{nil}},
			})
			patches.ApplyMethodSeq(b, "ContainerDumpFileAutoRemove", []OutputCell{
				{Values: Params{false}},
			})
			patches.ApplyFuncSeq(utils.CompressToDir, []OutputCell{
				{Values: Params{"test", nil}},
			})
			patches.ApplyFuncSeq(os.Remove, []OutputCell{
				{Values: Params{nil}},
			})
			defer patches.Reset()
			result, err := s.backup()
			So(err, ShouldBeNil)
			So(result, ShouldEqual, "test")
		})
		Convey("success condition 2", func() {
			patches := ApplyMethod(d, "GetContainerByName", func(_ *docker.Docker, name string) (types.Container, error) {
				return types.Container{}, nil
			})
			patches.ApplyMethodSeq(d, "RunCommandSync", []OutputCell{
				{Values: Params{"OK", nil}},
				{Values: Params{"OK", nil}},
			})
			patches.ApplyMethodSeq(b, "CheckResultError", []OutputCell{
				{Values: Params{false}},
			})
			patches.ApplyMethodSeq(b, "ContainerDumpFilePath", []OutputCell{
				{Values: Params{"test"}},
			})
			patches.ApplyFunc(filepath.Base, func(_ string) string {
				return "test"
			})
			patches.ApplyMethodSeq(d, "GetFileFromContainer", []OutputCell{
				{Values: Params{nil}},
			})
			patches.ApplyMethodSeq(b, "ContainerDumpFileAutoRemove", []OutputCell{
				{Values: Params{true}},
			})
			patches.ApplyFuncSeq(utils.CompressToDir, []OutputCell{
				{Values: Params{"test", nil}},
			})
			patches.ApplyFuncSeq(os.Remove, []OutputCell{
				{Values: Params{nil}},
			})
			defer patches.Reset()
			result, err := s.backup()
			So(err, ShouldBeNil)
			So(result, ShouldEqual, "test")
		})
		Convey("RunCommandSync error 2", func() {
			patches := ApplyMethod(d, "GetContainerByName", func(_ *docker.Docker, name string) (types.Container, error) {
				return types.Container{}, nil
			})
			patches.ApplyMethodSeq(d, "RunCommandSync", []OutputCell{
				{Values: Params{"OK", nil}},
				{Values: Params{"OK", errors.New("RunCommandSync")}},
			})
			patches.ApplyMethodSeq(b, "CheckResultError", []OutputCell{
				{Values: Params{false}},
			})
			patches.ApplyMethodSeq(b, "ContainerDumpFilePath", []OutputCell{
				{Values: Params{"test"}},
			})
			patches.ApplyFunc(filepath.Base, func(_ string) string {
				return "test"
			})
			patches.ApplyMethodSeq(d, "GetFileFromContainer", []OutputCell{
				{Values: Params{nil}},
			})
			patches.ApplyMethodSeq(b, "ContainerDumpFileAutoRemove", []OutputCell{
				{Values: Params{true}},
			})
			defer patches.Reset()
			_, err := s.backup()
			So(err, ShouldBeError)
			So(err.Error(), ShouldEqual, "RunCommandSync")
		})
	})
}

func TestSyncer_rClone(t *testing.T) {
	r := &RClone{}
	s := NewSyncer(&docker.Docker{}, SyncerOptions{
		RClone:    r,
		Backupper: &SqliteContainerBackupper{},
	}, RCloneOptions{})
	Convey("TestSyncer_rClone", t, func() {
		Convey("PushFile error", func() {
			patches := ApplyMethodSeq(r, "PushFile", []OutputCell{
				{Values: Params{"", errors.New("PushFile")}},
			})
			defer patches.Reset()
			err := s.rClone("", "", "", "", "")
			So(err, ShouldBeError)
			So(err.Error(), ShouldEqual, "PushFile")
		})
		Convey("Sync error", func() {
			patches := ApplyMethodSeq(r, "PushFile", []OutputCell{
				{Values: Params{"", nil}},
			})
			patches.ApplyMethodSeq(r, "Sync", []OutputCell{
				{Values: Params{errors.New("Sync")}},
			})
			defer patches.Reset()
			err := s.rClone("", "", "", "", "")
			So(err, ShouldBeError)
			So(err.Error(), ShouldEqual, "Sync")
		})
		Convey("success", func() {
			patches := ApplyMethodSeq(r, "PushFile", []OutputCell{
				{Values: Params{"", nil}},
			})
			patches.ApplyMethodSeq(r, "Sync", []OutputCell{
				{Values: Params{nil}},
			})
			defer patches.Reset()
			err := s.rClone("", "", "", "", "")
			So(err, ShouldBeNil)
		})
	})
}

func TestSyncer_Sync(t *testing.T) {
	s := NewSyncer(&docker.Docker{}, SyncerOptions{
		RClone:    &RClone{},
		Backupper: &SqliteContainerBackupper{},
	}, RCloneOptions{})
	Convey("TestSyncer_Sync", t, func() {
		Convey("backup error", func() {
			patches := ApplyPrivateMethod(s, "backup", func(_ *Syncer) (string, error) {
				return "", errors.New("backup")
			})
			defer patches.Reset()
			err := s.Sync()
			So(err, ShouldBeError)
			So(err.Error(), ShouldEqual, "backup")
		})
		Convey("rClone error", func() {
			patches := ApplyPrivateMethod(s, "backup", func(_ *Syncer) (string, error) {
				return "", nil
			})
			patches.ApplyPrivateMethod(s, "rClone", func(_ *Syncer, _, _, _, _, _ string) error {
				return errors.New("rClone")
			})
			defer patches.Reset()
			err := s.Sync()
			So(err, ShouldBeError)
			So(err.Error(), ShouldEqual, "rClone")
		})
		Convey("success", func() {
			patches := ApplyPrivateMethod(s, "backup", func(_ *Syncer) (string, error) {
				return "", nil
			})
			patches.ApplyPrivateMethod(s, "rClone", func(_ *Syncer, _, _, _, _, _ string) error {
				return nil
			})
			defer patches.Reset()
			err := s.Sync()
			So(err, ShouldBeNil)
		})
	})
}

func TestSyncer_CheckVersion(t *testing.T) {
	d := &docker.Docker{}
	s := NewSyncer(d, SyncerOptions{
		RClone:    &RClone{},
		Backupper: &SqliteContainerBackupper{},
	}, RCloneOptions{})
	Convey("TestSyncer_CheckVersion", t, func() {
		Convey("GetContainerByName error", func() {
			patches := ApplyMethodSeq(d, "GetContainerByName", []OutputCell{
				{Values: Params{types.Container{}, errors.New("GetContainerByName")}},
			})
			defer patches.Reset()
			_, err := s.CheckVersion()
			So(err, ShouldBeError)
			So(err.Error(), ShouldEqual, "GetContainerByName")
		})
		Convey("rClone error", func() {
			patches := ApplyMethodSeq(d, "GetContainerByName", []OutputCell{
				{Values: Params{types.Container{}, nil}},
			})
			patches.ApplyMethodSeq(d, "RunCommandSync", []OutputCell{
				{Values: Params{"", errors.New("RunCommandSync")}},
			})
			defer patches.Reset()
			_, err := s.CheckVersion()
			So(err, ShouldBeError)
			So(err.Error(), ShouldEqual, "RunCommandSync")
		})
		Convey("success", func() {
			patches := ApplyMethodSeq(d, "GetContainerByName", []OutputCell{
				{Values: Params{types.Container{}, nil}},
			})
			patches.ApplyMethodSeq(d, "RunCommandSync", []OutputCell{
				{Values: Params{"", nil}},
			})
			defer patches.Reset()
			_, err := s.CheckVersion()
			So(err, ShouldBeNil)
		})
	})
}
