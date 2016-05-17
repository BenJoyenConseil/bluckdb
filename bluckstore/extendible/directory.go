package extendible

import (
	page "github.com/BenJoyenConseil/bluckdb/bluckstore/extendible/page"
	"github.com/BenJoyenConseil/bluckdb/util"
	"os"
)

type Directory struct {
	globalDepth uint64
	pointerPageTable []page.Page
	file *os.File
}


func (self *Directory) getPage(key util.Hashable) page.Page {
	return self.pointerPageTable[self.extendibleHash(key)]
}

func (self *Directory) extendibleHash(key util.Hashable) int {
	return key.Hash() & (( 1 << self.globalDepth) -1)
}

func (self *Directory) Put(key, value string) error {
	selectedPage := self.getPage(util.String(key))
	fullError := selectedPage.Put(key, value)

	if fullError != nil {
		if self.globalDepth == selectedPage.LocalDepth() {
		 	self.pointerPageTable = append(self.pointerPageTable, self.pointerPageTable...)
			self.globalDepth += 1
		}

		newPage1 := page.New()
		newPage2 := page.New()

		for i, p := range self.pointerPageTable {
			if p == selectedPage {
				if (i >> selectedPage.LocalDepth()) & 1 == 1 {
					self.pointerPageTable[i] = newPage1
				} else {
					self.pointerPageTable[i] = newPage2
				}
			}
		}

		self.Put(key, value)
		iterator := page.NewRecordIterator(selectedPage)
		for iterator.HasNext() {
			record := iterator.Next()
			self.Put(string(record.Key()), string(record.Value()))
		}

		newPage1.SetLocalDepth(selectedPage.LocalDepth() + 1)
		newPage2.SetLocalDepth(selectedPage.LocalDepth() + 1)
	}

	return nil
}

func (self *Directory) Get(key string) (string, error) {
	return self.getPage(util.String(key)).Get(key)
}