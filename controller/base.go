package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

var controllerBox = make([]IController, 0)

type Controller struct {
	RouterMap map[string]gin.HandlerFunc
}

func (c *Controller) Success(ctx *gin.Context, msg string, data any) {
	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  msg,
		"data": data,
	})
}

func (c *Controller) Error(ctx *gin.Context, msg string) {
	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusBadRequest,
		"msg":  msg,
	})
}

func (c *Controller) getPath() string {
	return "/api"
}

func (c *Controller) Load(controller IController) {
	controllerBox = append(controllerBox, controller)
}

func (c *Controller) LoadRouter(url string, handlerFunc gin.HandlerFunc) {
	if len(c.RouterMap) == 0 {
		c.RouterMap = make(map[string]gin.HandlerFunc)
	}
	c.RouterMap[c.getPath()+url] = handlerFunc
}

func (c *Controller) GetRouterMap() map[string]gin.HandlerFunc {
	return c.RouterMap
}

type IController interface {
	LoadRouter(url string, handlerFunc gin.HandlerFunc)
	Load(controller IController)
	GetRouterMap() map[string]gin.HandlerFunc
}

func Run() {
	r := gin.Default()
	for _, c := range controllerBox {
		for url, handlerFunc := range c.GetRouterMap() {
			r.POST(url, handlerFunc)
		}
	}
	err := r.Run("127.0.0.1:8088")
	if err != nil {
		log.Fatalf("run err %v", err)
	}
}
