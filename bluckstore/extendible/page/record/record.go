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

type ByteRecordUnserializer struct {}


type ByteRecordSerializer struct {}

type Record interface {
	Bytes() []byte
	Payload() uint16
	Key() []byte
	Value() []byte
}

type RecordUnserializer interface {
	Unserialize(date []byte) Record
}

type RecordSerializer interface {
	Serialize(record Record) []byte
}

const HEADERS_BYTE_SIZE = 4
const uint16_BYTE_SIZE = 2

func (self *ByteRecord) Bytes() []byte {
	bytes := make([]byte, HEADERS_BYTE_SIZE)
	keyLen := uint16(len(self.key))
	valueLen :=  uint16(len(self.value))

	binary.LittleEndian.PutUint16(bytes[0 : uint16_BYTE_SIZE], keyLen)
	binary.LittleEndian.PutUint16(bytes[uint16_BYTE_SIZE : HEADERS_BYTE_SIZE], valueLen)
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

func New(key, value string) Record {
	return &ByteRecord{
		keyByteLen: uint16(len(key)),
		valueByteLen: uint16(len(value)),
		key: []byte(key),
		value: []byte(value),
	}
}

func (reader *ByteRecordUnserializer) Unserialize(data []byte) Record {

	keyLen := binary.LittleEndian.Uint16(data)
	valueLen := binary.LittleEndian.Uint16(data[2:])

	key := data[HEADERS_BYTE_SIZE : HEADERS_BYTE_SIZE + keyLen]
	value := data[HEADERS_BYTE_SIZE + keyLen : HEADERS_BYTE_SIZE + keyLen + valueLen]

	return &ByteRecord{
		keyByteLen: keyLen,
		valueByteLen: valueLen,
		key: key,
		value: value,
	}
}

func (writer *ByteRecordSerializer) Serialize(record Record) []byte {
	bytes := make([]byte, HEADERS_BYTE_SIZE)
	keyLen := uint16(len(record.Key()))
	valueLen :=  uint16(len(record.Value()))

	binary.LittleEndian.PutUint16(bytes[0:uint16_BYTE_SIZE], keyLen)
	binary.LittleEndian.PutUint16(bytes[uint16_BYTE_SIZE:HEADERS_BYTE_SIZE], valueLen)
	bytes = append(bytes[:], record.Key()[:]...)
	bytes = append(bytes[:], record.Value()[:]...)

	return bytes
}