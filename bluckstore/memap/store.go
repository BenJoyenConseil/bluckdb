package memap

import (
	"os"
	"github.com/edsrzf/mmap-go"
	"fmt"
	"encoding/binary"
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
		fmt.Println(err)
		err = nil
	}

	metaFile, err := os.OpenFile("/tmp/" + META_FILE_NAME, os.O_RDWR | os.O_CREATE, 0644)
	if err != nil {
		fmt.Println(err)
		err = nil
	} else {
		fStat, _ := metaFile.Stat()
		metaContent := make([]byte, fStat.Size())
		metaFile.Read(metaContent)

		gd, lastId, table := UnMarshallMeta(metaContent)
		if table != nil {
			store.dir = &Directory{
				table: table,
				gd: gd,
				lastPageId: lastId,
			}
		} else {
			store.dir = &Directory{
				gd: 0,
				table: make([]int, 1),
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



func UnMarshallMeta(data []byte) (gd uint, lastId int, table []int) {
	if len(data) <= 0 {
		return 0, 0, nil
	}
	gd = uint(binary.LittleEndian.Uint32(data))
	lastId = int(binary.LittleEndian.Uint32(data[4:]))
	for i := 8; i < len(data) ; i += 4 {
		id := int(binary.LittleEndian.Uint32(data[i:]))
		table = append(table, id)
	}
	return gd, lastId, table
}