package bluckstore

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/BenJoyenConseil/bluckdb/util"
)


func TestRecordPut_WhenRecordIsNIL_ShouldReturnNewRecordAndTrue(t *testing.T) {
	// Given
	var entry *Record
	var key util.String = "123"
	value := "Hello, world !"

	// When
	result, appended := entry.Put(key, value)

	// Then
	assert.True(t, appended)
	assert.Equal(t, result, &Record{key: key, value: value, next: nil})
}

func TestRecordPut_WhenContainsAlreadyThatKey_ShouldSetValueAndReturnFalse(t *testing.T) {
	// Given
	var key util.String = "123"
	value := "Hello, world !"
	entry := &Record{key: key, value: value, next: nil}


	// When
	result, appended := entry.Put(key, "Bye bye !")

	// Then
	assert.False(t, appended)
	assert.Equal(t, result, &Record{key: key, value: "Bye bye !", next: nil})
}

func TestRecordPut_WhenKeyCollision_ShouldFillNextRecord(t *testing.T) {
	// Given
	var entry *Record = &Record{util.String("321"), "Bye Bye world", nil}

	// When
	_, appended := entry.Put(util.String("123"), "Hello world")

	// Then
	assert.True(t, appended)
	assert.Equal(t, entry.next, &Record{util.String("123"), "Hello world", nil})
}

func TestRecordGet_WhenKeyIsEqual_ShouldReturnValue(t *testing.T) {
	// Given
	key := util.String("123")
	var entry *Record = &Record{util.String("123"), "Hello world", nil}

	// When
	has, result := entry.Get(key)

	// Then
	assert.True(t, has)
	assert.Equal(t, "Hello world", result)
}

func TestRecordGet_WhenKeyIsContainedInNextRecord_ShouldReturnValueFromNextRecord(t *testing.T) {
	// Given
	key := util.String("123")
	var entry *Record = &Record{util.String("321"), "Bye Bye world", &Record{util.String("123"), "Hello world", nil}}

	// When
	has, result := entry.Get(key)

	// Then
	assert.True(t, has)
	assert.Equal(t, "Hello world", result)
}

func TestRecordGet_WhenKeyIsNil_ShouldReturnFalseAndNil(t *testing.T) {
	// Given
	key := util.String("123")
	var entry *Record

	// When
	has, result := entry.Get(key)

	// Then
	assert.False(t, has)
	assert.Equal(t, nil, result)
}