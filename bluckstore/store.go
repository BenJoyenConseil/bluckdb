package bluckstore

import "github.com/BenJoyenConseil/bluckdb/bluckstore/mmap"

type KVStore interface {
	Get(k string) string
	Put(k, v string) error
	Open(absolutePath string)
	Close()
	DumpPage(pageId int) string
	Meta() *mmap.Directory
}
