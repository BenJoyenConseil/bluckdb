package memap

import (
	"os"
	"github.com/edsrzf/mmap-go"
)

type MmapKVStore struct {
	dir *Directory
}

const FILE_NAME = "data.db"

func New(pathDB string) *MmapKVStore {
	f, _ := os.OpenFile(pathDB + FILE_NAME, os.O_RDWR | os.O_CREATE | os.O_TRUNC, 0644)
	f.Write(make([]byte, 4096))
	dir := &Directory{
		dataFile: f,
		gd: 0,
		table: make([]int, 1),
	}
	dir.table[0] = 0
	dir.data, _ = mmap.Map(dir.dataFile, mmap.RDWR, 0)

	return &MmapKVStore{
		dir: dir,
	}
}

func (s *MmapKVStore) Get(k string) string {
	return s.dir.get(k)
}

func (s *MmapKVStore) Put(k, v string) {
	s.dir.put(k, v)
}

func (s *MmapKVStore) Close() {
	s.dir.data.Unmap()
	s.dir.dataFile.Close()
}

