package perception

import (
	"github.com/deckarep/golang-set"
	"sync"
	"time"
)

const (
	//请求时间差 (秒)
	Request_Time_Difference = 5
)

type PerceptionAgent struct {
	//读写锁  --  一把锁
	rwLock sync.RWMutex
	//当前Agent
	mySelf  *PNode
	rootMap map[string]*RsLine
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
	//顶点集合 <string> addr
	apexNodes mapset.Set

	//可到达定点 <string> addr  -- check服务填写
	reachNodes mapset.Set
}

//单例
var perceptionAgent *PerceptionAgent

//--------------------------------- 初始化---------------------------------------
//初始化
func init() {
	perceptionAgent = new(PerceptionAgent)
	perceptionAgent.mySelf = CreatePNode()
	perceptionAgent.rootMap = make(map[string]*RsLine)
}

/** 初始化 MySelfRsLines -- 供初始化时 Config使用
相当于
------------------------------------------
	MySelf
       A  -->  [A]
  /|\			|Node:A
------------------------------------------
   |			|Map
   |			|---[B]
   |				 |--(C,D,E)
   |			     |--(D,E)
   |			|---[C]
   |				 |--(D,E)
   |				 |--(F)
------------------------------------------
*/
func (this *PerceptionAgent) InitMySelfRsLinesSync(name string, ip string, port string) {
	this.rwLock.Lock()
	defer this.rwLock.Unlock()

	//RsLine init
	if this.rootMap[this.GetMySelfKey()] == nil {
		this.rootMap[this.GetMySelfKey()] = CreateRsLine()
	}
	//RsLine.InnerLine init
	if this.rootMap[this.GetMySelfKey()] == nil {
		this.rootMap[this.GetMySelfKey()] = CreateRsLine()
	}

	this.UpdateMySelfNode(name, ip, port)
	this.apexNodes = mapset.NewSet()
}

//---------------------------------获取方法集---------------------------------------
//获取当前 myself Node
func (this *PerceptionAgent) GetMySelf() *PNode {
	return this.mySelf
}

//获取当前 rootMap
func (this *PerceptionAgent) GetRootMap() map[string]*RsLine {
	return this.rootMap
}
//设置当前 rootMap
func (this *PerceptionAgent) SetRootMap(comeRootMap map[string]*RsLine) {
	this.rootMap = comeRootMap
}

//获取 MySelf Key
func (this *PerceptionAgent) GetMySelfKey() string {
	//return this.MySelf.GetPNodeIp()
	return this.mySelf.GetAddr()
}

//根据 Addrkey 获取 *RsLine
func (this *PerceptionAgent) GetRsLineByKey(addrKey string) *RsLine {
	return this.rootMap[addrKey]
}

//根据 Addrkey 获取 *RsLine ,如果 *RsLine 为nil则创建
func (this *PerceptionAgent) GetRsLineByKeyIfEmpty(addrKey string) *RsLine {
	if nil == this.GetRsLineByKey(addrKey) {
		this.rootMap[addrKey] = CreateRsLine()
	}
	return this.rootMap[addrKey]
}

func (this *PerceptionAgent) GetApexNodes() mapset.Set {
	return this.apexNodes
}
func (this *PerceptionAgent) SetApexNodes(apexNodes mapset.Set) {
	this.apexNodes.Clear()
	this.apexNodes.Add(apexNodes)
}
func (this *PerceptionAgent) AddToApexNodes(addr string) {
	this.apexNodes.Add(addr)
}

/** 对外 获取感知Agent */
func GetPerceptionAgent() *PerceptionAgent {
	return perceptionAgent
}
//---------------------------------复合更新---------------------------------------
//更新 MySelfNode 同时更新 RsLine 中的 对应的 MySelfNode
func (this *PerceptionAgent) UpdateMySelfNode(name string, ip string, port string) {
	//设置 PerceptionAgent.MySelf
	if nil == this.mySelf {
		this.mySelf = CreatePNode()
	}
	tn := time.Now()

	//记录原有Addr key
	srcAddr := this.GetMySelfKey()
	needChange := false

	this.mySelf.SetPNodeAndTime(name, ip, port, tn)
	nowAddr := this.GetMySelfKey()

	//IP or Port 做出修改，则对应Map value地址需要调整
	if srcAddr != nowAddr {
		this.rootMap[nowAddr] = this.rootMap[srcAddr]
		delete(this.rootMap, srcAddr)
		needChange = true
	}

	thisRsLine := this.GetRsLineByKeyIfEmpty(this.GetMySelfKey())

	//设置 PerceptionAgent.RootMap MySelf
	thisRsLine.GetKeyNode().SetPNodeAndTime(name, ip, port, tn)
	if needChange {
		thisRsLine.linealMap[nowAddr] = thisRsLine.linealMap[srcAddr]
		delete(thisRsLine.linealMap, srcAddr)
	}
}

//---------------------------------Controller 调用---------------------------------------
/** 对外 设置 感知Agent 用于 controller 调用 */
func (this *PerceptionAgent) SetPerceptionAgentSync(name string, ip string, port string) {
	this.rwLock.Lock()
	defer this.rwLock.Unlock()
	this.UpdateMySelfNode(name, ip, port)
}

//<=====


/** 对外 对方节点向己方合并 */
func (this *PerceptionAgent) RegisFromComeAgentSync(come *PerceptionAgent) {
	this.rwLock.Lock()
	defer this.rwLock.Unlock()

	if come == nil {
		return
	}

	//1 自关联注册 / 2 外来直系节点
	if this.GetMySelfKey() == come.GetMySelfKey() {
		//	自关联注册，比较时间后可直接更新
		this.updateMySelfAgent(come)
	} else {
		this.updateOtherAgent(come)
	}
}

//RegisFromComeAgent 调用 -- 自关联更新  -- 比较更新时间后进行更新
func (this *PerceptionAgent) updateMySelfAgent(come *PerceptionAgent) {
	//更新 MySelf
	if this.GetMySelf().GetUpTime().Unix() < come.GetMySelf().GetUpTime().Unix() {
		this.updatePerceptionAgent(come,0)
	}
}

//RegisFromComeAgent 调用 -- Other更新  -- 比较更新时间后进行更新
func (this *PerceptionAgent) updateOtherAgent(come *PerceptionAgent) {
	this.updatePerceptionAgent(come,1)
}

//更新感知节点
//come *PerceptionAgent 外来节点
//0 : 当前节点更新 / 1 : 外来节点更新
func (this *PerceptionAgent) updatePerceptionAgent(come *PerceptionAgent, selfOrOther int) {
	//当前节点更新
	if 0 == selfOrOther {
		//1 更新myself Node
		this.mySelf.SetPNode(come.mySelf.GetName(),come.mySelf.GetIp(),come.mySelf.GetPort())
		//2 更新

	} else  {
	//外来节点更新
	}
}

//<--
func (this *PerceptionAgent) updatePAgent(come *PerceptionAgent) {
	thisNodeSelf := this.GetMySelfKey()
	comeNodeSelf := come.GetMySelfKey()

	//当前节点更新
	if thisNodeSelf == comeNodeSelf {
		//1 更新myself Node
		this.mySelf.SetPNode(come.mySelf.GetName(),come.mySelf.GetIp(),come.mySelf.GetPort())
		//2 更新 RsLine.linealMap (myself + rsline)
		this.SetRootMap(come.GetRootMap())
		//3 更新定点
		this.SetApexNodes(come.GetApexNodes())
	} else {
	/**
	外来节点更新
		1、直属节点更新
						直属节点
					  /			\
		           无/            \有
		 设置到 linealMap		 事件判断
								/	    \
		 					   /		 \
				        无需更新		   需要更新(四种情况)：
											1、"我"不在
											2、"我"在行首
											3、"我"在行未
											4、"我"在行中

		2、Other区更新

		3、定点合并
	*/

	}
}

/** ------------  RelationshipLines 直系节点操作  ------------ */
//添加 场景  1、myself 2、myself 添加 B

/** RootMap 根据 rootKey 提取 转换为 RsLine */
func (this *PerceptionAgent) TransformRootMapToRsLine(rootKey string) *RsLine {
	result := CreateRsLine()
	//srcRsLine := this.RootMap[rootKey]

	//result.LinealNode = this.MySelf
	//if srcRsLine != nil {
	//	//临时用来存储 源定位中的 key
	//	keyList := &structure.ArrayList{}
	//	for k, v := range srcRsLine.LinealMap {
	//		if len(k) > 0 {
	//			tempList := structure.ArrayList{}
	//			tempList.Add(this.MySelf)
	//			for i:=0;i<v.Size();i++ {
	//				tempList.Add(v.Get(i))
	//			}
	//			keyList.Add(tempList)
	//		}
	//	}
	//	result.LinealMap[rootKey] = keyList
	//}
	return result
}

//
//func (this *PerceptionAgent) TransformRsLineToRootMap(rootKey string, rs *RsLine) {
//	//MySelf区更新
//	if this.GetMySelfKey_NoLock() == rootKey {
//		//1 为空直接拿来赋值  //2 对比更新 TODO
//		if this.RootMap[rootKey] == nil || rs.LinealNode.UpTime.Unix() > (this.MySelf.UpTime.Unix()+Request_Time_Difference){
//			this.MySelf.SetPNode(rs.LinealNode.GetPNodeName(),rs.LinealNode.GetPNodeIp(),rs.LinealNode.GetPNodePort())
//			rs.LinealNode.SetPNodeUpTime(this.MySelf.UpTime)
//			this.RootMap[rootKey] = rs
//		}
//	} else {
//		//	Other区 -- 更新原始
//		if nil == this.RootMap[rootKey] {
//			this.RootMap[rootKey] = CreateRsLine()
//		}
//
//		l1 := rs.LinealMap[rootKey]
//		for i:=0;i<l1.Size();i++ {
//			l2 := l1.Get(i).(structure.ArrayList)
//			node := l2.Get(0).(*PNode)
//			this.RootMap[rootKey].LinealMap[node.IP].Destroy()
//			if this.RootMap[rootKey].LinealMap[node.IP] == nil {
//				this.RootMap[rootKey].LinealMap[node.IP] = &structure.ArrayList{}
//			}
//			this.RootMap[rootKey].LinealNode = node
//			for j:=1;j<l2.Size();j++ {
//				this.RootMap[rootKey].LinealMap[node.IP].Add(l2.Get(j))
//			}
//		}
//	}
//
//}

/** 节点时间比对
时间差 = 当前时间 - 请求时刻对方节点当前时间 - 传输时间
使用时：本地时间 + 时间差 后再做时间运算
*/
func NodeTimeContrast(come *PerceptionAgent) int64 {
	comeTimeUnix := come.GetMySelf().GetLocalTime().Unix()
	localTimeunix := time.Now().Unix()

	var result int64
	result = localTimeunix - comeTimeUnix
	result -= Request_Time_Difference
	return result
}

