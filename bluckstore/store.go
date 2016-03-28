package bluckstore

import (
	"io/ioutil"
	"strings"
	"os"
	"bufio"
)

type KVStore interface {
	Get(k string) string
	Put(k, v string)
}

type MemKVStore struct {
	pairs map[string]string
}

type DiskKVStore struct {
	filename string
	count int
}

func (store *MemKVStore) Get(k string) string {
	return store.pairs[k]
}

func (store *MemKVStore) Put(k, v string) {
	store.pairs[k] = v
}

func NewMemStore() KVStore {
	return &MemKVStore{make(map[string]string)}
}

func (store *DiskKVStore) Get(k string) string {
	var value string
	body, _ := ioutil.ReadFile(store.filename)
	lines := strings.Split(string(body), "\n")

	for i := range lines {
		keyValuePair := strings.Split(lines[i], ":")
		if keyValuePair[0] == k {
			value = keyValuePair[1]
		}
	}
	if len(value) > 0 {
		return value
	}
	return "no such key"
}

func (store *DiskKVStore) Put(k, value string)  {
	f, err := os.OpenFile(store.filename, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	writer := bufio.NewWriter(f)
	writer.WriteString(k)
	writer.WriteString(":")
	writer.WriteString(value)
	writer.WriteString("\n")
	writer.Flush()
}

func NewDiskStore() KVStore {

	return &DiskKVStore{filename: "/tmp/data.blk", count:0}
}