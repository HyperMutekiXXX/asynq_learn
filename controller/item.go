package controller

import (
	"asynq_learn/model/item"
	"github.com/gin-gonic/gin"
)

type Item struct {
	Controller
}

func init() {
	c := &Item{}
	c.LoadRouter("/item/list", c.List)
	c.Load(c)
}

func (c *Item) List(ctx *gin.Context) {
	itemArr := make([]*item.Item, 0)
	for _, v := range item.ItemMap {
		itemArr = append(itemArr, v)
	}
	c.Success(ctx, "success", itemArr)
}
