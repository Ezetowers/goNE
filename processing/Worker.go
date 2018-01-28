package processing

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type Worker struct {
	waitGroup      *sync.WaitGroup
	id             int
	taskActionChan chan TaskAction
	finished       uint32
	logPrefix      string
}

func NewWorker(workerId int,
	waitGroup *sync.WaitGroup,
	taskActionChan chan TaskAction) *Worker {

	logPrefix := fmt.Sprintf("[WORKER N°%v]", workerId)
	worker := &Worker{
		waitGroup,
		workerId,
		taskActionChan,
		0,
		logPrefix,
	}
	return worker
}

/**
 * @brief      { function_description }
 *
 * @return     { description_of_the_return_value }
 */
func (self *Worker) Run() {
	Log.Noticef("%v Starting Worker loop", self.logPrefix)

	for self.finished == 0 {
		if taskAction, ok := <-self.taskActionChan; ok == true {
			self.processTaskAction(taskAction)
		} else {
			Log.Noticef("%v TaskActionChan has been closed.", self.logPrefix)
			break
		}
	}

	Log.Noticef("%v Ending Worker loop", self.logPrefix)
	self.waitGroup.Done()
}

/**
 * @brief      { function_description }
 *
 * @param      action  The action
 *
 * @return     { description_of_the_return_value }
 */
func (self *Worker) processTaskAction(action TaskAction) {
	// TODO:
	return
}

func (self *Worker) Finish() {
	atomic.StoreUint32(&self.finished, 1)
}
