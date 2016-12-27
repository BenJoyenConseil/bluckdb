package api

import (
	"github.com/BenJoyenConseil/bluckdb/bluckstore"
	"github.com/kataras/iris"
	"net/http"
	"strconv"
	"strings"
)

const (
	v1Path    = "/v1"
	dataPath  = "/v1/data"
	metaPath  = "/v1/meta"
	debugPath = "/v1/debug"
	idParam   = "id"
)

type RecordToJSON struct {
	Key string `json:"key"`
	Val string `json:"val"`
}

func AppendIrisHandlers(store *bluckstore.MultiStore) *iris.Framework {
	api := iris.New()

	api.Get("/", func(ctx *iris.Context) {

		ctx.JSON(http.StatusOK, struct{
			Version string `json:"version"`
			Message string `json:"message"`
		}{
			Version: "v0.1",
			Message: "You know, for fast persistency :)",
		})
	})

	apiV1 := api.Party(v1Path)


	apiV1.Get("/data/*randomName", func(ctx *iris.Context) {

		storePath := extractDynamicPath(dataPath, ctx.PathString())
		store := store.GetStore(storePath)

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
		store := store.GetStore(storePath)
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
		store := store.GetStore(storePath)
		ctx.JSON(http.StatusOK, store.Meta())
	})

	apiV1.Get("/debug/*randomName", func(ctx *iris.Context) {
		storePath := extractDynamicPath(debugPath, ctx.PathString())
		store := store.GetStore(storePath)
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
