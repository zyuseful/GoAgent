package config

import (
	"github.com/spf13/viper"
	"sync"
)

const (
	//server 配置
	server_port = "8999"
	server_ip   = "localhost"
)

type MyConfig struct {
	//锁
	lock sync.Mutex
	//server 配置
	server_port string
	server_ip   string
}

func (this *MyConfig) BuildMyConfig(viper *viper.Viper) {
	this.lock.Lock()
	defer this.lock.Unlock()
	
	this.server_port = ifUillStr(viper.GetString("server.port"), server_port)
	this.server_ip = ifUillStr(viper.GetString("server.ip"), server_ip)
}


func ifUillStr(v1 string, v2 string) string {
	if len(v1) > 0 {
		return v1
	}
	return v2
}
