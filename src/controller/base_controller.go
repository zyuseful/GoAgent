package controller

import (
	"bytes"
	"myagent/src/core/structure"
	MyWeb "myagent/src/core/web"
)

type BaseController struct {
	echoServer *MyWeb.EchoServer
	arr structure.ArrayList
}
type IBaseController interface {
	ControllerInit(myWeb *MyWeb.EchoServer)
}

//实现接口
func (this *BaseController) ControllerInit(myWeb *MyWeb.EchoServer) {
	this.echoServer = myWeb
}

func (this *BaseController) NullParam(info ...string) (string,bool) {
	var buffer bytes.Buffer
	size := len(info)
	var isEmpty bool = true

	if size <= 0 {
		return "",isEmpty
	}

	for i := 0; i < size; i++ {
		if (len(info[i]) < 0 ) {
			isEmpty = false
			buffer.WriteString(info[i])
			buffer.WriteString(" not found; ")
		}
	}
	return buffer.String(),isEmpty
}

func (this *BaseController) Success(result interface{}) WebResult {
	webResult := WebResult{
		Msg:    "",
		Result: result,
		Status: 0,
	}
	return webResult
}

func (this *BaseController) Failed(msg string,result interface{}) WebResult {
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

func (this *BaseController) GetMethods() (structure.ArrayList){
	return this.arr
}

//web 返回结果
type WebResult struct {
	Msg    string
	Result interface{} // result
	Status int         // 0 OK, 1 FAILURE
}
