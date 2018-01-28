package processing

type HandleInput func()
type HandleTimeout func(int64)

type TaskAction struct {
	handleInput   HandleInput
	handleTimeout HandleTimeout
	timeoutId     int64
}

func NewTaskAction(task ITask) *TaskAction {
	taskAction := &TaskAction{
		task.HandleInput,
		task.HandleTimeout,
		-1,
	}

	return taskAction
}

func (self *TaskAction) Execute() {
	if self.timeoutId > -1 {
		self.handleTimeout(self.TimeoutId())
	} else {
		self.handleInput()
	}
}

func (self *TaskAction) TimeoutId() int64 {
	return self.timeoutId
}

func (self *TaskAction) SetTimeoutId(timeoutId int64) {
	self.timeoutId = timeoutId
}
