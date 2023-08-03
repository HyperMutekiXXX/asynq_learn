package test

import (
	"asynq_learn/constant"
	"github.com/hibiken/asynq"
	"log"
	"testing"
)

func TestScheduled(t *testing.T) {
	// 创建scheduler
	scheduler := asynq.NewScheduler(asynq.RedisClientOpt{
		Addr:     constant.RedisAddr,
		Password: constant.RedisPass,
		DB:       0,
	}, nil)

	task := asynq.NewTask("example_task", nil)

	// 基于cron表达式
	entryID, err := scheduler.Register("* 0/30 0/1 * * ? ", task)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("registered an entry: %q\n", entryID)

	if err := scheduler.Run(); err != nil {
		log.Fatal(err)
	}
}
