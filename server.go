package main

import (
	"net/http"
	"fmt"
	"github.com/BenJoyenConseil/bluckdb/bluckstore/mmap"
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

	http.ListenAndServe(":2233", nil)
}
