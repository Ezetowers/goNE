package processing

import (
	"github.com/google/gopacket"
	"sync/atomic"
)

var atomicUniqueId int64

type ITask interface {
	HandleInput()
	HandleTimeout(timeoutId int64)
	UniqueId() int64
	EnqueuePacket(packet gopacket.Packet)
	DequeuePacket() gopacket.Packet
}

type Task struct {
	packetChan chan gopacket.Packet
	id         int64
	lastPacket gopacket.Packet
}

func NewTask() *Task {
	packetChan := make(chan gopacket.Packet)
	task := &Task{
		packetChan,
		atomic.AddInt64(&atomicUniqueId, 1),
		nil,
	}
	return task
}

func (self *Task) HandleInput() {
	// TODO:
}

func (self *Task) HandleTimeout(timeoutId int64) {
	// TODO:
}

func (self *Task) UniqueId() int64 {
	return self.id
}

func (self *Task) EnqueuePacket(packet gopacket.Packet) {
	self.packetChan <- packet
}

func (self *Task) DequeuePacket() gopacket.Packet {
	return <-self.packetChan
}
