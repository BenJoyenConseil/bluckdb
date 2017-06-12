package main

import (
	"flag"
	"os"
	"os/signal"

	"github.com/kataras/iris"

	"github.com/BenJoyenConseil/bluckdb/api"
	"github.com/BenJoyenConseil/bluckdb/bluckstore"

	"github.com/labstack/gommon/log"
)

const (
	logHeader = "${time_rfc3339}\t[BluckServer]\t${level}\t"
	logPrefix = `${time_rfc3339}\t${level}\t${prefix}\t${short_file}\t${line}`
)

func main() {
	var port string
	flag.StringVar(&port, "p", "2233", "Specify the host port to bind the database")
	flag.Parse()

	log.EnableColor()
	log.SetHeader(logHeader)
	log.SetPrefix(logPrefix)

	store := bluckstore.NewMmapMultiStore()
	go func() {
		sigchan := make(chan os.Signal, 10)
		signal.Notify(sigchan, os.Interrupt, os.Kill)
		<-sigchan
		log.Error("Program killed !")

		store.Close()

		os.Exit(0)
	}()

	log.Info("Launch the server to listen on port " + port)
	log.Info("Press ^Ĉ to exit")

	api := api.AppendIrisHandlers(store)

	api.Run(iris.Addr(":"+port), iris.WithoutBanner, iris.WithoutInterruptHandler)
}
