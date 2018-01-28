package main

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/Ezetowers/goNE/common"
	"github.com/Ezetowers/goNE/ne"
	"github.com/Ezetowers/goNE/processing"
	"github.com/Ezetowers/goNE/receiver"
	"github.com/op/go-logging"
)

// Config Init
type tomlConfig struct {
	common.BasicConfig
	Main mainSection
}

type mainSection struct {
	NetIface     string
	WorkersCount int
}

var MyConfig tomlConfig
var Log = logging.MustGetLogger("")

func initEnvironment() {
	common.InitConfig(&MyConfig)
	common.InitLogger(&MyConfig.BasicConfig)
	Log.Noticef("[MAIN] Environment correctly set. Starting %v program.\n",
		os.Args[0])
}

/*func main() {
	initEnvironment()

	dstSubnet := net.IPNet{
		net.ParseIP("192.168.1.0"),
		net.CIDRMask(24, 32),
	}

	srcSubnet := net.IPNet{
		net.ParseIP("10.0.0.0"),
		net.CIDRMask(8, 32),
	}

	neTask := processing.NewTask()
	pm := receiver.NewPacketMatcher(&dstSubnet, &srcSubnet, 0, neTask)

	neManager := ne.NewNeManager(MyConfig.Main.NetIface,
		MyConfig.Main.WorkersCount)
	neManager.AddPacketMatcher(pm)

	dstSubnet = net.IPNet{
		net.ParseIP("192.168.1.0"),
		net.CIDRMask(22, 32),
	}

	srcSubnet = net.IPNet{
		net.ParseIP("10.0.0.0"),
		net.CIDRMask(13, 32),
	}

	pm = receiver.NewPacketMatcher(&dstSubnet, &srcSubnet, 0, neTask)
	neManager.AddPacketMatcher(pm)

	handleSigintSignal(neManager)
	neManager.Start()
}*/

func main() {
	initEnvironment()

	neManager := ne.NewNeManager(MyConfig.Main.NetIface, MyConfig.Main.WorkersCount)
	addTestTask(neManager)

	handleSigintSignal(neManager)
	neManager.Start()
}

func addTestTask(neManager *ne.NeManager) {
	neTask := processing.NewTask()

	dstSubnet := net.IPNet{
		net.ParseIP("192.168.1.0"),
		net.CIDRMask(24, 32),
	}

	srcSubnet := net.IPNet{
		net.ParseIP("10.0.0.0"),
		net.CIDRMask(8, 32),
	}
	pm := receiver.NewPacketMatcher(&dstSubnet, &srcSubnet, 0, neTask)
	neManager.AddPacketMatcher(pm)

	dstSubnet = net.IPNet{
		net.ParseIP("192.168.1.0"),
		net.CIDRMask(22, 32),
	}

	srcSubnet = net.IPNet{
		net.ParseIP("10.0.0.0"),
		net.CIDRMask(13, 32),
	}

	pm = receiver.NewPacketMatcher(&dstSubnet, &srcSubnet, 0, neTask)
	neManager.AddPacketMatcher(pm)
}

func handleSigintSignal(neManager *ne.NeManager) {
	c := make(chan os.Signal, syscall.SIGINT)
	signal.Notify(c, os.Interrupt)

	go func() {
		// Block until we receive the SIGINT signal
		<-c
		Log.Noticef("[MAIN] SIGINT received. Proceed to finish program")
		neManager.Stop()
	}()
}
