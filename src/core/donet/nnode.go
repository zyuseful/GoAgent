package donet

import (
	"strings"
	"time"
)

const (
	//0 死 / 1 存活
	BIT_ActiveOrDeaded = 1 << 1
	//0 推测节点 / 1 可到达节点
	BIT_CheckThis = 1 << 2
)

/**
域网感知节点
*/
type DNode struct {
	//节点名称
	Name string
	IP   string
	Port string
	//节点地址
	Addr string
	//节点状态
	/**
	保留				   保留				0推测/1可达			0死/1活
	8					4					2					1

	*/
	State uint8
	//最后一次更新时间
	UpTime time.Time
	//当前节点系统时间  用于不同节点之间时间持平，请求回复时需将时间设置为最新
	LocalTime time.Time
}

//---------------------------------构造创建---------------------------------------
//创建 DNode
func CreateDNode() *DNode {
	result := &DNode{}
	return result
}

//创建 DNode
func CreateDNodeByParamsWithOutTime(name string, ip string, port string, addr string, state uint8) *DNode {
	now := time.Now()
	result := &DNode{
		name, ip, port, addr, state, now, now,
	}
	return result
}

//创建 DNode
func CreateDNodeByParams(name string, ip string, port string, addr string, state uint8, upTime time.Time, localTime time.Time) *DNode {
	result := &DNode{
		name, ip, port, addr, state, upTime, localTime,
	}
	return result
}

//---------------------------------基本 GetSet---------------------------------------
//获取节点Name
func (this *DNode) GetName() string {
	return this.Name
}

//设置节点IP
func (this *DNode) SetName(name string) {
	this.Name = name
}

//获取节点IP
func (this *DNode) GetIp() string {
	return this.IP
}

//设置节点IP
func (this *DNode) SetIp(ip string) {
	this.IP = ip
}

//获取节点Port
func (this *DNode) GetPort() string {
	return this.Port
}

//设置节点Port
func (this *DNode) SetPort(port string) {
	this.Port = port
}

//获取节点Addr
func (this *DNode) GetAddr() string {
	return this.Addr
}

//设置节点Addr
func (this *DNode) SetAddr(addr string) {
	this.Addr = addr
}

//获取节点upTime
func (this *DNode) SetUpTime(t time.Time) {
	this.UpTime = t
}

//设置节点upTime
func (this *DNode) GetUpTime() time.Time {
	return this.UpTime
}

//获取节点localTime
func (this *DNode) SetLocalTime(t time.Time) {
	this.LocalTime = t
}

//设置节点localTime
func (this *DNode) GetLocalTime() time.Time {
	return this.LocalTime
}

//获取节点state
func (this *DNode) SetState(s uint8) {
	this.State = s
}

//设置节点state
func (this *DNode) GetState() uint8 {
	return this.State
}

//---------------------------------状态设置---------------------------------------
//设置 当前节点 存活
func (this *DNode) SetStateActive() {
	this.State = this.State | BIT_ActiveOrDeaded
}

//设置 当前节点 死亡
func (this *DNode) SetStateDeaded() {
	this.State = this.State &^ BIT_ActiveOrDeaded
}

//设置 当前节点 推测检测
func (this *DNode) SetStateGuess() {
	this.State = this.State &^ BIT_CheckThis
}

//设置 当前节点 可到达
func (this *DNode) SetStateArrive() {
	this.State = this.State | BIT_CheckThis
}

//检测 节点是否存活
func (this *DNode) NNodeIsActive() bool {
	return this.State&BIT_ActiveOrDeaded > 0
}

//检测 节点是否可到达
func (this *DNode) NNodeIsReachable() bool {
	return this.State&BIT_CheckThis > 0
}

func (this *DNode) SetPNode(name string, ip string, port string) {
	tn := time.Now()
	this.SetPNodeAndTime(name, ip, port, tn)
}
func (this *DNode) SetPNodeAndTime(name string, ip string, port string, t time.Time) {
	if len(name) > 0 {
		this.Name = name
	}
	if len(ip) > 0 {
		this.IP = ip
	}
	if len(port) > 0 {
		this.Port = port
	}

	if len(ip) > 0 && len(port) > 0 {
		var builder strings.Builder
		builder.WriteString(ip)
		builder.WriteString(":")
		builder.WriteString(port)
		this.Addr = builder.String()
	}

	tn := time.Now()
	this.SetLocalTime(tn)
	this.SetUpTime(tn)
}
