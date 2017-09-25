package receiver

import (
	"github.com/google/gopacket"
	"sync"
	"sync/atomic"
)

type Dispatcher struct {
	waitGroup  *sync.WaitGroup
	packetChan chan gopacket.Packet
	finished   uint32
}

// Constructor
func NewDispatcher(waitGroup *sync.WaitGroup,
	packetChan chan gopacket.Packet) *Dispatcher {

	aDispatcher := &Dispatcher{
		waitGroup,
		packetChan,
		0,
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
