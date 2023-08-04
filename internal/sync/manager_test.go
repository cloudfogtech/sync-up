package sync

import (
	"errors"
	"reflect"
	"testing"

	. "github.com/agiledragon/gomonkey/v2"
	"github.com/cloudfogtech/sync-up/internal/docker"
	"github.com/cloudfogtech/sync-up/internal/env"
	. "github.com/smartystreets/goconvey/convey"
)

func TestNewManager(t *testing.T) {
	m := &Manager{}
	Convey("TestNewManager", t, func() {
		Convey("normal", func() {
			patches := ApplyPrivateMethod(reflect.TypeOf(m), "init", func(manager *Manager) error {
				return nil
			})
			defer patches.Reset()
			manager, err := NewManager()
			So(manager, ShouldNotBeNil)
			So(err, ShouldBeNil)
		})
		Convey("error", func() {
			patches := ApplyPrivateMethod(reflect.TypeOf(m), "init", func(manager *Manager) error {
				return errors.New("")
			})
			defer patches.Reset()
			manager, err := NewManager()
			So(manager, ShouldBeNil)
			So(err, ShouldBeError)
		})
	})
}

func TestManager_GetSyncers(t *testing.T) {
	Convey("TestManager_GetSyncers", t, func() {
		m := &Manager{syncer: map[string]*Syncer{
			"test": {},
		}}
		syncers := m.GetSyncers()
		So(len(syncers), ShouldEqual, 1)
	})
}

func TestManager_init(t *testing.T) {
	m := &Manager{syncer: map[string]*Syncer{
		"test": {},
	}}
	Convey("TestManager_init", t, func() {
		Convey("docker error", func() {
			patches := ApplyFunc(docker.NewDockerCli, func() (*docker.Docker, error) {
				return nil, errors.New("test")
			})
			defer patches.Reset()
			err := m.init()
			So(err, ShouldBeError)
		})
		Convey("serviceIds error", func() {
			patches := ApplyFunc(docker.NewDockerCli, func() (*docker.Docker, error) {
				return &docker.Docker{}, nil
			})
			defer patches.Reset()
			envManager := &env.Manager{}
			patches.ApplyFunc(env.NewManager, func() *env.Manager {
				return envManager
			})
			patches.ApplyMethod(envManager, "GetServiceIds", func() []string {
				return []string{}
			})
			err := m.init()
			So(err, ShouldBeError)
		})
		Convey("type not supported error", func() {
			patches := ApplyFunc(docker.NewDockerCli, func() (*docker.Docker, error) {
				return &docker.Docker{}, nil
			})
			defer patches.Reset()
			envManager := &env.Manager{}
			patches.ApplyFunc(env.NewManager, func() *env.Manager {
				return envManager
			})
			patches.ApplyMethod(envManager, "GetServiceIds", func() []string {
				return []string{
					"test",
				}
			})
			bpm := &ParserManager{}
			patches.ApplyFunc(NewParserManager, func(_ *env.Manager) *ParserManager {
				return bpm
			})
			patches.ApplyMethod(envManager, "GetEnvWithNilCheck", func(_ *env.Manager, tpl, id string) string {
				return "test"
			})
			patches.ApplyMethod(bpm, "NewBackupper", func(_ *ParserManager, t, id string) (Backupper, error) {
				return nil, errors.New("test")
			})
			err := m.init()
			So(err, ShouldBeError)
		})
		Convey("success", func() {
			patches := ApplyFunc(docker.NewDockerCli, func() (*docker.Docker, error) {
				return &docker.Docker{}, nil
			})
			defer patches.Reset()
			envManager := &env.Manager{}
			patches.ApplyFunc(env.NewManager, func() *env.Manager {
				return envManager
			})
			patches.ApplyMethod(envManager, "GetServiceIds", func() []string {
				return []string{
					"test",
				}
			})
			bpm := &ParserManager{}
			patches.ApplyFunc(NewParserManager, func(_ *env.Manager) *ParserManager {
				return bpm
			})
			patches.ApplyMethod(envManager, "GetEnvWithNilCheck", func(_ *env.Manager, tpl, id string) string {
				return "test"
			})
			patches.ApplyMethod(bpm, "NewBackupper", func(_ *ParserManager, t, id string) (Backupper, error) {
				return &PostgresContainerBackupper{}, nil
			})
			patches.ApplyMethod(m, "Add", func(_ *Manager, id string, syncer *Syncer) {
			})
			patches.ApplyFunc(NewSyncer, func(d *docker.Docker, syncerOptions SyncerOptions, rCloneOptions RCloneOptions) *Syncer {
				return &Syncer{}
			})
			patches.ApplyFunc(NewRClone, func(d *docker.Docker, containerName string) *RClone {
				return &RClone{}
			})
			patches.ApplyMethod(envManager, "GetEnvWithDefault", func(_ *env.Manager, _, _, _ string) string {
				return "test"
			})
			patches.ApplyMethod(envManager, "GetEnv", func(_ *env.Manager, _, _ string) string {
				return "test"
			})
			err := m.init()
			So(err, ShouldBeNil)
		})
	})
}

func TestManager_Add(t *testing.T) {
	Convey("TestManager_Add", t, func() {
		manager := Manager{
			syncer: make(map[string]*Syncer),
		}
		manager.Add("test", &Syncer{})
		So(len(manager.syncer), ShouldEqual, 1)
	})
}

func TestManager_Remove(t *testing.T) {
	Convey("TestManager_Remove", t, func() {
		manager := Manager{
			syncer: map[string]*Syncer{
				"test": {},
			},
		}
		manager.Remove("test")
		So(len(manager.syncer), ShouldEqual, 0)
	})
}

func TestManager_Get(t *testing.T) {
	Convey("TestManager_Remove", t, func() {
		manager := Manager{
			syncer: map[string]*Syncer{
				"test": {},
			},
		}
		So(manager.Get("test"), ShouldNotBeNil)
		So(manager.Get("test1"), ShouldBeNil)
	})
}
