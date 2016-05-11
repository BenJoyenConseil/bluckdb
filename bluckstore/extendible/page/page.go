package extendible

import (
	"errors"
	"github.com/BenJoyenConseil/bluckdb/bluckstore/extendible/page/record"
)

const (
	PAGE_DISK_SIZE = uint16(4096) // bytes
)

type PageDisk struct {
	localDepth uint64
	content []byte
	use uint16
	recordUnserializer extendible.RecordUnserializer
	recordSerializer extendible.RecordSerializer
}

type Page interface {
	Full(record extendible.Record) bool
	Put(key, value string) error
	Get(key string) (string, error)
	LocalDepth() uint64
	SetLocalDepth(num uint64)
	Content() []byte
}

func New() *PageDisk {
	return &PageDisk{
		localDepth: 0,
		content: make([]byte, PAGE_DISK_SIZE),
		use: 0,
		recordSerializer: &extendible.ByteRecordSerializer{},
		recordUnserializer: &extendible.ByteRecordUnserializer{},
	}
}

func (self *PageDisk) Full(record extendible.Record) bool {
	if record.Payload() + self.use > PAGE_DISK_SIZE {
		return true
	}
	return false
}

func (self *PageDisk) Put(key, value string) error {

	record := extendible.New(key, value)

	if self.Full(record) {
		return errors.New("The page is full. Need to split !")
	}
	recordSerialized := self.recordSerializer.Serialize((record))
	for i, b := range recordSerialized {
		self.content[int(self.use) + i] = b
	}
	self.use += uint16(len(recordSerialized))

	return nil
}

func (self *PageDisk) Get(key string) (string, error) {

	var record extendible.Record
	offset := 0

	for offset < len(self.content){
		record = self.recordUnserializer.Unserialize(self.content[offset:])
		offset += int(record.Payload())

		if string(record.Key()) == key {
			return string(record.Value()), nil
		}
	}

	return "", errors.New("Key not found : " + key)
}

func (self *PageDisk) Content() []byte  {
	return self.content
}

func (self *PageDisk) LocalDepth() uint64 {
	return self.localDepth
}

func (self *PageDisk) SetLocalDepth(depth uint64) {
	self.localDepth = depth
}