package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/jwambugu/golang-task-queues/pkg/queue"
	"log"
	"time"
)

type helloTask struct {
}

type HelloTaskPayload struct {
	queue.BaseTask

	Msg string `json:"msg,omitempty"`
}

func NewHelloTaskPayload(msg string) *HelloTaskPayload {
	return &HelloTaskPayload{
		BaseTask: queue.BaseTask{
			ProcessIn: []time.Duration{2 * time.Second, 4 * time.Second},
			Queue:     queue.Critical,
			TaskType:  queue.HelloTask,
		},
		Msg: msg,
	}
}

func (h *helloTask) Key() queue.TaskType {
	return queue.HelloTask
}

func (h *helloTask) GenericHandler(_ context.Context, t queue.Task) error {
	return nil
}

func (h *helloTask) RedisHandler(_ context.Context, t *asynq.Task) error {
	var payload *HelloTaskPayload

	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("tasks: unmarshal %v - %v", h.Key(), err)
	}

	log.Printf("[*] processing task: %v", payload.Msg)
	return nil
}

func NewHelloTask() queue.Worker {
	return &helloTask{}
}
