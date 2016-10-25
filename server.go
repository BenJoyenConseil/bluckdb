package main

import (
	"github.com/valyala/fasthttp"
	"github.com/kataras/iris"
	"github.com/BenJoyenConseil/bluckdb/bluckstore/mmap"
	"net/http"
	"github.com/labstack/gommon/log"
	"strconv"
	"os"
	"os/signal"
)

const (
	idParam = "id"
)

type record struct {
	Key string `json:"key"`
	Val string `json:"val"`
}

const (
	logHeader  = "${time_rfc3339}\t[BluckServer]\t${level}\t"
	logPrefix  = `${time_rfc3339}\t${level}\t${prefix}\t${short_file}\t${line}`
)

func main() {
	log.EnableColor()
	log.SetHeader(logHeader)
	log.SetPrefix(logPrefix)

	store := &mmap.MmapKVStore{}
	store.Open()
	defer store.Close()



	go func() {
		sigchan := make(chan os.Signal, 10)
		signal.Notify(sigchan, os.Interrupt)
		<-sigchan
		log.Error("Program killed !")

		store.Close()

		os.Exit(0)
	}()

	log.Info("Launch the server to listen on port 2233...")
	log.Info("Press ^Ĉ to exit")
	err := fasthttp.ListenAndServe(":2233", IrisHandler(store))
	if err != nil {
		log.Infof("Error occured while trying to run http server : %s", err.Error())
	}
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

	api.Get("/meta", func(ctx *iris.Context) {
		ctx.JSON(http.StatusOK, store.Dir)
	})

	api.Get("/debug", func(ctx *iris.Context) {
		pageId, _ := strconv.Atoi(ctx.Param("page_id"))
		ctx.WriteString(store.DumpPage(pageId))
	})

	//api.Delete(tableName, func(ctx *iris.Context) {
	//
	//})

	api.Build()
	return api.Router
}