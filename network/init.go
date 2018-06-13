package network

import (
	"github.com/op/go-logging"
)

// Ethernet Header constants
const (
	ETHERTYPE_IP   uint16 = 0x800
	ETHERTYPE_IPV6 uint16 = 0x86dd
	ETHERTYPE_ARP  uint16 = 0x806
	MAC_SIZE       int    = 6
)

// IP Header constants
const (
	IPPROTO_TCP    int = 6
	IPPROTO_UDP    int = 17
	IPPROTO_ICMP   int = 1
	IPPROTO_ICMPv6 int = 58
	IPV4_SIZE      int = 4
	IPV6_SIZE      int = 16
)

var Log *logging.Logger

func init() {
	Log = logging.MustGetLogger("")
}
