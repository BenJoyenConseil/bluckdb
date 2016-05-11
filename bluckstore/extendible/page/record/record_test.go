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
	expected := []byte{0x3, 0x0, 0x5, 0x0, '1', '2', '3', 'H', 'e', 'l', 'l', 'o'}
	assert.Equal(t, expected, result)
}


func TestRecordUnserializer(t *testing.T) {
	// Given
	data := []byte{0x3, 0x0, 0x5, 0x0, '1', '2', '3', 'H', 'e', 'l', 'l', 'o', '1'}
	recordUnserializer := &ByteRecordUnserializer{}

	// When
	result := recordUnserializer.Unserialize(data)

	// Then
	assert.Equal(t, []byte("123"), result.Key())
	assert.Equal(t, []byte("Hello"), result.Value())
}

func TestRecordSerializer(t *testing.T) {
	// Given
	record := &ByteRecord{
		key: []byte{'1', '2', '3'},
		value: []byte{'H', 'e', 'l', 'l', 'o'},
	}
	recordSerializer := &ByteRecordSerializer{}

	// When
	result := recordSerializer.Serialize(record)

	// Then
	expected := []byte{0x3, 0x0, 0x5, 0x0, '1', '2', '3', 'H', 'e', 'l', 'l', 'o'}
	assert.Equal(t, expected, result)
}