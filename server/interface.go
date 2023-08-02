package server

import (
	"asynq_learn/constant"
	"asynq_learn/model/order"
	"asynq_learn/tasks"
	"context"
	"encoding/json"
	"github.com/hibiken/asynq"
	"log"
)

var (
	server   *asynq.Server
	serveMux *asynq.ServeMux
)

func Start() {
	server = asynq.NewServer(
		asynq.RedisClientOpt{Addr: constant.RedisAddr, Password: constant.RedisPass},
		asynq.Config{
			// 最大并发量
			Concurrency: 10,
			// 设置队列，数字表示优先级
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
			// 错误处理
			ErrorHandler: asynq.ErrorHandlerFunc(handleError),
		},
	)
	serveMux = asynq.NewServeMux()
	serveMux.Handle(constant.OrderMission, tasks.NewOrderTask())
	err := server.Run(serveMux)
	if err != nil {
		log.Printf("err %v", err)
	}
}

func handleError(ctx context.Context, task *asynq.Task, err error) {
	o := &order.Order{}
	err = json.Unmarshal(task.Payload(), o)
	if err != nil {
		log.Printf("unmarshal err %v", err)
	}
	log.Printf("order_id: %v err %v", o.Id, err)
}
