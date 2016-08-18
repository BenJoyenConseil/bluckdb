package util

import (
	"testing"
	"github.com/stretchr/testify/assert"
)


func TestHash(t *testing.T) {
	// Given
	var key Key = "123"

	// When
	result := key.Hash()

	// Then
	assert.Equal(t, 1916298011, result)
}