package main

import (
	"github.com/valyala/fasthttp"
	"github.com/kataras/iris"
	"github.com/BenJoyenConseil/bluckdb/bluckstore/mmap"
	"net/http"
	"github.com/labstack/gommon/log"
)

const (
	idParam = "id"
)

type record struct {
	Key string `json:"key"`
	Val string `json:"val"`
}

func main() {
	store := &mmap.MmapKVStore{}
	store.Open()
	defer store.Close()

	log.Print("Server listening on port 2233")
	fasthttp.ListenAndServe(":2233", IrisHandler(store))
}

func IrisHandler(store *mmap.MmapKVStore) fasthttp.RequestHandler {
	api := iris.New()

	api.Get("/", func(ctx *iris.Context) {
		key := ctx.URLParam(idParam)
		r := &record{
			Key: key,
			Val: store.Get(key),
		}

		ctx.JSON(http.StatusOK, r)
	})

	api.Put("/", func(ctx *iris.Context) {
		key := ctx.URLParam(idParam)
		store.Put(key, string(ctx.PostBody()))
		ctx.SetStatusCode(http.StatusOK)
	})

	api.Get("/meta", func(ctx *iris.Context) {
		ctx.JSON(http.StatusOK, store.Dir)
	})

	api.Get("/debug", func(ctx *iris.Context) {

	})

	//api.Delete(tableName, func(ctx *iris.Context) {
	//
	//})

	api.Build()
	return api.Router
}