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

func (str String) Hash() int {
	h := fnv.New32a()
	h.Write([]byte(string(str)))
	return int(h.Sum32())
}
