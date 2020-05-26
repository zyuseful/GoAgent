package main

import (
	"fmt"
	MyCall "myagent/src/call"
	MyCmd "myagent/src/cmd"
)

var SERVER_PATH string

func main() {
	//SERVER_PATH = MyCommon.GetLocalPath()
	SERVER_PATH = "/Users/zys/go/src/myagent/conf"

	config := MyCall.ConfigInit(SERVER_PATH)
	ip := config.GetServerIP()
	port := config.GetServerPort()

	fmt.Println(ip,"  ",port)

	MyCmd.Execute()

}