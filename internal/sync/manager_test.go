package sync

import "testing"

func Test_SyncerManager(t *testing.T) {
	sm, _ := NewManager()
	for _, syncer := range sm.GetSyncers() {
		err := syncer.Sync()
		if err != nil {
			panic(err)
		}
	}
}
