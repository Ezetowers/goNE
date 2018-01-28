package processing

import (
	"sync"
)

type Scheduler struct {
	waitGroup      *sync.WaitGroup
	processor      SchedulerProcessor
	taskActionChan chan TaskAction
	workers        []*Worker
}

func NewScheduler(workersCount int) *Scheduler {
	// Create a wait group with the amount of workers to wait to finish
	aWaitGroup := new(sync.WaitGroup)
	aWaitGroup.Add(workersCount)

	// Channel used to pass taskActions to the Scheduler
	taskActionChan := make(chan TaskAction)

	aScheduler := &Scheduler{
		aWaitGroup,
		SchedulerProcessor{0},
		taskActionChan,
		make([]*Worker, workersCount),
	}

	// Create the workers
	for i := 0; i < workersCount; i++ {
		worker := NewWorker(i, aWaitGroup, taskActionChan)
		aScheduler.workers[i] = worker
	}

	return aScheduler
}

/**
 * @brief Start the Scheduler Workers. Every Worker is a goroutine
 *
 */
func (self *Scheduler) Start() {

	// Launch the workers
	for i := 0; i < len(self.workers); i++ {
		go self.workers[i].Run()
	}
}

/**
 * @brief Send the finish signal to the Workers and wait them
 * to finish gracefully
 */
func (self *Scheduler) Stop() {
	for i := 0; i < len(self.workers); i++ {
		self.workers[i].Finish()
	}

	// Close the Packet channel so the Workers can also finish
	close(self.taskActionChan)
	// Wait for all the workers to stop
	self.waitGroup.Wait()
}

/**
 * @brief      Adds a taskAction to the Workers processing channel
 *
 * @param      taskAction  Action to be enqueued
 *
 */
func (self *Scheduler) AddTaskAction(taskAction TaskAction) {
	self.taskActionChan <- taskAction
}
