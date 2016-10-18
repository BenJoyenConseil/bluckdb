package mmap

//
// Implementation of pattern Iterator to scan a Page
// The scan is made in reverse, from the last byte to the first (right to left).
//
type PageIterator struct {
	p       Page
	current int
}

//
// Return false when there is no other record to read.
// Return true if the cursor has found the beginning of the Page
//
func (it *PageIterator) hasNext() bool {
	if RECORD_TOTAL_HEADER_SIZE < it.current && it.current <= it.p.use() {
		return true
	}

	return false
}

//
// Return the next pointer to the byte array casted in a ByteRecord struct (provides methods).
// Use an cursor to know the current Record index in the Page.
// Update the cursor (current) to minus the payload of the Record
//
func (it *PageIterator) next() ByteRecord {
	r := ByteRecord(it.p[:it.current])
	it.current -= int(r.KeyLen()+r.ValLen()) + int(RECORD_TOTAL_HEADER_SIZE)
	return r
}
