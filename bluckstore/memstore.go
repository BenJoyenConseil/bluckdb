package bluckstore

import "github.com/BenJoyenConseil/bluckdb/util"


type MemKVStore struct {
	hashmap *HashMap
}

func (store *MemKVStore) Get(k string) string {
	return store.hashmap.Get(util.String(k)).(string)
}

func (store *MemKVStore) Put(k, v string) {
	store.hashmap.Put(util.String(k), v)
}

func NewMemStore() KVStore {
	store := &MemKVStore{}
	store.hashmap = NewHashMap()
	return store
}