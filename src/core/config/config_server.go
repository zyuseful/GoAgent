package config

import (
	"github.com/spf13/viper"
	MyCommon "myagent/src/core/common"
	"strings"
	"sync"
)

const (
	//----------------配置文件----------------
	//----server 配置----
	Server_Port        = "8999"
	Server_IP          = "localhost"
	Server_UploadSpace = "./uploadFiles"
	//Server_UsingHttps  = false //是否使用https 默认使用http
)

type MyConfig struct {
	//锁
	lock sync.Mutex
	//--------server 配置--------
	//服务port
	server_port string
	//服务IP
	server_ip string
	//服务文件上传
	server_uploadSpace string

	//是否使用https : true -> https,则需要cert+privateKey /false -> http,无需cert、privateKey
	server_usingHttps bool
	//是否使用 echo自动tls : autoHttps = true 使用 echo的自动 https / autoHttps = false 使用自己私钥、证书
	server_autoHttps  bool
	//使用https cert 文件
	server_httpsCertFile string
	//使用https 私钥 文件
	server_privateKeyFile string

	//----------------引用配置----------------
	//格式 ${MTCF__Server_Port} ,${MTCF__Server_IP}
	//MT是前缀，防止与系统环境变量冲突 CF是配置文件前缀
	referenceVariableMap map[string]string
}

func (this *MyConfig) BuildMyConfig() {
	this.lock.Lock()
	defer this.lock.Unlock()

	//初始化环境变量可变参数
	if this.referenceVariableMap == nil {
		this.referenceVariableMap = make(map[string]string)
	}

	//服务port
	this.server_port = ifUillStr(viper.GetString("server.port"), Server_Port)
	//服务IP
	this.server_ip = ifUillStr(viper.GetString("server.ip"), Server_IP)

	//服务文件上传
	this.server_uploadSpace = ifUillStr(viper.GetString("server.uploadSpace"), Server_UploadSpace)

	//是否使用https
	this.server_usingHttps = viper.GetBool("server.usingHttps")
	//是否使用 echo自动tls
	this.server_autoHttps = viper.GetBool("server.autoHttps")
	//cert file
	this.server_httpsCertFile = viper.GetString("server.httpsCertFile")
	//private key file
	this.server_privateKeyFile = viper.GetString("server.privateKeyFile")


	//初始化环境变量可变参数 -- 配置
	this.referenceVariableMap["${MTCF__Server_UploadSpace}"] = this.server_uploadSpace
	this.referenceVariableMap["${MTPG__Run_Space}"] = MyCommon.GetLocalPath()
}

//----------------Server 部分----------------
/** 获取服务 ip */
func (this *MyConfig) GetServerIP() string {
	return this.server_ip
}
/** 获取服务 port */
func (this *MyConfig) GetServerPort() string {
	return this.server_port
}
/** 获取上传地址 */
func (this *MyConfig) GetServerUploadSpace() string {
	return this.server_uploadSpace
}
/** 是否使用https */
func (this *MyConfig) GetServerUsingHttps() bool {
	return this.server_usingHttps
}
/** 是否使用 echo自动tls */
func (this *MyConfig) GetServerAutoHttps() bool {
	return this.server_autoHttps
}
/** https cert 文件 */
func (this *MyConfig) GetServerHttpsCertFile() string {
	return this.server_httpsCertFile
}
/** https privateKey 文件 */
func (this *MyConfig) GetServerHttpsPrivateKeyFile() string {
	return this.server_privateKeyFile
}

//----------------引用变量  递归----------------
func (this *MyConfig) TransformReferenceVariableStr(cmd string) string {
	if len(cmd) <= 0 {
		return cmd
	}

	contains := strings.Contains(cmd, "${MT")
	if contains {
		for k, v := range this.referenceVariableMap {
			if strings.Contains(cmd, k) {
				cmd = strings.Replace(cmd, k, v, -1)
			}
		}
	}

	return cmd
}

/**
获取所有 程序配置 环境变量
*/
func (this *MyConfig) GetALLReferenceVariableStr() []string {
	keys := make([]string, 0, len(this.referenceVariableMap))
	for k := range this.referenceVariableMap {
		keys = append(keys, k)
	}
	return keys
}

//----------------方法工具----------------
func ifUillStr(v1 string, v2 string) string {
	if len(v1) > 0 {
		return v1
	}
	return v2
}
