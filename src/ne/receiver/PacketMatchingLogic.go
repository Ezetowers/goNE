package receiver

import (
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
		dstPort uint16) (task processing.ITask, err error)
}

type SimplePacketMatchingLogic struct {
	packetMatchers []PacketMatcher
}

/**
 * @brief      { function_description }
 *
 * @return     { description_of_the_return_value }
 */
func (s *SimplePacketMatchingLogic) LogicName() string {
	return "SimplePacketMatchingLogic"
}

/**
 * @brief      Adds a packet matcher.
 *
 * @param      pm    { parameter_description }
 *
 * @return     { description_of_the_return_value }
 */
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

/**
 * @brief      Gets the task.
 *
 * @param      srcAddress  The source address
 * @param      dstAddress  The destination address
 * @param      dstPort     The destination port
 *
 * @return     The task.
 */
func (s *SimplePacketMatchingLogic) GetTask(srcAddress *net.IP,
	dstAddress *net.IP,
	dstPort uint16) (task processing.ITask, ok bool) {

	for _, pm := range s.packetMatchers {
		if pm.MatchOnlySrc(srcAddress, dstAddress, dstPort) {
			return pm.Task, true
		}
	}

	return nil, false
}

/**
 * @brief      { function_description }
 *
 * @return     { description_of_the_return_value }
 */
func (s *SimplePacketMatchingLogic) Dump() {
	Log.Infof("[SIMPLE_PM_LOGIC] Packet Matcher Container content:\n")
	for i, pm := range s.packetMatchers {
		Log.Infof("[SIMPLE_PM_LOGIC] PM N°%v: %v\n", i, pm.String())
	}
}
