package call

import (
	"fmt"
	"github.com/robfig/cron/v3"
	MyCommon "myagent/src/core/common"
	MyConfig "myagent/src/core/config"
	MyPerceptionCore "myagent/src/core/perception"
)

//agent 初始化，影响 perception
func AgentInit(currentConfig *MyConfig.MyConfig,cron *cron.Cron) {
	ip := currentConfig.GetServerIP()
	port := currentConfig.GetServerPort()
	name := currentConfig.GetServerName()
	if len(ip) > 0  &&  ("localhost"!=ip || "127.0.0.1" != ip)  {
		ip = MyCommon.GetLocalIp()
	}

	//初始化 Myself Node
	MyPerceptionCore.GetPerceptionAgent().SetPerceptionAgent(name,ip,port)
	//初始化 Myself RsLines
	MyPerceptionCore.GetPerceptionAgent().InitMySelfRsLines()
	//初始化


	//如果使用 agent check 则添加定时任务 TODO
	if currentConfig.GetServerCheckAgent() {
		forCron := currentConfig.GetServerCheckTimeForCron()
		fmt.Println(forCron)
		//cron.AddFunc(forCron, func() {

		//})
	}
}