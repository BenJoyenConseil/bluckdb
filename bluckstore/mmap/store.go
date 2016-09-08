package memap

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/edsrzf/mmap-go"
	"io/ioutil"
	"os"
)

type MmapKVStore struct {
	Dir *Directory
}

const FILE_NAME = "data.db"
const META_FILE_NAME = "db.meta"
const DB_DIRECTORY = "/tmp/"

func New() *MmapKVStore {
	os.Remove(DB_DIRECTORY + FILE_NAME)
	os.Remove(DB_DIRECTORY + META_FILE_NAME)
	store := &MmapKVStore{}
	store.Open()
	return store
}

func (store *MmapKVStore) Open() {
	f, err := os.OpenFile(DB_DIRECTORY+FILE_NAME, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("OpenFile DB error : " + err.Error())
		err = nil
	}

	metaFile, err := os.OpenFile(DB_DIRECTORY+META_FILE_NAME, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("OpenFile Meta error : " + err.Error())
		err = nil
	} else {

		buf, err := ioutil.ReadFile(DB_DIRECTORY + META_FILE_NAME)
		if err != nil {
			fmt.Println("ReadFile error : " + err.Error())
			err = nil
		}
		dec := gob.NewDecoder(bytes.NewBuffer(buf))
		store.Dir = &Directory{}
		err = dec.Decode(&store.Dir)
		if err != nil {
			fmt.Println("Decoding error : " + err.Error())
			store.Dir = &Directory{
				Gd:    0,
				Table: make([]int, 1),
			}
			f.Write(make([]byte, 4096))
			metaFile.Write(store.Dir.serializeMeta())
		}
	}

	store.Dir.metaFile = metaFile
	store.Dir.dataFile = f
	store.Dir.data, err = mmap.Map(store.Dir.dataFile, mmap.RDWR, 0644)
	if err != nil {
		fmt.Println(err)
		err = nil
	}
}

func (s *MmapKVStore) Get(k string) string {
	return s.Dir.get(k)
}

func (s *MmapKVStore) Put(k, v string) {
	s.Dir.put(k, v)
}

func (s *MmapKVStore) Close() {
	s.Dir.data.Unmap()
	s.Dir.dataFile.Close()
	s.Dir.metaFile.Close()
}

func (s *MmapKVStore) Rm() {
	os.Remove(DB_DIRECTORY + FILE_NAME)
	os.Remove(DB_DIRECTORY + META_FILE_NAME)
}
