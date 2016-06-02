package memap

import "github.com/edsrzf/mmap-go"

type FileMapped struct {
	content mmap.MMap
}

func (f *FileMapped) Put(key, value string) {
	f.content = append(f.content, []byte(key + value)...)
	f.content.Flush()
}
