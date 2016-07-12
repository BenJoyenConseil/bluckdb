package memap

import (
	"testing"
	"github.com/BenJoyenConseil/bluckdb/util"
	"github.com/stretchr/testify/assert"
	"encoding/binary"
	"strconv"
	"os"
	"github.com/edsrzf/mmap-go"
	"io/ioutil"
)

func TestDirectory_ExtendibleHash(t *testing.T) {
	// Given
	d := &Directory{
		gd: 4, // 4 bytes of the key.Hash will be used to distribute keys
	}
	key := "321" //   0011
	key2 := "123" //  1011

	// When
	result := d.extendibleHash(util.String(key))
	result2 := d.extendibleHash(util.String(key2))

	// Then
	assert.Equal(t, 3, result)
	assert.Equal(t, 11, result2)
}

func TestDirectory_GetPage(t *testing.T) {
	// Given
	d := &Directory{
		gd: 4, // means the table size is 2^4 length
		table: make([]int, 16),
		data: make([]byte, 4096 * 4),
	}
	key := "123" //   1011
	d.table[11] = 2

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
		gd: 4, // means we take 4 significant bytes of the hash result
		table: make([]int, 16),
		data: make([]byte, 4096 * 4),
	}
	key := "123" //   1011
	d.table[11] = 2
	fileOffset := 2 * 4096
	binary.LittleEndian.PutUint16(d.data[fileOffset + 4094:], 9)
	binary.LittleEndian.PutUint16(d.data[fileOffset:], 3)
	binary.LittleEndian.PutUint16(d.data[fileOffset + 2:], 2)
	copy(d.data[fileOffset + 4:], []byte(key))
	copy(d.data[fileOffset + 7:], []byte("Hi"))

	// When
	result := d.get(key)

	// Then
	assert.Equal(t, "Hi", result)
}

func TestDirectory_Expand(t *testing.T) {
	// Given
	d := &Directory{
		gd: 3,
		table: make([]int, 8),
	}

	// When
	d.expand()

	// Then
	assert.Equal(t, 16, len(d.table))
	assert.Equal(t, 4, int(d.gd))
}

func TestDirectory_Split(t *testing.T) {
	// Given
	d := &Directory{
		gd: 1,
	}
	page := Page(make([]byte, 4096))
	page.setLd(0)
	fillPage(page, 2)

	// When
	p1, p2 := d.split(page)

	// Then
	assert.Equal(t, "key0value yop yop", string(p1[4:21]))
	assert.Equal(t, 21, p1.use())
	assert.Equal(t, "key1value yop yop", string(p2[4:21]))
	assert.Equal(t, 21, p2.use())
}

func TestDirectory_NextPageId(t *testing.T) {
	// Given
	d := &Directory{
		lastPageId: 29920,
	}

	// When
	result := d.nextPageId()

	// Then
	assert.Equal(t, 29921, result)
}

func TestDirectoryReplace(t *testing.T)  {
	// Given
	dir := &Directory{
		table:[]int{0, 1, 3, 2, 0, 1, 3, 2},
		gd: 2,
		lastPageId: 4,
	}

	// When
	r1, r2 := dir.replace(2, 2)

	// Then
	assert.Equal(t, 2, r1)
	assert.Equal(t, 5, r2)
}

func TestDirectorySerializeMeta(t *testing.T)  {
	// Given
	dir := &Directory{
		table:[]int{0, 1},
		gd: 1,
		lastPageId: 1,
	}

	// When
	result := dir.serializeMeta()

	// Then
	assert.Equal(t, []byte{0x1, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0,}, result)
	assert.Equal(t, 16, len(result))
}

func TestDirectory_Put(t *testing.T) {
	// Given
	fPath := "/tmp/test.db"
	f, _ := os.OpenFile(fPath, os.O_RDWR | os.O_CREATE | os.O_TRUNC, 0644)
	defer f.Close()
	f.Write(make([]byte, 4096))
	dir := &Directory{
		dataFile: f,
		gd: 0,
		table: make([]int, 1),
	}
	dir.table[0] = 0
	dir.data, _ = mmap.Map(dir.dataFile, mmap.RDWR, 0)
	defer dir.data.Unmap()
	var page Page = Page(dir.data[0:4096])
	fillPage(page, 5)

	// When
	dir.put("123", "Yolo !")

	// Then
	buf := make([]byte, 9)
	dir.dataFile.ReadAt(buf, 109)
	assert.Equal(t, "123Yolo !", string(buf))
	os.Remove(fPath)
}

func TestDirectory_PutShouldIncrementLD_WhenPageIsFull(t *testing.T) {
	// Given
	fPath := "/tmp/test.db"
	f, _ := os.OpenFile(fPath, os.O_RDWR | os.O_CREATE | os.O_TRUNC, 0644)
	defer f.Close()
	f.Write(make([]byte, 4096))
	dir := &Directory{
		dataFile: f,
		gd: 0,
		table: make([]int, 1),
	}
	dir.table[0] = 0
	dir.data, _ = mmap.Map(dir.dataFile, mmap.RDWR, 0)
	defer dir.data.Unmap()
	var page Page = Page(dir.data[0:4096])
	fillPage(page, 5)
	binary.LittleEndian.PutUint16(page[4094:], 4090)    // full

	// When
	dir.put("123", "Yolo !")

	// Then
	assert.Equal(t, 8192, len(dir.data))
	assert.Equal(t, 1, Page(dir.data[:4096]).ld())
	assert.Equal(t, 1, Page(dir.data[4096:8192]).ld())
	os.Remove(fPath)
}

func TestDirectory_Put_INT(t *testing.T) {
	// Given
	fPath := "/tmp/test.db"
	metaFPath := "/tmp/metaTest.db"
	f, _ := os.OpenFile(fPath, os.O_RDWR | os.O_CREATE | os.O_TRUNC, 0644)
	metaF, _ := os.OpenFile(metaFPath, os.O_RDWR | os.O_CREATE | os.O_TRUNC, 0644)
	defer f.Close()
	defer metaF.Close()
	f.Write(make([]byte, 4096))
	dir := &Directory{
		dataFile: f,
		gd: 0,
		table: make([]int, 1),
		metaFile: metaF,
	}
	dir.table[0] = 0
	dir.data, _ = mmap.Map(dir.dataFile, mmap.RDWR, 0)
	defer dir.data.Unmap()

	// When
	for i := 0; i < 4000; i++ {
		dir.put("key" + strconv.Itoa(i), "Yolo !")
	}

	// Then
	assert.Equal(t, []int{0, 1, 2, 3, 5, 7, 6, 4, 15, 12, 8, 14, 10, 9, 13, 11, 0, 22, 2, 3, 18, 7, 25, 19, 20, 24, 16, 21, 10, 17, 23, 11}, dir.table)

	os.Remove(fPath)
	os.Remove(metaFPath)
}

func TestDirectory_Put_ShouldWriteMetaFileOnDisk_WhenTableIsExpanded(t *testing.T) {
	// Given
	fPath := "/tmp/test.db"
	metaFPath := "/tmp/metaTest.db"
	f, _ := os.OpenFile(fPath, os.O_RDWR | os.O_CREATE | os.O_TRUNC, 0644)
	metaF, _ := os.OpenFile(metaFPath, os.O_RDWR | os.O_CREATE | os.O_TRUNC, 0644)
	defer f.Close()
	defer metaF.Close()
	f.Write(make([]byte, 4096))
	dir := &Directory{
		dataFile: f,
		gd: 0,
		table: make([]int, 1),
		metaFile: metaF,
	}
	dir.table[0] = 0
	dir.data, _ = mmap.Map(dir.dataFile, mmap.RDWR, 0)
	defer dir.data.Unmap()

	// When
	for i := 0; i < 1000; i++ {

		dir.put("key" + strconv.Itoa(i), "Yolo !")
	}

	// Then
	result, _ := ioutil.ReadFile(metaFPath)
	assert.NotEmpty(t, result)

	os.Remove(fPath)
	os.Remove(metaFPath)
}


func TestDirectory_Put_SameKey(t *testing.T) {
	// Given
	fPath := "/tmp/test.db"
	f, _ := os.OpenFile(fPath, os.O_RDWR | os.O_CREATE | os.O_TRUNC, 0644)
	defer f.Close()
	f.Write(make([]byte, 4096))
	dir := &Directory{
		dataFile: f,
		gd: 0,
		table: make([]int, 1),
	}
	dir.table[0] = 0
	dir.data, _ = mmap.Map(dir.dataFile, mmap.RDWR, 0)
	defer dir.data.Unmap()

	// When
	for i := 0; i < 5; i++ {

		dir.put("key", "Yolo ! " + strconv.Itoa(i))
	}

	// Then BOUM !!!
	assert.Equal(t, "Yolo ! 4", dir.get("key"))
	os.Remove(fPath)
}

func fillPage(page Page, numberOfRecord int) {
	for i := 0; i < numberOfRecord; i++{
		itoa := strconv.Itoa(i)
		page.put("key" + itoa, "value yop yop")
	}
}
