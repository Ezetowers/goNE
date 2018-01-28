package receiver

import (
	"sync"
	"sync/atomic"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

type Sniffer struct {
	netIface      string
	pcapHandle    *pcap.Handle
	waitGroup     *sync.WaitGroup
	packetChan    chan gopacket.Packet
	finished      uint32
	packetCounter uint32
}

// Constructor
func NewSniffer(netIface string,
	waitGroup *sync.WaitGroup,
	packetChan chan gopacket.Packet) *Sniffer {

	// Put a timeout of one second to avoid a deadlock in case of a ctrl+c exit
	handle, err := pcap.OpenLive(netIface, 0xFFFF, true, pcap.BlockForever)
	if err != nil {
		panic(err)
	}

	aSniffer := &Sniffer{
		netIface,
		handle,
		waitGroup,
		packetChan,
		0,
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

	// Close the Packet channel so the Dispatcher can also finish
	close(sniffer.packetChan)
	sniffer.waitGroup.Done()
	Log.Noticef("[SNIFFER] Ending Sniffing loop")
}

func (sniffer *Sniffer) handlePacket(packet gopacket.Packet) {
	Log.Debugf("[SNIFFER] Packet arrived. Counter: %v", sniffer.packetCounter)
	sniffer.packetCounter += 1
	sniffer.packetChan <- packet
}

func (sniffer *Sniffer) Finish() {
	atomic.StoreUint32(&sniffer.finished, 1)
	sniffer.pcapHandle.Close()
}
