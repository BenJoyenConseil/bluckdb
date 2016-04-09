package bluckstore

import (
	"io/ioutil"
	"strings"
	"os"
	"bufio"
	"strconv"
	"crypto/md5"
	"encoding/binary"
)

type DiskKVStore struct {
	filename string
	numPartition int
}

func (store *DiskKVStore) Get(k string) string {
	var value string
	body, _ := ioutil.ReadFile(store.filename + strconv.Itoa(consistentHash(k, store.numPartition)) + ".blk")
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

	partition := strconv.Itoa(consistentHash(k, store.numPartition))
	f, err := os.OpenFile(store.filename + partition + ".blk", os.O_APPEND|os.O_RDWR, 0600)
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

	fileNameTemplate := "/tmp/data"
	extension := ".blk"
	numPartition := 10

	for i := 0; i < numPartition; i++ {
		f, _ := os.OpenFile(fileNameTemplate + strconv.Itoa(i) + extension, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0600)
		f.Close()
	}
	return &DiskKVStore{filename: fileNameTemplate, numPartition:numPartition}
}

func consistentHash(k string, numPartition int) int{
	h := md5.New()
	h.Write([]byte(k))

	r := binary.LittleEndian.Uint32(h.Sum(nil))
	i := int(r) % numPartition
	return i
}
