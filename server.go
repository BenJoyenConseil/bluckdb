package main

import (
	"net/http"
	"fmt"
	mmap "github.com/BenJoyenConseil/bluckdb/bluckstore/mmap"
	"encoding/json"
	"encoding/gob"
	"bytes"
	"io/ioutil"
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
		buf, err := ioutil.ReadFile(mmap.DB_DIRECTORY + mmap.META_FILE_NAME)
		if err != nil {
			fmt.Fprint(w, "ReadFile META FILE error : " + err.Error())
			return
		}
		dec := gob.NewDecoder(bytes.NewBuffer(buf))
		dir := &mmap.Directory{}
		err = dec.Decode(&dir)
		if err != nil {
			fmt.Fprint(w, "Decoding META file error : " + err.Error())
			return
		}
		jsonEnc := json.NewEncoder(w)
		jsonEnc.Encode(dir)
	})

	http.HandleFunc("/debug", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		r.ParseForm()
		pageId, err := strconv.Atoi(r.Form.Get("page_id"))
		if err != nil {
			fmt.Fprint(w, "Unable to parse page_id : " + err.Error())
			return
		}

		f, err := os.Open(mmap.DB_DIRECTORY + mmap.FILE_NAME)
		if err != nil {
			fmt.Fprint(w, "ReadFile DATA FILE error : " + err.Error())
			return
		}
		buff := make([]byte, 4096)
		f.ReadAt(buff, int64(pageId * 4096))
		fmt.Fprint(w, string(buff))
	})

	http.ListenAndServe(":2233", nil)
}
