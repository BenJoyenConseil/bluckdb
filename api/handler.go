package api

import (
	"io/ioutil"
	"strings"

	"github.com/BenJoyenConseil/bluckdb/bluckstore"

	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
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

func AppendIrisHandlers(store *bluckstore.MultiStore) *iris.Application {
	api := iris.New()

	api.Get("/", func(ctx context.Context) {

		ctx.JSON(struct {
			Version string `json:"version"`
			Message string `json:"message"`
		}{
			Version: "v0.1",
			Message: "You know, for fast persistency :)",
		})
	})

	apiV1 := api.Party(v1Path)

	apiV1.Get("/data/{randomName:path}", func(ctx context.Context) {
		println("----------->" + ctx.Path())
		storePath := extractDynamicPath(dataPath, ctx.Path())
		store := store.GetStore(storePath)

		key := ctx.URLParam(idParam)

		store.Lock()
		v := store.Get(key)
		store.Unlock()

		r := &RecordToJSON{
			Key: key,
			Val: v,
		}

		ctx.JSON(r)
	})

	apiV1.Put("/data/{randomName:path}", func(ctx context.Context) {
		storePath := extractDynamicPath(dataPath, ctx.Path())
		store := store.GetStore(storePath)
		key := ctx.URLParam(idParam)

		store.Lock()
		body, _ := ioutil.ReadAll(ctx.Request().Body)
		err := store.Put(key, string(body))
		store.Unlock()

		if err != nil {
			type error struct {
				Message string `json:"message"`
			}
			msg := &error{Message: err.Error()}
			ctx.StatusCode(iris.StatusRequestEntityTooLarge)
			ctx.JSON(msg)
			return
		}
	})

	apiV1.Get("/meta/{randomName:path}", func(ctx context.Context) {

		storePath := extractDynamicPath(metaPath, ctx.Path())
		store := store.GetStore(storePath)
		ctx.JSON(store.Meta())
	})

	apiV1.Get("/debug/{randomName:path}", func(ctx context.Context) {
		storePath := extractDynamicPath(debugPath, ctx.Path())
		store := store.GetStore(storePath)
		pageId, _ := ctx.Params().GetInt("page_id")
		p := store.DumpPage(pageId) // this panics on raw installations
		ctx.WriteString(p)
	})

	//api.Delete(tableName, func(ctx context.Context) {
	//
	//})

	return api
}

func extractDynamicPath(fixedPath string, fullPath string) string {
	return strings.TrimPrefix(fullPath, fixedPath)
}
