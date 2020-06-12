package controller

import (
	"bytes"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"io"
	"io/ioutil"
	"myagent/src/core/structure"
	MyWeb "myagent/src/core/web"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	HTTP = "http://"
	HTTPS = "https://"

	//提示信息
	AgentNotYetOpen = "尚未开启Agent功能"
)

type BaseController struct {
	echoServer *MyWeb.EchoServer
	arr        structure.ArrayList
}
type IBaseController interface {
	ControllerInit(myWeb *MyWeb.EchoServer)
}


//实现接口
func (this *BaseController) ControllerInit(myWeb *MyWeb.EchoServer) {
	this.echoServer = myWeb
}

func (this *BaseController) NullParam(info ...string) (string, bool) {
	var buffer bytes.Buffer
	size := len(info)
	var isEmpty bool = false

	if size <= 0 {
		return "", true
	}

	for i := 0; i < size; i++ {
		if len(info[i]) <= 0 {
			isEmpty = isEmpty || true
			buffer.WriteString("params [")
			buffer.WriteString(strconv.Itoa(i))
			buffer.WriteString("]")
			buffer.WriteString(" not found; ")
		}
	}
	return buffer.String(), isEmpty
}

//获取配置文件中的默认文件上传路径
func (this *BaseController) GetConfigUploadFilePath() string {
	//配置文件中的上传地址
	space := this.echoServer.GetConfig().GetServerUploadSpace()
	if len(space) <= 0 {
		return ""
	}

	var upFilePath strings.Builder
	upFilePath.WriteString(space)
	upFilePath.WriteString("/")

	return upFilePath.String()
}

/**
context echo.Context  echo 服务
paramStr 请求参数名称
paramType 来源:
	0 and !(1,2)   default  --  context.QueryParam(paramStr) 请求样例: curl -X GET http://localhost:1323\?name\=Joe
		 		   				context.FormValue(paramStr)  请求样例: curl -X POST http://localhost:1323 -d 'name=Joe'
	1 返回 context.QueryParam(paramStr)
	2 返回 context.FormValue(paramStr)
*/
func (this *BaseController) GetParamFormEchoContext(context echo.Context, paramStr string, paramType int) string {
	var paramRs string
	switch paramType {
	case 1:
		paramRs = context.QueryParam(paramStr)
	case 2:
		paramRs = context.FormValue(paramStr)
	default:
		paramRs = this.IfNullStr(context.QueryParam(paramStr), context.FormValue(paramStr))
	}

	//如果获取的长度大于 0 -- 进行引用判断
	if len(paramRs) > 0 {
		paramRs = this.echoServer.GetConfig().TransformReferenceVariableStr(paramRs)
	}

	return paramRs
}

/**
多行命令处理 拼接为一条命令 （windows、linux支持) 如：cmd1 && cmd2 && cmd3
*/
func (this *BaseController) MultilineCmdDispose(cmds string) string {
	if len(cmds) > 0 {
		return strings.Replace(cmds, "\r\n", " && ", -1)
	}
	return cmds
}

func (this *BaseController) IfNullStr(s1 string, s2 string) string {
	if len(s1) <= 0 {
		return s2
	} else {
		return s1
	}
}

func (this *BaseController) Success(result interface{}) WebResult {
	webResult := WebResult{
		Msg:    "",
		Result: result,
		Status: 0,
	}
	return webResult
}

func (this *BaseController) Failed(msg string, result interface{}) WebResult {
	webResult := WebResult{
		Msg:    msg,
		Result: result,
		Status: 1,
	}
	return webResult
}

func (this *BaseController) AddMethod(method string) {
	this.arr.Add(method)
}

func (this *BaseController) GetMethods() structure.ArrayList {
	return this.arr
}

func (this *BaseController) GetHttpOrHttpsStr() string {
	usinghttps := this.echoServer.GetConfig().GetServerUsingHttps()
	if usinghttps {
		return HTTPS
	}else {
		return HTTP
	}
}

/** 发送GET请求
// url：         请求地址
// response：    请求返回的内容
*/
func Get(url string,outTimeSecond int64) string {

	// 超时时间：5秒
	client := &http.Client{Timeout: time.Duration(outTimeSecond) * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	var buffer [512]byte
	result := bytes.NewBuffer(nil)
	for {
		n, err := resp.Body.Read(buffer[0:])
		result.Write(buffer[0:n])
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
	}

	return result.String()
}

/** 发送 Post请求
// url：         请求地址
// data：        POST请求提交的数据
// contentType： 请求体格式，如：application/json
// content：     请求放回的内容
*/
func (this *BaseController) Post(url string, data interface{}, contentType string, outTimeSecond int64) string {
	// 超时时间：5秒
	client := &http.Client{Timeout: time.Duration(outTimeSecond) * time.Second}
	jsonStr, _ := json.Marshal(data)
	if len(contentType) <=0 {
		contentType = "application/json"
	}
	resp, err := client.Post(url, contentType, bytes.NewBuffer(jsonStr))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	result, _ := ioutil.ReadAll(resp.Body)
	return string(result)
}

//web 返回结果
type WebResult struct {
	Msg    string
	Result interface{} // result
	Status int         // 0 OK, 1 FAILURE
}
