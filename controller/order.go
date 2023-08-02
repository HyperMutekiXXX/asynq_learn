package controller

import (
	"asynq_learn/client"
	"asynq_learn/constant"
	"asynq_learn/model/item"
	"asynq_learn/model/order"
	"asynq_learn/tasks"
	"fmt"
	"github.com/gin-gonic/gin"
	"sync/atomic"
)

type Order struct {
	Controller
}

func init() {
	o := &Order{}
	o.LoadRouter("/order/list", o.List)
	o.LoadRouter("/order/buy", o.Buy)
	o.LoadRouter("/order/pay", o.Pay)
	o.LoadRouter("/order/cancel", o.Cancel)
	o.Load(o)
}

func (o *Order) Buy(ctx *gin.Context) {
	tmp := &order.Order{}
	err := ctx.Bind(tmp)
	if err != nil {
		o.Error(ctx, err.Error())
		return
	}

	itemData, ok := item.ItemMap[tmp.ItemId]
	if !ok {
		o.Error(ctx, fmt.Sprintf("商品ID[%v]不存在", tmp.ItemId))
		return
	}
	if itemData.Num == 0 {
		o.Error(ctx, fmt.Sprintf("商品[%v]没货", itemData.Name))
		return
	}
	atomic.AddInt64(&itemData.Num, -1)
	id := order.GetAtId()
	tmp.Id = id
	tmp.Status = constant.PayWait
	order.OrderMap[id] = tmp
	task := tasks.NewOrderTask()
	err = client.EnqueueContext(ctx, task.GenOrderTask(tmp))
	if err != nil {
		o.Error(ctx, err.Error())
		return
	}
	o.Success(ctx, "请支付", tmp)
}

func (o *Order) List(ctx *gin.Context) {
	orderArr := make([]*order.Order, 0)
	for _, v := range order.OrderMap {
		orderArr = append(orderArr, v)
	}
	o.Success(ctx, "success", orderArr)
}

func (o *Order) Pay(ctx *gin.Context) {
	tmp := &order.Order{}
	err := ctx.Bind(tmp)
	if err != nil {
		o.Error(ctx, err.Error())
		return
	}
	orderData := order.OrderMap[tmp.Id]
	if orderData.Status != 1 {
		o.Error(ctx, fmt.Sprintf("订单ID[%v],已取消或者已支付", orderData.Id))
		return
	}
	orderData.Status = constant.PayOk
	o.Success(ctx, "支付成功", nil)
}

func (o *Order) Cancel(ctx *gin.Context) {
	tmp := &order.Order{}
	err := ctx.Bind(tmp)
	if err != nil {
		o.Error(ctx, err.Error())
		return
	}
	orderData := order.OrderMap[tmp.Id]
	if orderData.Status != 1 {
		o.Error(ctx, fmt.Sprintf("订单ID[%v],已取消或者已支付", orderData.Id))
		return
	}
	orderData.Status = constant.PayCancel
	itemData, ok := item.ItemMap[orderData.ItemId]
	if !ok {
		o.Error(ctx, fmt.Sprintf("商品ID[%v]不存在", orderData.ItemId))
		return
	}
	atomic.AddInt64(&itemData.Num, 1)

	o.Success(ctx, "取消成功", nil)
}
