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
		current: 0,
	}

	// When
	result := it.hasNext()

	// Then
	assert.False(t, result)
}

func TestHasNext_shouldReturnTrueWhenCurrentIsLowerThanPageUse(t *testing.T) {
	// Given
	it := &PageIterator{
		p: Page(make([]byte, 4096)),
		current: 0,
	}
	binary.LittleEndian.PutUint16(it.p[4094:], 1)

	// When
	result := it.hasNext()

	// Then
	assert.True(t, result)
}

func TestNext(t *testing.T) {
	// Given
	it := &PageIterator{
		p: Page(make([]byte, 4096)),
		current: 0,
	}
	binary.LittleEndian.PutUint16(it.p[4094:], 4)
	binary.LittleEndian.PutUint16(it.p[0:], 1)
	binary.LittleEndian.PutUint16(it.p[2:], 1)
	it.p[4] = 'H'
	it.p[5] = 'i'

	// When
	rKey, rVal := it.next()

	// Then
	assert.Equal(t, "H", rKey)
	assert.Equal(t, "i", rVal)
	assert.Equal(t, 6, it.current)
}
