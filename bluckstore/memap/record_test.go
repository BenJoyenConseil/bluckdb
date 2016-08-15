package memap

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"encoding/binary"
)

var r ByteRecord = ByteRecord([]byte{
	'1', '2', '3', '4', '5',
	'Y', 'o', 'l', 'o',
	0x4, 0x0,
	0x5, 0x0,
})


func TestByteRecord_Key(t *testing.T) {
	// Given
	
	// When
	result := r.Key()

	// Then
	assert.Equal(t, "12345", string(result))
}

func TestByteRecord_Val(t *testing.T) {
	// Given

	// When
	result := r.Val()

	// Then
	assert.Equal(t, "Yolo", string(result))
}

func TestByteRecord_KeyLen(t *testing.T) {
	// Given

	// When
	result := r.KeyLen()

	// Then
	assert.Equal(t, uint16(5), result)
}

func TestByteRecord_ValLen(t *testing.T) {
	// Given

	// When
	result := r.ValLen()

	// Then
	assert.Equal(t, uint16(4), result)
}

func TestByteRecord_SetKey(t *testing.T) {
	// Given
	r = ByteRecord(make([]byte, 14))

	// When
	r.Write("654321", "Yolo")

	// Then
	assert.Equal(t, "654321", string(r[0 : 6]))
	assert.Equal(t, uint16(4), binary.LittleEndian.Uint16(r[10:]))
	assert.Equal(t, uint16(6), binary.LittleEndian.Uint16(r[12:]))
}