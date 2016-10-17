package memap

import (
	"bytes"
	"encoding/binary"
	"errors"
)

//
// A page is an array of 4096 bytes (the same size as the SSD Hard Drive default configuration on Linux)
//
// from 4092 to 4094, 2 bytes to store local depth (ld)
// from 4094 to 4096, 2 bytes to store the number of bytes used
//
type Page []byte

const PAGE_SIZE = 4096
const PAGE_USE_OFFSET = 4094
const PAGE_LOCAL_DEPTH_OFFSET = 4092

//
// Read the Trailer of the Page where Use is stored (4094).
// Do LittleEndian unserialization on a 2 bytes slice
//
func (p Page) use() int {
	return int(binary.LittleEndian.Uint16(p[PAGE_USE_OFFSET:]))
}

func (p Page) rest() int {
	return PAGE_LOCAL_DEPTH_OFFSET - p.use()
}

func (p Page) ld() int {
	return int(binary.LittleEndian.Uint16(p[PAGE_LOCAL_DEPTH_OFFSET:]))
}

func (p Page) setLd(v int) {
	binary.LittleEndian.PutUint16(p[PAGE_LOCAL_DEPTH_OFFSET:], uint16(v))
}

func (p Page) get(k string) (v string, err error) {
	it := &PageIterator{
		current: p.use(),
		p:       p,
	}
	l := uint16(len(k))
	for it.hasNext() {
		r := it.next()

		if l == r.KeyLen() {
			if bytes.Compare(r.Key(), []byte(k)) == 0 {
				return string(r.Val()), nil
			}
		}
	}

	return "", errors.New("Key not found")
}

//
// Put writes key and value on disk, within the Page boundaries.
// It checks if the page's available space is more than the record size.
// If so, it writes record after the last offset given by page.use()
// If not, it returns an error.
//
func (p Page) put(k, v string) error {
	payload := len(k) + len(v) + RECORD_TOTAL_HEADER_SIZE

	// TODO : should p.rest() be keeped in memory to skip the task of unserialization ?
	if p.rest() >= payload {

		use := p.use()
		r := ByteRecord(p[use : use + payload])
		r.Write(k, v)
		binary.LittleEndian.PutUint16(p[PAGE_USE_OFFSET:], uint16(use + payload))

		return nil
	}

	return errors.New("The page is full.")
}
