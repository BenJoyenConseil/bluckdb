package memap

import (
	"encoding/binary"
)

type Record interface {
	Key() []byte
	Val() []byte
	KeyLen() uint16
	ValLen() uint16
}

type RecordWrtier interface {
	Write(key, val string)
}

type ByteRecord []byte

func (r ByteRecord) Key() []byte {
	return r[ : r.KeyLen()]
}

func (r ByteRecord) Val() []byte {
	return r[r.KeyLen() : r.KeyLen() + r.ValLen()]
}

func (r ByteRecord) KeyLen() uint16 {
	l := len(r)
	return binary.LittleEndian.Uint16(r[l - 2:])
}

func (r ByteRecord) ValLen() uint16 {
	l := len(r)
	return binary.LittleEndian.Uint16(r[l - 4:])
}

func (r ByteRecord) Write(key, val string) {
	content := []byte(key + val)

	for i:= 0; i < len(content); i++ {
		r[i] = content[i]
	}
	binary.LittleEndian.PutUint16(r[len(content) :], uint16(len(val)))
	binary.LittleEndian.PutUint16(r[len(content) + 2 : ], uint16(len(key)))
}