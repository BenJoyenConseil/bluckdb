package main

import (
	"github.com/kataras/iris"
	"github.com/labstack/gommon/log"
	"os"
	"os/signal"
	"sync"
	"github.com/BenJoyenConseil/bluckdb/bluckstore"
)

const (
	idParam = "id"
)

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
		stores: make(map[string]bluckstore.ThreadSafeStore),
		lock: &sync.RWMutex{},
	}

	go func() {
		sigchan := make(chan os.Signal, 10)
		signal.Notify(sigchan, os.Interrupt, os.Kill)
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

	api.Listen(":2233")
}



type server struct {
	stores map[string]bluckstore.ThreadSafeStore
	lock *sync.RWMutex
}


func (server *server) getStore(path string) bluckstore.ThreadSafeStore {
	server.lock.Lock()
	if server.stores[path] == nil {
		log.Infof("First time using the store instance in path %s", path)
		s := bluckstore.NewStore()

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
