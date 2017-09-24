package main

import (
	"common"
	"github.com/op/go-logging"
	"ne"
	"os"
	"os/signal"
	"syscall"
)

// Config Init
type tomlConfig struct {
	common.BasicConfig
	Main mainSection
}

type mainSection struct {
	NetIface string
}

var MyConfig tomlConfig
var Log = logging.MustGetLogger("")

func initEnvironment() {
	common.InitConfig(&MyConfig)
	common.InitLogger(&MyConfig.BasicConfig)
	Log.Noticef("Environment correctly set. Starting %v program.\n",
		os.Args[0])
}

func main() {
	initEnvironment()
	neManager := ne.NewNeManager(MyConfig.Main.NetIface)
	handleSigintSignal(neManager)
	neManager.Start()
}

func handleSigintSignal(neManager *ne.NeManager) {
	c := make(chan os.Signal, syscall.SIGINT)
	signal.Notify(c, os.Interrupt)

	go func() {
		<-c
		Log.Noticef("[MAIN] SIGINT received. Proceed to finish program")
		neManager.Stop()
	}()
}
