package controller

import (
	"github.com/labstack/echo/v4"
	"io"
	MyOS "myagent/src/core/myos"
	MyWeb "myagent/src/core/web"
	"net/http"
	"os"
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
	CmdController__SyncCmd           = "/CmdController/SyncCmd"
	CmdController__UploadFileSyncCmd = "/CmdController/UploadFileSyncCmd"
	CmdController__AllCmdEnv         = "/CmdController/AllCmdEnv"
)

type CmdController struct {
	BaseController
}

//实现接口
func (this *CmdController) ControllerInit(myWeb *MyWeb.EchoServer) {
	this.BaseController.ControllerInit(myWeb)

	//将所需服务在这里自行维护
	this.SyncCmd()
	this.SyncCmdsAndFile()
	this.AllCmdEnvCmd()
}

//同步命令
func (this CmdController) SyncCmd() {
	this.echoServer.GetEcho().GET(CmdController__SyncCmd, func(context echo.Context) error {
		paramCmd := this.GetParamFormEchoContext(context, "cmd", 0)
		paramTool := this.GetParamFormEchoContext(context, "tool", 0)
		paramStr, b := this.NullParam(paramCmd, paramTool)
		if b {
			return context.JSON(http.StatusOK, this.Failed(paramStr, nil))
		}
		paramCmd = this.MultilineCmdDispose(paramCmd)
		result := &MyOS.SystemCommandResult{}
		result.DoCommandWithCmdStr(paramTool, paramCmd)
		return context.JSON(http.StatusOK, this.Success(result))
	})
}

//展示所有ENV命令
func (this CmdController) AllCmdEnvCmd() {
	this.echoServer.GetEcho().GET(CmdController__AllCmdEnv, func(context echo.Context) error {
		envArr := this.BaseController.echoServer.GetConfig().GetALLReferenceVariableStr()
		return context.JSON(http.StatusOK, this.Success(envArr))
	})
}

//上传脚本、文件后执行命令
func (this CmdController) SyncCmdsAndFile() {
	this.echoServer.GetEcho().POST(CmdController__UploadFileSyncCmd, func(context echo.Context) error {
		paramCmd := this.GetParamFormEchoContext(context, "cmd", 0)
		paramTool := this.GetParamFormEchoContext(context, "tool", 0)

		paramStr, b := this.NullParam(paramCmd, paramTool)
		if b {
			return context.JSON(http.StatusOK, this.Failed(paramStr, nil))
		}
		paramCmd = this.MultilineCmdDispose(paramCmd)

		//配置文件中的上传地址
		file, err := context.FormFile("file")

		if err != nil {
			return err
		}
		src, err := file.Open()
		if err != nil {
			return err
		}
		defer src.Close()

		// Destination
		dst, err := os.Create(this.GetConfigUploadFilePath() + file.Filename)
		if err != nil {
			return err
		}
		defer dst.Close()

		// Copy
		if _, err = io.Copy(dst, src); err != nil {
			return err
		}

		result := &MyOS.SystemCommandResult{}
		result.DoCommandWithCmdStr(paramTool, paramCmd)
		return context.JSON(http.StatusOK, this.Success(result))
	})
}
