package memap

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"os"
	"strconv"
	"testing"
)

func TestStorePut_shouldReOpen_UsingMeta(t *testing.T) {
	// Given
	os.Remove("/tmp/data.db")
	os.Remove("/tmp/db.meta")
	store := MmapKVStore{}
	store.Open()
	store.Put("KEY", "VALUE")
	store.Close()

	// When
	store.Open()

	// Then
	assert.Equal(t, 4096, len(store.Dir.data))
	assert.Equal(t, "/tmp/data.db", store.Dir.dataFile.Name())
	assert.Equal(t, []int{0}, store.Dir.Table)
	assert.Equal(t, 0, int(store.Dir.Gd))
	assert.Equal(t, "KEYVALUE", string(store.Dir.data[0:8]))
	store.Close()
}

func BenchmarkMmapPut(b *testing.B) {
	store := &MmapKVStore{}
	store.Rm()
	store.Open()
	defer store.Close()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		store.Put(strconv.Itoa(n), "mec, elle est où ma caisse ??")
	}

}

func setup() {
	store := &MmapKVStore{}
	store.Rm()
	store.Open()
	size := 1000000
	for i := 0; i < size; i++ {
		store.Put(strconv.Itoa(i), "mec, elle est où ma caisse ??")
	}
	store.Close()
}

func BenchmarkMmapGet(b *testing.B) {
	setup()
	store := &MmapKVStore{}
	store.Open()
	defer store.Close()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		store.Get(strconv.Itoa(rand.Intn(1000000 - 1)))
	}
}
