package memap

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/BenJoyenConseil/bluckdb/util"
	"github.com/edsrzf/mmap-go"
	"github.com/stretchr/testify/assert"
	"os"
	"strconv"
	"testing"
)

func TestDirectory_ExtendibleHash(t *testing.T) {
	// Given
	d := &Directory{
		Gd: 4, // 4 bytes of the key.Hash will be used to distribute keys
	}
	key := "321"  //   0011
	key2 := "123" //  1011

	// When
	result := d.extendibleHash(util.Key(key))
	result2 := d.extendibleHash(util.Key(key2))

	// Then
	assert.Equal(t, 3, result)
	assert.Equal(t, 11, result2)
}

func TestDirectory_GetPage(t *testing.T) {
	// Given
	d := &Directory{
		Gd:    4, // means the table size is 2^4 length
		Table: make([]int, 16),
		data:  make([]byte, 4096*4),
	}
	key := "123" //   1011
	d.Table[11] = 2

	// When
	page, idPage := d.getPage(key)

	// Then
	assert.Equal(t, 2, idPage)
	assert.Equal(t, 8192, cap(page))
	assert.Equal(t, 4096, len(page))
}

func TestDirectory_Get(t *testing.T) {
	// Given
	d := &Directory{
		Gd:    4, // means we take 4 significant bytes of the hash result
		Table: make([]int, 16),
		data:  make([]byte, 4096*4),
	}
	d.Table[11] = 2
	pageOffset := 2 * 4096
	binary.LittleEndian.PutUint16(d.data[pageOffset+4094:], 9) // use

	binary.LittleEndian.PutUint16(d.data[pageOffset+7:], 3) // keyLen
	binary.LittleEndian.PutUint16(d.data[pageOffset+5:], 2) // valLen
	key := "123"                                            //   1011
	copy(d.data[pageOffset:], []byte(key))
	copy(d.data[pageOffset+3:], []byte("Hi"))

	// When
	result := d.get(key)

	// Then
	assert.Equal(t, "Hi", result)
}

func TestDirectory_Expand(t *testing.T) {
	// Given
	d := &Directory{
		Gd:    3,
		Table: make([]int, 8),
	}

	// When
	d.expand()

	// Then
	assert.Equal(t, 16, len(d.Table))
	assert.Equal(t, 4, int(d.Gd))
}

func TestDirectory_Split(t *testing.T) {
	// Given
	d := &Directory{
		Gd: 1,
	}
	page := Page(make([]byte, 4096))
	page.setLd(0)
	fillPage(page, 2)

	// When
	p1, p2 := d.split(page)

	// Then
	assert.Equal(t, "key0value yop yop", string(p1[0:17]))
	assert.Equal(t, 21, p1.use())
	assert.Equal(t, "key1value yop yop", string(p2[0:17]))
	assert.Equal(t, 21, p2.use())
}

func TestDirectory_Split_ShouldSkipRecordWhenHasBeenAlreadyRead(t *testing.T) {
	// Given
	d := &Directory{
		Gd: 1,
	}
	page := Page(make([]byte, 4096))
	page.setLd(0)
	fillPage(page, 2)
	page.put("key0", "value updated")

	// When
	p1, p2 := d.split(page)

	// Then
	assert.Equal(t, "key0value updated", string(p1[0:17]))
	assert.Equal(t, 21, p1.use())
	assert.Equal(t, "key1value yop yop", string(p2[0:17]))
	assert.Equal(t, 21, p2.use())
}

func TestDirectory_NextPageId(t *testing.T) {
	// Given
	d := &Directory{
		LastPageId: 29920,
	}

	// When
	result := d.nextPageId()

	// Then
	assert.Equal(t, 29921, result)
}

func TestDirectoryReplace(t *testing.T) {
	// Given
	dir := &Directory{
		Table:      []int{0, 1, 3, 2, 0, 1, 3, 2},
		Gd:         2,
		LastPageId: 4,
	}

	// When
	r1, r2 := dir.replace(2, 2)

	// Then
	assert.Equal(t, 2, r1)
	assert.Equal(t, 5, r2)
}

func TestDirectory_Put(t *testing.T) {
	// Given
	fPath := "/tmp/test.db"
	f, _ := os.OpenFile(fPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	defer f.Close()
	f.Write(make([]byte, 4096))
	dir := &Directory{
		dataFile: f,
		Gd:       0,
		Table:    make([]int, 1),
	}
	dir.Table[0] = 0
	dir.data, _ = mmap.Map(dir.dataFile, mmap.RDWR, 0)
	defer dir.data.Unmap()
	var page Page = Page(dir.data[0:4096])
	fillPage(page, 5)

	// When
	dir.put("123", "Yolo !")

	// Then
	result := make([]byte, 9)
	dir.dataFile.ReadAt(result, 105)
	assert.Equal(t, "123Yolo !", string(result))
	os.Remove(fPath)
}

func TestDirectory_PutShouldIncrementLD_WhenPageIsFull(t *testing.T) {
	// Given
	fPath := "/tmp/test.db"
	f, _ := os.OpenFile(fPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	defer f.Close()
	f.Write(make([]byte, 4096))
	dir := &Directory{
		dataFile: f,
		Gd:       0,
		Table:    make([]int, 1),
	}
	dir.Table[0] = 0
	dir.data, _ = mmap.Map(dir.dataFile, mmap.RDWR, 0)
	defer dir.data.Unmap()

	var page Page = Page(dir.data[0:4096])
	fillPage(page, 182) // set Page to 4076

	// When
	dir.put("12345678", "Yolo !")

	// Then
	assert.Equal(t, 8192, len(dir.data))
	assert.Equal(t, 1, int(dir.Gd))
	assert.Equal(t, 1, Page(dir.data[:4096]).ld())
	assert.Equal(t, 1, Page(dir.data[4096:8192]).ld())
	os.Remove(fPath)
}

func TestDirectory_Put_INT(t *testing.T) {
	// Given
	fPath := "/tmp/test.db"
	f, _ := os.OpenFile(fPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	defer f.Close()
	f.Write(make([]byte, 4096))
	dir := &Directory{
		dataFile: f,
		Gd:       0,
		Table:    make([]int, 1),
	}
	dir.Table[0] = 0
	dir.data, _ = mmap.Map(dir.dataFile, mmap.RDWR, 0)
	defer dir.data.Unmap()

	// When
	for i := 0; i < 4000; i++ {
		dir.put("key"+strconv.Itoa(i), "Yolo !")
	}

	// Then
	assert.Equal(t, []int{0, 1, 2, 3, 5, 7, 6, 4, 15, 12, 8, 14, 10, 9, 13, 11, 0, 22, 2, 3, 18, 7, 25, 19, 20, 24, 16, 21, 10, 17, 23, 11}, dir.Table)

	os.Remove(fPath)
}

func TestDirectory_Put_SameKey(t *testing.T) {
	// Given
	fPath := "/tmp/test.db"
	f, _ := os.OpenFile(fPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	defer f.Close()
	f.Write(make([]byte, 4096))
	dir := &Directory{
		dataFile: f,
		Gd:       0,
		Table:    make([]int, 1),
	}
	dir.Table[0] = 0
	dir.data, _ = mmap.Map(dir.dataFile, mmap.RDWR, 0)
	defer dir.data.Unmap()

	// When
	for i := 0; i < 5; i++ {

		dir.put("key", "Yolo ! "+strconv.Itoa(i))
	}

	// Then BOUM !!!
	assert.Equal(t, "Yolo ! 4", dir.get("key"))
	os.Remove(fPath)
}

func fillPage(page Page, numberOfRecord int) {
	for i := 0; i < numberOfRecord; i++ {
		itoa := strconv.Itoa(i)
		page.put("key"+itoa, "value yop yop")
	}
}

func TestDirectory_String(t *testing.T) {
	// Given
	dir := &Directory{
		Gd:         1,
		Table:      []int{0, 1, 0, 1},
		LastPageId: 1,
	}

	// When
	result := new(bytes.Buffer)
	fmt.Fprint(result, dir)

	// Then BOUM !!!
	assert.Equal(t, "{\"table\":[0,1,0,1],\"globalDepth\":1,\"LastPageId\":1}\n", result.String())
}

func BenchmarkDirectory_Put(b *testing.B) {

}

/*
func TestDirectory_Gc(t *testing.T) {
	// Given
	fPath := "/tmp/test.db"
	metaFPath := "/tmp/metaTest.db"
	os.Remove(fPath)
	os.Remove(metaFPath)
	f, _ := os.OpenFile(fPath, os.O_RDWR | os.O_CREATE | os.O_TRUNC, 0644)
	metaF, _ := os.OpenFile(metaFPath, os.O_RDWR | os.O_CREATE | os.O_TRUNC, 0644)
	defer f.Close()
	defer metaF.Close()
	f.Write(make([]byte, 4096))
	dir := &Directory{
		dataFile: f,
		Gd: 0,
		Table: make([]int, 1),
		metaFile: metaF,
	}
	dir.Table[0] = 0
	dir.data, _ = mmap.Map(dir.dataFile, mmap.RDWR, 0)
	defer dir.data.Unmap()

	// When
	for i := 0; i < 1000; i++ {

		dir.put("key" + strconv.Itoa(rand.Intn(20)), "Yolo ! " + strconv.Itoa(i))
	}

	// Then BOUM !!!
	assert.Equal(t, 1, len(dir.Table))
	assert.Equal(t, 0, dir.Gd)
	assert.Equal(t, 4096, len(dir.data))
}*/
