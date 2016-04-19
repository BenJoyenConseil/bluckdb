package util

import (
	"testing"
	"github.com/stretchr/testify/assert"
)


func TestHash(t *testing.T) {
	// Given
	var key String = "123"

	// When
	result := key.Hash()

	// Then
	assert.Equal(t, 1916298011, result)
}

func TestEquals_WhenStringContentsAreEqual_ShouldReturnTrue(t *testing.T) {
	// Given
	var key String = "123"

	// When
	result := key.Equals(String("123"))

	// Then
	assert.True(t, result)
}

func TestEquals_WhenStringContentsAreNotEqual_ShouldReturnFalse(t *testing.T) {
	// Given
	var key String = "123"

	// When
	result := key.Equals(String("124"))

	// Then
	assert.False(t, result)
}