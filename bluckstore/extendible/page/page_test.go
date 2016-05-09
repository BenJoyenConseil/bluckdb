package extendible

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/BenJoyenConseil/bluckdb/bluckstore/extendible/page/record"
)

type stubRecord struct {payload uint16; key []byte; value []byte}
func (self *stubRecord) Bytes() []byte {return append(self.key, self.value...)}
func (self *stubRecord) Payload() uint16 {return self.payload}
func (self *stubRecord) Key() []byte {return self.key}
func (self *stubRecord) Value() []byte {return self.value}

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

func TestNew(t *testing.T) {
	// Given

	// When
	result := New()

	// Then
	assert.Equal(t, uint64(0), result.localDepth)
	assert.Equal(t, uint16(0), result.use)
	assert.Equal(t, 4096, len(result.content))
}

type stubRecordUnSerializer struct {}
func (self *stubRecordUnSerializer) Unserialize(data []byte) extendible.Record {
	return &stubRecord{
		payload: uint16(8),
		key: data[0:3],
		value: data[3:8],
	}
}

func TestGet(t *testing.T) {
	// Given
	page := &Page{
		content: []byte{},
	}
	pair1 := []byte("121Hello")
	pair2 := []byte("122Hello")
	pair3 := []byte("123Hello")
	pair4 := []byte("124Hello")
	page.content = append(page.content, pair1...)
	page.content = append(page.content, pair2...)
	page.content = append(page.content, pair3...)
	page.content = append(page.content, pair4...)
	page.recordUnserializer = &stubRecordUnSerializer{}

	// When
	result, _ := page.Get("123")

	// Then
	assert.Equal(t, "Hello", result)
}

type stubRecordSerializer struct {}
func (self *stubRecordUnSerializer) Serialize(record extendible.Record) []byte{return append(append([]byte("lol"), record.Key()...), record.Value()...)}

func TestPut(t *testing.T) {
	// Given
	page := &Page{
		content: []byte{},
	}
	page.recordSerializer = &stubRecordUnSerializer{}

	// When

	error := page.Put("121", "Hello")
	error = page.Put("122", "Hello")
	error = page.Put("123", "Hello")
	result1 := page.content[0:11]
	result2 := page.content[11:22]
	result3 := page.content[22:]

	// Then
	assert.Equal(t, nil, error)
	assert.Equal(t, "lol121Hello", string(result1))
	assert.Equal(t, "lol122Hello", string(result2))
	assert.Equal(t, "lol123Hello", string(result3))
}

func TestPut_Overflow(t *testing.T) {
	// Given
	page := &Page{
		use: 4095,
	}

	// When
	error := page.Put("123", "Hello")

	// Then
	assert.NotNil(t, error)
}