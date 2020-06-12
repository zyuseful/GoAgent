package perception

import (
	"myagent/src/core/structure"
	"sync"
	"time"
)

const (
	//请求时间差 (秒)
	Request_Time_Difference = 5
)

type PerceptionAgent struct {
	//读写锁  --  一把锁
	RWLock sync.RWMutex
	//当前Agent
	MySelf  *PNode
	RootMap map[string]*RsLine
	/**
	  	计算后的关系
	  		----------------------------------------------------------------------------------------------
	  	   | 							     RootMap map[string]*RsLine			     			     	  |
	  	   |----------------------------------------------------------------------------------------------|
	  	   |	            |   RootMap   |       RsLine.linealMap(map[string]*structure.ArrayList)		  |
	  	   |----------------------------------------------------------------------------------------------|
	  	   |      区域       |		                        关系图										  |
	  	   |----------------------------------------------------------------------------------------------|
	  	   |     			|	    	  |		 B		|	C	,	D	,	E	,	F	,	G			  |
	  	   |     MySelf		|	   A 	  |		 B		|   C	,	F	,	G							  |
	       |     			|	    	  |		 F		|	G											  |
	  	   |----------------------------------------------------------------------------------------------|
	  	   |     Other		|	   B	  |      		|
	  		----------------------------------------------------------------------------------------------
	*/
}

var perceptionAgent *PerceptionAgent

//初始化
func init() {
	perceptionAgent = new(PerceptionAgent)
}

//---------------------获取方法集---------------------------
func (this *PerceptionAgent) GetMySelfKey_NoLock() string {
	return this.MySelf.GetPNodeIp()
}

/** 获取当前Agent 自己的 ADDR */
func (this *PerceptionAgent) GetMySelfADDR_NoLock() string {
	return this.MySelf.ADDR
}

/** 获取myself RsLine 无锁 */
func (this *PerceptionAgent) GetMySelfRsLines_NoLock() *RsLine {
	this.CreateWhenMySelfRsLineIsNull_NoLock()
	return this.RootMap[this.GetMySelfKey_NoLock()]
}

/** 检查myself RsLine 是否为空，若为空则创建 */
func (this *PerceptionAgent) CreateWhenMySelfRsLineIsNull_NoLock() {
	if this.RootMap[this.GetMySelfADDR_NoLock()] == nil {
		this.RootMap[this.GetMySelfADDR_NoLock()] = CreateRsLine()
	}
}

/** 对外 设置 感知Agent 用于 controller 调用 */
func (this *PerceptionAgent) SetPerceptionAgent(name string, ip string, port string) {
	this.RWLock.Lock()
	defer this.RWLock.Unlock()
	this.MySelf.SetPNode(name, ip, port)
}

/** 对外 获取感知Agent */
func GetPerceptionAgent() *PerceptionAgent {
	return perceptionAgent
}

/** ------------  RelationshipLines 直系节点操作  ------------ */
//添加 场景  1、myself 2、myself 添加 B
func (this *PerceptionAgent) UpdatePerceptionAgentRsLineSync(come *PerceptionAgent) {
	this.RWLock.Lock()
	this.contrastAndUpdatePerceptionAgent(come)
	defer this.RWLock.Unlock()
}

/**  */
func (this *PerceptionAgent) contrastAndUpdatePerceptionAgent(come *PerceptionAgent) {
	//首先判断 是我自己还是外来
	thisSelfKey := this.GetMySelfKey_NoLock()
	comeSelfKey := come.GetMySelfKey_NoLock()

	//---------------- 1、myself MySelf区更新 ----------------
	if thisSelfKey == comeSelfKey {
		this.MySelf = come.MySelf
		this.RootMap = come.RootMap
	} else {
		//---------------- 2、other 区更新 TODO ----------------
		/** 直接替换来 */
		//不存在节点 -- 直接拿来
		if nil == this.GetMySelfRsLines_NoLock() {
			this.RootMap[comeSelfKey] = come.GetMySelfRsLines_NoLock()
		} else {
			//other区 -> myself区
			var nodeTimeContrast int64
			nodeTimeContrast = NodeTimeContrast(come)

			//时间比对符合更新
			if this.MySelf.UpTime.Unix() < (come.MySelf.UpTime.Unix() + nodeTimeContrast + Request_Time_Difference) {
				this.RootMap[comeSelfKey] = come.GetMySelfRsLines_NoLock()
				t := time.Now()
				this.RootMap[comeSelfKey].LinealNode.SetPNodeLocalTime(t)
				this.RootMap[comeSelfKey].LinealNode.SetPNodeUpTime(t)
			}
		}
	}
	this.comeAgentAppendToThisAgent(come)
}

func (this *PerceptionAgent) comeAgentAppendToThisAgent(come *PerceptionAgent) {
	keyArr := structure.ArrayList{}
	for kv,_:=range this.RootMap {
		if len(kv) > 0 {
			keyArr.Add(kv)
		}
	}

	for i:=0;i<keyArr.Size();i++ {
		getKey := keyArr.Get(i).(string)
		rsline := come.TransformRootMapToRsLine(getKey)
		this.TransformRsLineToRootMap(getKey,rsline)
	}
}

/** RootMap 根据 rootKey 提取 转换为 RsLine */
func (this *PerceptionAgent) TransformRootMapToRsLine(rootKey string) *RsLine {
	result := CreateRsLine()
	srcRsLine := this.RootMap[rootKey]

	result.LinealNode = this.MySelf
	if srcRsLine != nil {
		//临时用来存储 源定位中的 key
		keyList := &structure.ArrayList{}
		for k, v := range srcRsLine.LinealMap {
			if len(k) > 0 {
				tempList := structure.ArrayList{}
				tempList.Add(this.MySelf)
				for i:=0;i<v.Size();i++ {
					tempList.Add(v.Get(i))
				}
				keyList.Add(tempList)
			}
		}
		result.LinealMap[rootKey] = keyList
	}
	return result
}
func (this *PerceptionAgent) TransformRsLineToRootMap(rootKey string, rs *RsLine) {
	//MySelf区更新
	if this.GetMySelfKey_NoLock() == rootKey {
		//TODO 更新判断
		this.MySelf = rs.LinealNode
	}

	srcRsLine := this.RootMap[rootKey]
	if nil == srcRsLine {
		this.RootMap[rootKey] = CreateRsLine()
	}

	l1 := rs.LinealMap[rootKey]
	for i:=0;i<l1.Size();i++ {
		l2 := l1.Get(i).(structure.ArrayList)
		node := l2.Get(0).(*PNode)
		this.RootMap[rootKey].LinealMap[node.IP].Destroy()
		if this.RootMap[rootKey].LinealMap[node.IP] == nil {
			this.RootMap[rootKey].LinealMap[node.IP] = &structure.ArrayList{}
		}
		this.RootMap[rootKey].LinealNode = node
		for j:=1;j<l2.Size();j++ {
			this.RootMap[rootKey].LinealMap[node.IP].Add(l2.Get(j))
		}
	}
}

/** 节点时间比对
时间差 = 当前时间 - 请求时刻对方节点当前时间 - 传输时间
使用时：本地时间 + 时间差 后再做时间运算
*/
func NodeTimeContrast(come *PerceptionAgent) int64 {
	comeTimeUnix := come.MySelf.GetPNodeLocalTime().Unix()
	localTimeunix := time.Now().Unix()

	var result int64
	result = (localTimeunix - comeTimeUnix)
	result -= Request_Time_Difference
	return result
}

/** 初始化 MySelfRsLines
相当于
	|A | - | - |
*/
func (this *PerceptionAgent) InitMySelfRsLines() {
	this.RWLock.Lock()
	defer this.RWLock.Unlock()

	//RsLine init
	if this.RootMap[this.GetMySelfKey_NoLock()] == nil {
		this.RootMap[this.GetMySelfKey_NoLock()] = CreateRsLine()
	}
	//RsLine.InnerLine init
	if this.RootMap[this.GetMySelfKey_NoLock()] == nil {
		this.RootMap[this.GetMySelfKey_NoLock()] = CreateRsLine()
	}

	//TODO
	//this.RootMap[this.GetMySelfKey_NoLock()].SetRsLine_LinealNode(time.Now())
}
