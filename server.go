package main

import (
	"github.com/bjc/bluckdb/bluckstore"
	"net/http"
	"fmt"
)

func main() {

	store := bluckstore.NewMemStore()
	http.HandleFunc("/get/", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		fmt.Fprint(w, r.Form.Get("key"))
		fmt.Fprint(w, store.Get(r.Form.Get("key")))
	})

	http.HandleFunc("/put/", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		store.Put(r.Form.Get("key"), r.Form.Get("value"))
	})

	http.ListenAndServe(":8080", nil)
}