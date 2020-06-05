package call

import (
	MyWebController "myagent/src/controller"
	MyConfig "myagent/src/core/config"
	MyWebCore "myagent/src/core/web"
)

func NetServers(echoServer *MyWebCore.EchoServer) {
	var fileController MyWebController.IBaseController
	fileController = &MyWebController.FileController{}

	var cmdController MyWebController.IBaseController
	cmdController = &MyWebController.CmdController{}

	var testController MyWebController.IBaseController
	testController = &MyWebController.TestController{}

	initAndSetNets(echoServer, fileController, cmdController,testController)
}

func NetInit(currentConfig *MyConfig.MyConfig) *MyWebCore.EchoServer {
	echoServer := &MyWebCore.EchoServer{}
	echoServer.Init()

	echoServer.IP = currentConfig.GetServerIP()
	echoServer.PORT = currentConfig.GetServerPort()
	echoServer.SetConfig(currentConfig)

	return echoServer
}

func initAndSetNets(echoServer *MyWebCore.EchoServer, ss ...MyWebController.IBaseController) {
	if ss == nil || len(ss) <= 0 {
		return
	}
	for i := 0; i < len(ss); i++ {
		if nil == ss[i] {
			continue
		}
		ss[i].ControllerInit(echoServer)
	}
}
