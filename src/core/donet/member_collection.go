package donet

import (
	"myagent/src/core/structure"
	"strings"
)

type MemberCollection struct {
	//存放 *DNode 元素
	MemberElements *structure.ArrayList
}

func CreateMemberCollection() *MemberCollection {
	result := &MemberCollection{}
	result.MemberElements = &structure.ArrayList{}
	return result
}

func (this *MemberCollection) AddMember(dNode *DNode) {
	contains := this.IfContains(dNode)
	if !contains {
		this.MemberElements.Add(dNode)
	}
}

//查看是否包含相同addr 的 DNode
func (this *MemberCollection) IfContains(dNode *DNode) bool {
	result := false
	var node *DNode
	for i := 0; i < this.MemberElements.Size(); i++ {
		node = this.MemberElements.Get(i).(*DNode)
		if strings.Compare(node.GetAddr(), dNode.GetAddr()) == 0 {
			result = true
			break
		}
	}
	return result
}

//查看是否包含相同addr 的 DNode -- 只要存在一个就返回
func (this *MemberCollection) IfContainsArr(dNode ...*DNode) bool {
	result := false
	var node *DNode
	var i, j int

	for i = 0; i < len(dNode); i++ {
		for j = i + 1; j < len(dNode); j++ {
			if dNode[i].GetAddr() == dNode[j].GetAddr() {
				result = true
				return result
			}
		}
	}

	for i = 0; i < this.MemberElements.Size(); i++ {
		for j = 0; j < len(dNode); j++ {
			node = this.MemberElements.Get(i).(*DNode)
			if node.GetAddr() == dNode[j].GetAddr() {
				result = true
				break
			}
		}
	}
	return result
}

/**
查看是否包含所有 的 DNode -- 只要存在一个就返回

*/
func (this *MemberCollection) ContainsAllArr(dNode ...*DNode) (bool, int) {
	var node *DNode
	var i, j, first int

	if nil == dNode || len(dNode) <= 0 {
		return false, -1
	}

	//如果输入长度大于成员长度 false
	if this.MemberElements.Size() < len(dNode) {
		return false, -1
	}

	//首先找到第一个元素
	first = -1
	for i = 0; i < this.MemberElements.Size(); i++ {
		node = this.Get(i)
		if node.GetAddr() == dNode[0].GetAddr() {
			first = i
			break
		} else {
			continue
		}
	}

	if first == -1 {
		return false, -1
	}

	//一对一比较
	for i = 0; i < len(dNode); i++ {
		if this.MemberElements.Size() <= i+first {
			if nil != this.MemberElements.Get(first+i) {
				node = this.MemberElements.Get(first + i).(*DNode)
				if dNode[i].GetAddr() == node.GetAddr() {
					j++
				}
			}
		}
	}

	return len(dNode) == j, first
}

func (this *MemberCollection) Get(index int) *DNode {
	if this.MemberElements.IsEmpty() {
		return nil
	}

	return this.MemberElements.Get(index).(*DNode)
}

func (this *MemberCollection) ToDNodeArr() []*DNode {
	nodes := make([]*DNode, this.MemberElements.Size())
	for i := 0; i < this.MemberElements.Size(); i++ {
		nodes[i] = this.Get(i)
	}
	return nodes
}
