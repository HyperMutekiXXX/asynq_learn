package client

import (
	"asynq_learn/constant"
	"context"
	"github.com/hibiken/asynq"
	"log"
)

var client *asynq.Client

func init() {
	client = asynq.NewClient(asynq.RedisClientOpt{Addr: constant.RedisAddr})
}

func EnqueueContext(ctx context.Context, task *asynq.Task) error {
	info, err := client.EnqueueContext(ctx, task)
	if err != nil {
		log.Printf("could not enqueue task: %v", err)
		return err
	}
	log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)
	return err
}
