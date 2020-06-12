package crontab

import (
	"github.com/robfig/cron/v3"
)

var c *cron.Cron

func Init() {
	if nil == c {
		//https://godoc.org/github.com/robfig/cron
		//添加Seconds是对标准cron规范的最常见修改
		c = cron.New(cron.WithSeconds())
	}

	//这里默认使用秒
	if nil == c {
		c = cron.New(cron.WithSeconds())
	}
}

func Start() {
	c.Start()
}

func GetCron() *cron.Cron {
	return c
}
