package receiver

import (
	"sync"
	"sync/atomic"
	"time"
)

type Dispatcher struct {
	waitGroup *sync.WaitGroup
	finished  uint32
}

// Constructor
func NewDispatcher(waitGroup *sync.WaitGroup) *Dispatcher {
	aDispatcher := &Dispatcher{
		waitGroup,
		0,
	}

	return aDispatcher
}

func (dispatcher *Dispatcher) Run() {
	Log.Noticef("[DISPATCHER] Starting Dispatching loop")

	for dispatcher.finished == 0 {
		time.Sleep(time.Second)
	}

	Log.Noticef("[DISPATCHER] Ending Dispatching loop")
	dispatcher.waitGroup.Done()
}

func (dispatcher *Dispatcher) Finish() {
	atomic.StoreUint32(&dispatcher.finished, 1)
}
