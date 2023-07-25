package sync

import (
	"errors"
	"fmt"
	"github.com/cloudfogtech/sync-up/internal/common"
	"github.com/cloudfogtech/sync-up/internal/docker"
	"github.com/cloudfogtech/sync-up/internal/env"
)

type Manager struct {
	syncer map[string]*Syncer
}

func NewManager() (*Manager, error) {
	sm := &Manager{
		syncer: make(map[string]*Syncer),
	}
	err := sm.init()
	if err != nil {
		return nil, err
	}
	return sm, nil
}

func (sm *Manager) GetSyncers() []*Syncer {
	syncers := make([]*Syncer, 0)
	for _, syncer := range sm.syncer {
		syncers = append(syncers, syncer)
	}
	return syncers
}

func (sm *Manager) init() error {
	d, err := docker.NewDockerCli()
	if err != nil {
		return err
	}
	m := env.NewManager()
	serviceIds := m.GetServiceIds()
	if len(serviceIds) < 1 {
		return errors.New("no SyncUp instance definition")
	}
	bpm := NewParserManager(m)
	for _, id := range serviceIds {
		t := m.GetEnvWithNilCheck(env.TypeTpl, id)
		if !bpm.CheckType(t) {
			return errors.New(fmt.Sprintf("'%s' is not a supported type", t))
		}
		backupper := bpm.NewBackupper(t, id)
		sm.Add(id, NewSyncer(d, SyncerOptions{
			Id:        id,
			RClone:    NewRClone(d, m.GetEnvWithNilCheck(env.RCTpl, id)),
			Backupper: backupper,
			LocalDir:  m.GetEnvWithDefault(env.LocalDirPathTpl, id, common.TempDirPath),
			Cron:      m.GetEnv(env.CronTpl, id),
		}, RCloneOptions{
			Dir:            m.GetEnvWithNilCheck(env.RCDirPathTpl, id),
			BandwidthLimit: m.GetEnv(env.RCBandwidthLimitTpl, id),
			RemoteName:     m.GetEnvWithNilCheck(env.RCRemoteNameTpl, id),
			RemotePath:     m.GetEnvWithNilCheck(env.RCRemotePathTpl, id),
		}))
	}
	return nil
}

func (sm *Manager) Add(id string, syncer *Syncer) {
	sm.syncer[id] = syncer
}

func (sm *Manager) Remove(id string) {
	delete(sm.syncer, id)
}

func (sm *Manager) Get(id string) *Syncer {
	return sm.syncer[id]
}
