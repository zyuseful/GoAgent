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

	this.TestWorld()
}

func (this TestController) HelloWorld() {
	this.echoServer.GetEcho().POST("/Test/Hello", func(context echo.Context) error {
		return context.JSON(http.StatusOK, this.Success("Hello"))
	})
}

func (this TestController) TestWorld() {
	this.echoServer.GetEcho().GET("/Test/Hello1", func(context echo.Context) error {

		//post := this.Post("http://www.baidu.com", "asdf", "", -1)
		post := this.Post("http://localhost:8091/Test/Hello", "asdf", "", 0)
		return context.JSON(http.StatusOK, post)
		//post := this.Post(concatURL, MyPerceptionCore.GetPerceptionAgent(), "", 0)
	})
}
