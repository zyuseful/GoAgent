package controller

import (
	"github.com/labstack/echo/v4"
	MyOS "myagent/src/core/myos"
	MyWeb "myagent/src/core/web"
	"net/http"
)

/**
Contrller 写法
1、定义一个结构体(类),其中包含(继承)BaseController
2、实现接口 func (this *CmdController) ControllerInit(myWeb *MyWeb.EchoServer) ...
3、controller中的 对外服务方法写法如下  fun (this CmdController) DoSomething() ...
   在初始化网络部分服务时会反射这种写法的方法进行执行,服务方法不需要形参
4、将执行web 接口服务在 func (this *CmdController) ControllerInit(myWeb *MyWeb.EchoServer) ... 中注册
5、如果非网络服务方法：1、不建议在这里编写；2、可以使用 fun (this *CmdController) DoSomething() ... 形式进行编写
 */

//url 常量
const (
	CmdController__SyncCmd="/CmdController/SyncCmd"
)

type CmdController struct {
	BaseController
}

//实现接口
func (this *CmdController) ControllerInit(myWeb *MyWeb.EchoServer) {
	this.BaseController.ControllerInit(myWeb)

	//将所需服务在这里自行维护
	this.SyncCmd()
}

//同步命令
func (this CmdController) SyncCmd() {
	this.echoServer.GetEcho().GET(CmdController__SyncCmd, func(context echo.Context) error {
		paramCmd := context.QueryParam("Cmd")
		paramTool := context.QueryParam("Tool")
		paramStr, b := this.NullParam(paramCmd,paramTool)
		if !b {
			return context.JSON(http.StatusOK,this.Failed(paramStr,nil))
		}
		result := &MyOS.SystemCommandResult{}
		result.DoCommandWithCmdStr(paramTool,paramCmd)
		return context.JSON(http.StatusOK,this.Success(result))
	})
}

