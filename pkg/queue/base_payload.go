package queue

import "time"

type TaskType string

// Payload queues a Task.
//
// LastError stores the last error encountered whn running the Job
// Queue indicates the Queue the Job should be executed on.
// RunIn indicates when to run the Job.
// TaskType is an identifier for the Job.
type Payload interface {
	LastError() error
	Queue() Queue
	RunIn() time.Duration
	TaskType() TaskType
}

type Job struct {
	LastErr   error         `json:"last_error,omitempty"`
	OnQueue   Queue         `json:"on_queue,omitempty"`
	ProcessIn time.Duration `json:"process_in,omitempty"`
	Task      TaskType      `json:"task,omitempty"`
}

func (b *Job) LastError() error {
	return b.LastErr
}

func (b *Job) Queue() Queue {
	return b.OnQueue
}

func (b *Job) RunIn() time.Duration {
	return b.ProcessIn
}

func (b *Job) TaskType() TaskType {
	return b.Task
}
