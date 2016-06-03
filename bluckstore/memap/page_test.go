package memap

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"encoding/binary"
	"errors"
)

func TestUse(t *testing.T) {
	// Given
	var p Page = make([]byte, 4096)
	p[4094] = 0x3
	p[4095] = 0x0

	// When
	result := p.use()

	// Then
	assert.Equal(t, 3, result)
}

func TestAdd(t *testing.T) {
	// Given
	var p Page = make([]byte, 4096)

	// When
	p.add('Y')

	// Then
	assert.Equal(t, byte('Y'), p[0])
}

func TestAdd_shouldIncrementUse(t *testing.T) {
	// Given
	var p Page = make([]byte, 4096)

	// When
	p.add('Y')

	// Then
	assert.Equal(t, 1, p.use())
}

func TestAddMany(t *testing.T) {
	// Given
	var p Page = make([]byte, 4096)

	// When
	p.add([]byte{'Y', 'o', 'l', 'o'}...)
	result := string(p[0:4])

	// Then
	assert.Equal(t, "Yolo", result)
}

func TestRest_DefaultIs4092(t *testing.T) {
	// Given
	var p Page = make([]byte, 4096)

	// When
	result := p.rest()

	// Then
	assert.Equal(t, 4092, result)
}

func TestRest(t *testing.T) {
	// Given
	var p Page = make([]byte, 4096)
	// use = 4000
	binary.LittleEndian.PutUint16(p[4094:], uint16(4000))

	// When
	result := p.rest()

	// Then
	assert.Equal(t, 92, result)
}

func TestLd(t *testing.T) {
	// Given
	var p Page = make([]byte, 4096)
	binary.LittleEndian.PutUint16(p[4092:], uint16(16))

	// When
	result := p.ld()

	// Then
	assert.Equal(t, 16, result)
}

func TestSetLd(t *testing.T) {
	// Given
	var p Page = make([]byte, 4096)
	ld := 30

	// When
	p.setLd(ld)
	result := binary.LittleEndian.Uint16(p[4092:])

	// Then
	assert.Equal(t, 30, int(result))
}

func TestGet(t *testing.T) {
	// Given
	var p Page = make([]byte, 4096)
	binary.LittleEndian.PutUint16(p[4094:], 14) // use

	// insert a record
	k := "key1"
	v := "Yolo !"
	binary.LittleEndian.PutUint16(p[4:], 4) // length of key
	binary.LittleEndian.PutUint16(p[6:], 6) // length of value
	copy(p[8:12], k)
	copy(p[12:18], v)
	// end record

	// When
	result := p.get(k)

	// Then
	assert.Equal(t, "Yolo !", string(result))
}

func TestPut_UseShouldBeIncrementedWithThePayloadOfTheNewRecord(t *testing.T) {
	// Given
	var p Page = make([]byte, 4096)
	k := "key1"   // len (=2) + key (=4)   = 6 bytes
	v := "Yolo !" // len (=2) + value (=6) = 8 bytes

	// When
	p.put(k, v)

	// Then
	assert.Equal(t, 14, p.use())
}

func TestPut_(t *testing.T) {
	// Given
	var p Page = make([]byte, 4096)
	k := "key1"
	v := "Yolo !"

	// When
	p.put(k, v)

	// Then
	lenKey := []byte(p[0:2])
	assert.Equal(t, []byte{0x4, 0x0}, lenKey) // {0x4, 0x0} : LittleEndian style
	lenVal := []byte(p[2:4])
	assert.Equal(t, []byte{0x6, 0x0}, lenVal)
	rKey := []byte(p[4:8])
	assert.Equal(t, []byte{'k', 'e', 'y', '1'}, rKey)
	rVal := []byte(p[8:14])
	assert.Equal(t, []byte{'Y', 'o', 'l', 'o', ' ', '!'}, rVal)
}

func TestPut_shouldReturnAnErrorWhenRestOfPageIsLowerThanRecordPayload(t *testing.T) {
	// Given
	var p Page = make([]byte, 4096)
	binary.LittleEndian.PutUint16(p[4094:], 4080)
	k := "key1"
	v := "Yolo !"

	// When
	result := p.put(k, v)

	// Then
	assert.Equal(t, errors.New("The page is full. use = 4080"), result)
}