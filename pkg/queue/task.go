package queue

import "time"

type TaskType string

func (t TaskType) String() string {
	return string(t)
}

// Task queues a job to be executed.
//
// # LastError stores the last error encountered whn running the Task
//
// OnQueue indicates the Queue the Task should be executed on.
//
// RunIn indicates when to run the Task.
//
// Type is an identifier for the Task.
type Task interface {
	LastError() error
	OnQueue() Queue
	RunIn() []time.Duration
	Type() TaskType
}

type BaseTask struct {
	LastErr   error           `json:"last_error,omitempty"`
	ProcessIn []time.Duration `json:"process_in,omitempty"`
	Queue     Queue           `json:"queue,omitempty"`
	TaskType  TaskType        `json:"task_type,omitempty"`
}

func (b *BaseTask) LastError() error {
	return b.LastErr
}

func (b *BaseTask) OnQueue() Queue {
	return b.Queue
}

func (b *BaseTask) RunIn() []time.Duration {
	return b.ProcessIn
}

func (b *BaseTask) Type() TaskType {
	return b.TaskType
}

// NewTask creates a new Task running on the Default Queue
func NewTask(t TaskType) BaseTask {
	return BaseTask{
		Queue:    Default,
		TaskType: t,
	}
}

// NewTaskOnQueue creates a Task to be run on the provided Queue
func NewTaskOnQueue(q Queue, t TaskType) BaseTask {
	return BaseTask{
		Queue:    q,
		TaskType: t,
	}
}
