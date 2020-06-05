package web

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/crypto/acme/autocert"
	MyConfig "myagent/src/core/config"
)

type EchoServer struct {
	IP   string
	PORT string
	echo *echo.Echo
	myConfig *MyConfig.MyConfig
}

func (this *EchoServer) GetEcho() *echo.Echo {
	return this.echo
}
func (this *EchoServer) GetConfig() *MyConfig.MyConfig {
	return this.myConfig
}
func (this *EchoServer) SetConfig(myConfig *MyConfig.MyConfig) {
	this.myConfig = myConfig
}

func (this *EchoServer) Init() {
	this.echo = echo.New()
}

func (this *EchoServer) Start() {
	//ip + port 处理
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

	//http、https处理
	usingHttps := this.myConfig.GetServerUsingHttps()
	//是否使用 echo自动tls
	authHttps := this.myConfig.GetServerAutoHttps()
	//https cert 文件
	certFile := this.myConfig.GetServerHttpsCertFile()
	//https privateKey 文件
	privateKey := this.myConfig.GetServerHttpsPrivateKeyFile()


	if !usingHttps {
		//http
		this.echo.Start(ipAndPort)
	} else {
		if authHttps {
			//https auto tls
			// Cache certificates
			this.echo.AutoTLSManager.Cache = autocert.DirCache("/var/www/.cache")
			this.echo.Use(middleware.Recover())
			this.echo.Use(middleware.Logger())
			//this.echo.StartAutoTLS(ipAndPort)
			this.echo.StartAutoTLS(":443")
		}else {
			//https private tls
			this.echo.StartTLS(ipAndPort,certFile,privateKey)
		}
	}


}

