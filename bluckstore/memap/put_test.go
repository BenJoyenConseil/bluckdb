package memap

import (
	"testing"
	"encoding/binary"
	"fmt"
	"os"
	"github.com/edsrzf/mmap-go"
	"strconv"
	"github.com/BenJoyenConseil/bluckdb/util"
	"errors"
	"github.com/stretchr/testify/assert"
)

type Page []byte

func (p Page) use() int {
	return int(binary.LittleEndian.Uint16(p[4094 : 4096]))
}

func (p Page) add(v...byte) {
	copy(p[p.use():], v)
	use := make([]byte, 2)
	binary.LittleEndian.PutUint16(use, uint16(p.use() + len(v)))
	copy(p[4094:], use)
}

func (p Page) put(k, v string) error{
	headerSize := 4
	lens := make([]byte, headerSize)
	lenK := len(k)
	lenV := len(v)
	binary.LittleEndian.PutUint16(lens, uint16(lenK))
	binary.LittleEndian.PutUint16(lens[2:], uint16(lenV))
	payload := lenK + lenV + headerSize
	if p.left() >= payload {
		p.add(lens...)
		p.add([]byte(k)...)
		p.add([]byte(v)...)
		return nil
	} else {
		return errors.New("The page is full. use = " + strconv.Itoa(p.use()))
	}
}

func (p Page) left() int {
	return 4092 - p.use()
}

func (p Page) ld() int {
	return int(binary.LittleEndian.Uint16(p[4092:]))
}

func (p Page) setLd(v int) {
	binary.LittleEndian.PutUint16(p[4092:], uint16(v))
}

func (p Page) get(k string) string {

	for i := 0; i <= p.use(); {

		lenKey := int(binary.LittleEndian.Uint16(p[i : i + 2]))
		lenVal := int(binary.LittleEndian.Uint16(p[i + 2 : i + 4]))

		currentKey := string(p[i + 4 : i + 4 + lenKey])
		if currentKey == k {
			return string(p[ i + 4 + lenKey : i + 4 + lenKey + lenVal])
		}else {
			i += 4 + lenKey + lenVal
		}
	}
	return ""
}

type PageIterator struct {
	p Page
	current int
}

func (it *PageIterator) next() (k, v string) {
	lenKey := int(binary.LittleEndian.Uint16(it.p[it.current : it.current + 2]))
	lenVal := int(binary.LittleEndian.Uint16(it.p[it.current + 2 : it.current + 4]))

	key := string(it.p[it.current + 4 : it.current + 4 + lenKey])
	value := string(it.p[ it.current + 4 + lenKey : it.current + 4 + lenKey + lenVal])

	it.current += lenKey + lenVal + 4

	return key, value
}

func (it *PageIterator) hasNext() bool {
	if it.current < it.p.use() {
		return true
	}

	return false
}

func TestPut(t *testing.T){
	f, err  := os.OpenFile("/tmp/data.db", os.O_RDWR | os.O_CREATE | os.O_TRUNC, 0644)
	f.Write(make([]byte, 4096))
	m, err := mmap.Map(f, mmap.RDWR, 0)
	defer f.Close()
	if err != nil {
		fmt.Println(err)
	}
	defer m.Unmap()

	var page Page = Page(m[0:4096])
	fill(page)
	fmt.Println(page.left())
	fmt.Println(page.get("key180"))
}

func fill(page Page) {
	for i := 0; i < 185; i++{
		itoa := strconv.Itoa(i)
		page.put("key" + itoa, "value yop yop")
	}
}

type Directory struct {
	table []int
	data mmap.MMap
	gd uint
	dataFile *os.File
	lastPageId int
}

func (dir *Directory) getPageId(k util.Hashable) int {
	return k.Hash() & (( 1 << dir.gd) -1)
}

func (dir *Directory) getPage(k util.Hashable) (Page, int) {
	id := k.Hash() & (( 1 << dir.gd) -1)
	offset := dir.table[id] * 4096
	return Page(dir.data[offset : offset + 4096]), id
}

func (dir *Directory) get(k string) string {
	p, _ := dir.getPage(util.String(k))
	return p.get(k)
}

func (dir *Directory) expand() {
	dir.table = append(dir.table, dir.table...)
	dir.gd ++
}

func (dir *Directory) split(page Page) (p1, p2 Page) {
	p1 = make([]byte, 4096)
	p2 = make([]byte, 4096)

	it := &PageIterator{p: page, current: 0}

	for it.hasNext() {
		k, v := it.next()
		h := util.String(k).Hash() & (( 1 << dir.gd) -1)
		if (h >> uint(page.ld())) & 1 == 1 {
			p2.put(k, v)
		} else {
			p1.put(k, v)
		}
	}
	return p1, p2

}

func (dir *Directory) replace(obsoletePageId int, ld uint) (p1, p2 int) {
	p1Id := obsoletePageId
	p2Id := dir.nextPageId()

	for i := 0; i < len(dir.table); i++ {
		if obsoletePageId != dir.table[i] {
			continue
		}
		if (i >> ld) & 1 == 1 {
			dir.table[i] = p2Id
		} else {
			dir.table[i] = p1Id
		}
	}
	return p1Id, p2Id

}

func TestReplace(t *testing.T)  {
	// Given
	dir := &Directory{
		table:[]int{0, 1, 3, 2, 0, 1, 3, 2},
		gd: 2,
		lastPageId: 4,
	}

	// When
	r1, r2 := dir.replace(2, 2)

	// Then
	assert.Equal(t, 2, r1)
	assert.Equal(t, 5, r2)
}

func (dir *Directory) nextPageId() int {
	dir.lastPageId ++
	return dir.lastPageId
}

func (dir *Directory) put(key, value string) {
	page, id := dir.getPage(util.String(key))
	err := page.put(key, value)

	if err != nil {
		if uint(page.ld()) == dir.gd {
			dir.expand()
		}
		if uint(page.ld()) < dir.gd {

			p1, p2 := dir.split(page)
			id1, id2 := dir.replace(dir.table[id], uint(page.ld()))
			p1.setLd(page.ld() + 1)
			p2.setLd(page.ld() + 1)

			dir.dataFile.WriteAt(p1, int64(id1 * 4096))
			dir.dataFile.WriteAt(p2, int64(id2 * 4096))
			dir.data.Unmap()
			dir.data, _ = mmap.Map(dir.dataFile, mmap.RDWR, 0)
			dir.put(key, value)
		}

	}
}

func BenchmarkMemapPut(b *testing.B){
	f, err := os.OpenFile("/tmp/data.db", os.O_RDWR | os.O_CREATE | os.O_TRUNC, 0644)
	defer f.Close()
	if err != nil {
		fmt.Println(err)
	}
	f.Write(make([]byte, 4096))

	// init
	dir := &Directory{
		dataFile: f,
		gd: 0,
		table: make([]int, 1),
	}
	dir.table[0] = 0
	dir.data, _ = mmap.Map(dir.dataFile, mmap.RDWR, 0)
	defer dir.data.Unmap()
	// given
	var page Page = Page(dir.data[0:4096])
	fill(page)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {

		dir.put("yolo !! " + strconv.Itoa(i), "mec, elle est où ma caisse ??")
	}
	fmt.Println(dir.gd)

}

func BenchmarkMemapGet(b *testing.B){
	f, err := os.OpenFile("/tmp/data.db", os.O_RDWR | os.O_CREATE | os.O_TRUNC, 0644)
	defer f.Close()
	if err != nil {
		fmt.Println(err)
	}
	f.Write(make([]byte, 4096))

	// init
	dir := &Directory{
		dataFile: f,
		gd: 0,
		table: make([]int, 1),
	}
	dir.table[0] = 0
	dir.data, _ = mmap.Map(dir.dataFile, mmap.RDWR, 0)
	defer dir.data.Unmap()
	// given
	var page Page = Page(dir.data[0:4096])
	fill(page)
	for i := 0; i < b.N; i++ {

		dir.put("yolo !! " + strconv.Itoa(i), "mec, elle est où ma caisse ??")
	}


	b.ResetTimer()
	for i := 0; i < b.N; i++ {

		dir.get("yolo !! " + strconv.Itoa(i))
	}
}