package queue

import "time"

// Task represents a Job to be executed
type Task string

func (t Task) String() string {
	return string(t)
}

// TaskRunner queues a Task.
//
// ProcessIn indicates when to run the Task.
// Queue indicates the level the Task should be executed on.
// Task indicates the
type TaskRunner interface {
	ProcessIn() time.Duration
	Priority() Queue
	Task() Task
}

// Job represents a Task that will be executed
type Job struct {
	Payload      Task          `json:"payload,omitempty"`
	RunAfter     time.Duration `json:"run_after,omitempty"`
	WithPriority Queue         `json:"with_priority,omitempty"`
}

func (j *Job) ProcessIn() time.Duration {
	return j.RunAfter
}

func (j *Job) Priority() Queue {
	return j.WithPriority
}

func (j *Job) Task() Task {
	return j.Payload
}
