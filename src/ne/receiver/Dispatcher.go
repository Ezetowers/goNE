package receiver

import (
	"github.com/google/gopacket"
	"sync"
	"sync/atomic"
)

type Dispatcher struct {
	waitGroup      *sync.WaitGroup
	packetChan     chan gopacket.Packet
	finished       uint32
	pmLogic        SimplePacketMatchingLogic
	processingLock *sync.Mutex
}

// Constructor
func NewDispatcher(waitGroup *sync.WaitGroup,
	packetChan chan gopacket.Packet) *Dispatcher {
	processingLock := new(sync.Mutex)

	aDispatcher := &Dispatcher{
		waitGroup,
		packetChan,
		0,
		SimplePacketMatchingLogic{
			[]PacketMatcher{},
		},
		processingLock,
	}

	return aDispatcher
}

func (dispatcher *Dispatcher) Run() {
	Log.Noticef("[DISPATCHER] Starting Dispatching loop")

	for dispatcher.finished == 0 {
		if packet, ok := <-dispatcher.packetChan; ok == true {
			Log.Debugf("[DISPATCHER] Packet received. Proceed to process it")
			// Iterate over all layers, printing out each layer type
			for _, layer := range packet.Layers() {
				Log.Debugf("[SNIFFER] PACKET LAYER:", layer.LayerType())
			}
		} else {
			Log.Noticef("[DISPATCHER] PacketChan has been closed. ")
			break
		}
	}

	Log.Noticef("[DISPATCHER] Ending Dispatching loop")
	dispatcher.waitGroup.Done()
}

func (dispatcher *Dispatcher) Finish() {
	atomic.StoreUint32(&dispatcher.finished, 1)
}

func (dispatcher *Dispatcher) AddPacketMatcher(pm *PacketMatcher) error {
	dispatcher.processingLock.Lock()
	defer dispatcher.processingLock.Unlock()
	response := dispatcher.pmLogic.AddPacketMatcher(pm)
	dispatcher.pmLogic.Dump()
	return response
}

func (dispatcher *Dispatcher) RemovePacketMatcher(pm *PacketMatcher) error {
	dispatcher.processingLock.Lock()
	defer dispatcher.processingLock.Unlock()
	response := dispatcher.pmLogic.RemovePacketMatcher(pm)
	dispatcher.pmLogic.Dump()
	return response
}
