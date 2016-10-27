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
	"sync"
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

	log.EnableColor()
	log.SetHeader(logHeader)
	log.SetPrefix(logPrefix)

	server := &server{
		stores: make(map[string]LockableKVStore),
		lock: &sync.RWMutex{},
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
		store.Lock()
		r := &record{
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

type server struct {
	stores map[string]LockableKVStore
	lock *sync.RWMutex
}

type LockableStore struct {
	lock *sync.RWMutex
	store bluckstore.KVStore
}
func (s *LockableStore) Get(k string) string {return s.store.Get(k)}
func (s *LockableStore) Put(k, v string) error {return s.store.Put(k, v)}
func (s *LockableStore) Open(abs string) { s.store.Open(abs)}
func (s *LockableStore) Close() { s.store.Close()}
func (s *LockableStore) Meta() *mmap.Directory{return s.store.Meta()}
func (s *LockableStore) DumpPage(id int) string {return s.store.DumpPage(id)}
func (s *LockableStore) Lock() { s.lock.Lock()}
func (s *LockableStore) Unlock() { s.lock.Unlock()}

func (server *server) getStore(path string) LockableKVStore {
	server.lock.Lock()
	if server.stores[path] == nil {
		log.Infof("Store %s not existing Creating a instance", path)
		s := &LockableStore{
			store: &mmap.MmapKVStore{},
			lock: &sync.RWMutex{},
		}

		s.Open(path)
		log.Debugf("Open %s", path)
		server.stores[path] = s
	}
	server.lock.Unlock()
	return server.stores[path]
}

func (server *server) close() {
	for _, store := range server.stores {
		store.Close()
	}
}


type LockableKVStore interface {
	bluckstore.KVStore
	Lock()
	Unlock()
}
