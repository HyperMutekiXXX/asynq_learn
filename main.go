package main

import (
	"asynq_learn/controller"
	"asynq_learn/server"
)

func main() {
	// 协程启动asynq服务端
	go server.Start()
	controller.Run()
}
