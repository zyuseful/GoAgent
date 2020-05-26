package config

import (
	"github.com/spf13/viper"
	"sync"
)

const (
	//server 配置
	Server_Port = "8999"
	Server_IP   = "localhost"
)

type MyConfig struct {
	//锁
	lock sync.Mutex
	//server 配置
	server_port string
	server_ip   string
}

func (this *MyConfig) BuildMyConfig() {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.server_port = ifUillStr(viper.GetString("server.port"), Server_Port)
	this.server_ip = ifUillStr(viper.GetString("server.ip"), Server_IP)
}

//----------------Server 部分----------------
func (this *MyConfig) GetServerIP() string  {
	return this.server_ip
}
func (this *MyConfig) GetServerPort() string  {
	return this.server_port
}



func ifUillStr(v1 string, v2 string) string {
	if len(v1) > 0 {
		return v1
	}
	return v2
}
