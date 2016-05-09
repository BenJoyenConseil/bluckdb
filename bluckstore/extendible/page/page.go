package extendible

import (
	"errors"
	"github.com/BenJoyenConseil/bluckdb/bluckstore/extendible/page/record"
)

const (
	PAGE_SIZE = uint16(4096) // byte
)

type Page struct {
	localDepth uint64
	content []byte
	use uint16
	recordReader extendible.RecordReader
}

func New() *Page{
	return &Page{
		localDepth: 0,
		content: make([]byte, PAGE_SIZE),
		use: 0,
	}
}

func (self *Page) Full(record extendible.Record) bool {
	if record.Payload() + self.use > PAGE_SIZE {
		return true
	}
	return false
}

func (self *Page) Put(key, value string) error {
	return errors.New("Not implemented")
}

func (self * Page) Get(key string) (string, error) {

	var record extendible.Record
	offset := 0

	for offset < len(self.content){
		record = self.recordReader.Read(self.content[offset:])
		offset += int(record.Payload())
		if string(record.Key()) == key {
			return string(record.Value()), nil
		}
	}

	return "", errors.New("Key not found : " + key)
}