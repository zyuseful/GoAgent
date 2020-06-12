package call

import (
	"github.com/robfig/cron/v3"
	MyCrontab "myagent/src/core/crontab"
)

//定时任务初始化 -- Agent 定时检查
func CronInit() *cron.Cron{
	MyCrontab.Init()
	return MyCrontab.GetCron()
}
