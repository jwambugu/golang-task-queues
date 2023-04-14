package queue_impl

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/jwambugu/golang-task-queues/pkg/queue"
	"log"
)

type redisDistributor struct {
	client *asynq.Client
}

func (r *redisDistributor) Close() error {
	return r.client.Close()
}

func (r *redisDistributor) Enqueue(ctx context.Context, task queue.Task) error {
	var (
		taskType = task.Type().String()
		onQueue  = task.OnQueue()
	)

	if taskType == "" {
		return errors.New("queue_impl: redis - task type is required")
	}

	if onQueue == "" {
		onQueue = queue.Default
	}

	payload, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("queue_impl: redis marshal - %v", err)
	}

	var (
		t           = asynq.NewTask(taskType, payload)
		enqueueOpts = []asynq.Option{
			asynq.Queue(onQueue.String()),
			asynq.MaxRetry(5),
		}
	)

	for _, duration := range task.RunIn() {
		enqueueOpts = append(enqueueOpts, asynq.ProcessIn(duration))
	}

	log.Printf("%+v", enqueueOpts)

	info, err := r.client.EnqueueContext(ctx, t, enqueueOpts...)
	if err != nil {
		return fmt.Errorf("queue_impl: redis enqueue - %v", err)
	}

	log.Printf("queue_impl: redis - %+v", info)
	return nil
}

func NewRedisQueue(opts asynq.RedisClientOpt) queue.Queuer {
	client := asynq.NewClient(opts)
	return &redisDistributor{
		client: client,
	}
}

type redisProcessor struct {
	srv *asynq.Server
}

func (r *redisProcessor) Run(workers ...queue.Worker) error {
	mux := asynq.NewServeMux()
	for _, worker := range workers {
		mux.HandleFunc(worker.Key().String(), worker.RedisHandler)
	}
	return r.srv.Run(mux)
}

func NewRedisProcessor(redisOpts asynq.RedisClientOpt) queue.Processor {
	srv := asynq.NewServer(
		redisOpts,
		asynq.Config{
			Concurrency: 10,
			Queues:      queue.Queues,
		},
	)
	return &redisProcessor{
		srv: srv,
	}
}
