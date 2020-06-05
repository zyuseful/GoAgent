package controller

import (
	"github.com/labstack/echo/v4"
	MyWeb "myagent/src/core/web"
	"net/http"
)

type TestController struct {
	BaseController
}

//实现接口
func (this *TestController) ControllerInit(myWeb *MyWeb.EchoServer) {
	this.BaseController.ControllerInit(myWeb)

	//将所需服务在这里自行维护
	this.HelloWorld()
}

func (this TestController) HelloWorld() {
	this.echoServer.GetEcho().GET("/Test/Hello", func(context echo.Context) error {
		return context.JSON(http.StatusOK, this.Success("Hello"))
	})
}