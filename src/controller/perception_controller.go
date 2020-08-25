package controller

import (
	"encoding/json"
	"fmt"
	MyCommon "myagent/src/core/common"
	MyDnet "myagent/src/core/donet"
	MyWeb "myagent/src/core/web"
	"net/http"

	"github.com/labstack/echo/v4"
)

//url 常量
const (
	/** 设置节点 */
	PerceptionController__SetLocalPerception = "/PerceptionController/SetLocalPerception"
	/** 注册节点 */
	PerceptionController__RegisterTo = "/PerceptionController/RegisterTo"
	/** 接收注册节点 */
	PerceptionController__RegisterCome = "/PerceptionController/RegisterCome"
	/** 获取我的感知信息 */
	PerceptionController__ResultMyPerception = "/PerceptionController/ResultMyPerception"
)

/** 感知服务 负责节点互加、路径计算 */
type PerceptionController struct {
	BaseController
}

//实现接口
func (this *PerceptionController) ControllerInit(myWeb *MyWeb.EchoServer) {
	this.BaseController.ControllerInit(myWeb)

	//将所需服务在这里自行维护
	/** 获取我的感知信息 */
	this.ResultMyPerception()
	/** 设置节点 */
	this.SetLocalPerception()
	/** 注册节点 由当前节点出发添加其他节点 */
	this.RegisterTo()
	/** 接收注册节点信息 */
	this.RegisterCome()
}

/** 获取我的感知信息 */
func (this *PerceptionController) ResultMyPerception() {
	this.echoServer.GetEcho().GET(PerceptionController__ResultMyPerception, func(context echo.Context) error {
		return context.JSON(http.StatusOK, this.Success(MyDnet.GetPerceptionDNET()))
	})
}

/** 设置节点 */
func (this *PerceptionController) SetLocalPerception() {
	this.echoServer.GetEcho().GET(PerceptionController__SetLocalPerception, func(context echo.Context) error {
		name := this.GetParamFormEchoContext(context, "name", 0)
		ip := this.GetParamFormEchoContext(context, "ip", 0)

		MyDnet.GetPerceptionDNET().UpdateRootNode(name, ip, "")
		MyDnet.GetPerceptionDNET().RootNode.SetStateArrive()
		MyDnet.GetPerceptionDNET().RootNode.SetStateActive()
		ss := MyDnet.GetPerceptionDNET()
		//zy
		fmt.Println(ss)
		return context.JSON(http.StatusOK, this.Success(MyDnet.GetPerceptionDNET()))
	})
}

/** 注册节点 从A节点操作向B注册*/
func (this *PerceptionController) RegisterTo() {
	this.echoServer.GetEcho().POST(PerceptionController__RegisterTo, func(context echo.Context) error {
		if !this.echoServer.GetConfig().GetServerCheckAgent() {
			return context.JSON(http.StatusOK, this.Failed(AgentNotYetOpen, nil))
		}

		toIP := this.GetParamFormEchoContext(context, "toIP", 0)
		toPort := this.GetParamFormEchoContext(context, "toPort", 0)
		paramStr, b := this.NullParam(toIP, toPort)
		if b {
			return context.JSON(http.StatusOK, this.Failed(paramStr, nil))
		}

		ss := MyDnet.GetPerceptionDNET()
		//zy
		fmt.Println(ss)

		//同步处理 A -> B
		concatURL := MyCommon.StringConcat(this.GetHttpOrHttpsStr(), toIP, ":", toPort, PerceptionController__RegisterCome)
		//post := this.Post(concatURL, MyDnet.GetPerceptionDNET(), "", 0)
		post := this.Post(concatURL, ss, "", 0)
		result := new(MyDnet.PerceptionDNET)
		json.Unmarshal([]byte(post), result)
		//fmt.Println(post)
		//fmt.Println(result)

		//接收返回 Agent进行计算
		MyDnet.GetPerceptionDNET().MergePerceptionDNETSync(result)

		//给A 调用者发送消息
		return context.JSON(http.StatusOK, this.Success(MyDnet.GetPerceptionDNET()))
	})
}

/** 接收注册 从A节点操作向B注册 这里的A就是注册来源*/
func (this *PerceptionController) RegisterCome() {
	this.echoServer.GetEcho().POST(PerceptionController__RegisterCome, func(context echo.Context) error {
		if !this.echoServer.GetConfig().GetServerCheckAgent() {
			return context.JSON(http.StatusOK, this.Failed(AgentNotYetOpen, nil))
		}

		agent := new(*MyDnet.PerceptionDNET)
		if err := context.Bind(agent); err != nil {
			//TODO
			return context.JSON(http.StatusOK, this.Failed("TODO", nil))
		}

		//B -> A
		//MyDnet.GetPerceptionDNET().MergePerceptionDNETSync(agent)
		return context.JSON(http.StatusOK, MyDnet.GetPerceptionDNET())
	})
}
