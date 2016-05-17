package extendible

import (
	"testing"
	"github.com/BenJoyenConseil/bluckdb/util"
	"github.com/stretchr/testify/assert"
	page "github.com/BenJoyenConseil/bluckdb/bluckstore/extendible/page"
	record "github.com/BenJoyenConseil/bluckdb/bluckstore/extendible/page/record"
	"errors"
)


type stubPage struct{content string; full bool; localDepth uint64}

func (stub *stubPage) Get(key string) (string, error){
	if key == "123" {return "Hello", nil}
	return "", errors.New("Not implemented !!")
}
func (stub *stubPage) Put(key, value string) error {
	if key == "full" {
		return errors.New("the page is full")
	}
	return nil
}
func (stub *stubPage) Full(record record.Record) bool {return stub.full}
func (stub *stubPage) LocalDepth() uint64 {return stub.localDepth}
func (stub *stubPage) SetLocalDepth(num uint64) {stub.localDepth = num}
func (stub *stubPage) Content() []byte { return []byte(stub.content) }


func TestGetPage(t *testing.T) {
	// Given
	dir := &Directory{
		pointerPageTable: make([]page.Page, 2),
		globalDepth: 1,
	}
	key := util.String("123")
	var pointerPage = &stubPage{}
	dir.pointerPageTable[1] = pointerPage

	// When
	result := dir.getPage(key)

	// Then
	var expected = pointerPage
	assert.True(t, expected == result)
}

func TestGet(t *testing.T) {
	// Given
	dir := &Directory{
		pointerPageTable: make([]page.Page, 2),
		globalDepth: 1,
	}
	dir.pointerPageTable[1] = &stubPage{}

	// When
	result, err := dir.Get("123")

	// Then
	assert.Equal(t, "Hello", result)
	assert.Nil(t, err)
}

func TestExtendibleHashing (t *testing.T) {
	// Given
	dir := &Directory{ globalDepth: 0}

	// When
	result := dir.extendibleHash(util.String("key"))

	// Then
	assert.Equal(t, 0, result)
}

func TestExtendibleHashing_whenGlobalDepthIsEqualTo3_shouldReturn4 (t *testing.T) {
	// Given
	dir := &Directory{ globalDepth: 3}

	// When
	result := dir.extendibleHash(util.String("key"))

	// Then
	assert.Equal(t, 4, result)
}


func TestExtendibleHashing_whenGlobalDepthIsEqualTo6_shouldReturn44 (t *testing.T) {
	// Given
	dir := &Directory{ globalDepth: 6}

	// When
	result := dir.extendibleHash(util.String("key"))

	// Then
	assert.Equal(t, 44, result)
}

/*
To split a bucket/page j when inserting record with search-key value Kj :
Update the second half of the bucket/page address table entries originally pointing to j, to point to z
 */
func TestPut_WhenPageIsFull_WhenLocalDepthIsLowerthanGlobal_ShouldIncrementLocalDepth(t *testing.T) {
	// Given
	dir := &Directory{
		pointerPageTable: make([]page.Page, 2),
		globalDepth: 1,
	}
	p1 := &stubPage{
		localDepth: 0,
		full: true,
	}
	dir.pointerPageTable[0] = p1
	dir.pointerPageTable[1] = p1
	value := ""
	key := "full"

	// When
	dir.Put(key, value)

	// Then
	assert.Equal(t, 1, int(dir.pointerPageTable[0].LocalDepth()))
	assert.Equal(t, 1, int(dir.pointerPageTable[1].LocalDepth()))
}

func TestPut_WhenPageIsFull_WhenLocalDepthIsLowerthanGlobal_ShouldCreateNewPage(t *testing.T) {
	// Given
	dir := &Directory{
		pointerPageTable: make([]page.Page, 2),
		globalDepth: 1,
	}
	p1 := &stubPage{
		localDepth: 0,
		full: true,
	}
	dir.pointerPageTable[0] = p1
	dir.pointerPageTable[1] = p1
	value := ""
	key := "full"

	// When
	dir.Put(key, value)

	// Then
	assert.NotEqual(t, p1, dir.pointerPageTable[0])
	assert.NotEqual(t, p1, dir.pointerPageTable[1])
}

func TestPut_WhenPageIsFull_WhenLocalDepthIsLowerthanGlobal_ShouldRebalanceRecord(t *testing.T) {
	// Given
	record1 := string([]byte{0x2, 0x0, 0x0, 0x0, 'r', '1'})
	record2 := string([]byte{0x2, 0x0, 0x0, 0x0, 'r', '2'})
	record3 := string([]byte{0x2, 0x0, 0x0, 0x0, 'r', '3'})
	dir := &Directory{
		pointerPageTable: make([]page.Page, 2),
		globalDepth: 1,
	}
	p1 := &stubPage{
		localDepth: 0,
		full: true,
		content: record1 + record2 + record3,
	}
	dir.pointerPageTable[0] = p1
	dir.pointerPageTable[1] = p1
	key := "full"
	value := "r4"

	// When
	dir.Put(key, value)

	// Then
	assert.Equal(t, "\x04\x00\x02\x00fullr4\x02\x00\x00\x00r1\x02\x00", string(dir.pointerPageTable[0].Content()[:18]))
	assert.Equal(t, "\x02\x00\x00\x00r2", string(dir.pointerPageTable[1].Content()[:6]))
}

func TestPut_WhenPageIsFull_AndLocalDepthEqualToGlobal_ShouldIncrementGlobalDepth(t *testing.T) {
	// Given
	depth := uint64(1)
	dir := &Directory{
		pointerPageTable: make([]page.Page, 2),
		globalDepth: depth,
	}
	p1 := &stubPage{
		localDepth: depth,
		full: true,
	}

	p2 := &stubPage{
		localDepth: depth,
		full: true,
	}
	dir.pointerPageTable[0] = p1
	dir.pointerPageTable[1] = p2
	value := ""
	key := "full"

	// When
	dir.Put(key, value)

	// Then
	assert.Equal(t, 2, int(dir.globalDepth))
}

func TestPut_WhenPageIsFull_AndLocalDepthEqualToGlobal_ShouldMultiplyByTwoThePointerPageTable(t *testing.T) {
	// Given
	depth := uint64(1)
	dir := &Directory{
		pointerPageTable: make([]page.Page, 2),
		globalDepth: depth,
	}
	p1 := &stubPage{
		localDepth: depth,
		full: true,
	}

	p2 := &stubPage{
		localDepth: depth,
		full: true,
	}
	dir.pointerPageTable[0] = p1
	dir.pointerPageTable[1] = p2
	value := ""
	key := "full"

	// When
	dir.Put(key, value)

	// Then
	assert.Equal(t, 4, len(dir.pointerPageTable))
	pointerAtIndex1 := dir.pointerPageTable[1]
	pointerAtIndex3 := dir.pointerPageTable[3]
	assert.True(t, pointerAtIndex1 == pointerAtIndex3)
}