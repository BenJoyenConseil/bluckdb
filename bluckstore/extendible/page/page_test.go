package extendible

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

type stubRecord struct {payload uint16}
func (self *stubRecord) Bytes() []byte {return nil}
func (self *stubRecord) Payload() uint16 {return self.payload}

func TestFull_When_Record_Payload_Plus_USE_areHigherTo_PAGE_SIZE(t *testing.T) {
	// Given
	page := &Page{
		use: PAGE_SIZE - 13,
	}

	// When
	result := page.Full(&stubRecord{payload: 14})

	// Then
	assert.True(t, result)
}

func TestFull_When_Record_Payload_Plus_USE_areLowerTo_PAGE_SIZE(t *testing.T) {
	// Given
	page := &Page{
		use: PAGE_SIZE - 20,
	}

	// When
	result := page.Full(&stubRecord{payload: 14})

	// Then
	assert.False(t, result)
}

func TestFull_When_Record_Payload_Plus_USE_areEqualTo_PAGE_SIZE(t *testing.T) {
	// Given
	page := &Page{
		use: PAGE_SIZE - 20,
	}

	// When
	result := page.Full(&stubRecord{payload: 20})

	// Then
	assert.False(t, result)
}