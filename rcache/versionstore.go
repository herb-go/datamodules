package rcache

import (
	"sync"

	"github.com/herb-go/herbdata/kvdb"
)

type VersionStore interface {
	SetVersion(key []byte, version []byte) error
	GetVersion(Key []byte) ([]byte, error)
}

type SyncmapVersionStore struct {
	store sync.Map
}

func (s *SyncmapVersionStore) SetVersion(key []byte, version []byte) error {
	s.store.Store(key, version)
	return nil
}
func (s *SyncmapVersionStore) GetVersion(key []byte) ([]byte, error) {
	data, ok := s.store.Load(key)
	if !ok {
		return nil, kvdb.ErrKeyNotFound
	}
	return data.([]byte), nil
}
func NewSyncmapVersionStore() *SyncmapVersionStore {
	return &SyncmapVersionStore{}
}
