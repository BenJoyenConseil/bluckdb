package memap

import "encoding/binary"

type PageIterator struct {
	p Page
	current int
}

func (it *PageIterator) hasNext() bool {
	if it.current < it.p.use() {
		return true
	}

	return false
}

func (it *PageIterator) next() (k, v string) {
	lenKey := int(binary.LittleEndian.Uint16(it.p[it.current : it.current + 2]))
	lenVal := int(binary.LittleEndian.Uint16(it.p[it.current + 2 : it.current + 4]))

	key := string(it.p[it.current + 4 : it.current + 4 + lenKey])
	value := string(it.p[ it.current + 4 + lenKey : it.current + 4 + lenKey + lenVal])

	it.current += lenKey + lenVal + 4

	return key, value
}
