package web

import "github.com/labstack/echo/v4"

type EchoServer struct {
	IP   string
	PORT string
	echo *echo.Echo
}

func (this *EchoServer) GetEcho() *echo.Echo {
	return this.echo
}
func (this *EchoServer) Init() {
	this.echo = echo.New()
}

func (this *EchoServer) Start() {
	ipAndPort := ""
	if (len(this.IP)<=0) {
		ipAndPort=":"
	}else {
		ipAndPort=this.IP +":"
	}

	if (len(this.PORT)<=0) {
		ipAndPort+="8999"
	}else {
		ipAndPort+=this.PORT
	}

	this.echo.Start(ipAndPort)
}

