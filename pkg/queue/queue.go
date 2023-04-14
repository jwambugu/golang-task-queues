package queue

import (
	"context"
	"github.com/hibiken/asynq"
)

// Queue represents the queue which a Task runs on.
type Queue string

const (
	Critical Queue = "critical"
	Default  Queue = "default"
	Low      Queue = "low"
)

func (p Queue) String() string {
	return string(p)
}

var Queues = map[string]int{
	Critical.String(): 6, // processed 60% of the time
	Default.String():  3, // processed 30% of the time
	Low.String():      1, // processed 10% of the time
}

// Worker processes a Task based on its configuration.
//
// Key is a unique identifier for the Task being executed
// Handler processes the given Task
type Worker interface {
	Key() TaskType
	GenericHandler(ctx context.Context, t Task) error
	RedisHandler(ctx context.Context, t *asynq.Task) error
}

// Queuer provides methods for managing the queue
//
// Close closes the underlying connection for the given driver
// Enqueue adds the Task to the Queue
type Queuer interface {
	Close() error
	Enqueue(ctx context.Context, task Task) error
}

// Processor starts and manages all the Worker handling a Task
//
// Run starts the task processing and blocks until an os signal to exit the program is received.
// Once it receives a signal, it gracefully shuts down all active workers and other goroutines to process the tasks.
type Processor interface {
	Run(workers ...Worker) error
}
