package bluckstore



type KVStore interface {
	Get(k string) string
	Put(k, v string)
}

