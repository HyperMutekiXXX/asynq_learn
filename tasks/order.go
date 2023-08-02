package tasks

import (
	"asynq_learn/constant"
	"asynq_learn/model/item"
	"asynq_learn/model/order"
	"context"
	"encoding/json"
	"github.com/hibiken/asynq"
	"log"
	"sync/atomic"
	"time"
)

type OrderTask struct {
}

func (o *OrderTask) ProcessTask(ctx context.Context, task *asynq.Task) error {
	payload := task.Payload()
	orderData := &order.Order{}
	err := json.Unmarshal(payload, orderData)
	if err != nil {
		log.Printf("unmarshal err %v", err)
		return err
	}

	tmp := order.OrderMap[orderData.Id]
	if tmp.Status != constant.PayWait {
		return nil
	}

	tmp.Status = constant.PayTimeOut

	itemData := item.ItemMap[tmp.ItemId]
	atomic.AddInt64(&itemData.Num, 1)

	log.Printf("订单[%v],商品[%v]已超时", tmp.Id, itemData.Name)

	return err
}

func NewOrderTask() *OrderTask {
	return &OrderTask{}
}

func (o *OrderTask) GenOrderTask(order *order.Order) *asynq.Task {
	marshal, err := json.Marshal(order)
	if err != nil {
		log.Printf("marshal err %v", err)
		return nil
	}
	delay := time.Minute * 1
	return asynq.NewTask(constant.OrderMission, marshal, asynq.ProcessIn(delay), asynq.MaxRetry(10), asynq.Timeout(delay))
}
