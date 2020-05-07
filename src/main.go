package main

import (
	"myagent/src/cmd"
	MyLog "myagent/src/log"
	"myagent/src/myos"
)

func main1() {
	ss := &myos.SystemCommandResult{}
	ss.DoCommandWithCmdStr(myos.LINUX_TOOL_BASH, "ls ~ && echo 'asdf' > ~/asdf.log")
	MyLog.Info(ss)
}

func main2() {
	ss := &myos.SystemFileResult{}
	ss.ListFiles("/Users/zys")
	for _, fs := range ss.Result {
		//fmt.Println(fs.FileName,"  ",fs.IsDir,"  ",fs.FileSize,"  ",fs.UpdateTime,"  ",fs.FilePath)
		MyLog.Info(fs)
	}
	MyLog.Info(ss.StartTime.Sub(ss.EndTime))
}

func main() {
	cmd.Execute()
}
