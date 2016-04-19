package bluckstore

import "github.com/BenJoyenConseil/bluckdb/util"

type HashMap struct {
	table []*Entry

	size int
}

const BUCKET_NUMER int = 8

func NewHashMap() *HashMap{
	return &HashMap{table:make([]* Entry, BUCKET_NUMER), size:0}
}

func (self * HashMap) bucket(key util.Hashable) int {
	return key.Hash() % len(self.table)
}

func (self * HashMap) expand() {
	table := self.table
	self.table = make([]*Entry, len(table) * 2)

	for _, E := range table {
		for e := E; e != nil ; e = e.next {
			self.Put(E.key, E.value)
		}
	}
}

func (self * HashMap) Put(key util.Hashable, value interface{}) {
	bucket := self.bucket(key)
	var appended bool

	self.table[bucket], appended = self.table[bucket].Put(key, value)
	if appended {
		self.size++
	}
	if self.size * 2 > len(self.table) {
		self.expand()
	}
}

func (self * HashMap) Get(key util.Hashable) interface{} {
	bucket := self.bucket(key)

	if has, value := self.table[bucket].Get(key); has {
		return value
	}
	return "key not found"
}