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
	binary.LittleEndian.PutUint16(p[4094:], 16) // use

	// insert a record
	k := "key1"
	v := "Yolo !"
	binary.LittleEndian.PutUint16(p[0:], 4) // length of key
	binary.LittleEndian.PutUint16(p[2:], 6) // length of value
	copy(p[4:], k)
	copy(p[8:], v)
	// end record

	// When
	result := p.get(k)

	// Then
	assert.Equal(t, "Yolo !", string(result))
}

func TestGet_ShouldReturnEmptyStringWhenKeyDoesntExist(t *testing.T) {
	// Given
	var p Page = make([]byte, 4096)
	binary.LittleEndian.PutUint16(p[4094:], 14) // use

	// insert a record
	k := "key1"
	v := "Yolo !"
	binary.LittleEndian.PutUint16(p[0:], 4) // length of key
	binary.LittleEndian.PutUint16(p[2:], 6) // length of value
	copy(p[6:], k)
	copy(p[8:], v)
	// end record

	// When
	result := p.get("key321")

	// Then
	assert.Empty(t, string(result))
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

func TestShift_ShouldShiftRecordsAfterTheRemovedOne_UntilPageUse(t *testing.T) {
	// Given
	var p Page = make([]byte, 4096)
	binary.LittleEndian.PutUint16(p[4094:], 31)
	k := "key"
	k2 := "key2"
	v := "Yolo !!"
	v2 := "Yolo 2 !!"
	binary.LittleEndian.PutUint16(p[0:], 3) // length of key
	binary.LittleEndian.PutUint16(p[2:], 7) // length of value
	copy(p[6:], k)
	copy(p[8:], v)
	binary.LittleEndian.PutUint16(p[14:], 4) // length of key
	binary.LittleEndian.PutUint16(p[16:], 9) // length of value
	copy(p[18:], k2)
	copy(p[22:], v2)

	offset := 0
	size := 14
	// When
	p.shift(offset, size)

	// Then
	assert.Equal(t, "key2", string(p[4:8]))
	assert.Equal(t, "Yolo 2 !!", string(p[8:17]))
}

func TestShift_ShouldMinusPageUseWithTheRemovedRecordPayload(t *testing.T) {
	// Given
	var p Page = make([]byte, 4096)
	binary.LittleEndian.PutUint16(p[4094:], 31)

	offset := 0
	size := 14
	// When
	p.shift(offset, size)

	// Then
	assert.Equal(t, 17, p.use())
}

func TestFind_ShouldReturnThe_OffsetOfRecord_lenKey_lenValue(t *testing.T) {
	// Given
	var p Page = make([]byte, 4096)
	binary.LittleEndian.PutUint16(p[4094:], 28) // use

	// insert a record
	k := "key1"
	v := "Yolo !"
	binary.LittleEndian.PutUint16(p[0:], 4) // length of key
	binary.LittleEndian.PutUint16(p[2:], 6) // length of value
	copy(p[4:], k)
	copy(p[8:], v)
	// end record
	// insert a record
	k2 := "key2"
	binary.LittleEndian.PutUint16(p[14:], 4) // length of key
	binary.LittleEndian.PutUint16(p[16:], 6) // length of value
	copy(p[18:], k2)
	copy(p[22:], v)
	// end record

	// When
	offset, lenK, lenV := p.find("key2")

	// Then
	assert.Equal(t, 14, offset)
	assert.Equal(t, 4, lenK)
	assert.Equal(t, 6, lenV)
}

func TestRemovet(t *testing.T) {
	// Given
	var p Page = make([]byte, 4096)
	binary.LittleEndian.PutUint16(p[4094:], 28) // use

	// insert a record
	k := "key1"
	v := "Yolo !"
	binary.LittleEndian.PutUint16(p[0:], 4) // length of key
	binary.LittleEndian.PutUint16(p[2:], 6) // length of value
	copy(p[4:], k)
	copy(p[8:], v)
	// end record
	// insert a record
	k2 := "key2"
	binary.LittleEndian.PutUint16(p[14:], 4) // length of key
	binary.LittleEndian.PutUint16(p[16:], 6) // length of value
	copy(p[18:], k2)
	copy(p[22:], v)
	// end record

	// When
	p.remove("key2")

	// Then
	assert.Equal(t, 14, p.use())
}