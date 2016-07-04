package bluckstore

import (
	"testing"
	"strconv"
	"os"
	"github.com/BenJoyenConseil/bluckdb/util"
	"math/rand"
)

func BenchmarkPutDiskKVStore(b *testing.B) {
	for f := 0; f < BUCKET_NUMER; f++ {
		os.Remove(buildPartitionFilePathString(util.String(strconv.Itoa(f))))
	}
	store := NewDiskStore()
	size := 1000000


	for n := 0; n < b.N; n++ {
		store.Put(strconv.Itoa(rand.Intn(size - 1)), "mec, elle est où ma caisse ??")
	}
}


func BenchmarkGetDiskKVStore(b *testing.B) {
	// setup
	for f := 0; f < BUCKET_NUMER; f++ {
		os.Remove(buildPartitionFilePathString(util.String(strconv.Itoa(f))))
	}
	store := NewDiskStore()
	size := 1000000
	for i := 0; i < size; i++ {
		store.Put(strconv.Itoa(i), "mec, elle est où ma caisse ??")
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		store.Get(strconv.Itoa(rand.Intn(size - 1)))
	}
}
