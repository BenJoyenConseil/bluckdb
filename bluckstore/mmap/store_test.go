package memap

import (
    "os"
    "strconv"
    "github.com/stretchr/testify/assert"
    "testing"
    "math/rand"
)

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
    assert.Equal(t, []int{0}, store.dir.Table)
    assert.Equal(t, 0, int(store.dir.Gd))
    assert.Equal(t, "KEYVALUE", string(store.dir.data[0:8]))
    store.Close()
}

func BenchmarkMmapPut(b *testing.B) {
    /*devNull, _ := os.Open(os.DevNull)
    store := &MmapKVStore{
        dir: &Directory{
            Table: make([]int, 1),
            dataFile: devNull,
            metaFile: devNull,
            Gd: 0,
            data: make([]byte, 4096),
        },
    }*/
    store := &MmapKVStore{}
    store.Rm()
    store.Open()
    defer store.Close()

    b.ResetTimer()
    for n := 0; n < b.N; n++ {
        store.Put(strconv.Itoa(n), "mec, elle est où ma caisse ??")
    }

}

func BenchmarkHashmapPut(b *testing.B) {
    store := make(map[string]string)

    b.ResetTimer()
    for n := 0; n < b.N; n++ {
        store[strconv.Itoa(n)] = "mec, elle est où ma caisse ??"
    }
}

func setup(){
    store := &MmapKVStore{}
    store.Rm()
    store.Open()
    size := 10000
    for i := 0; i < size; i++ {
        store.Put(strconv.Itoa(i), "mec, elle est où ma caisse ??")
    }
    store.Close()
}

func BenchmarkMmapGet(b *testing.B) {
    //setup()
    store := &MmapKVStore{}
    store.Open()
    defer store.Close()

    b.ResetTimer()
    for n := 0; n < b.N; n++ {
        store.Get(strconv.Itoa(rand.Intn(10000 - 1)))
    }
}

func BenchmarkHashmapGet(b *testing.B) {
    store := make(map[string]string)
    size := 10000
    for i := 0; i < size; i++ {
        store[strconv.Itoa(i)] = "mec, elle est où ma caisse ??"
    }
    devNull, _ := os.Open(os.DevNull)
    defer devNull.Close()

    b.ResetTimer()
    for n := 0; n < b.N; n++ {
        devNull.Write([]byte(store[strconv.Itoa(rand.Intn(10000 - 1))]))
    }
}
