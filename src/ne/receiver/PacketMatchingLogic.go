package receiver

import (
	"errors"
	"ne/processing"
	"net"
	"sort"
)

type PacketMatchingLogic interface {
	LogicName() string
	AddPacketMatcher() error
	RemovePacketPacket() error
	GetTask(srcAddress net.IP,
		dstAddress net.IP,
		dstPort uint32) (task processing.ITask, err error)
}

type SimplePacketMatchingLogic struct {
	packetMatchers []PacketMatcher
}

func (s *SimplePacketMatchingLogic) LogicName() string {
	return "SimplePacketMatchingLogic"
}

func (s *SimplePacketMatchingLogic) AddPacketMatcher(pm *PacketMatcher) error {
	// TODO: Check duplicates
	s.packetMatchers = append(s.packetMatchers, *pm)
	Log.Infof("[SIMPLE_PM_LOGIC] Added Packet Matcher: %v", pm)

	sort.Slice(s.packetMatchers[:], func(i, j int) bool {
		return PacketMatcherGreaterThan(&s.packetMatchers[i], &s.packetMatchers[j])
	})

	return nil
}

func (s *SimplePacketMatchingLogic) RemovePacketMatcher(pm *PacketMatcher) error {
	// TODO:
	return nil
}

func (s *SimplePacketMatchingLogic) GetTask(srcAddress *net.IP,
	dstAddress *net.IP,
	dstPort uint32) (task processing.ITask, err error) {

	for _, pm := range s.packetMatchers {
		if pm.MatchOnlySrc(srcAddress, dstAddress, dstPort) {
			return pm.Task, nil
		}
	}

	return nil, errors.New("Packet Matcher does not belong to any stored network")
}

func (s *SimplePacketMatchingLogic) Dump() {
	Log.Infof("[SIMPLE_PM_LOGIC] Packet Matcher Container content:\n")
	for i, pm := range s.packetMatchers {
		Log.Infof("[SIMPLE_PM_LOGIC] PM N°%v: %v\n", i, pm.String())
	}
}
