package processing

import (
	"github.com/google/gopacket"
	"github.com/op/go-logging"
)

var Log = logging.MustGetLogger("")

type ITask interface {
	HandleInput()
	// Timeout in milliseconds
	HandleTimeout(timeout uint64)
	UniqueId() int64
	EnqueuePacket(packet gopacket.Packet)
	DequeuePacket() gopacket.Packet
}

type Task struct {
	packetChan chan gopacket.Packet
}

func NewTask() *Task {
	packetChan := make(chan gopacket.Packet)
	task := &Task{
		packetChan,
	}

	return task
}

func (t *Task) HandleInput() {
	// TODO:
}

func (t *Task) HandleTimeout(timeout uint64) {
	// TODO:
}

func (t *Task) UniqueId() int64 {
	// TODO:
	return -1
}

func (t *Task) EnqueuePacket(packent gopacket.Packet) {
	// TODO:
}

func (t *Task) DequeuePacket() gopacket.Packet {
	// TODO:
	return <-t.packetChan
}
