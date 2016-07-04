package bluckstore

import (
	"testing"
	"strconv"
	"math/rand"
)

func BenchmarkPutMemKVStore(b *testing.B) {
	store := NewMemStore()
	size := 1000000

	for n := 0; n < b.N; n++ {
		store.Put(strconv.Itoa(rand.Intn(size - 1)), "mec, elle est où ma caisse ??")
	}
}


func BenchmarkGetMemKVStore(b *testing.B) {
	// setup
	store := NewMemStore()
	size := 1000000
	for i := 0; i < size; i++ {
		store.Put(strconv.Itoa(i), "mec, elle est où ma caisse ??")
	}

	// bench
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		store.Get(strconv.Itoa(rand.Intn(size - 1)))
	}
}
