package call

import (
	"fmt"
	"github.com/robfig/cron/v3"
	MyCommon "myagent/src/core/common"
	MyConfig "myagent/src/core/config"
	MyDnet "myagent/src/core/donet"
)

//agent 初始化，影响 perception
func AgentInit(currentConfig *MyConfig.MyConfig, cron *cron.Cron) {
	ip := currentConfig.GetServerIP()
	port := currentConfig.GetServerPort()
	name := currentConfig.GetServerName()
	if len(ip) > 0 && ("localhost" != ip || "127.0.0.1" != ip) {
		ip = MyCommon.GetLocalIp()
	}

	//初始化 Myself Node
	//MyDnet.SetPerceptionDNET(MyDnet.CreatePerceptionDNET(MyDnet.CreateDNodeByParamsWithOutTime(name,ip,port,ip+":"+port,MyDnet.BIT_ActiveOrDeaded|MyDnet.BIT_CheckThis)))
	//MyDnet.SetPerceptionDNET(MyDnet.CreatePerceptionDNET(MyDnet.CreateDNodeByParamsWithOutTime(name,ip,port,ip+":"+port,MyDnet.BIT_ActiveOrDeaded|MyDnet.BIT_CheckThis)))

	root := MyDnet.CreateDNodeByParamsWithOutTime(name, ip, port, ip+":"+port, MyDnet.BIT_ActiveOrDeaded|MyDnet.BIT_CheckThis)
	MyDnet.ThisPerceptionDNET = MyDnet.CreatePerceptionDNET(root)

	//如果使用 agent check 则添加定时任务 TODO
	if currentConfig.GetServerCheckAgent() {
		forCron := currentConfig.GetServerCheckTimeForCron()
		fmt.Println(forCron)
		//cron.AddFunc(forCron, func() {

		//})
	}
}
