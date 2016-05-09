package extendible

import (
	"testing"
	"github.com/BenJoyenConseil/bluckdb/util"
	"github.com/stretchr/testify/assert"
	"github.com/BenJoyenConseil/bluckdb/bluckstore/extendible/page"
)

func TestGetPage(t *testing.T) {
	// Given
	dir := &Directory{
		pointerPageTable: make([]*extendible.Page, 2),
		globalDepth: 1,
	}
	key := util.String("123")
	key2 := util.String("122")
	var pointerTable *extendible.Page = &extendible.Page{}
	dir.pointerPageTable[1] = pointerTable

	// When
	result := dir.getPage(key)
	result2 := dir.getPage(key2)

	// Then
	var expected *extendible.Page = pointerTable
	assert.True(t, expected == result)
	assert.Nil(t, result2)
}
