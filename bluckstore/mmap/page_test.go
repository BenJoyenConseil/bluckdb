package mmap

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestRest_DefaultIsPAGE_LOCAL_DEPTH_OFFSET(t *testing.T) {
	// Given
	var p Page = make([]byte, PAGE_SIZE)

	// When
	result := p.rest()

	// Then
	assert.Equal(t, PAGE_LOCAL_DEPTH_OFFSET, result)
}

func TestRest(t *testing.T) {
	// Given
	var p Page = make([]byte, PAGE_SIZE)
	// use = 4000
	binary.LittleEndian.PutUint16(p[PAGE_USE_OFFSET:], uint16(4000))

	// When
	result := p.rest()

	// Then
	assert.Equal(t, 92, result)
}

func TestLd(t *testing.T) {
	// Given
	var p Page = make([]byte, PAGE_SIZE)
	binary.LittleEndian.PutUint16(p[PAGE_LOCAL_DEPTH_OFFSET:], uint16(16))

	// When
	result := p.ld()

	// Then
	assert.Equal(t, 16, result)
}

func TestSetLd(t *testing.T) {
	// Given
	var p Page = make([]byte, PAGE_SIZE)
	ld := 30

	// When
	p.setLd(ld)
	result := binary.LittleEndian.Uint16(p[PAGE_LOCAL_DEPTH_OFFSET:])

	// Then
	assert.Equal(t, 30, int(result))
}

func TestGet(t *testing.T) {
	// Given
	var p Page = make([]byte, PAGE_SIZE)
	binary.LittleEndian.PutUint16(p[PAGE_USE_OFFSET:], 14) // use

	// insert a record
	k := "key1"
	v := "Yolo !"
	binary.LittleEndian.PutUint16(p[12:], 4) // length of key
	binary.LittleEndian.PutUint16(p[10:], 6) // length of value
	copy(p[0:], k)
	copy(p[4:], v)
	// end record

	// When
	result, err := p.Get("key1")

	// Then
	assert.Equal(t, "Yolo !", string(result))
	assert.Nil(t, err)
}

func TestGet_ShouldReturnEmptyStringWhenKeyDoesntExist(t *testing.T) {
	// Given
	var p Page = make([]byte, PAGE_SIZE)
	binary.LittleEndian.PutUint16(p[PAGE_USE_OFFSET:], 14) // use

	// insert a record
	k := "key1"
	v := "Yolo !"
	binary.LittleEndian.PutUint16(p[12:], 4) // length of key
	binary.LittleEndian.PutUint16(p[10:], 6) // length of value
	copy(p[0:], k)
	copy(p[4:], v)
	// end record

	// When
	result, err := p.Get("key321")

	// Then
	assert.Empty(t, string(result))
	assert.Error(t, err)
}

func TestPut_UseShouldBeIncrementedWithThePayloadOfTheNewRecord(t *testing.T) {
	// Given
	var p Page = make([]byte, PAGE_SIZE)
	k := "key1"   // len (=2) + key (=4)   = 6 bytes
	v := "Yolo !" // len (=2) + value (=6) = 8 bytes

	// When
	p.Put(k, v)

	// Then
	assert.Equal(t, 14, p.Use())
}

func TestPut_(t *testing.T) {
	// Given
	var p Page = make([]byte, PAGE_SIZE)
	k := "key1"
	v := "Yolo !"

	// When
	p.Put(k, v)
	p.Put(k, "Yolo updated !")

	// Then
	lenKey := binary.LittleEndian.Uint16(p[12:14])
	assert.Equal(t, uint16(4), lenKey) // {0x4, 0x0} : LittleEndian

	lenVal := binary.LittleEndian.Uint16(p[10:12])
	assert.Equal(t, uint16(6), lenVal)

	rKey := string(p[0:4])
	assert.Equal(t, "key1", rKey)

	rVal := string(p[4:10])
	assert.Equal(t, "Yolo !", rVal)

	lenVal = binary.LittleEndian.Uint16(p[32:34])
	assert.Equal(t, uint16(14), lenVal)
	rVal = string(p[18:32])
	assert.Equal(t, "Yolo updated !", rVal)
}

func TestPut_shouldReturnAnErrorWhenRestOfPageIsLowerThanRecordPayload(t *testing.T) {
	// Given
	var p Page = make([]byte, PAGE_SIZE)
	binary.LittleEndian.PutUint16(p[PAGE_USE_OFFSET:], 4080)
	k := "key1"
	v := "Yolo !"

	// When
	result := p.Put(k, v)

	// Then
	assert.Equal(t, errors.New("The page is full."), result)
}

func BenchmarkPagePut(b *testing.B) {

	for n := 0; n < b.N; n++ {
		var p Page = make([]byte, PAGE_SIZE)
		for i := 0; i < 100; i++ {
			p.Put("key", "mec, elle est où ma caisse ??")
		}
	}
	fmt.Println(b.N)
}

func BenchmarkPageGet(b *testing.B) {
	var p Page = make([]byte, PAGE_SIZE)
	for n := 0; n < 150; n++ {
		p.Put("key"+strconv.Itoa(n), "mec, elle est où ma caisse ??")
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		p.Get("key" + strconv.Itoa(n%150))
	}
	fmt.Println(b.N)
}
