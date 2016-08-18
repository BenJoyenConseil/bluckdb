package memap

import (
	"encoding/binary"
	"errors"
	"strconv"
	"bytes"
)

/*
 A page is an array of 4096 bytes (the same size as the SSD Hard Drive default configuration on Linux)

 from 4092 to 4094, 2 bytes to store local depth (ld)
 from 4094 to 4096, 2 bytes to store the number of bytes used
 */
type Page []byte

func (p Page) use() int {
	return int(binary.LittleEndian.Uint16(p[4094 : 4096]))
}

func (p Page) rest() int {
	return 4092 - p.use()
}

func (p Page) ld() int {
	return int(binary.LittleEndian.Uint16(p[4092:]))
}

func (p Page) setLd(v int) {
	binary.LittleEndian.PutUint16(p[4092:], uint16(v))
}

func (p Page) get(k string) (v string, err error) {
	it := &PageIterator{
		current: p.use(),
		p: p,
	}

	for it.hasNext() {
		r := it.next()
		if bytes.Compare(r.Key(), []byte(k)) == 0 {
			return string(r.Val()), nil
		}
	}

	return "", errors.New("Key not found")
}

func (p Page) put(k, v string) error{
	payload := len(k) + len(v) + RECORD_TOTAL_HEADER_SIZE
	if p.rest() >= payload {
		r := ByteRecord(p[p.use() : p.use() + payload])
		r.Write(k, v)
		binary.LittleEndian.PutUint16(p[4094:],  uint16(p.use() +payload))
		return nil
	} else {
		return errors.New("The page is full. use = " + strconv.Itoa(p.use()))
	}
}