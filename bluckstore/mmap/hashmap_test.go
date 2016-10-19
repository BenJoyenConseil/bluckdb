package mmap

import (
	"math/rand"
	"os"
	"strconv"
	"testing"
)

func BenchmarkHashmapGet(b *testing.B) {
	store := make(map[string]string)
	size := 10000
	for i := 0; i < size; i++ {
		store[strconv.Itoa(i)] = "mec, elle est où ma caisse ??"
	}
	devNull, _ := os.Open(os.DevNull)
	defer devNull.Close()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		devNull.Write([]byte(store[strconv.Itoa(rand.Intn(10000-1))]))
	}
}

func BenchmarkHashmapPut(b *testing.B) {
	store := make(map[string]string)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		store[strconv.Itoa(n)] = "mec, elle est où ma caisse ??"
	}
}
