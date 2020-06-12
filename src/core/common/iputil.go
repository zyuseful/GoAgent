package common

import (
	"fmt"
	"net"
)

//获取本机ip
func GetLocalIp() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println("get local ip failed")
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
