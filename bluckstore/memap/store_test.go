package memap

import (
	"testing"
	"fmt"
	"os"
	"github.com/edsrzf/mmap-go"
	"strconv"
)


func fill(page Page) {
	for i := 0; i < 185; i++{
		itoa := strconv.Itoa(i)
		page.put("key" + itoa, "value yop yop")
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