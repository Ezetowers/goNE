package network

import (
	"github.com/op/go-logging"
)

// Ethernet Header constants
const ETHERTYPE_IP uint16 = 0x800
const ETHERTYPE_IPV6 uint16 = 0x86dd
const ETHERTYPE_ARP uint16 = 0x806

// IPv4 Header constants
const IPPROTO_TCP int = 6
const IPPROTO_UDP int = 17
const IPPROTO_ICMP int = 1
const IPPROTO_ICMPv6 int = 58

var Log *logging.Logger

func init() {
	Log = logging.MustGetLogger("")
}
