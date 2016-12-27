package main

import (
	"github.com/kataras/iris"
	"strconv"
	"strings"
	"net/http"
)

type RecordToJSON struct {
	Key string `json:"key"`
	Val string `json:"val"`
}

func IrisHandler(server *server) *iris.Framework {
	api := iris.New()

	apiV1 := api.Party(v1Path)

	apiV1.Get("/data/*randomName", func(ctx *iris.Context) {

		storePath := extractDynamicPath(dataPath, ctx.PathString())
		store := server.getStore(storePath)

		key := ctx.URLParam(idParam)
		store.Lock()
		r := &RecordToJSON{
			Key: key,
			Val: store.Get(key),
		}
		store.Unlock()

		ctx.JSON(http.StatusOK, r)
	})

	apiV1.Put("/data/*randomName", func(ctx *iris.Context) {

		storePath := extractDynamicPath(dataPath, ctx.PathString())
		store := server.getStore(storePath)
		key := ctx.URLParam(idParam)

		store.Lock()
		err := store.Put(key, string(ctx.PostBody()))
		store.Unlock()

		if err != nil {
			type error struct {
				Message string `json:"message"`
			}
			msg := &error{Message: err.Error()}
			ctx.JSON(http.StatusRequestEntityTooLarge, msg)
		}

		ctx.SetStatusCode(http.StatusOK)
	})

	apiV1.Get("/meta/*randomName", func(ctx *iris.Context) {

		storePath := extractDynamicPath(metaPath, ctx.PathString())
		store := server.getStore(storePath)
		ctx.JSON(http.StatusOK, store.Meta())
	})

	apiV1.Get("/debug/*randomName", func(ctx *iris.Context) {
		storePath := extractDynamicPath(debugPath, ctx.PathString())
		store := server.getStore(storePath)
		pageId, _ := strconv.Atoi(ctx.Param("page_id"))
		ctx.WriteString(store.DumpPage(pageId))
	})

	//api.Delete(tableName, func(ctx *iris.Context) {
	//
	//})

	return api
}

func extractDynamicPath(fixedPath string, fullPath string) string {
	return strings.TrimPrefix(fullPath, fixedPath)
}
