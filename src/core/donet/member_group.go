package donet

import (
	"myagent/src/core/structure"
)

/**
MemberGroup -- 成员集合
{
	memberLeader:*DNode	-- 组长
	memberGroup:*structure.ArrayList -- 组员（集合）
			|_____{A,B,C,D}
			|_____{A,D}
			|_____{C}
}
*/
type MemberGroup struct {
	//组长
	MemberLeader *DNode
	//成员组 -- 成员是一个个集合 结合存放的是 *MemberCollection
	MemberGroup *structure.ArrayList
}

func (this *MemberGroup) SetMemberGroup(mp *structure.ArrayList) {
	this.MemberGroup = mp
}

//--------------------------------创建 MemberCollection--------------------------------
//创建 成员集合
func (this *MemberGroup) createMemberGroup() *MemberGroup {
	var result *MemberGroup
	result = &MemberGroup{}
	return result
}

//创建 对外调用
func CreateMemberGroup(memberLeader *DNode) *MemberGroup {
	var result *MemberGroup
	collection := result.createMemberGroup()
	if nil != memberLeader {
		//1 设置组长
		collection.MemberLeader = memberLeader
		//2 初始化成员组
		collection.MemberGroup = &structure.ArrayList{}
	}
	return collection
}

//--------------------------------操作 MemberCollection--------------------------------
//创建成员集合
func (this *MemberGroup) createMemberArr() *structure.ArrayList {
	result := &structure.ArrayList{}
	return result
}

//批量加入member成员 （MemberGroup.MemberLeader 不会发生改变）
/*
func (this *MemberGroup) addMembers(dNode ...*DNode)  {
	//当dNode 长度为1 --
	if len(dNode) <= 0 {
		return
	}

	//如果长度为0 -- 判断所有 成员集合的第一个节点是否包含
	if len(dNode) == 1 {
		group := this.MemberGroup
		//是否需要添加
		needAdd := true
		for i:=0;i<group.Size();i++ {
			memberList := group.Get(i).(*MemberCollection)
			if nil != memberList && !memberList.MemberElements.IsEmpty() {
				has, hasIndex := memberList.ContainsAllArr(dNode[0])
				if has && hasIndex == 0 {
					needAdd = false
					break
				}
			}
		}

		if needAdd {
			collection := CreateMemberCollection()
			collection.AddMember(dNode[0])
			this.MemberGroup.Add(collection)
		}
	} else {
	//如果长度 > 0 则
	//1. 如果所有成员中都不存在 -- 新创建并增加
	//2. 如果存在判断位置
		needAdd := true
		group := this.MemberGroup
		for i:=0;i<group.Size();i++ {
			memberList := group.Get(i).(*MemberCollection)
			if nil != memberList && !memberList.MemberElements.IsEmpty() {
				has, hasIndex := memberList.ContainsAllArr(dNode...)
				if has && hasIndex == 0 {
					needAdd = false
					break
				}
			}
		}

		if needAdd {
			collection := CreateMemberCollection()
			collection.MemberElements.AppendTo(0,dNode)
			this.MemberGroup.Add(collection)
		}
	}
}
*/
func (this *MemberGroup) addMembers(dNode []*DNode) {
	//当dNode 长度为1 --
	if len(dNode) <= 0 {
		return
	}

	//如果长度为0 -- 判断所有 成员集合的第一个节点是否包含
	if len(dNode) == 1 {
		group := this.MemberGroup
		//是否需要添加
		needAdd := true
		for i := 0; i < group.Size(); i++ {
			memberList := group.Get(i).(*MemberCollection)
			if nil != memberList && !memberList.MemberElements.IsEmpty() {
				has, hasIndex := memberList.ContainsAllArr(dNode[0])
				if has && hasIndex == 0 {
					needAdd = false
					break
				}
			}
		}

		if needAdd {
			collection := CreateMemberCollection()
			collection.AddMember(dNode[0])
			this.MemberGroup.Add(collection)
		}
	} else {
		//如果长度 > 0 则
		//1. 如果所有成员中都不存在 -- 新创建并增加
		//2. 如果存在判断位置
		needAdd := true
		group := this.MemberGroup
		for i := 0; i < group.Size(); i++ {
			memberList := group.Get(i).(*MemberCollection)
			if nil != memberList && !memberList.MemberElements.IsEmpty() {
				has, hasIndex := memberList.ContainsAllArr(dNode...)
				if has && hasIndex == 0 {
					needAdd = false
					break
				}
			}
		}

		if needAdd {
			collection := CreateMemberCollection()
			if nil != dNode {
				for i := 0; i < len(dNode); i++ {
					collection.AddMember(dNode[i])
				}
			}
			this.MemberGroup.Add(collection)
		}
	}
}

/**
找出所需节点
*/
func (this *MemberGroup) FindNodeFromMemberGroup(dNode *DNode) structure.ArrayList {
	result := structure.ArrayList{}
	group := this.MemberGroup
	for i := 0; i < group.Size(); i++ {
		memberList := group.Get(i).(*MemberCollection)
		if nil != memberList && !memberList.MemberElements.IsEmpty() {
			arr, i2 := memberList.ContainsAllArr(dNode)
			if arr {
				var indexArr [2]int
				indexArr[0] = i
				indexArr[1] = i2
				result.Add(indexArr)
			}
		}
	}
	return result
}

/**
通过 MemberGroup 的 index 获取 成员集合
*/
func (this *MemberGroup) GetMemberCollectionByIndex(index int) *MemberCollection {
	return this.MemberGroup.Get(index).(*MemberCollection)
}

/**
将MemberGroup 降级处理，转换成 MemberCollection 集合
*/
func (this *MemberGroup) ConvertMemberGroupToMemberCollectionArrayLists() *structure.ArrayList {
	result := &structure.ArrayList{}
	for i := 0; i < this.MemberGroup.Size(); i++ {
		get := this.MemberGroup.Get(i).(*MemberCollection)
		newMemberCollection := CreateMemberCollection()
		newMemberCollection.AddMember(this.MemberLeader)
		for j := 0; j < get.MemberElements.Size(); j++ {
			newMemberCollection.AddMember(get.Get(j))
		}
		result.Add(newMemberCollection)
	}
	return result
}
