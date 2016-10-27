package mmap

import (
	"bytes"
	"encoding/gob"
	"errors"
	"github.com/labstack/gommon/log"
	"io/ioutil"
	"os"
)

type MmapKVStore struct {
	Dir  *Directory
	Path string
}

const (
	FILE_NAME         = "bluck.data"
	META_FILE_NAME    = "bluck.meta"
	DB_DEFAULT_FOLDER = "/tmp/"
)

//
// Open create the datafile and the metadata file if they do not exist.
// Else if they exist, it loads from the disk and mmap the datafile.
//
func (store *MmapKVStore) Open(absolutePath string) {

	store.Path = absolutePath
	err := os.MkdirAll(store.Path, 0777)
	if err != nil {
		log.Errorf("Mkdir %s impossible : %s", absolutePath, err.Error())
	}

	dataFileName := absolutePath + FILE_NAME
	f, err := os.OpenFile(dataFileName, os.O_RDWR, 0644)
	defer f.Close()

	if err != nil {
		log.Warnf("OpenFile DB : %s", err.Error())
		log.Infof("Try to create the file : %s", dataFileName)

		f, _ = os.Create(dataFileName)
		store.Dir = &Directory{
			Gd:    0,
			Table: make([]int, 1),
		}
		f.Write(make([]byte, PAGE_SIZE))

		// TODO : flush metadata

	} else {
		log.Infof("Datafile %s detected", dataFileName)
		meta, err := ioutil.ReadFile(absolutePath + META_FILE_NAME)

		if err != nil {
			log.Errorf("Error while trying to OpenFile Metadata : %s", err.Error())
			panic(err)
		} else {
			store.Dir = DecodeMeta(bytes.NewBuffer(meta))
		}
	}

	store.Dir.dataFile = f
	store.Dir.mmapDataFile()
}

func (s *MmapKVStore) Get(k string) string {
	return s.Dir.get(k)
}

//
// Put inserts data into the memory which is mapped on disk, but does not persit the metadata
// If the store crashes in an inconsistent way (metadata != data), you need to use the recovery tool (RestoreMETA func)
//
func (s *MmapKVStore) Put(k, v string) error {
	if len(k)+len(v) > 2045 {
		return errors.New("The record is too long")
	}

	s.Dir.put(k, v)
	return nil
}

//
// Close must be called at the end of the connection with the store, to persist safely the data and to persist metadata on disk.
// The persistence of the metadata is only done in this method.
// Usage : defer store.Close()
//
func (s *MmapKVStore) Close() {
	s.Dir.data.Flush()
	s.Dir.data.Unmap()
	s.Dir.dataFile.Close()

	metaFile, err := os.OpenFile(s.Path+META_FILE_NAME, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	defer metaFile.Close()
	if err != nil {
		log.Error(err)
		err = nil
	}

	written, err := metaFile.WriteAt(EncodeMeta(s.Dir).Bytes(), 0)
	if err != nil {
		log.Errorf("Metadata not persisted (bytes written %d) : %s", written, err.Error())
	} else {
		log.Info("Store closing : metadata are persisted")
	}
}

func (s *MmapKVStore) Rm() {
	os.Remove(s.Path + FILE_NAME)
	os.Remove(s.Path + META_FILE_NAME)
}

func (s *MmapKVStore) RestoreMETA() {
	/*

		stat, _ := f.Stat()
		numBuckets := FindBucketNumber(stat.Size())
		tableSize := NextPowerOfTwo(uint(numBuckets))

		store.Dir = &Directory{
			Gd: FindTwoToPowerOfN(uint(numBuckets)),
			Table: make([]int, tableSize),
			LastPageId: int(numBuckets) - 1,
		}

		for i := 0; i < int(tableSize); i ++ {
			if i >= int(numBuckets) {
				store.Dir.Table[i] = 0
			} else {
				store.Dir.Table[i] = i
			}
		}
	*/
}

func FindBucketNumber(fileSize int64) int64 {
	return fileSize / int64(PAGE_SIZE)
}

func FindTwoToPowerOfN(v uint) uint {
	for i := uint(1); ; i++ {
		if (v >> i) <= 0 {
			return i
			break
		}
	}
	return 0 // never
}

func NextPowerOfTwo(v uint) uint {
	v--
	v |= v >> 1
	v |= v >> 2
	v |= v >> 4
	v |= v >> 8
	v |= v >> 16
	v++
	return v
}

func DecodeMeta(buff *bytes.Buffer) *Directory {
	dec := gob.NewDecoder(buff)
	var dir Directory
	dec.Decode(&dir)
	return &dir
}

func EncodeMeta(dir *Directory) *bytes.Buffer {
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	enc.Encode(&dir)
	return &buff
}

func (s *MmapKVStore) DumpPage(pageId int) string {
	start := pageId * 4096
	end := start + 4096
	return string(s.Dir.data[start:end])
}

func (s *MmapKVStore) Meta() *Directory {
	return s.Dir
}
