package bluckstore

import (
	"testing"
	"fmt"
	"strconv"
)

func BenchmarkPutDiskKVStore(b *testing.B) {
	store := NewDiskStore()

	for i := 0; i < b.N; i++ {
		id := strconv.Itoa(i)
		store.Put(id, "hello world " + id)
		fmt.Println(id, "hello world n°" + id)
	}
}


func BenchmarkGetDiskKVStore(b *testing.B) {
	store := NewDiskStore()

	for i := 0; i < b.N; i++ {
		id := strconv.Itoa(i)
		result := store.Get(id)
		fmt.Println(id + " : " + result)
	}
}