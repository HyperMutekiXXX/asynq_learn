package test

import (
	"asynq_learn/constant"
	"context"
	"encoding/json"
	"github.com/hibiken/asynq"
	"log"
	"testing"
	"time"
)

const (
	TaskHandleFunc = "Hello:HandleFunc"
	TaskHandler    = "Hello:Handler"
)

type Hello struct {
	Msg string `json:"msg"`
}

// 生产者发送信息进队列
func TestHelloClient(t *testing.T) {
	client := asynq.NewClient(asynq.RedisClientOpt{
		Addr:     constant.RedisAddr,
		Password: constant.RedisPass,
		DB:       0,
	})
	defer client.Close()
	hello := &Hello{Msg: "Hello World!"}
	payload, err := json.Marshal(hello)
	if err != nil {
		log.Fatalf("unmarshal err %v", err)
	}
	handlerTask := asynq.NewTask(TaskHandler, payload, asynq.ProcessIn(time.Minute*1))
	handleFuncTask := asynq.NewTask(TaskHandleFunc, payload)

	// 进队列
	info, err := client.Enqueue(handlerTask)
	if err != nil {
		log.Printf("could not enqueue task: %v", err)
	}
	log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)

	info1, err := client.Enqueue(handleFuncTask)
	if err != nil {
		log.Printf("could not enqueue task: %v", err)
	}
	log.Printf("enqueued task: id=%s queue=%s", info1.ID, info1.Queue)

}

// HelloTask 实现asynq.Handler接口
type HelloTask struct {
}

func (h *HelloTask) ProcessTask(ctx context.Context, task *asynq.Task) error {
	hello := &Hello{}
	err := json.Unmarshal(task.Payload(), hello)
	if err != nil {
		log.Printf("unmarshal err %v", err)
		return err
	}
	log.Printf(hello.Msg)
	return err
}

// HelloWorldTask 使用HandleFunc的方式
func HelloWorldTask(ctx context.Context, task *asynq.Task) error {
	hello := &Hello{}
	err := json.Unmarshal(task.Payload(), hello)
	if err != nil {
		log.Printf("unmarshal err %v", err)
		return err
	}
	log.Printf(hello.Msg)
	return err
}

// ErrorFunc 错误处理函数
func ErrorFunc(ctx context.Context, task *asynq.Task, err error) {
	log.Printf("err %v", err)
}

// 消费者监听队列并消费信息
func TestHelloServer(t *testing.T) {
	srv := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:     constant.RedisAddr,
			Password: constant.RedisPass,
			DB:       0,
		},
		asynq.Config{
			// 并发量
			Concurrency: 10,
			// 设置队列及其权重
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
			// 设置错误处理
			ErrorHandler: asynq.ErrorHandlerFunc(ErrorFunc),
		},
	)

	mux := asynq.NewServeMux()
	// HandleFunc方式
	mux.HandleFunc(TaskHandleFunc, HelloWorldTask)
	//Handler方式
	mux.Handle(TaskHandler, &HelloTask{})

	// 启动消费者
	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}
