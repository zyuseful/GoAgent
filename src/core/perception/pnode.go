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
	Name string
	IP   string
	PORT string
	//节点地址
	ADDR string
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
	State uint8
	//最后一次更新时间
	UpTime time.Time
	//当前节点系统时间  用于不同节点之间时间持平，请求回复时需将时间设置为最新
	LocalTime time.Time
}

func CreatePNode() *PNode {
	result := &PNode{}
	return result
}

func (this *PNode) SetPNode(name string, ip string, port string) {
	if len(name) > 0 {
		this.Name = name
	}
	if len(ip) > 0 {
		this.IP = ip
	}
	if len(port) > 0 {
		this.PORT = port
	}

	if len(ip) > 0 && len(port) > 0 {
		this.ADDR = ip + ":" + port
	}
	this.SetThisPNodeActive()

	tn := time.Now()
	this.SetPNodeLocalTime(tn)
	this.SetPNodeUpTime(tn)
}

func (this *PNode) SetPNodeUpTime(t time.Time) {
	this.UpTime = t
}
func (this *PNode) GetPNodeUpTime() time.Time {
	return this.UpTime
}

func (this *PNode) SetPNodeLocalTime(t time.Time) {
	this.LocalTime = t
}
func (this *PNode) GetPNodeLocalTime() time.Time {
	return this.LocalTime
}

func (this *PNode) SetThisPNodeActive() {
	//设置状态为： 本机检测 活
	this.State = this.State | BIT_ActiveOrDeaded
	this.State = this.State | BIT_CheckThis
}
func (this *PNode) SetThisPNodeDeaded() {
	//设置状态为： 本机检测 活
	this.State = this.State &^ BIT_ActiveOrDeaded
	this.State = this.State | BIT_CheckThis
}
func (this *PNode) SetCheckComePNode() {
	this.State = this.State | BIT_CheckCome
}
func (this *PNode) SetNoCheckComePNode() {
	this.State = this.State &^ BIT_CheckCome
}
func (this *PNode) SetCheckComePNodeActive() {
	this.State = this.State | BIT_CheckCome
	this.SetThisPNodeActive()
}
func (this *PNode) SetCheckComePNodeDeaded() {
	this.State = this.State | BIT_CheckCome
	this.SetThisPNodeDeaded()
}

func (this *PNode) SetPNodeState(comeCheck bool, thisCheck bool, isActive bool) {
	if isActive {
		this.State = this.State | BIT_ActiveOrDeaded
	} else {
		this.State = this.State &^ BIT_ActiveOrDeaded
	}

	if thisCheck {
		this.State = this.State | BIT_CheckThis
	} else {
		this.State = this.State &^ BIT_CheckThis
	}

	if comeCheck {
		this.State = this.State | BIT_CheckCome
	} else {
		this.State = this.State &^ BIT_CheckCome
	}
}

func (this *PNode) CheckThisPNodeIsActive() bool {
	return this.State&BIT_ActiveOrDeaded > 0
}

func (this *PNode) CheckComePNode() bool {
	return this.State&BIT_CheckCome > 0
}

func (this *PNode) GetPNodeIp() string {
	return this.IP
}

func (this *PNode) GetPNodeAddr() string {
	return this.ADDR
}

func (this *PNode) GetPNodePort() string {
	return this.PORT
}
func (this *PNode) GetPNodeName() string {
	return this.PORT
}
