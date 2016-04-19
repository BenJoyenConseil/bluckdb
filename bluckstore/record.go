package bluckstore


import "github.com/BenJoyenConseil/bluckdb/util"

type Entry struct {
	key util.Hashable

	value interface{}
	next *Entry
}

func (self *Entry) Put(key util.Hashable, value interface{}) (entry *Entry, appended bool) {
	if self == nil {
		return &Entry{key: key, value: value, next: nil}, true
	} else if self.key.Equals(key) {
		self.value = value
		return self, false
	} else {
		self.next, appended = self.next.Put(key, value)
		return self, appended
	}
}

func (self *Entry) Get(key util.Hashable) (has bool, value interface{}) {
	if self == nil {
		return false, nil
	} else if self.key.Equals(key) {
		return true, self.value
	} else {
		return self.next.Get(key)
	}
}