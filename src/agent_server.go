package main

import (
	MyCall "myagent/src/call"
	MyCommon "myagent/src/core/common"
	MyConfig "myagent/src/core/config"
	MyLog "myagent/src/core/log"
	MyOS "myagent/src/core/myos"
	MyWebCore "myagent/src/core/web"
)

func main1() {
	ss := &MyOS.SystemCommandResult{}
	ss.DoCommandWithCmdStr(MyOS.LINUX_TOOL_BASH, "ls ~ && echo 'asdf' > ~/asdf.log")
	MyLog.Info(ss)
}

func main2() {
	ss := &MyOS.SystemFileResult{}
	ss.ListFiles("/Users/zys")
	for _, fs := range ss.Result {
		//fmt.Println(fs.FileName,"  ",fs.IsDir,"  ",fs.FileSize,"  ",fs.UpdateTime,"  ",fs.FilePath)
		MyLog.Info(fs)
	}
	MyLog.Info(ss.StartTime.Sub(ss.EndTime))
}

var RUN_PATH string
var CURRENT_CONFIG *MyConfig.MyConfig
var ECHO_SERVER *MyWebCore.EchoServer


func main() {
	RUN_PATH = MyCommon.GetLocalPath()
	RUN_PATH = "/Users/zys/go/src/myagent/conf"

	CURRENT_CONFIG = MyCall.ConfigInit(RUN_PATH)
	ECHO_SERVER = MyCall.NetInit(CURRENT_CONFIG)
	MyCall.NetServers(ECHO_SERVER)
	ECHO_SERVER.Start()
}
