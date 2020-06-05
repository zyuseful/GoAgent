package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	EchoS "myagent/src/core/web"
	"net/http"
	//"github.com/spf13/viper"
	common "myagent/src/core/common"
)

func main_afero1() {
	path := common.GetLocalPath()
	fmt.Println("当前路径：",path)
}

func main_afero2() {
	server := EchoS.EchoServer{}
	server.Init()
	server.GetEcho().GET("/", func(context echo.Context) error {
		return context.String(http.StatusOK,"Hello World")
	})
	server.Start()
}

func main() {
	path := common.GetLocalPath()
	path+="/conf/"

	viper.SetConfigName("server") // name of config file (without extension)
	viper.SetConfigType("yaml") // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("/Users/zys/go/src/myagent/conf")               // optionally look for config in the working directory
	viper.AddConfigPath("$HOME")               // optionally look for config in the working directory
	viper.AddConfigPath(".")               // optionally look for config in the working directory
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}
