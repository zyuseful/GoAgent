package donet

import (
	"fmt"
	"myagent/src/core/structure"
	"sync"
	"time"
)

/**
           |----------------------------------------------------------------------------------------------|
	  	   |     			|	    	  |		 B		|	C	,	D	,	E	,	F	,	G			  |
	  	   |     myself区  	|	   A 	  |		 B		|   C	,	F	,	G							  |
	       |     			|	    	  |		 F		|	G											  |
	  	   |----------------------------------------------------------------------------------------------|
	  	   |     other区  	|	   B	  |      C 		|												  |
	  		----------------------------------------------------------------------------------------------
*/
type PerceptionDNET struct {
	RWLock sync.RWMutex
	//根起点
	RootNode *DNode
	//所属域 --> {"数据网","A子网","B子网"...}
	Domains *structure.ArrayList
	//成员map --> 成员Node Addr : 成员集合
	MemberMap map[string]*MemberGroup
	//成员节点+锁
	//MemberNodeAndLock map[string]*MemberNode
}

var ThisPerceptionDNET *PerceptionDNET

func GetPerceptionDNET() *PerceptionDNET {
	//perceptionDNET.RWLock.RLock()
	//defer perceptionDNET.RWLock.RUnlock()
	return ThisPerceptionDNET
}

//---------------------------------构造创建---------------------------------------
//创建 Perception_DNET 感知域网
func CreatePerceptionDNET(rootNode *DNode) *PerceptionDNET {
	var result *PerceptionDNET
	result = &PerceptionDNET{}
	//设置根节点
	result.RootNode = rootNode
	//所属域
	result.Domains = &structure.ArrayList{}
	//初始成员map
	result.MemberMap = make(map[string]*MemberGroup)
	//初始化成员节点+锁
	//result.MemberNodeAndLock = make(map[string]*MemberNode)
	return result
}
func CreatePerceptionDNETNoP(rootNode *DNode) PerceptionDNET {
	var result PerceptionDNET
	result = PerceptionDNET{}
	//设置根节点
	result.RootNode = rootNode
	//所属域
	result.Domains = &structure.ArrayList{}
	//初始成员map
	result.MemberMap = make(map[string]*MemberGroup)
	//初始化成员节点+锁
	//result.MemberNodeAndLock = make(map[string]*MemberNode)
	return result
}

//--------------------------------基本操作---------------------------------------
//更新节点
func (this *PerceptionDNET) UpdateRootNode(name string, ip string, port string) {
	this.RWLock.Lock()
	this.RWLock.Unlock()

	if len(name) > 0 {
		this.RootNode.SetName(name)
	}

	changeAddr := false
	if len(ip) > 0 {
		this.RootNode.SetIp(ip)
		changeAddr = true
	}
	if len(port) > 0 {
		this.RootNode.SetPort(port)
		changeAddr = true
	}
	if changeAddr {
		this.RootNode.SetAddr(this.RootNode.GetIp() + ":" + this.RootNode.GetPort())
	}
}

/**
增加直联节点
注意：此操作会检查是否添加过该节点
*/
func (this *PerceptionDNET) AddDirectConnectionNode(cNode *DNode) {
	//节点地址 IP:PORT -- 唯一信息
	addr := cNode.GetAddr()
	//无该直系节点 -- 添加创建
	if nil == this.MemberMap[addr] {
		//成员map
		this.MemberMap[addr] = CreateMemberGroup(cNode)
		//成员节点+锁
		//this.MemberNodeAndLock[addr] = CreateMemberNode(cNode)
	}
}

/**
增加直联节点(同步)
注意：此操作会检查是否添加过该节点
*/
func (this *PerceptionDNET) AddDirectConnectionNodeSync(cNode *DNode) {
	this.RWLock.Lock()
	defer this.RWLock.Unlock()
	//节点地址 IP:PORT -- 唯一信息
	addr := cNode.GetAddr()
	//无该直系节点 -- 添加创建
	if nil == this.MemberMap[addr] {
		//成员map
		this.MemberMap[addr] = CreateMemberGroup(cNode)
		//成员节点+锁
		//this.MemberNodeAndLock[addr] = CreateMemberNode(cNode)
	}
}

/**
是否存在该直系节点
*/
func (this *PerceptionDNET) ExistDirectConnectionNodeSync(addr string) bool {
	this.RWLock.Lock()
	defer this.RWLock.Unlock()
	return nil != this.MemberMap[addr]
}

//--------------------------------获取---------------------------------------
//线程安全:获取成员leader（直系节点） Addr
func (this *PerceptionDNET) MemberLeadersSync() structure.ArrayList {
	this.RWLock.Lock()
	defer this.RWLock.Unlock()
	leaders := this.MemberLeaders()
	return leaders
}

//线程不安全:获取成员leader（直系节点） Addr
func (this *PerceptionDNET) MemberLeaders() structure.ArrayList {
	set := structure.ArrayList{}
	for k, _ := range this.MemberMap {
		set.Add(k)
	}
	return set
}

//--------------------------------合并操作---------------------------------------
/**
合并感知节点
*/
func (this *PerceptionDNET) MergePerceptionDNETSync(come *PerceptionDNET) {
	if nil == come || nil == come.RootNode {
		return
	}

	thisRoot := this.RootNode.GetAddr()
	comeRoot := come.RootNode.GetAddr()

	//自我合并 or 外来合并
	if thisRoot == comeRoot {
		this.mergePerceptionMySelf(come)
	} else {
		this.mergePerceptionOthers(come)
	}

	/**
	自我去重 -- 优化 TODO

	1 A->B->C D
	2       D
	3       C D E F
	4       D G F
	5    E->F

	以如上关系而言
	2 和 4 表达的含义是一致的，因此可以省略 2 保留 4
	同理1 和 3
	*/

}

/**
自我合并
比较root 节点的更新时间进行合并 以便后期拓展持久化、中心化
*/
func (this *PerceptionDNET) mergePerceptionMySelf(come *PerceptionDNET) {
	thisUnix := this.RootNode.GetUpTime().Unix()
	comeUnix := come.RootNode.GetUpTime().Unix()
	//如果时间更新，说明是后变更，则直接整体替换现有感知网络
	//后期升级可以做成员合并，并记录合并记录，持久化或中心+持久化
	//自我合并不需要考虑时间差值问题，因为是自身的修改，机器不变、时钟不变
	if comeUnix >= thisUnix {
		this = come
	}
}

/**
外部合并
*/
func (this *PerceptionDNET) mergePerceptionOthers(come *PerceptionDNET) {
	//获取本agent 直系节点成员set集合
	thisleaderList := this.MemberLeaders()
	comeRootaddr := come.RootNode.GetAddr()
	//查看对方是否已备自己添加在直属节点中
	_, contains := thisleaderList.Contains(comeRootaddr)

	//1 如果没有对方直系节点 -- 拿来 --> 去除自己
	//2 如果存在对方直系节点 -- 更新
	if !contains {
		//1.1  添加对方Root到直系节点中
		now := time.Now()
		addNode := CreateDNodeByParams(come.RootNode.GetName(), come.RootNode.GetIp(), come.RootNode.GetPort(), come.RootNode.GetAddr(), 0, now, now)
		addNode.SetStateActive()
		addNode.SetStateArrive()
		this.AddDirectConnectionNode(addNode)
	}

	//1.2 查看对方直属成员 -- 排除自己(delete) --because直属节点信息一定来自自己
	leaders := come.MemberLeaders()
	i := 0
	for i = 0; i < leaders.Size(); i++ {
		comeLeaderAddrStr := leaders.Get(i)
		//如果自己存在于对方的直属节点，则将自己剔除,方便后面将对方降级处理
		if this.RootNode.GetAddr() == comeLeaderAddrStr {
			delete(come.MemberMap, this.RootNode.GetAddr())
		}
	}
	//1.3 查看对方直属节点排除自己后 其成员组中 是否存在自己
	//1.3.1 将对方降级处理
	comeMemberCollection := come.ConvertToMemberGroup()
	//1.3.2 查看成员集合中是否存在自己
	comeGroupFindIndexArr := comeMemberCollection.FindNodeFromMemberGroup(this.RootNode)

	//该成员组中 发现 自己
	if !comeGroupFindIndexArr.IsEmpty() {
		for i := 0; i < comeGroupFindIndexArr.Size(); i++ {
			tempForFoundMe := comeGroupFindIndexArr.Get(i).([2]int)
			//我所在的member group index
			groupMembersIndex := tempForFoundMe[0]
			//我所在的 member collection index
			memberCollectionIndex := tempForFoundMe[1]
			//come 成员集合
			tempForMemberCollection := comeMemberCollection.GetMemberCollectionByIndex(groupMembersIndex)
			//判断 我 在对方成员集合中的位置
			//
			//计算我自己
			if memberCollectionIndex > 0 {
				beforeMe := tempForMemberCollection.Get(memberCollectionIndex - 1)
				_, b := thisleaderList.Contains(beforeMe.GetAddr())
				//我的直属节点未找到，说明可加入到直属节点中
				if !b {
					this.AddDirectConnectionNode(beforeMe)
				}
			}
		}
	} else {
		comeMemberGroup := comeMemberCollection.MemberGroup
		for i := 0; i < comeMemberGroup.Size(); i++ {
			comeMemberCollection := comeMemberGroup.Get(i).(*MemberCollection)
			this.MemberMap[comeRootaddr].addMembers(comeMemberCollection.ToDNodeArr())
		}
		comeMemberCollection.MemberGroup.Size()
	}
}

//--------------------------------转化操作---------------------------------------
/**
将 PerceptionDNET 降级转换为 MemberCollection
*/
func (this *PerceptionDNET) ConvertToMemberGroup() *MemberGroup {

	copyDNode := CreateDNodeByParamsWithOutTime(this.RootNode.GetName(), this.RootNode.GetIp(), this.RootNode.GetPort(), this.RootNode.GetAddr(), 0)
	//1 创建 MemberCollection 并设置 组长
	resultMemberGroup := CreateMemberGroup(copyDNode)
	//2 向成员组转化
	for _, v := range this.MemberMap {
		collection := CreateMemberCollection()
		resultMemberGroup.MemberGroup.Add(collection)
		collection.MemberElements.Add(v.MemberLeader)
		for i := 0; i < v.MemberGroup.Size(); i++ {
			memberCollection := v.MemberGroup.Get(i).(*MemberCollection)
			for j := 0; j < memberCollection.MemberElements.Size(); j++ {
				if nil != memberCollection.MemberElements.Get(j) {
					collection.MemberElements.Add(memberCollection.MemberElements.Get(j))
				}
			}
		}
	}

	return resultMemberGroup
}

//--------------------------------添加操作---------------------------------------
/**
添加依赖链
创建MemberGroup
	填充MemberLeader
	填充MemberCollection
*/
func (this *PerceptionDNET) AddDNodesArrayListForDNet(nodes *structure.ArrayList) {
	if nil == nodes || nodes.IsEmpty() {
		return
	}
	firstNode := nodes.Get(0).(*DNode)

	this.AddDirectConnectionNode(firstNode)

	size := nodes.Size()
	//第一个元素作为直属节点，成员从第二节点开始算起
	dNodes := make([]*DNode, size-1)
	for i := 1; i < size; i++ {
		node := nodes.Get(i).(*DNode)
		dNodes[i-1] = node
	}
	this.MemberMap[firstNode.GetAddr()].addMembers(dNodes)
}

//--------------------------------打印操作---------------------------------------
func (this *PerceptionDNET) Print() {
	rootLine := true
	memberleaderLine := true

	split := "->"
	rootLen := len(this.RootNode.GetAddr()) + len(split)
	memberGroupLeader := 0
	for _, v := range this.MemberMap {
		memberGroupleader := v.MemberLeader
		memberGroup := v.MemberGroup

		//打印 root节点
		if rootLine {
			fmt.Print(this.RootNode.GetAddr(), split)
			rootLine = false
		} else {
			for n := 0; n < rootLen; n++ {
				fmt.Print(" ")
			}
		}

		//打印memberGroupLeader
		memberleaderLine = true
		fmt.Print(memberGroupleader.GetAddr(), split)
		memberGroupLeader = len(memberGroupleader.GetAddr()) + len(split)

		if nil != memberGroup && !memberGroup.IsEmpty() {
			//MemberGroup
			for i := 0; i < memberGroup.Size(); i++ {
				if nil == memberGroup.Get(i) {
					continue
				}
				if memberleaderLine {
					memberleaderLine = false
				} else {
					for n := 0; n < memberGroupLeader+rootLen; n++ {
						fmt.Print(" ")
					}
				}

				memberCollection := memberGroup.Get(i).(*MemberCollection)
				//MemberCollection
				for j := 0; j < memberCollection.MemberElements.Size(); j++ {
					dNode := memberCollection.Get(j)
					if nil != dNode {
						fmt.Print(dNode.GetAddr(), " ")
					}
				}
				fmt.Println()
			}
		} else {
			fmt.Println()
		}
		//fmt.Println()
	}
}

//--------------------------------节点查找---------------------------------------
//TODO
func (this *PerceptionDNET) FindWayToDNode(addr string) {

}
