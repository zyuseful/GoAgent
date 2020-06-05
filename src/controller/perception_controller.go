package controller

import (
	"github.com/labstack/echo/v4"
	MyOS "myagent/src/core/myos"
	MyWeb "myagent/src/core/web"
	"net/http"
)

//url 常量
const (
	PerceptionController__SetLocalName  = "/PerceptionController/SetLocalName"
	PerceptionController__RegisterAgent = "/PerceptionController/RegisterAgent"
)

/** 感知服务 负责节点互加、路径计算 */
type PerceptionController struct {
	BaseController
}

//实现接口
func (this *PerceptionController) ControllerInit(myWeb *MyWeb.EchoServer) {
	this.BaseController.ControllerInit(myWeb)

	//将所需服务在这里自行维护
}

/** 设置节点名称 */
func (this *PerceptionController) SetLocalName() {
	this.echoServer.GetEcho().GET(PerceptionController__SetLocalName, func(context echo.Context) error {
		name := this.GetParamFormEchoContext(context, "name", 0)


	})
}

/** 注册节点 */
func (this *PerceptionController) RegisterAgent() {
	this.echoServer.GetEcho().GET(PerceptionController__RegisterAgent, func(context echo.Context) error {
		paramCmd := this.GetParamFormEchoContext(context, "url", 0)
		paramTool := this.GetParamFormEchoContext(context, "tool", 0)
		paramStr, b := this.NullParam(paramCmd, paramTool)
		if !b {
			return context.JSON(http.StatusOK, this.Failed(paramStr, nil))
		}
		paramCmd = this.MultilineCmdDispose(paramCmd)
		result := &MyOS.SystemCommandResult{}
		result.DoCommandWithCmdStr(paramTool, paramCmd)
		return context.JSON(http.StatusOK, this.Success(result))
	})
}
