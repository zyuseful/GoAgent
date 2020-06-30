package perception

import "time"

const (
	BIT_ActiveOrDeaded = 1 << 1
	BIT_CheckThis      = 1 << 2
	BIT_CheckCome      = 1 << 3
)

//感知节点
type PNode struct {
	//节点名称
	name string
	ip   string
	port string
	//节点地址
	addr string
	//节点状态
	/**
	8	4	2	1
			   0/1		--死/活
		   0/1			--来源检测检测/来源未检
	   0/1				--本机未检测/本机检测

			1   0		--来源判断死
			1   1       --来源判断活
		0	1	0		--本机未检测 来源机死
		0	1	1		--本机未检测 来源机活
		1	1	0		--本机已检测 来源机死
		1	1	1		--本机已检测 来源机活
	*/
	state uint8
	//最后一次更新时间
	upTime time.Time
	//当前节点系统时间  用于不同节点之间时间持平，请求回复时需将时间设置为最新
	localTime time.Time
}

//---------------------------------构造创建---------------------------------------
//创建 PNode
func CreatePNode() *PNode {
	result := &PNode{}
	return result
}
//创建 PNode
func CreatePNodeByParams(name string, ip string, port string, addr string,state uint8,upTime time.Time,localTime time.Time) *PNode{
	result := &PNode{
		name,ip,port,addr,state,upTime,localTime,
	}
	return result
}

//---------------------------------基本 GetSet---------------------------------------
//获取节点Name
func (this *PNode) GetName() string {
	return this.name
}
//设置节点IP
func (this *PNode) SetName(name string) {
	this.name = name
}

//获取节点IP
func (this *PNode) GetIp() string {
	return this.ip
}
//设置节点IP
func (this *PNode) SetIp(ip string) {
	this.ip = ip
}

//获取节点Port
func (this *PNode) GetPort() string {
	return this.port
}
//设置节点Port
func (this *PNode) SetPort(port string) {
	this.port = port
}

//获取节点Addr
func (this *PNode) GetAddr() string {
	return this.addr
}
//设置节点Addr
func (this *PNode) SetAddr(addr string) {
	this.addr = addr
}

//获取节点upTime
func (this *PNode) SetUpTime(t time.Time) {
	this.upTime = t
}
//设置节点upTime
func (this *PNode) GetUpTime() time.Time {
	return this.upTime
}

//获取节点localTime
func (this *PNode) SetLocalTime(t time.Time) {
	this.localTime = t
}
//设置节点localTime
func (this *PNode) GetLocalTime() time.Time {
	return this.localTime
}

//获取节点state
func (this *PNode) SetState(s uint8) {
	this.state = s
}
//设置节点state
func (this *PNode) GetState() uint8 {
	return this.state
}

//---------------------------------状态设置---------------------------------------
//
//设置 当前节点 存活 + 检测
func (this *PNode) SetThisPNodeActive() {
	//设置状态为： 本机检测 活
	this.state = this.state | BIT_ActiveOrDeaded
	this.state = this.state | BIT_CheckThis
}
//设置 当前节点 死亡 + 检测
func (this *PNode) SetThisPNodeDeaded() {
	//设置状态为： 本机检测 活
	this.state = this.state &^ BIT_ActiveOrDeaded
	this.state = this.state | BIT_CheckThis
}

func (this *PNode) SetCheckComePNode() {
	this.state = this.state | BIT_CheckCome
}
func (this *PNode) SetNoCheckComePNode() {
	this.state = this.state &^ BIT_CheckCome
}
func (this *PNode) SetCheckComePNodeActive() {
	this.state = this.state | BIT_CheckCome
	this.SetThisPNodeActive()
}
func (this *PNode) SetCheckComePNodeDeaded() {
	this.state = this.state | BIT_CheckCome
	this.SetThisPNodeDeaded()
}

func (this *PNode) CheckThisPNodeIsActive() bool {
	return this.state&BIT_ActiveOrDeaded > 0
}

func (this *PNode) CheckComePNode() bool {
	return this.state&BIT_CheckCome > 0
}


func (this *PNode) SetPNodeState(comeCheck bool, thisCheck bool, isActive bool) {
	if isActive {
		this.state = this.state | BIT_ActiveOrDeaded
	} else {
		this.state = this.state &^ BIT_ActiveOrDeaded
	}

	if thisCheck {
		this.state = this.state | BIT_CheckThis
	} else {
		this.state = this.state &^ BIT_CheckThis
	}

	if comeCheck {
		this.state = this.state | BIT_CheckCome
	} else {
		this.state = this.state &^ BIT_CheckCome
	}
}
func (this *PNode) GetPNodeState() (bool,bool,bool) {
	var sCome,sThis,sIsAlive bool

	sIsAlive = this.state&BIT_ActiveOrDeaded > 0
	sThis = this.state&BIT_CheckThis > 0
	sCome = this.state&BIT_ActiveOrDeaded > 0

	return sCome,sThis,sIsAlive
}

//---------------------------------状态设置---------------------------------------
func (this *PNode) SetPNode(name string, ip string, port string) {
	tn := time.Now()
	this.SetPNodeAndTime(name,ip,port,tn)
}
func (this *PNode) SetPNodeAndTime(name string, ip string, port string,t time.Time) {
	if len(name) > 0 {
		this.name = name
	}
	if len(ip) > 0 {
		this.ip = ip
	}
	if len(port) > 0 {
		this.port = port
	}

	if len(ip) > 0 && len(port) > 0 {
		this.addr = ip + ":" + port
	}
	this.SetThisPNodeActive()

	tn := time.Now()
	this.SetLocalTime(tn)
	this.SetUpTime(tn)
}





