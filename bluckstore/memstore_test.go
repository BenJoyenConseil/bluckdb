package bluckstore

import (
	"testing"
	"strconv"
	"math/rand"
)

func BenchmarkPutMemKVStore(b *testing.B) {
	store := NewMemStore()

	for i := 0; i < b.N; i++ {
		id := strconv.Itoa(i)
		store.Put(id, "hello world " + id)
		//fmt.Println(id, "hello world n°" + id)
	}
}


func BenchmarkGetMemKVStore(b *testing.B) {
	// setup
	store := NewMemStore()
	for i := 0; i < b.N; i++ {
		id := strconv.Itoa(rand.Intn(b.N))
		store.Put(id, "hello world " + id)
	}

	// bench
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		id := strconv.Itoa(i)
		store.Get(id)
	}
}
