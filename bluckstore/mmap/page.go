package mmap

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/labstack/gommon/log"
)

//
// A page is an array of 4096 bytes (the same size as the SSD Hard Drive default configuration on Linux)
//
// from 4092 to 4094, 2 bytes to store local depth (ld)
// from 4094 to 4096, 2 bytes to store the number of bytes used
//
type Page []byte

const (
	PAGE_SIZE               = 4096
	PAGE_USE_OFFSET         = 4094
	PAGE_LOCAL_DEPTH_OFFSET = 4092
)

//
// Read the Trailer of the Page where Use is stored (4094).
// Do LittleEndian deserialization on a 2 bytes slice
//
func (p Page) Use() int {
	return int(binary.LittleEndian.Uint16(p[PAGE_USE_OFFSET:]))
}

func (p Page) rest() int {
	return PAGE_LOCAL_DEPTH_OFFSET - p.Use()
}

func (p Page) ld() int {
	return int(binary.LittleEndian.Uint16(p[PAGE_LOCAL_DEPTH_OFFSET:]))
}

func (p Page) setLd(v int) {
	binary.LittleEndian.PutUint16(p[PAGE_LOCAL_DEPTH_OFFSET:], uint16(v))
}

//
// Get iterates through the page using PageIterator with cursor setted to p.Use().
// It compares the key length first, if it matches it compares all bytes and then returns the value.
// It begins to watch the last record because the page is append only, so the last version of the value for a key is more likely to be near the end
//
func (p Page) Get(k string) (v string, err error) {
	it := &PageIterator{
		current: p.Use(),
		p:       p,
	}
	l := uint16(len(k))
	for it.HasNext() {
		r := it.Next()

		if l == r.keyLen() {
			if bytes.Compare(r.key(), []byte(k)) == 0 {
				return string(r.val()), nil
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
// There is no lookup in this function to know if the key already exists.
// It is append only, so the duplicates will be garbage collected during the next split of the page.
//
func (p Page) Put(k, v string) error {
	payload := len(k) + len(v) + RECORD_TOTAL_HEADER_SIZE

	// TODO : should p.rest() be keeped in memory to skip the task of deserialization ?
	if p.rest() >= payload {

		use := p.Use()
		r := ByteRecord(p[use : use+payload])
		r.Write(k, v)
		binary.LittleEndian.PutUint16(p[PAGE_USE_OFFSET:], uint16(use+payload))

		return nil
	}

	return errors.New("The page is full.")
}

func (p Page) Gc() Page {
	tmp := Page(make([]byte, PAGE_SIZE))
	tmp.setLd(p.ld())
	lookup := make(map[string]bool)
	it := &PageIterator{p: p, current: p.Use()}

	for it.HasNext() {

		r := it.Next()
		k := string(r.key())
		if _, ok := lookup[k]; ok {
			// this record is skipped because a younger version exists, garbage collection of older version
			continue
		} else {
			lookup[k] = true
			tmp.Put(k, string(r.val()))
		}
	}
	log.Debugf("p.rest=%d tmp.rest=%d", p.rest(), tmp.rest())
	return tmp
}
