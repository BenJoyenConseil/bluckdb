package mmap

import (
	"bytes"
	"encoding/gob"
	"github.com/edsrzf/mmap-go"
	"github.com/labstack/gommon/log"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"os"
	"strconv"
	"testing"
)

func TestStoreOpen_shouldReOpen_UsingMeta(t *testing.T) {
	// Given
	os.Remove(DB_DEFAULT_FOLDER + FILE_NAME)
	os.Remove(DB_DEFAULT_FOLDER + META_FILE_NAME)
	store := MmapKVStore{}
	store.Open(DB_DEFAULT_FOLDER)
	store.Put("KEY", "VALUE")
	store.Close()

	// When
	store.Open(DB_DEFAULT_FOLDER)

	// Then
	assert.Equal(t, PAGE_SIZE, len(store.Dir.data))
	assert.Equal(t, DB_DEFAULT_FOLDER+FILE_NAME, store.Dir.dataFile.Name())
	assert.Equal(t, []int{0}, store.Dir.Table)
	assert.Equal(t, 0, int(store.Dir.Gd))
	assert.Equal(t, "KEYVALUE", string(store.Dir.data[0:8]))
	store.Close()
}

func TestMmapKVStore_Close_ShouldCreateMetadataFile_UsingTheStoreFolderAttribute(t *testing.T) {
	// Given
	store := &MmapKVStore{
		Dir: &Directory{
			Gd:         2,
			LastPageId: 4,
		},
		Path: "/tmp/",
	}

	// When
	store.Close()

	// Then
	f, err := os.Open(store.Path + META_FILE_NAME)
	assert.NotNil(t, f)
	assert.Nil(t, err)
	assert.Nil(t, os.Remove(store.Path+"bluck.meta"))
	os.RemoveAll("/tmp/bluckdb")
}

func TestDecodeMeta(t *testing.T) {
	// Given
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	enc.Encode(Directory{
		Gd:         1,
		LastPageId: 2,
		Table:      []int{0, 1, 0, 2},
	})

	// When
	result := DecodeMeta(&buff)

	// Then
	assert.Equal(t, 2, result.LastPageId)
	assert.Equal(t, []int{0, 1, 0, 2}, result.Table)
	assert.Equal(t, 1, int(result.Gd))
}

func TestEncodeMeta(t *testing.T) {
	// Given
	dir := &Directory{
		Gd:         1,
		LastPageId: 2,
		Table:      []int{0, 1, 0, 2},
	}

	// When
	result := EncodeMeta(dir)

	// Then
	var r Directory
	dec := gob.NewDecoder(result)
	dec.Decode(&r)
	assert.Equal(t, *dir, r)
}

func TestMmapKVStore_RestoreMETA_shouldReOpen_UsingFileStatToBuildMeta(t *testing.T) {
	// Given
	os.Remove(DB_DEFAULT_FOLDER + FILE_NAME)
	f, _ := os.OpenFile(DB_DEFAULT_FOLDER + FILE_NAME, os.O_RDWR|os.O_CREATE, 0644)
	f.Write(make([]byte, 16384))
	f.Close()
	store := &MmapKVStore{}

	// When
	store.Open(DB_DEFAULT_FOLDER)

	// Then
	assert.Equal(t, 16384, len(store.Dir.data))
	assert.Equal(t, []int{0, 1, 2, 0}, store.Dir.Table)
	assert.Equal(t, 2, int(store.Dir.Gd))
	assert.Equal(t, 2, store.Dir.LastPageId)
	store.Close()
}

func TestMmapKVStore_Open_shouldCreateNewFileWhenNotExisting_UsingFSFolder(t *testing.T) {
	// Given
	fsFolder := "/tmp/bluckdb/exotic/path/"
	os.RemoveAll(fsFolder)
	store := &MmapKVStore{}

	// When
	store.Open(fsFolder)

	// Then
	assert.Equal(t, PAGE_SIZE, len(store.Dir.data))
	assert.Equal(t, []int{0}, store.Dir.Table)
	assert.Equal(t, 0, int(store.Dir.Gd))
	assert.Equal(t, 0, store.Dir.LastPageId)
	assert.Nil(t, os.Remove(fsFolder+"bluck.data"))
	store.Close()
	os.RemoveAll("/tmp/bluckdb")
}

func TestNextPowerOfTwoNom(t *testing.T) {
	// Given
	numBuckets := uint(15921)

	// When
	result := NextPowerOfTwo(numBuckets)

	// Then
	assert.Equal(t, int(16384), int(result))
}

func TestFindTwoToPowerOfN(t *testing.T) {
	// Given
	numBuckets := uint(15921)

	// When
	result := FindTwoToPowerOfN(numBuckets)

	// Then
	assert.Equal(t, uint(14), result)
}

func TestFindBucketNumber(t *testing.T) {
	// Given
	fileSize := int64(65216512)

	// When
	result := FindBucketNumber(fileSize)

	// Then
	assert.Equal(t, int64(15922), result)
}

func TestMmapKVStore_DumpPage(t *testing.T) {
	// Given
	store := &MmapKVStore{
		Dir: &Directory{
			data:  mmap.MMap(make([]byte, 4096)),
			Table: make([]int, 1),
		},
	}
	copy(store.Dir.data, "12345salut!")

	// When
	result := store.DumpPage(0)

	// Then
	assert.Contains(t, result, "12345salut!")
}

func TestMmapKVStore_Put_WhenRecordPayloadIsToBig(t *testing.T) {
	// Given
	store := &MmapKVStore{}

	// When
	err := store.Put("1234", string(make([]byte, 2046)))
	err = store.Put("1234", string(make([]byte, 2046)))

	// Then
	assert.NotNil(t, err)
}

func BenchmarkMmapPut(b *testing.B) {
	p := "/tmp/bluckdb/bench/"
	log.SetLevel(log.OFF)
	store := &MmapKVStore{}
	store.Rm(p)
	store.Open(p)
	defer store.Close()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		store.Put(strconv.Itoa(n), "mec, elle est où ma caisse ??")
	}

}

func BenchmarkMmapRangePut(b *testing.B) {
	p := "/tmp/bluckdb/bench/"
	log.SetLevel(log.OFF)
	store := &MmapKVStore{}
	store.Rm(p)
	store.Open(p)
	defer store.Close()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for i := 0; i < 1000000; i++ {

			store.Put(strconv.Itoa(i), "mec, elle est où ma caisse ??")
		}
	}

}

func setup() {
	p := "/tmp/bluckdb/bench/"
	store := &MmapKVStore{}
	store.Rm(p)
	store.Open(p)
	size := 10000
	for i := 0; i < size; i++ {
		store.Put(strconv.Itoa(i), "mec, elle est où ma caisse ??")
	}
	store.Close()
}

func BenchmarkMmapGet(b *testing.B) {
	log.SetLevel(log.OFF)
	b.StopTimer()
	setup()
	store := &MmapKVStore{}
	store.Open("/tmp/bluckdb/bench/")
	defer store.Close()

	b.ResetTimer()
	b.StartTimer()
	for n := 0; n < b.N; n++ {
		store.Get(strconv.Itoa(rand.Intn(1000000 - 1)))
	}
}
