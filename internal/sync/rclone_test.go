package sync

import (
	"errors"
	"path/filepath"
	"testing"

	. "github.com/agiledragon/gomonkey/v2"
	"github.com/cloudfogtech/sync-up/internal/docker"
	"github.com/docker/docker/api/types"
	. "github.com/smartystreets/goconvey/convey"
)

func TestNewRClone(t *testing.T) {
	Convey("TestNewRClone", t, func() {
		rclone := NewRClone(&docker.Docker{}, "test")
		So(rclone, ShouldNotBeNil)
		So(rclone.containerName, ShouldEqual, "test")
	})
}

func TestRClone_PushFile(t *testing.T) {
	Convey("TestRClone_PushFile", t, func() {
		d := &docker.Docker{}
		rClone := NewRClone(d, "test")
		Convey("GetContainerByName error", func() {
			patches := ApplyMethod(d, "GetContainerByName", func(_ *docker.Docker, name string) (types.Container, error) {
				return types.Container{}, errors.New("test")
			})
			defer patches.Reset()
			_, err := rClone.PushFile(RClonePushOptions{})
			So(err, ShouldBeError)
		})
		Convey("SendFileToContainer error", func() {
			patches := ApplyMethod(d, "GetContainerByName", func(d *docker.Docker, name string) (types.Container, error) {
				return types.Container{}, nil
			})
			patches.ApplyMethod(d, "SendFileToContainer", func(d *docker.Docker, container types.Container, filePath, dest string) error {
				return errors.New("test")
			})
			defer patches.Reset()
			_, err := rClone.PushFile(RClonePushOptions{})
			So(err, ShouldBeError)
		})
		Convey("success", func() {
			patches := ApplyMethod(d, "GetContainerByName", func(d *docker.Docker, name string) (types.Container, error) {
				return types.Container{}, nil
			})
			patches.ApplyMethod(d, "SendFileToContainer", func(d *docker.Docker, container types.Container, filePath, dest string) error {
				return nil
			})
			patches.ApplyFunc(filepath.Base, func(path string) string {
				return "test"
			})
			defer patches.Reset()
			path, err := rClone.PushFile(RClonePushOptions{})
			So(err, ShouldBeNil)
			So(path, ShouldNotBeNil)
		})
	})
}

func TestRClone_Sync(t *testing.T) {
	Convey("TestRClone_Sync", t, func() {
		d := &docker.Docker{}
		rClone := NewRClone(d, "test")
		Convey("GetContainerByName error", func() {
			patches := ApplyMethod(d, "GetContainerByName", func(_ *docker.Docker, name string) (types.Container, error) {
				return types.Container{}, errors.New("test")
			})
			defer patches.Reset()
			err := rClone.Sync(RCloneSyncOptions{})
			So(err, ShouldBeError)
		})
		Convey("RunCommandSync error 1", func() {
			patches := ApplyMethod(d, "GetContainerByName", func(d *docker.Docker, name string) (types.Container, error) {
				return types.Container{}, nil
			})
			patches.ApplyMethod(d, "RunCommandSync", func(d *docker.Docker, container types.Container, commands []string) (string, error) {
				return "", errors.New("test")
			})
			defer patches.Reset()
			err := rClone.Sync(RCloneSyncOptions{
				FilePath:       "/test/filePath",
				Progress:       true,
				BandwidthLimit: "1M",
				RemoteName:     "test",
				RemotePath:     "/test/tt",
			})
			So(err, ShouldBeError)
		})
		Convey("RunCommandSync error 2", func() {
			patches := ApplyMethod(d, "GetContainerByName", func(d *docker.Docker, name string) (types.Container, error) {
				return types.Container{}, nil
			})
			patches.ApplyMethodSeq(d, "RunCommandSync", []OutputCell{
				{Values: Params{"ok", nil}},
				{Values: Params{"", errors.New("test")}},
			})
			defer patches.Reset()
			err := rClone.Sync(RCloneSyncOptions{
				FilePath:       "/test/filePath",
				Progress:       true,
				BandwidthLimit: "1M",
				RemoteName:     "test",
				RemotePath:     "/test/tt",
			})
			So(err, ShouldBeError)
		})
		Convey("success", func() {
			patches := ApplyMethod(d, "GetContainerByName", func(d *docker.Docker, name string) (types.Container, error) {
				return types.Container{}, nil
			})
			patches.ApplyMethodSeq(d, "RunCommandSync", []OutputCell{
				{Values: Params{"ok", nil}},
				{Values: Params{"ok", nil}},
			})
			defer patches.Reset()
			err := rClone.Sync(RCloneSyncOptions{
				FilePath:       "/test/filePath",
				Progress:       true,
				BandwidthLimit: "1M",
				RemoteName:     "test",
				RemotePath:     "/test/tt",
			})
			So(err, ShouldBeNil)
		})
	})
}
