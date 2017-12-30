package ne

import (
	"github.com/google/gopacket"
	"ne/receiver"
	"sync"
)

type NeManager struct {
	sniffer    *receiver.Sniffer
	dispatcher *receiver.Dispatcher
	waitGroup  *sync.WaitGroup
}

// Constructor
func NewNeManager(netIface string) *NeManager {
	aWaitGroup := new(sync.WaitGroup)

	// Channel used to pass packets from the Sniffer to the Dispatcher
	packetChan := make(chan gopacket.Packet)

	aNeManager := &NeManager{
		receiver.NewSniffer(netIface, aWaitGroup, packetChan),
		receiver.NewDispatcher(aWaitGroup, packetChan),
		aWaitGroup,
	}

	return aNeManager
}

func (neManager *NeManager) Start() {
	// Start every NE process as a goroutine and add them to the wait group
	neManager.waitGroup.Add(2)

	Log.Noticef("[NE_MANAGER] Starting NeManager event loop")
	go neManager.sniffer.Run()
	go neManager.dispatcher.Run()
	Log.Noticef("[SNIFFER] Ending Sniffing loop")

	neManager.waitGroup.Wait()
	Log.Noticef("[NE_MANAGER] Program stopped. Manager goroutines ended succesfully. ")
}

func (neManager *NeManager) Stop() {
	Log.Noticef("[NE_MANAGER] Stopping NeManager event loop. Wait goroutines to finish")
	neManager.sniffer.Finish()
	neManager.dispatcher.Finish()
}

func (neManager *NeManager) AddPacketMatcher(packetMatcher *receiver.PacketMatcher) error {
	return neManager.dispatcher.AddPacketMatcher(packetMatcher)
}

func (neManager *NeManager) RemovePacketMatcher(packetMatcher *receiver.PacketMatcher) error {
	return neManager.dispatcher.RemovePacketMatcher(packetMatcher)
}
