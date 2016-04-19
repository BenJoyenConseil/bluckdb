package bluckstore

import (
	"testing"
	"strconv"
	"math/rand"
	"os"
)

func BenchmarkPutDiskKVStore(b *testing.B) {
	store := NewDiskStore()

	for i := 0; i < b.N; i++ {
		id := strconv.Itoa(rand.Intn(i + 1))
		store.Put(id, "hello world " + id)
		//fmt.Println(id, "hello world n°" + id)
	}
}


func BenchmarkGetDiskKVStore(b *testing.B) {
	// setup
	for f := 0; f < 10; f++ {
		os.Remove("/tmp/data" + strconv.Itoa(f) + ".blk")
	}
	store := NewDiskStore()
	for i := 0; i < b.N; i++ {
		id := strconv.Itoa(rand.Intn(b.N))
		store.Put(id, "hello world " + id)
	}

	// start bench
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		id := strconv.Itoa(i)
		store.Get(id)
		//fmt.Println(id + " : " + result)
	}
}
