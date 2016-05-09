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
	return "", errors.New("Not implemented")
}