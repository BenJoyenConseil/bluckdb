package memap

import (
	"encoding/binary"
	"errors"
	"strconv"
)

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
	for i := 0; i < p.use(); {

		lenKey := int(binary.LittleEndian.Uint16(p[i : i + 2]))
		lenVal := int(binary.LittleEndian.Uint16(p[i + 2 : i + 4]))

		currentKey := string(p[i + 4 : i + 4 + lenKey])
		if currentKey == k {
			return string(p[ i + 4 + lenKey : i + 4 + lenKey + lenVal])
		}else {
			i += 4 + lenKey + lenVal
		}
	}
	return ""
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
