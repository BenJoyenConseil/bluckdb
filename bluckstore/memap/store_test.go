package memap

import (
	"os"
	"strconv"
	"github.com/stretchr/testify/assert"
	"testing"
	"io/ioutil"
)



func TestOpen(t *testing.T) {
	// Given
	content := []byte{0x1, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x2, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x2, 0x0, 0x0, 0x0,}
	ioutil.WriteFile("/tmp/data.db", make([]byte, 8192), 0644)
	ioutil.WriteFile("/tmp/db.meta", []byte(content), 0644)
	store := MmapKVStore{}
    defer store.Close()

    // When
    store.Open()

	// Then
	assert.Equal(t, 8192, len(store.dir.data))
	assert.Equal(t, "/tmp/data.db", store.dir.dataFile.Name())
	assert.Equal(t, []int{2, 1, 2}, store.dir.table)
	assert.Equal(t, 1, int(store.dir.gd))
    os.Remove("/tmp/data.db")
    os.Remove("/tmp/db.meta")
}

func TestStorePut_shouldReOpen_UsingMeta(t *testing.T) {
    // Given
    os.Remove("/tmp/data.db")
    os.Remove("/tmp/db.meta")
    store := MmapKVStore{}
    store.Open()
    store.Put("KEY", "VALUE")
    store.Close()

    // When
    store.Open()

    // Then
    assert.Equal(t, 4096, len(store.dir.data))
    assert.Equal(t, "/tmp/data.db", store.dir.dataFile.Name())
    assert.Equal(t, []int{0}, store.dir.table)
    assert.Equal(t, 0, int(store.dir.gd))
    assert.Equal(t, "KEYVALUE", string(store.dir.data[4:12]))
    store.Close()
}

func TestStoreUnMarshallMeta(t *testing.T) {
    // Given
    data := []byte{0x1, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0,}

    // When
    gd, lastId, table := UnMarshallMeta(data)

    // Then
    assert.Equal(t, uint(1), gd)
    assert.Equal(t, 1, lastId)
    assert.Equal(t, []int{0, 1}, table)
}

func fill(page Page) {
	for i := 0; i < 185; i++ {
		itoa := strconv.Itoa(i)
		page.put("key"+itoa, "value yop yop")
	}
}

func BenchmarkMemapPut(b *testing.B) {
    store := New()
    defer store.Close()


    b.ResetTimer()
	for i := 0; i < b.N; i++ {

		store.Put(strconv.Itoa(i), "mec, elle est où ma caisse ??")
	}

}

func BenchmarkMemapGet(b *testing.B) {
    store := New()
    defer store.Close()

	for i := 0; i < b.N; i++ {

        store.Put(strconv.Itoa(i), "mec, elle est où ma caisse ??")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {

        store.Get("yolo !! " + strconv.Itoa(i))
	}
}
