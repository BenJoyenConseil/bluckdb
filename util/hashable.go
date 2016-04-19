package util

type Hashable interface {
	Equals(h Hashable) bool
	Hash() int
}
