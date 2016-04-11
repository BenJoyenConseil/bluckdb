package bluckstore

import (
	"testing"
	"github.com/stretchr/testify/assert"
)


func TestHashMapBucket(t *testing.T) {
	// Given
	hashMap := NewHashMap()

	// When
	result := hashMap.bucket(String("123"))
	result2 := hashMap.bucket(String("124"))


	// Then
	assert.Equal(t, 3, result)
	assert.Equal(t, 2, result2)
}

func TestHashMapExpand(t *testing.T) {
	// Given
	hashMap := NewHashMap()
	oldEntry := &Entry{String("123"), "some value", nil}
	hashMap.table[1] = oldEntry

	// When
	hashMap.expand()

	// Then
	assert.Equal(t, 16, len(hashMap.table))
	assert.Contains(t, hashMap.table, oldEntry)
}

func TestHashMapPut(t *testing.T) {
	// Given
	hashMap := NewHashMap()

	// When
	hashMap.Put(String("123"), "Hello world")

	// Then
	assert.Equal(t, &Entry{String("123"), "Hello world", nil}, hashMap.table[3])
}

func TestHashMapPut_WhenTableLenIsLessThan_2xActuelSize_ShouldExpandTable(t *testing.T) {
	// Given
	hashMap := NewHashMap()
	hashMap.size = 6

	// When
	hashMap.Put(String("123"), "Hello world")

	// Then
	bucketNumberX2 := 16
	assert.Equal(t, bucketNumberX2, len(hashMap.table))
}

func TestHashMapGet(t *testing.T) {
	// Given
	hashmap := NewHashMap()
	hashmap.table[3] = &Entry{String("123"), "Hello world", nil}

	// When
	result := hashmap.Get(String("123"))

	// Then
	assert.Equal(t, "Hello world", result)
}

func TestHashMapGet_WhenKeyIsEqualToTheSecondEntry_ShouldReturnValueOfTheSecondEntry(t *testing.T) {
	// Given
	hashmap := NewHashMap()
	second := &Entry{String("123"), "Hello world", nil}
	hashmap.table[3] = &Entry{String("122"), "some value", second}

	// When
	result := hashmap.Get(String("123"))

	// Then
	assert.Equal(t, "Hello world", result)
}