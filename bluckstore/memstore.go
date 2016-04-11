package bluckstore



type MemKVStore struct {
	hashmap *HashMap
}

func (store *MemKVStore) Get(k string) string {
	return store.hashmap.Get(String(k)).(string)
}

func (store *MemKVStore) Put(k, v string) {
	store.hashmap.Put(String(k), v)
}

func NewMemStore() KVStore {
	store := &MemKVStore{}
	store.hashmap = NewHashMap()
	return store
}