package bluckstore



type MemKVStore struct {
	pairs map[string]string
}



func (store *MemKVStore) Get(k string) string {
	return store.pairs[k]
}

func (store *MemKVStore) Put(k, v string) {
	store.pairs[k] = v
}

func NewMemStore() KVStore {
	return &MemKVStore{make(map[string]string)}
}