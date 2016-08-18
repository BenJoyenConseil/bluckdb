package util

import (
	"hash/fnv"
)

type Key []byte

// TODO : use Murmur3 ?
var h = fnv.New32a()

func (k Key) Hash() int {
	h.Reset()
	h.Write(k)
	return int(h.Sum32())
}
