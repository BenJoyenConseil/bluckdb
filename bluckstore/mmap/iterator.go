package memap

type PageIterator struct {
	p       Page
	current int
}

func (it *PageIterator) hasNext() bool {
	if RECORD_TOTAL_HEADER_SIZE < it.current && it.current <= it.p.use() {
		return true
	}

	return false
}

func (it *PageIterator) next() ByteRecord {
	r := ByteRecord(it.p[:it.current])
	it.current -= int(r.KeyLen()+r.ValLen()) + int(RECORD_TOTAL_HEADER_SIZE)
	return r
}
