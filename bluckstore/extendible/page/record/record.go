package extendible

import (
	"encoding/binary"
)

type ByteRecord struct {
	keyByteLen uint16
	valueByteLen uint16
	key []byte
	value []byte
}
type ByteRecordReader struct {}

type Record interface {
	Bytes() []byte
	Payload() uint16
	Key() []byte
	Value() []byte
}

type RecordReader interface {
	Read(date []byte) Record
}



func (self *ByteRecord) Bytes() []byte {
	bytes := make([]byte, 4)
	keyLen := uint16(len(self.key))
	valueLen :=  uint16(len(self.value))

	binary.LittleEndian.PutUint16(bytes[0:2], keyLen)
	binary.LittleEndian.PutUint16(bytes[2:4], valueLen)
	bytes = append(bytes[:], self.key[:]...)
	bytes = append(bytes[:], self.value[:]...)

	return bytes
}

func (self *ByteRecord) Payload() uint16 {
	return uint16(2 + 2 + len(self.key) + len(self.value))
}

func (self *ByteRecord) Key() []byte {
	return self.key
}

func (self *ByteRecord) Value() []byte {
	return self.value
}

func (self *ByteRecordReader) Read(data []byte) *ByteRecord {
	var keyPlusValueLen uint16 = 4

	keyLen := binary.LittleEndian.Uint16(data)
	valueLen := binary.LittleEndian.Uint16(data[2:])

	key := data[keyPlusValueLen : keyPlusValueLen + keyLen]
	value := data[keyPlusValueLen + keyLen : keyPlusValueLen + keyLen + valueLen]

	return &ByteRecord{
		keyByteLen: keyLen,
		valueByteLen: valueLen,
		key: key,
		value: value,
	}
}
