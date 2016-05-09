package extendible

import (
	"testing"
	"github.com/BenJoyenConseil/bluckdb/util"
	"github.com/stretchr/testify/assert"
	page "github.com/BenJoyenConseil/bluckdb/bluckstore/extendible/page"
	record "github.com/BenJoyenConseil/bluckdb/bluckstore/extendible/page/record"
	"errors"
)


type stubPage struct{content []byte}
func (stub *stubPage) Get(key string) (string, error){
	if key == "123" {return "Hello", nil}
	return "Stub", errors.New("Stub !!")
}
func (stub *stubPage) Put(key, value string) error {return errors.New("Stub !")}
func (stub *stubPage) Full(record record.Record) bool{return true}

func TestGetPage(t *testing.T) {
	// Given
	dir := &Directory{
		pointerPageTable: make([]page.Page, 2),
		globalDepth: 1,
	}
	key := util.String("123")
	key2 := util.String("122")
	var pointerTable *stubPage = &stubPage{}
	dir.pointerPageTable[1] = pointerTable

	// When
	result := dir.getPage(key)
	result2 := dir.getPage(key2)

	// Then
	var expected *stubPage = pointerTable
	assert.True(t, expected == result)
	assert.Nil(t, result2)
}

func TestGet(t *testing.T) {
	// Given
	dir := &Directory{
		pointerPageTable: make([]page.Page, 2),
		globalDepth: 1,
	}
	var pointerTable *stubPage = &stubPage{}
	dir.pointerPageTable[1] = pointerTable

	// When
	result, err := dir.Get("123")

	// Then
	assert.Equal(t, "Hello", result)
	assert.Nil(t, err)
}