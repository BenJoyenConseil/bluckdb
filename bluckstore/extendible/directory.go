package extendible

import (
	"github.com/BenJoyenConseil/bluckdb/bluckstore/extendible/page"
	"errors"
	"github.com/BenJoyenConseil/bluckdb/util"
)

type Directory struct {
	globalDepth uint64
	pointerPageTable []*extendible.Page
}

func (self *Directory) getPage(key util.Hashable) *extendible.Page {
	return self.pointerPageTable[key.Hash() & (( 1 << self.globalDepth) -1)]
}

func (self *Directory) Put(key, value string) error {
	return errors.New("Not implemented")
}

func (self *Directory) Get(key string) (string, error) {
	return "", errors.New("Not implemented")
}