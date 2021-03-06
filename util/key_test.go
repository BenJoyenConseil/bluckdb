package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHash(t *testing.T) {
	// Given
	var key Key = []byte("123")

	// When
	result := key.Hash()

	// Then
	assert.Equal(t, 1916298011, result)
}
