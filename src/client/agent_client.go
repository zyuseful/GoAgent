package main

import (
	"flag"
	"fmt"
	MyConfig "myagent/src/core/config"
	"os"
)

const (
	default_ip   = MyConfig.Server_IP
	default_port = MyConfig.Server_Port
)

var (
	// help
	h bool

	// server
	ip   string
	port string
)

func init() {
	flag.BoolVar(&h, "h", false, "this help")

	flag.StringVar(&ip, "ip", "", "set myagent server ip for connect")
	flag.StringVar(&port, "port", "", "set myagent server port for connect")

	// 改变默认的 Usage，flag包中的Usage 其实是一个函数类型。这里是覆盖默认函数实现，具体见后面Usage部分的分析
	flag.Usage = usage
}

func usage() {
	fmt.Fprintf(os.Stderr, `myagent client:
Usage:agent_client [-h] [-ip server ip] [-port server port]

Options:
`)
	flag.PrintDefaults()
}

func main() {
	mainStop := make(chan bool, 1)
	if len(ip) <=0 || len(port) <=0 {
		fmt.Println("ip or port is null, you can use 'help' to show and set ip,port !")
		mainStop<-true
	}else {
		go run(mainStop)
	}
	<-mainStop
}

func run(stop chan bool) {
	var str string
	fmt.Printf("请输入内容：")
	for {
		str = ""
		fmt.Scanf("%s", &str)
		switch str {
		case "":
			fmt.Println()
			break
		case "quit":
			exit(stop,"Bye!")
			break
		case "help":
			usage()
			exit(stop,"Bye!")
			break
		default:
			fmt.Println(str)
		}
	}
}

func exit(stop chan bool,info string) {
	fmt.Println(info)
	stop <- true
}
