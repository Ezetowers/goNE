package ne

import (
	"sync"

	"github.com/Ezetowers/goNE/processing"
	"github.com/Ezetowers/goNE/receiver"
	"github.com/google/gopacket"
)

type NeManager struct {
	sniffer    *receiver.Sniffer
	dispatcher *receiver.Dispatcher
	scheduler  *processing.Scheduler
	waitGroup  *sync.WaitGroup
}

// Constructor
func NewNeManager(netIface string,
	workersCount int) *NeManager {
	aWaitGroup := new(sync.WaitGroup)

	// Channel used to pass packets from the Sniffer to the Dispatcher
	packetChan := make(chan gopacket.Packet)

	aNeManager := &NeManager{
		receiver.NewSniffer(netIface, aWaitGroup, packetChan),
		receiver.NewDispatcher(aWaitGroup, packetChan),
		processing.NewScheduler(workersCount),
		aWaitGroup,
	}

	return aNeManager
}

func (self *NeManager) Start() {
	// Start every NE process as a goroutine and add them to the wait group
	self.waitGroup.Add(2)

	Log.Noticef("[NE_MANAGER] Starting NeManager event loop")
	go self.sniffer.Run()
	go self.dispatcher.Run()
	self.scheduler.Start()

	self.waitGroup.Wait()
	Log.Noticef("[NE_MANAGER] Program stopped. Manager goroutines ended succesfully. ")
}

func (self *NeManager) Stop() {
	Log.Noticef("[NE_MANAGER] Stopping NeManager event loop. Wait goroutines to finish")
	self.scheduler.Stop()
	self.sniffer.Finish()
	self.dispatcher.Finish()
}

func (neManager *NeManager) AddPacketMatcher(packetMatcher *receiver.PacketMatcher) error {
	return neManager.dispatcher.AddPacketMatcher(packetMatcher)
}

func (neManager *NeManager) RemovePacketMatcher(packetMatcher *receiver.PacketMatcher) error {
	return neManager.dispatcher.RemovePacketMatcher(packetMatcher)
}
