package controller

import (
	"asynq_learn/model/user"
	"github.com/gin-gonic/gin"
)

type User struct {
	Controller
}

func init() {
	c := &User{}
	c.LoadRouter("/user/list", c.List)
	c.Load(c)
}

func (c *User) List(ctx *gin.Context) {
	userArr := make([]*user.User, 0)
	for _, v := range user.UserMap {
		userArr = append(userArr, v)
	}
	c.Success(ctx, "success", userArr)
}
