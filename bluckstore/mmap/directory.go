package mmap

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/BenJoyenConseil/bluckdb/util"
	"github.com/edsrzf/mmap-go"
	"github.com/labstack/gommon/log"
	"os"
	"syscall"
)

type Directory struct {
	Table      []int `json:"table"`
	data       mmap.MMap
	Gd         uint `json:"globalDepth"`
	dataFile   *os.File
	LastPageId int `json:"LastPageId"`
}

func (dir *Directory) extendibleHash(k util.Hashable) int {
	return k.Hash() & ((1 << dir.Gd) - 1)
}

func (dir *Directory) getPage(k string) (Page, int, error) {
	hash := dir.extendibleHash(util.Key(k))

	if hash > len(dir.Table)-1 {
		return nil, -1, fmt.Errorf("hash (%d) out of the Table array %d", hash, len(dir.Table))
	}

	id := dir.Table[hash]
	offset := id * PAGE_SIZE

	if offset+PAGE_SIZE > len(dir.data) {
		return nil, -1, fmt.Errorf("offset (%d) out of the data array %d", offset+PAGE_SIZE, len(dir.data))
	}
	return Page(dir.data[offset : offset+PAGE_SIZE]), id, nil
}

func (dir *Directory) get(k string) string {
	p, _, err := dir.getPage(k)
	if err != nil {
		return ""
	}

	val, _ := p.Get(k)
	return val
}

func (dir *Directory) expand() {
	dir.Table = append(dir.Table, dir.Table...)
	dir.Gd++
}

func (dir *Directory) increaseSize() {

	stats, _ := dir.dataFile.Stat()
	size := stats.Size()
	dir.dataFile.WriteAt(make([]byte, size), int64(size))
	dir.mmapDataFile()
}

func (dir *Directory) split(page Page) (p1, p2 Page) {

	p1 = make([]byte, PAGE_SIZE)
	p2 = make([]byte, PAGE_SIZE)

	it := &PageIterator{p: page, current: page.Use()}

	for it.HasNext() {

		r := it.Next()
		k := string(r.key())

		h := util.Key(k).Hash() & ((1 << dir.Gd) - 1)

		if (h>>uint(page.ld()))&1 == 1 {
			p2.Put(k, string(r.val()))
		} else {
			p1.Put(k, string(r.val()))
		}
	}
	return p1, p2
}

func (dir *Directory) nextPageId() int {
	dir.LastPageId++
	return dir.LastPageId
}

func (dir *Directory) replace(obsoletePageId int, ld uint) (p1Id, p2Id int) {

	p1Id = obsoletePageId
	p2Id = dir.nextPageId()

	for i := 0; i < len(dir.Table); i++ {

		if obsoletePageId != dir.Table[i] {
			continue
		}
		if (i>>ld)&1 == 1 {
			dir.Table[i] = p2Id
		} else {
			dir.Table[i] = p1Id
		}
	}
	return p1Id, p2Id

}

func (dir *Directory) put(key, value string) {
	page, id, _ := dir.getPage(key)
	err := page.Put(key, value)

	if err != nil {
		// TODO : log trace err
		if uint(page.ld()) == dir.Gd {
			dir.expand()
		}
		if uint(page.ld()) < dir.Gd {
			cleaned := page.Gc()
			copy(dir.data[id*PAGE_SIZE:id*PAGE_SIZE+PAGE_SIZE], cleaned)
			if err2 := cleaned.Put(key, value); err2 == nil {
				log.Debugf("Yavait de la place")
				return
			}
			log.Debugf("C'est quand même blindé")
			p1, p2 := dir.split(cleaned)

			id1, id2 := dir.replace(id, uint(page.ld()))

			p1.setLd(page.ld() + 1)
			p2.setLd(page.ld() + 1)

			if id2*PAGE_SIZE >= len(dir.data) {
				log.Debugf("dir.data is to small : increase ! id2=%d, id2*PageSize=%d, lenDirData=%d", id2, id2*PAGE_SIZE, len(dir.data))
				dir.increaseSize()
			}
			copy(dir.data[id1*PAGE_SIZE:id1*PAGE_SIZE+PAGE_SIZE], p1)
			copy(dir.data[id2*PAGE_SIZE:id2*PAGE_SIZE+PAGE_SIZE], p2)

			dir.put(key, value)
		}

	}
}

func (dir *Directory) mmapDataFile() {

	var err error
	dir.data.Unmap()
	dir.data, err = mmap.Map(dir.dataFile, mmap.RDWR|syscall.MAP_POPULATE, 0644)

	log.Debugf("mmap data size = %d bytes", len(dir.data))
	if err != nil {
		log.Error(err)
	}
}

func (dir *Directory) String() string {
	w := new(bytes.Buffer)
	jsonEnc := json.NewEncoder(w)
	jsonEnc.Encode(dir)
	return w.String()
}
