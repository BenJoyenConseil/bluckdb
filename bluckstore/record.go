package bluckstore


import "github.com/BenJoyenConseil/bluckdb/util"

type Record struct {
	key util.Hashable

	value interface{}
	next *Record
}

func (self *Record) Put(key util.Hashable, value interface{}) (entry *Record, appended bool) {
	if self == nil {
		return &Record{key: key, value: value, next: nil}, true
	} else if self.key.Equals(key) {
		self.value = value
		return self, false
	} else {
		self.next, appended = self.next.Put(key, value)
		return self, appended
	}
}

func (self *Record) Get(key util.Hashable) (has bool, value interface{}) {
	if self == nil {
		return false, nil
	} else if self.key.Equals(key) {
		return true, self.value
	} else {
		return self.next.Get(key)
	}
}