package bluckstore

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"fmt"
	"strconv"
)

func TestGet(t *testing.T) {
	// Given
	assert := assert.New(t)
	store := &MemKVStore{hashmap:map[string]string{"123" : "hello world"}}

	// When
	result := store.Get("123")

	// Then
	assert.Equal("hello world", result)
}

func TestPut(t *testing.T) {
	// Given
	assert := assert.New(t)
	store := &MemKVStore{make(map[string]string)}

	// When
	store.Put("123", "hello world")

	// Then
	assert.Equal(store.hashmap["123"], "hello world")
}

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