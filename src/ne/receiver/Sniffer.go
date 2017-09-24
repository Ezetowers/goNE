package receiver

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/op/go-logging"
	"sync"
	"sync/atomic"
)

var Log = logging.MustGetLogger("")

type Sniffer struct {
	netIface   string
	pcapHandle *pcap.Handle
	waitGroup  *sync.WaitGroup
	finished   uint32
}

// Constructor
func NewSniffer(netIface string, waitGroup *sync.WaitGroup) *Sniffer {
	// Put a timeout of one second to avoid a deadlock in case of a ctrl+c exit
	handle, err := pcap.OpenLive(netIface, 0xFFFF, true, 1)
	if err != nil {
		panic(err)
	}

	aSniffer := &Sniffer{
		netIface,
		handle,
		waitGroup,
		0,
	}

	return aSniffer
}

func (sniffer *Sniffer) Run() {
	handle := sniffer.pcapHandle
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	Log.Noticef("[SNIFFER] Starting Sniffing loop")
	for packet := range packetSource.Packets() {
		if sniffer.finished != 0 {
			break
		}
		sniffer.handlePacket(packet)
	}

	sniffer.waitGroup.Done()
	Log.Noticef("[SNIFFER] Ending Sniffing loop")
}

func (sniffer *Sniffer) handlePacket(packet gopacket.Packet) {
	Log.Debugf("[SNIFFER] Packet arrived")
	// Iterate over all layers, printing out each layer type
	for _, layer := range packet.Layers() {
		Log.Debugf("[SNIFFER] PACKET LAYER:", layer.LayerType())
	}
}

func (sniffer *Sniffer) Finish() {
	atomic.StoreUint32(&sniffer.finished, 1)
}
