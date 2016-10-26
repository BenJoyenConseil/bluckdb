package main

import (
	"github.com/BenJoyenConseil/bluckdb/bluckstore"
	"github.com/BenJoyenConseil/bluckdb/bluckstore/mmap"
	"github.com/kataras/iris"
	"github.com/labstack/gommon/log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
)

const (
	idParam = "id"
)

type record struct {
	Key string `json:"key"`
	Val string `json:"val"`
}

const (
	logHeader = "${time_rfc3339}\t[BluckServer]\t${level}\t"
	logPrefix = `${time_rfc3339}\t${level}\t${prefix}\t${short_file}\t${line}`
	v1Path    = "/v1"
	dataPath  = "/v1/data"
	metaPath  = "/v1/meta"
	debugPath = "/v1/debug"
)

func main() {
	log.SetLevel(log.DEBUG)

	log.EnableColor()
	log.SetHeader(logHeader)
	log.SetPrefix(logPrefix)

	server := &server{
		stores: make(map[string]bluckstore.KVStore),
	}

	go func() {
		sigchan := make(chan os.Signal, 10)
		signal.Notify(sigchan, os.Interrupt)
		<-sigchan
		log.Error("Program killed !")

		server.close()

		os.Exit(0)
	}()

	log.Info("Launch the server to listen on port 2233...")
	log.Info("Press ^Ĉ to exit")

	api := IrisHandler(server)

	api.Set(iris.OptionDisableBanner(true))
	api.Set(iris.OptionDisablePathCorrection(true))
	//api.Set(iris.OptionMaxConnsPerIP(1))
	api.Listen(":2233")
}

func IrisHandler(server *server) *iris.Framework {
	api := iris.New()

	apiV1 := api.Party(v1Path)

	apiV1.Get("/data/*randomName", func(ctx *iris.Context) {
		storePath := extractDynamicPath(dataPath, ctx.PathString())
		store := server.getStore(storePath)

		key := ctx.URLParam(idParam)
		r := &record{
			Key: key,
			Val: store.Get(key),
		}

		ctx.JSON(http.StatusOK, r)
	})

	apiV1.Put("/data/*randomName", func(ctx *iris.Context) {
		storePath := extractDynamicPath(dataPath, ctx.PathString())
		store := server.getStore(storePath)
		key := ctx.URLParam(idParam)
		err := store.Put(key, string(ctx.PostBody()))

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

type server struct {
	stores map[string]bluckstore.KVStore
}

func (server *server) getStore(path string) bluckstore.KVStore {
	if server.stores[path] == nil {
		s := &mmap.MmapKVStore{}
		s.Open(path)
		log.Debugf("Open %s", path)
		server.stores[path] = s
	}
	return server.stores[path]
}

func (server *server) close() {
	for _, store := range server.stores {
		store.Close()
	}
}
