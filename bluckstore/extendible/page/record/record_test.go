package extendible

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestPayload(t *testing.T) {
	// Given
	record := &ByteRecord{
		key: make([]byte, 3),
		value: make([]byte, 9),
	}

	// When
	result := record.Payload()

	// Then
	assert.Equal(t, uint16(16), result)
}

func TestBytes(t *testing.T) {
	// Given
	record := &ByteRecord{
		key: []byte{'1', '2', '3'},
		value: []byte{'H', 'e', 'l', 'l', 'o'},
	}

	// When
	result := record.Bytes()

	// Then
	var lenKeyEquals3 byte = 0x3
	var lenValueEquals5 byte = 0x5
	assert.Equal(t, []byte{lenKeyEquals3, lenValueEquals5,'1', '2', '3', 'H', 'e', 'l', 'l', 'o'}, result)
}