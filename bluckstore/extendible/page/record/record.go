package extendible

type ByteRecord struct {
	keyByteLen uint16
	valueByteLen uint16
	key []byte
	value []byte
}

type Record interface {
	Bytes() []byte
	Payload() uint16
}

func (self *ByteRecord) Bytes() []byte {
	bytes := []byte{
		byte(uint16(len(self.key))),
		byte(uint16(len(self.value))),
	}
	bytes = append(bytes[:], self.key[:]...)
	bytes = append(bytes[:], self.value[:]...)

	return bytes
}

func (self *ByteRecord) Payload() uint16 {
	return uint16(2 + 2 + len(self.key) + len(self.value))
}
