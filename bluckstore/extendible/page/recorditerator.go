package extendible

import (
	"github.com/BenJoyenConseil/bluckdb/bluckstore/extendible/page/record"
)


type RecordIterator struct {
	content []byte
	unserializer extendible.RecordUnserializer
	current extendible.Record
}

func NewRecordIterator(page Page) *RecordIterator {
	return &RecordIterator{
		content: page.Content(),
		unserializer: &extendible.ByteRecordUnserializer{},
		current: nil,
	}
}

func (it *RecordIterator) HasNext() bool {
	if len(it.content) > 0 {
		if it.content[0] == 0x00 {
			return false
		}
		return true
	}
	return false
}

func (it *RecordIterator) Next() extendible.Record{
	if it.content == nil || len(it.content) == 0 {
		return nil
	}

	record := it.unserializer.Unserialize(it.content)
	it.content = it.content[int(record.Payload()):]
	return record
}
