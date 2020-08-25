package donet

import "sync"

type MemberNode struct {
	MemberNode *DNode
	RWLock     sync.RWMutex
}

func CreateMemberNode(memberNode *DNode) *MemberNode {
	node := &MemberNode{}
	node.MemberNode = memberNode
	return node
}
