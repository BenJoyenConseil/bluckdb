package main

import (
	"fmt"
	mmap "github.com/BenJoyenConseil/bluckdb/bluckstore/mmap"
	"net/http"
	"os"
	"strconv"
)

func main() {

	store := &mmap.MmapKVStore{}
	store.Open()
	defer store.Close()
	http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		r.ParseForm()
		fmt.Fprint(w, r.Form.Get("key"))
		fmt.Fprint(w, " : ")
		fmt.Fprint(w, store.Get(r.Form.Get("key")))
	})

	http.HandleFunc("/put", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		r.ParseForm()
		store.Put(r.Form.Get("key"), r.Form.Get("value"))
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "BluckDB KV store. Yolo !")
	})

	http.HandleFunc("/meta", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, store.Dir)
	})

	http.HandleFunc("/debug", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		r.ParseForm()
		pageId, err := strconv.Atoi(r.Form.Get("page_id"))
		if err != nil {
			fmt.Fprint(w, "Unable to parse page_id : "+err.Error())
			return
		}

		f, err := os.Open(mmap.DB_DIRECTORY + mmap.FILE_NAME)
		if err != nil {
			fmt.Fprint(w, "ReadFile DATA FILE error : "+err.Error())
			return
		}
		buff := make([]byte, mmap.PAGE_SIZE)
		f.ReadAt(buff, int64(pageId * mmap.PAGE_SIZE))
		fmt.Fprint(w, string(buff))
	})

	http.ListenAndServe(":2233", nil)
}
