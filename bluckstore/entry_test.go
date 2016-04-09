package bluckstore

import (
	"testing"
	"github.com/stretchr/testify/assert"
)


func TestEntryPut_WhenEntryIsNIL_ShouldReturnNewEntry(t *testing.T) {
	// Given
	var entry *Entry
	var key String = "123"
	value := "Hello, world !"

	// When
	result, appended := entry.Put(key, value)

	// Then
	assert.True(t, appended)
	assert.Equal(t, result, &Entry{key: key, value: value, next: nil})
}

func TestEntryPut_WhenContainsAllReadyKey_ShouldSetValueAndReturnFalse(t *testing.T) {
	// Given
	var key String = "123"
	value := "Hello, world !"
	entry := &Entry{key: key, value: value, next: nil}


	// When
	result, appended := entry.Put(key, "Bye bye !")

	// Then
	assert.False(t, appended)
	assert.Equal(t, result, &Entry{key: key, value: "Bye bye !", next: nil})
}

func TestEntryPut_WhenKeyCollision_ShouldFillNextEntry(t *testing.T) {
	// Given
	var entry *Entry = &Entry{String("321"), "Bye Bye world", nil}

	// When
	_, appended := entry.Put(String("123"), "Hello world")

	// Then
	assert.True(t, appended)
	assert.Equal(t, entry.next, &Entry{String("123"), "Hello world", nil})
}

func TestEntryGet_WhenKeyIsEqual_ShouldReturnValue(t *testing.T) {
	// Given
	key := String("123")
	var entry *Entry = &Entry{String("123"), "Hello world", nil}

	// When
	has, result := entry.Get(key)

	// Then
	assert.True(t, has)
	assert.Equal(t, "Hello world", result)
}

func TestEntryGet_WhenKeyIsContainedInNextEntry_ShouldReturnValueFromNextEntry(t *testing.T) {
	// Given
	key := String("123")
	var entry *Entry = &Entry{String("321"), "Bye Bye world", &Entry{String("123"), "Hello world", nil}}

	// When
	has, result := entry.Get(key)

	// Then
	assert.True(t, has)
	assert.Equal(t, "Hello world", result)
}

func TestEntryGet_WhenKeyIsNil_ShouldReturnFalseAndNil(t *testing.T) {
	// Given
	key := String("123")
	var entry *Entry

	// When
	has, result := entry.Get(key)

	// Then
	assert.False(t, has)
	assert.Equal(t, nil, result)
}