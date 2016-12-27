package bluckstore

import (
	"github.com/labstack/gommon/log"
	"sync"
)

type MultiStore struct {
	stores map[string]ThreadSafeStore
	lock   *sync.RWMutex
}

func NewMmapMultiStore() *MultiStore {
	return &MultiStore{
		stores: make(map[string]ThreadSafeStore),
		lock:   &sync.RWMutex{},
	}
}

func (server *MultiStore) GetStore(path string) ThreadSafeStore {
	server.lock.Lock()
	if server.stores[path] == nil {
		log.Infof("First time using the store instance in path %s", path)
		s := NewMmapStore()

		s.Open(path)
		log.Debugf("Open %s", path)
		server.stores[path] = s
	}
	server.lock.Unlock()
	return server.stores[path]
}

func (server *MultiStore) Close() {
	for _, store := range server.stores {
		store.Close()
	}
}
