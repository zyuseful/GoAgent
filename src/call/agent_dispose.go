package call

import (
	"github.com/robfig/cron/v3"
	MyConfig "myagent/src/core/config"
	MyPerceptionCore "myagent/src/core/perception"
	MyCommon "myagent/src/core/common"
	"time"
)

//agent 初始化，影响 perception
func AgentInit(currentConfig *MyConfig.MyConfig,cron *cron.Cron) {
	ip := currentConfig.GetServerIP()
	port := currentConfig.GetServerPort()
	name := currentConfig.GetServerName()
	if len(ip) > 0  &&  "localhost"!=ip && "127.0.0.1" != ip  {
		MyPerceptionCore.GetPerceptionAgent().MySelf.IP = ip
	} else {
		MyPerceptionCore.GetPerceptionAgent().MySelf.IP = MyCommon.GetLocalIp()
	}

	//初始化 Myself Node
	MyPerceptionCore.GetPerceptionAgent().SetPerceptionAgent(name,ip,port)
	//初始化 Myself RsLines
	MyPerceptionCore.GetPerceptionAgent().InitMySelfRsLines()
	//初始化
	MyPerceptionCore.GetPerceptionAgent().LocalTime = time.Now()
	MyPerceptionCore.GetPerceptionAgent().AgentUpTime = MyPerceptionCore.GetPerceptionAgent().LocalTime

	//如果使用 agent check 则添加定时任务 TODO
	if currentConfig.GetServerCheckAgent() {
		forCron := currentConfig.GetServerCheckTimeForCron()
		cron.AddFunc(forCron, func() {
			
		})
	}
}