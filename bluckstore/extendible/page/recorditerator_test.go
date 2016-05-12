package extendible

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/BenJoyenConseil/bluckdb/bluckstore/extendible/page/record"
)

func TestNext(t *testing.T) {
	// Given
	iterator := &RecordIterator{
		content: append([]byte{0x3, 0x0, 0x5, 0x0}, []byte("keyvalue")...),
		unserializer: &extendible.ByteRecordUnserializer{},
	}

	// When
	result := iterator.Next()

	// Then
	assert.Equal(t, "key", string(result.Key()))
	assert.Equal(t, "value", string(result.Value()))
	assert.Empty(t, iterator.content)
}

func TestNext_ContentIsEmptyl(t *testing.T) {
	// Given
	iterator := &RecordIterator{
		content: make([]byte, 0),
		unserializer: &extendible.ByteRecordUnserializer{},
	}

	// When
	result := iterator.Next()

	// Then
	assert.Nil(t, result)
}

func TestNext_ContentIsNil(t *testing.T) {
	// Given
	iterator := &RecordIterator{
		content: nil,
	}

	// When
	result := iterator.Next()

	// Then
	assert.Nil(t, result)
}


func TestHasNext_ContentIsNotInitialized(t *testing.T) {
	// Given
	iterator := &RecordIterator{
		content: make([]byte, 12),
	}

	// When
	result := iterator.HasNext()

	// Then
	assert.False(t, result)
}