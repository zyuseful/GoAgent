package call

import (
	"fmt"
	"github.com/spf13/viper"
	MyConfig "myagent/src/core/config"
)

func ConfigInit(runPath string) (*MyConfig.MyConfig) {
	viper.SetConfigName("server")   // name of config file (without extension)
	viper.SetConfigType("yaml")     // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(runPath) // optionally look for config in the working directory
	viper.AddConfigPath("$HOME")    // call multiple times to add many search paths
	viper.AddConfigPath(".")        // optionally look for config in the working directory
	err := viper.ReadInConfig()     // Find and read the config file
	if err != nil {                 // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	currentConfig := &MyConfig.MyConfig{}
	//完成全局变量配置值变更
	currentConfig.BuildMyConfig()
	return currentConfig
}