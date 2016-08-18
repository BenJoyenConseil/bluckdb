package util

import (
	"hash/fnv"
)

type String string

func (self String) Equals(other Hashable) bool {
	if o, ok := other.(String); ok {
		return self == o
	} else {
		return false
	}
}

// TODO : use Murmur3 ?
var h = fnv.New32a()

func (str String) Hash() int {
	h.Reset()
	h.Write([]byte(str))
	return int(h.Sum32())
}
