package memap

import (
	"encoding/binary"
)

const (
	RECORD_HEADER_SIZE = 2
	RECORD_TOTAL_HEADER_SIZE = 4
)

type Record interface {
	Key() []byte
	Val() []byte
	KeyLen() uint16
	ValLen() uint16
}

type RecordWriter interface {
	Write(key, val string)
}

type ByteRecord []byte

func (r ByteRecord) Key() []byte {
	l := len(r)
	return r[l - RECORD_TOTAL_HEADER_SIZE - int(r.KeyLen()) - int(r.ValLen()) : l - RECORD_TOTAL_HEADER_SIZE - int(r.ValLen())]
}

func (r ByteRecord) Val() []byte {
	l := len(r)
	return r[l - RECORD_TOTAL_HEADER_SIZE - int(r.ValLen()) : l - RECORD_TOTAL_HEADER_SIZE]
}

func (r ByteRecord) KeyLen() uint16 {
	l := len(r)
	return binary.LittleEndian.Uint16(r[l - RECORD_HEADER_SIZE:])
}

func (r ByteRecord) ValLen() uint16 {
	l := len(r)
	return binary.LittleEndian.Uint16(r[l - RECORD_TOTAL_HEADER_SIZE:])
}

/*

 Record order = [key, value, lenVal, lenKey]

*/
func (r ByteRecord) Write(key, val string) {
	total := len(key) + len(val)
	copy(r[:], key)
	copy(r[len(key):], val)
	binary.LittleEndian.PutUint16(r[total :], uint16(len(val)))
	binary.LittleEndian.PutUint16(r[total + RECORD_HEADER_SIZE : ], uint16(len(key)))
}