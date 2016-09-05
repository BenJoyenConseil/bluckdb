package memap

import (
	"github.com/edsrzf/mmap-go"
	"os"
	"github.com/BenJoyenConseil/bluckdb/util"
	"bytes"
	"encoding/gob"
	"fmt"
)

type Directory struct {
	Table []int
	data mmap.MMap
	Gd uint
	dataFile *os.File
	metaFile *os.File
	LastPageId int
}


func (dir *Directory) extendibleHash(k util.Hashable) int {
	return k.Hash() & (( 1 << dir.Gd) -1)
}

func (dir *Directory) getPage(k string) (Page, int) {
	id := dir.Table[dir.extendibleHash(util.Key(k))]
	offset := id * 4096
	return Page(dir.data[offset : offset + 4096]), id
}

func (dir *Directory) get(k string) string {
	p, _ := dir.getPage(k)
	val,_ := p.get(k)
	return  val
}

func (dir *Directory) expand() {
	dir.Table = append(dir.Table, dir.Table...)
	dir.Gd ++
}

func (dir *Directory) split(page Page) (p1, p2 Page) {
	lookup := make(map[string]bool)
	p1 = make([]byte, 4096)
	p2 = make([]byte, 4096)

	it := &PageIterator{p: page, current: page.use()}

	for it.hasNext() {
		r := it.next()
		k := string(r.Key())
		if _, ok := lookup[k]; ok {
			// this record is skipped because a younger version exists
			continue
		} else {
			lookup[k] = true
		}
		h := util.Key(k).Hash() & (( 1 << dir.Gd) -1)

		if (h >> uint(page.ld())) & 1 == 1 {
			p2.put(k, string(r.Val()))
		} else {
			p1.put(k, string(r.Val()))
		}
	}
	return p1, p2
}

func (dir *Directory) nextPageId() int {
	dir.LastPageId ++
	return dir.LastPageId
}

func (dir *Directory) replace(obsoletePageId int, ld uint) (p1, p2 int) {
	p1Id := obsoletePageId
	p2Id := dir.nextPageId()

	for i := 0; i < len(dir.Table); i++ {
		if obsoletePageId != dir.Table[i] {
			continue
		}
		if (i >> ld) & 1 == 1 {
			dir.Table[i] = p2Id
		} else {
			dir.Table[i] = p1Id
		}
	}
	return p1Id, p2Id

}


func (dir *Directory) put(key, value string) {
	page, id := dir.getPage(key)
	err := page.put(key, value)

	if err != nil {
		// TODO : log trace err
		if uint(page.ld()) == dir.Gd {
			dir.expand()
		}
		if uint(page.ld()) < dir.Gd {

			p1, p2 := dir.split(page)
			id1, id2 := dir.replace(id, uint(page.ld()))
			p1.setLd(page.ld() + 1)
			p2.setLd(page.ld() + 1)

			dir.dataFile.WriteAt(p1, int64(id1 * 4096))
			dir.dataFile.WriteAt(p2, int64(id2 * 4096))
			dir.data.Unmap()
			dir.data, err = mmap.Map(dir.dataFile, mmap.RDWR, 0)
			if err != nil {
				fmt.Println(err)
			}
			dir.put(key, value)
			dir.metaFile.WriteAt(dir.serializeMeta(), 0)

		}

	}
}

func (dir *Directory) serializeMeta() []byte {
	var meta bytes.Buffer
	enc := gob.NewEncoder(&meta)

	enc.Encode(dir)
	return meta.Bytes()
}