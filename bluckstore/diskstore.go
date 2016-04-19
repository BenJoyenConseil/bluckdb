package bluckstore

import (
	"io/ioutil"
	"strings"
	"os"
	"bufio"
	"strconv"
	"github.com/BenJoyenConseil/bluckdb/util"
)

const extension = ".blk"
const filename = "data"
const path = "/tmp/"

type DiskKVStore struct {}

func (store *DiskKVStore) Get(k string) string {
	var value string
	file := partitionFile(util.String(k))
	body, _ := ioutil.ReadFile(file)
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

	file := partitionFile(util.String(k))
	f, err := os.OpenFile(file, os.O_APPEND|os.O_RDWR, 0600)
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

	for i := 0; i < BUCKET_NUMER; i++ {
		file := partitionFile(util.String(strconv.Itoa(i)))
		f, _ := os.OpenFile(file, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0600)
		f.Close()
	}
	return &DiskKVStore{}
}


func partitionFile(k util.String) string{
	bucketId := strconv.Itoa(k.Hash() % BUCKET_NUMER)
	return path + filename + bucketId + extension
}
