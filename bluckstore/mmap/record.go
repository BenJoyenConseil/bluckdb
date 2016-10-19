package mmap

import (
	"encoding/binary"
)

const (
	RECORD_HEADER_SIZE       = 2
	RECORD_TOTAL_HEADER_SIZE = 4
)

//
// This interface is for documentation purpose only.
// A Record is composed of two byte arrays for key and value data, and footers that contain meta info for variable length data.
// In theory, A key or a value could not excead 16384 bytes length (uint16 is used to store length).
// But at the moment, there is no logic to handle a record length higher than a Page (4096 bytes)
//
type Record interface {
	key() []byte
	val() []byte
	keyLen() uint16
	valLen() uint16
}

type RecordWriter interface {
	Write(key, val string)
}

//
// A Record that is mapped on a byte slice
//
type ByteRecord []byte

func (r ByteRecord) key() []byte {
	l := len(r)
	return r[l-RECORD_TOTAL_HEADER_SIZE-int(r.keyLen())-int(r.valLen()) : l-RECORD_TOTAL_HEADER_SIZE-int(r.valLen())]
}

func (r ByteRecord) val() []byte {
	l := len(r)
	return r[l-RECORD_TOTAL_HEADER_SIZE-int(r.valLen()) : l-RECORD_TOTAL_HEADER_SIZE]
}

func (r ByteRecord) keyLen() uint16 {
	l := len(r)
	return binary.LittleEndian.Uint16(r[l-RECORD_HEADER_SIZE:])
}

func (r ByteRecord) valLen() uint16 {
	l := len(r)
	return binary.LittleEndian.Uint16(r[l-RECORD_TOTAL_HEADER_SIZE:])
}

//
// Record order = [key, value, lenVal, lenKey]
// Write does it fast because it skips "String to Byte Slice conversion",
// using copy(dst, src[]) special case with string
//
func (r ByteRecord) Write(key, val string) {
	lenKey := uint16(len(key))
	lenVal := uint16(len(val))
	total := lenKey + lenVal
	copy(r[:], key)
	copy(r[lenKey:], val)
	binary.LittleEndian.PutUint16(r[total:], lenVal)
	binary.LittleEndian.PutUint16(r[total+RECORD_HEADER_SIZE:], lenKey)
}
