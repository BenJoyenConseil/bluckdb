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

func (store *DiskKVStore) Get(key string) string {
	file := buildPartitionFilePathString(util.String(key))
	body, _ := ioutil.ReadFile(file)
	lines := strings.Split(string(body), "\n")

	for i := len(lines) -1; i >= 0 ; i-- {
		keyValuePair := strings.Split(lines[i], ":")
		if keyValuePair[0] == key {
			return keyValuePair[1]
		}
	}
	return "no such key"
}

func (store *DiskKVStore) Put(k, value string)  {

	file := buildPartitionFilePathString(util.String(k))
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
		file := buildPartitionFilePathString(util.String(strconv.Itoa(i)))
		f, _ := os.OpenFile(file, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0600)
		f.Close()
	}
	return &DiskKVStore{}
}


func buildPartitionFilePathString(key util.String) string{
	bucketId := strconv.Itoa(key.Hash() % BUCKET_NUMER)
	return path + filename + bucketId + extension
}
