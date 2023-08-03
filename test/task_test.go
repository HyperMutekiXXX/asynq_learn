package test

import (
	"asynq_learn/constant"
	"asynq_learn/tasks"
	"context"
	"encoding/json"
	"github.com/hibiken/asynq"
	"log"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: constant.RedisAddr})
	defer client.Close()

	// ------------------------------------------------------
	// Example 1: Enqueue task to be processed immediately.
	//            Use (*Client).Enqueue method.
	// ------------------------------------------------------

	task, err := tasks.NewEmailDeliveryTask(42, "some:template:id")
	if err != nil {
		log.Fatalf("could not create task: %v", err)
	}
	info, err := client.Enqueue(task)
	if err != nil {
		log.Fatalf("could not enqueue task: %v", err)
	}
	log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)

	// ------------------------------------------------------------
	// Example 2: Schedule task to be processed in the future.
	//            Use ProcessIn or ProcessAt option.
	// ------------------------------------------------------------

	info, err = client.Enqueue(task, asynq.ProcessIn(24*time.Hour))
	if err != nil {
		log.Fatalf("could not scheduled task: %v", err)
	}
	log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)

	// ----------------------------------------------------------------------------
	// Example 3: Set other options to tune task processing behavior.
	//            Options include MaxRetry, Queue, Timeout, Deadline, Unique etc.
	// ----------------------------------------------------------------------------

	task, err = tasks.NewImageResizeTask("https://example.com/myassets/image.jpg")
	if err != nil {
		log.Fatalf("could not create task: %v", err)
	}
	info, err = client.Enqueue(task, asynq.MaxRetry(10), asynq.Timeout(3*time.Minute))
	if err != nil {
		log.Fatalf("could not enqueue task: %v", err)
	}
	log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)
}

func TestServer(t *testing.T) {
	srv := asynq.NewServer(
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
		},
	)

	// mux maps a type to a handler
	mux := asynq.NewServeMux()
	mux.HandleFunc(tasks.TypeEmailDelivery, tasks.HandleEmailDeliveryTask)
	mux.Handle(tasks.TypeImageResize, &tasks.ImageProcessor{})
	// ...register other handlers...

	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}

func TestDddClientLow(t *testing.T) {
	client := asynq.NewClient(asynq.RedisClientOpt{
		Addr: constant.RedisAddr,
	})

	payload, err := json.Marshal(map[string]interface{}{"order_id": 2})
	if err != nil {
		log.Fatalf("marshal err %v", err)
	}
	task := asynq.NewTask("update:status", payload, asynq.Queue("low"))
	delay := time.Minute * 1
	info, err := client.Enqueue(task, asynq.ProcessIn(delay))
	if err != nil {
		log.Fatalf("could not enqueue task: %v", err)
	}
	log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)
}

func TestDddClient(t *testing.T) {
	client := asynq.NewClient(asynq.RedisClientOpt{
		Addr: constant.RedisAddr,
	})

	payload, err := json.Marshal(map[string]interface{}{"order_id": 1})
	if err != nil {
		log.Fatalf("marshal err %v", err)
	}
	delay := time.Minute * 1
	task := asynq.NewTask("update:status", payload, asynq.ProcessIn(delay))

	info, err := client.Enqueue(task)
	if err != nil {
		log.Fatalf("could not enqueue task: %v", err)
	}
	log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)
}

func TestDddServer(t *testing.T) {
	srv := asynq.NewServer(
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
			ErrorHandler: &ErrorTestHandler{},
		},
	)
	mux := asynq.NewServeMux()
	mux.HandleFunc("update:status", func(ctx context.Context, task *asynq.Task) error {
		res := map[string]interface{}{}
		err := json.Unmarshal(task.Payload(), &res)
		if err != nil {
			log.Fatalf("unmarshal err %v", err)
			return err
		}
		log.Printf("order id : %v", res["order_id"])
		return nil
	})
	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}

type ErrorTestHandler struct {
}

func (h *ErrorTestHandler) HandleError(ctx context.Context, task *asynq.Task, err error) {
	//TODO implement me
	panic("implement me")
}
