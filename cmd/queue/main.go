package main

import (
	"context"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/jwambugu/golang-task-queues/pkg/queue_impl"
	"github.com/jwambugu/golang-task-queues/pkg/tasks"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const redisAddr = "0.0.0.0:6379"

func main() {
	redisQueue := queue_impl.NewRedisQueue(asynq.RedisClientOpt{
		Addr: redisAddr,
	})

	//goland:noinspection GoUnhandledErrorResult
	defer redisQueue.Close()

	ctx := context.Background()
	x := tasks.NewHelloTaskPayload("Hello World")
	log.Printf("%+v", x.OnQueue())
	err := redisQueue.Enqueue(ctx, tasks.NewHelloTaskPayload("Hello World"))
	if err != nil {
		log.Fatalf("queue hello: %v", err)
	}

	errChan := make(chan error)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		redisProcessor := queue_impl.NewRedisProcessor(asynq.RedisClientOpt{
			Addr: redisAddr,
		})

		errChan <- redisProcessor.Run(
			tasks.NewHelloTask(),
		)
	}()

	log.Fatalf("exit %v", <-errChan)
}
