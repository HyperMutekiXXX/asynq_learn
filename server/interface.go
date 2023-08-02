package server

import (
	"asynq_learn/constant"
	"asynq_learn/tasks"
	"github.com/hibiken/asynq"
	"log"
)

var (
	server   *asynq.Server
	serveMux *asynq.ServeMux
)

//func init() {
//	go start()
//}

func Start() {
	server = asynq.NewServer(
		asynq.RedisClientOpt{Addr: constant.RedisAddr},
		asynq.Config{
			// Specify how many concurrent workers to use
			Concurrency: 10,
			// Optionally specify multiple queues with different priority.
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
			// See the godoc for other configuration options

			// 错误处理
			//ErrorHandler: &ErrorTestHandler{},
		},
	)
	serveMux = asynq.NewServeMux()
	serveMux.Handle(constant.OrderMission, tasks.NewOrderTask())
	err := server.Run(serveMux)
	if err != nil {
		log.Printf("err %v", err)
	}
}
