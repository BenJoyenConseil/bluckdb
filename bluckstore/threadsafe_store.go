package bluckstore

import (
	"github.com/BenJoyenConseil/bluckdb/bluckstore/mmap"
	"sync"
)

func NewMmapStore() ThreadSafeStore {
	return &store{
		store: &mmap.MmapKVStore{},
		lock:  &sync.RWMutex{},
	}
}

type store struct {
	lock  *sync.RWMutex
	store KVStore
}

func (s *store) Get(k string) string    { return s.store.Get(k) }
func (s *store) Put(k, v string) error  { return s.store.Put(k, v) }
func (s *store) Open(abs string)        { s.store.Open(abs) }
func (s *store) Close()                 { s.store.Close() }
func (s *store) Meta() *mmap.Directory  { return s.store.Meta() }
func (s *store) DumpPage(id int) string { return s.store.DumpPage(id) }

func (s *store) Lock()   { s.lock.Lock() }
func (s *store) Unlock() { s.lock.Unlock() }

type ThreadSafeStore interface {
	KVStore
	Lock()
	Unlock()
}
