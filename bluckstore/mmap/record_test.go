package mmap

import (
	"encoding/binary"
	"github.com/stretchr/testify/assert"
	"testing"
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
	result := r.key()

	// Then
	assert.Equal(t, "12345", string(result))
}

func TestByteRecord_Key_WhenSliceIsBiggerThanRecord(t *testing.T) {
	// Given
	var r ByteRecord = ByteRecord([]byte{
		0, 0, 0, 0, 0, 0,
		'1', '2', '3', '4', '5',
		'Y', 'o', 'l', 'o',
		0x4, 0x0,
		0x5, 0x0,
	})

	// When
	result := r.key()

	// Then
	assert.Equal(t, "12345", string(result))
}

func TestByteRecord_Val(t *testing.T) {
	// Given

	// When
	result := r.val()

	// Then
	assert.Equal(t, "Yolo", string(result))
}

func TestByteRecord_Val_WhenSliceIsBiggerThanRecord(t *testing.T) {
	// Given
	var r ByteRecord = ByteRecord([]byte{
		0, 0, 0, 0, 0, 0,
		'1', '2', '3', '4', '5',
		'Y', 'o', 'l', 'o',
		0x4, 0x0,
		0x5, 0x0,
	})

	// When
	result := r.val()

	// Then
	assert.Equal(t, "Yolo", string(result))
}

func TestByteRecord_KeyLen(t *testing.T) {
	// Given

	// When
	result := r.keyLen()

	// Then
	assert.Equal(t, uint16(5), result)
}

func TestByteRecord_ValLen(t *testing.T) {
	// Given

	// When
	result := r.valLen()

	// Then
	assert.Equal(t, uint16(4), result)
}

func TestByteRecord_Write(t *testing.T) {
	// Given
	r = ByteRecord(make([]byte, 14))

	// When
	r.Write("654321", "Yolo")

	// Then
	assert.Equal(t, "654321", string(r[0:6]))
	assert.Equal(t, uint16(4), binary.LittleEndian.Uint16(r[10:]))
	assert.Equal(t, uint16(6), binary.LittleEndian.Uint16(r[12:]))
}

func BenchmarkByteRecord_Write(b *testing.B) {

	k := "key123123123"
	v := "yolo i am the value of the key 123123123 !! yolo !"
	r := ByteRecord(make([]byte, len(k)+len(v)+RECORD_TOTAL_HEADER_SIZE))

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		r.Write(k, v)
	}
}

type AppendRecord []byte

func (r AppendRecord) Write(key, val string) {
	lenKey := uint16(len(key))
	lenVal := uint16(len(val))
	total := lenKey + lenVal
	binary.LittleEndian.PutUint16(r[total:], lenVal)
	binary.LittleEndian.PutUint16(r[total+RECORD_HEADER_SIZE:], lenKey)
	r = append(r[0:0], (key + val)...)
}

func BenchmarkByteRecord_Write_WitheAppend(b *testing.B) {
	k := "key123123123"
	v := "yolo i am the value of the key 123123123 !! yolo !"
	r := AppendRecord(make([]byte, (len(k) + len(v) + RECORD_TOTAL_HEADER_SIZE)))

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		r.Write(k, v)
	}
}
