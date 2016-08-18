package memap

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"encoding/binary"
)

func TestHasNext(t *testing.T) {
	// Given
	it := &PageIterator{
		p: Page(make([]byte, 4096)),
		current: 4,
	}
	binary.LittleEndian.PutUint16(it.p[4094:], 4)

	// When
	result := it.hasNext()

	// Then
	assert.False(t, result)
}

func TestHasNext_shouldReturnTrueWhenCurrentIsHigherThan_TOTALHEADERSIZE(t *testing.T) {
	// Given
	it := &PageIterator{
		p: Page(make([]byte, 4096)),
		current: 5,
	}
	binary.LittleEndian.PutUint16(it.p[4094:], 5)

	// When
	result := it.hasNext()

	// Then
	assert.True(t, result)
}

func TestHasNext_ShouldReturnTrueWhenCurrentIsEqualToPageUse(t *testing.T) {
	// Given
	it := &PageIterator{
		p: Page(make([]byte, 4096)),
		current: 5,
	}
	binary.LittleEndian.PutUint16(it.p[4094:], 5)

	// When
	result := it.hasNext()

	// Then
	assert.True(t, result)
}

func TestNext(t *testing.T) {
	// Given
	it := &PageIterator{
		p: Page(make([]byte, 4096)),
		current: 12,
	}
	binary.LittleEndian.PutUint16(it.p[4094:], 12)
	it.p[0] = 'H'
	it.p[1] = 'i'
	binary.LittleEndian.PutUint16(it.p[2:], 1)
	binary.LittleEndian.PutUint16(it.p[4:], 1)
	it.p[6] = 'Y'
	it.p[7] = 'o'
	binary.LittleEndian.PutUint16(it.p[8:], 1)
	binary.LittleEndian.PutUint16(it.p[10:], 1)

	// When
	r := it.next()
	assert.Equal(t, 6, it.current)
	r2 := it.next()

	// Then
	assert.Equal(t, "Y", string(r.Key()))
	assert.Equal(t, "o", string(r.Val()))
	assert.Equal(t, "H", string(r2.Key()))
	assert.Equal(t, "i", string(r2.Val()))
	assert.Equal(t, 0, it.current)
}
