package memap

import (
"encoding/binary"
"errors"
"strconv"
)
const TOTAL_HEADERS_SIZE = 4

/*
 A page is an array of 4096 bytes (the same size as the SSD Hard Drive default configuration on Linux)
 from 4092 to 4094, 2 bytes for number of bytes rest
 from 4094 to 4096, 2 bytes for number of bytes used
 */
type Page []byte

func (p Page) use() int {
	return int(binary.LittleEndian.Uint16(p[4094 : 4096]))
}

func (p Page) add(v...byte) {
	copy(p[p.use():], v)
	use := make([]byte, 2)
	binary.LittleEndian.PutUint16(use, uint16(p.use() + len(v)))
	copy(p[4094:], use)
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

func (p Page) get(k string) string {
	offset, lenK, lenV := p.find(k)

	return string(p[offset + TOTAL_HEADERS_SIZE + lenK : offset + TOTAL_HEADERS_SIZE + lenK + lenV])
}

func (p Page) put(k, v string) error{
	headerSize := 4
	lens := make([]byte, headerSize)
	lenK := len(k)
	lenV := len(v)
	binary.LittleEndian.PutUint16(lens, uint16(lenK))
	binary.LittleEndian.PutUint16(lens[2:], uint16(lenV))
	payload := lenK + lenV + headerSize
	if p.rest() >= payload {
		p.add(lens...)
		p.add([]byte(k)...)
		p.add([]byte(v)...)
		return nil
	} else {
		return errors.New("The page is full. use = " + strconv.Itoa(p.use()))
	}
}

func (p Page) remove(k string) {
	offset, lenK, lenV := p.find(k)
	p.shift(offset, TOTAL_HEADERS_SIZE + lenK + lenV)
}

func (p Page) shift(offset, size int) {
	for i := offset + size; i < p.use(); i++ {
		p[i - size] = p[i]
	}
	binary.LittleEndian.PutUint16(p[4094:], uint16(p.use() - size))
}

func (p Page) find(k string) (offset, lenK, lenV int) {
	l := len(k)

	for i := 0; i < p.use(); {

		lenK = int(binary.LittleEndian.Uint16(p[i : i + 2]))

		if l == lenK {
			currentKey := string(p[i + TOTAL_HEADERS_SIZE : i + TOTAL_HEADERS_SIZE + lenK])
			lenV = int(binary.LittleEndian.Uint16(p[i + 2 : i + TOTAL_HEADERS_SIZE]))
			if currentKey == k {
				return i, lenK, lenV
			} else {
				i += TOTAL_HEADERS_SIZE + lenK + lenV
			}
		} else {
			i += TOTAL_HEADERS_SIZE + lenK + lenV
		}
	}
	return 0, 0, 0
}