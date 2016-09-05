package memap

import (
	"os"
	"github.com/edsrzf/mmap-go"
	"fmt"
	"io/ioutil"
	"encoding/gob"
	"bytes"
)

type MmapKVStore struct {
	dir *Directory
}

const FILE_NAME = "data.db"
const META_FILE_NAME = "db.meta"

func New() *MmapKVStore {
	os.Remove("/tmp/" + FILE_NAME)
	os.Remove("/tmp/" + META_FILE_NAME)
	store := &MmapKVStore{}
	store.Open()
	return store
}

func (store *MmapKVStore) Open() {
	f, err := os.OpenFile("/tmp/" + FILE_NAME, os.O_RDWR | os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("OpenFile DB error : " + err.Error())
		err = nil
	}

	metaFile, err := os.OpenFile("/tmp/" + META_FILE_NAME, os.O_RDWR | os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("OpenFile Meta error : " + err.Error())
		err = nil
	} else {

		buf, err := ioutil.ReadFile("/tmp/" + META_FILE_NAME)
		if err != nil {
			fmt.Println("ReadFile error : " + err.Error())
			err = nil
		}
		dec := gob.NewDecoder(bytes.NewBuffer(buf))
		store.dir = &Directory{}
		err = dec.Decode(&store.dir)
		if err != nil {
			fmt.Println("Decoding error : " + err.Error())
			store.dir = &Directory{
				Gd: 0,
				Table: make([]int, 1),
			}
			f.Write(make([]byte, 4096))
			metaFile.Write(store.dir.serializeMeta())
		}
	}


	store.dir.metaFile = metaFile
	store.dir.dataFile = f
	store.dir.data, err = mmap.Map(store.dir.dataFile, mmap.RDWR, 0644)
	if err != nil {
		fmt.Println(err)
		err = nil
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
	s.dir.metaFile.Close()
}

func (s *MmapKVStore) Rm()  {
	os.Remove("/tmp/" + FILE_NAME)
	os.Remove("/tmp/" + META_FILE_NAME)
}