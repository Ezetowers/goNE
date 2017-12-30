package net

import (
	"encoding/binary"
	"fmt"
)

// Ethernet
type EthHeader struct {
	Hdest   []byte
	Hsource []byte
	Hproto  uint16
}

// Constructor: Don't copy the data, just create a slice pointing to
// the data received. We want to make this operation as efficient as possible
func NewEthHeader(data []byte) *EthHeader {
	aEthHeader := &EthHeader{
		data[0:6],
		data[6:12],
		binary.BigEndian.Uint16(data[12:16]),
	}

	return aEthHeader
}

func (self *EthHeader) String() string {
	hdest := NewMac(self.Hdest)
	hsource := NewMac(self.Hsource)
	return fmt.Sprintf("[ETH_LAYER] HSource: %v - HDest: %v - HProto: 0x%.4x",
		hsource.String(),
		hdest.String(),
		self.Hproto)
}

// IPv4
// type IPv4Header struct {
// }

// func newIPv4Header(data []byte, protocol uint16) {

// }

// func (self *IPv4Header) String() string {

// }
