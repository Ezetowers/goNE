package receiver

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
	"ne/net"
	"sync"
	"sync/atomic"
)

type IPAddrProcessing func(packet gopacket.Packet, data *IPHeaderInfo) error
type TransportHeaderProcessing func(packet gopacket.Packet, data *TransportHeaderInfo) error

type Dispatcher struct {
	waitGroup      *sync.WaitGroup
	packetChan     chan gopacket.Packet
	finished       uint32
	pmLogic        SimplePacketMatchingLogic
	processingLock *sync.Mutex
	ipMap          map[uint16]IPAddrProcessing
	transportMap   map[int]TransportHeaderProcessing
}

type IPHeaderInfo struct {
	srcAddr  net.IPAddress
	dstAddr  net.IPAddress
	protocol int
}

type TransportHeaderInfo struct {
	port uint16
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
		make(map[uint16]IPAddrProcessing),
		make(map[int]TransportHeaderProcessing),
	}

	// Fill the ipMap with the proper callbacks
	aDispatcher.ipMap[net.ETHERTYPE_IP] = aDispatcher.processIPv4Packet
	aDispatcher.ipMap[net.ETHERTYPE_IPV6] = aDispatcher.processIPv6Packet

	// Fill the transportMap with the proper
	aDispatcher.transportMap[net.IPPROTO_UDP] = aDispatcher.processUDPHeader

	return aDispatcher
}

/**
 * @brief      { function_description }
 *
 * @return     { description_of_the_return_value }
 */
func (self *Dispatcher) Run() {
	Log.Noticef("[DISPATCHER] Starting Dispatching loop")

	for self.finished == 0 {
		if packet, ok := <-self.packetChan; ok == true {
			self.processPacket(packet)
		} else {
			Log.Noticef("[DISPATCHER] PacketChan has been closed. ")
			break
		}
	}

	Log.Noticef("[DISPATCHER] Ending Dispatching loop")
	self.waitGroup.Done()
}

/**
 * @brief      { function_description }
 *
 * @param      packet  The packet
 *
 * @return     { description_of_the_return_value }
 */
func (self *Dispatcher) processPacket(packet gopacket.Packet) {

	ethHeader := net.NewEthHeader(packet.Layer(layers.LayerTypeEthernet).LayerContents())

	// if ethHeader.Hproto == net.ETHERTYPE_IP || ethHeader.Hproto == net.ETHERTYPE_IPV6 {
	if ethHeader.Hproto == net.ETHERTYPE_IP || ethHeader.Hproto == net.ETHERTYPE_IPV6 {
		Log.Debugf("[DISPATCHER] Packet received. Proceed to process it")
		Log.Debugf("[DISPATCHER] %v", ethHeader)

		// Process the Network Layer packet
		var ipHeaderInfo IPHeaderInfo
		if err := self.ipMap[ethHeader.Hproto](packet, &ipHeaderInfo); err != nil {
			Log.Errorf("[DISPATCHER] Packet IPv4 Header could not be parsed. Error: %v", err)
			return
		}

		// FIXME: For the moment we'll only process UDP packets
		if ipHeaderInfo.protocol == net.IPPROTO_UDP {
			Log.Debugf("[DISPATCHER] IP header info. SrcAddr: %v - DstAddr: %v - TransportProtocol: %v",
				ipHeaderInfo.srcAddr.String(),
				ipHeaderInfo.dstAddr.String(),
				ipHeaderInfo.protocol)

			var headerInfo TransportHeaderInfo
			if err := self.transportMap[ipHeaderInfo.protocol](packet, &headerInfo); err != nil {
				Log.Errorf("[DISPATCHER] Packet Transport Header could not be parsed. Error: %v", err)
				return
			}

		}

	}
}

/**
 * @brief      Retrieves the data to process the IPv4 packet received
 *
 * @param      packet  The packet that holds the data
 * @param      data    Struct pointer where the data will be stored
 *
 * @return     nil if the IPv4 header data could be retrieved from the packet
 */
func (self *Dispatcher) processIPv4Packet(packet gopacket.Packet,
	data *IPHeaderInfo) error {

	ipv4Header, err := ipv4.ParseHeader(packet.Layer(layers.LayerTypeIPv4).LayerContents())
	if err != nil {
		return err
	}

	data.srcAddr = *net.NewIPAddress(ipv4Header.Src[12:16])
	data.dstAddr = *net.NewIPAddress(ipv4Header.Dst[12:16])
	data.protocol = ipv4Header.Protocol
	return nil
}

/**
 * @brief      Retrieves the data to process the IPv4 packet received
 *
 * @param      packet  The packet that holds the data
 * @param      data    Struct pointer where the data will be stored
 *
 * @return     nil if the IPv4 header data could be retrieved from the packet
 */
func (self *Dispatcher) processIPv6Packet(packet gopacket.Packet,
	data *IPHeaderInfo) error {

	ipv6Header, err := ipv6.ParseHeader(packet.Layer(layers.LayerTypeIPv6).LayerContents())
	if err != nil {
		return err
	}

	data.srcAddr = *net.NewIPAddress(ipv6Header.Src)
	data.dstAddr = *net.NewIPAddress(ipv6Header.Dst)
	data.protocol = ipv6Header.NextHeader

	return nil
}

/**
 * @brief      { function_description }
 *
 * @return     { description_of_the_return_value }
 */
func (self *Dispatcher) processUDPHeader(packet gopacket.Packet,
	data *TransportHeaderInfo) error {

	udpHeader := packet.Layer(layers.LayerTypeUDP).LayerContents()
	Log.Debugf("[DISPATCHER] Processing UDP Packet: %v", udpHeader)
	return nil
}

/**
 * @brief      { function_description }
 *
 * @return     { description_of_the_return_value }
 */
func (self *Dispatcher) Finish() {
	atomic.StoreUint32(&self.finished, 1)
}

/**
 * @brief      Adds a packet matcher.
 *
 * @param      pm    { parameter_description }
 *
 * @return     { description_of_the_return_value }
 */
func (self *Dispatcher) AddPacketMatcher(pm *PacketMatcher) error {
	self.processingLock.Lock()
	defer self.processingLock.Unlock()
	response := self.pmLogic.AddPacketMatcher(pm)
	self.pmLogic.Dump()
	return response
}

/**
 * @brief      Removes a packet matcher.
 *
 * @param      pm    { parameter_description }
 *
 * @return     { description_of_the_return_value }
 */
func (self *Dispatcher) RemovePacketMatcher(pm *PacketMatcher) error {
	self.processingLock.Lock()
	defer self.processingLock.Unlock()
	response := self.pmLogic.RemovePacketMatcher(pm)
	self.pmLogic.Dump()
	return response
}
