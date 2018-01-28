package receiver

import (
	"fmt"
	"net"

	"github.com/Ezetowers/goNE/processing"
)

type PacketMatcher struct {
	// Range of addresses valid as destination address
	dstSubnet net.IPNet
	srcSubnet net.IPNet
	dstPort   uint16
	Task      processing.ITask
}

/*func (p *PacketMatcher) Equals(rhsPacket *PacketMatcher) bool {
	if p.srcSubnet == rhsPacket.srcSubnet {
		if p.dstSubnet == rhsPacket.dstSubnet {
			return p.dstPort == rhsPacket.dstPort
		}
	}

	return false
}*/

func NewPacketMatcher(dstSubnet *net.IPNet,
	srcSubnet *net.IPNet,
	dstPort uint16,
	task processing.ITask) *PacketMatcher {
	return &PacketMatcher{
		*dstSubnet,
		*srcSubnet,
		dstPort,
		task,
	}
}

func (p *PacketMatcher) MatchOnlySrc(srcAddress *net.IP,
	dstAddress *net.IP,
	dstPort uint16) bool {

	return p.dstPort == dstPort &&
		p.srcSubnet.Contains(*srcAddress) &&
		p.dstSubnet.Contains(*dstAddress)

}

func (p *PacketMatcher) MatchSrcAndPort(srcAddress *net.IP, dstPort uint16) bool {
	return p.dstPort == dstPort && p.srcSubnet.Contains(*srcAddress)
}

func (p *PacketMatcher) Match(srcAddress *net.IP) bool {
	return p.srcSubnet.Contains(*srcAddress)
}

/* Function to order the PacketMatcher from the General Network to the Particular one */
func PacketMatcherGreaterThan(p1 *PacketMatcher, p2 *PacketMatcher) bool {

	dstMaskComparation := compareIPMasks(&p1.dstSubnet.Mask, &p2.dstSubnet.Mask)
	srcMaskComparation := compareIPMasks(&p1.srcSubnet.Mask, &p2.srcSubnet.Mask)

	if dstMaskComparation == 0 {
		if srcMaskComparation == 0 {
			return false
		}
		return srcMaskComparation > 1
	}

	return dstMaskComparation > 1
}

/**
 * @brief Compare the IPMasks.
 *
 * @param      m1    The m 1
 * @param      m2    The m 2
 *
 * @return     0  if m1 == m2
 * 			   1  if m1 > m2
 * 			   -1 if m1 < m2
 */
func compareIPMasks(m1 *net.IPMask, m2 *net.IPMask) int8 {
	ones_mask1, _ := m1.Size()
	ones_mask2, _ := m2.Size()

	if ones_mask1 == ones_mask2 {
		return 0
	} else if ones_mask1 > ones_mask2 {
		return 1
	} else {
		return -1
	}
}

// String Interface
func (p *PacketMatcher) String() string {
	ones_dst, _ := p.dstSubnet.Mask.Size()
	ones_src, _ := p.srcSubnet.Mask.Size()

	return fmt.Sprintf("[PACKET_MATCHER] Dst Subnet: %v/%v - Src Subnet: %v/%v - DstPort: %v\n",
		p.dstSubnet.IP.String(),
		ones_dst,
		p.srcSubnet.IP.String(),
		ones_src,
		p.dstPort)
}
