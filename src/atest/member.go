package main

import (
	"fmt"
	"myagent/src/core/donet"
	"myagent/src/core/structure"
)

func mains1() {
	collection := donet.CreateMemberCollection()
	nodeA := donet.CreateDNodeByParamsWithOutTime("A", "127.0.0.1", "5555", "127.0.0.1:5555", 0)
	nodeB := donet.CreateDNodeByParamsWithOutTime("B", "127.0.0.2", "5555", "127.0.0.2:5555", 0)
	nodeC := donet.CreateDNodeByParamsWithOutTime("C", "127.0.0.3", "5555", "127.0.0.3:5555", 0)
	nodeD := donet.CreateDNodeByParamsWithOutTime("D", "127.0.0.4", "5555", "127.0.0.4:5555", 0)

	collection.AddMember(nodeA)
	collection.AddMember(nodeB)
	collection.AddMember(nodeC)
	collection.AddMember(nodeD)

	fmt.Println(collection.ContainsAllArr(nodeB, nodeC))

	//fmt.Println(collection.Contains(nodeA))
}
func mains0() {
	node := donet.CreateDNodeByParamsWithOutTime("A", "127.0.0.1", "5555", "127.0.0.1:5555", 0)
	dnet := donet.CreatePerceptionDNET(node)

	nodeB := donet.CreateDNodeByParamsWithOutTime("B", "127.0.0.1", "5555", "127.0.0.1:5555", 0)
	memberCollection1 := donet.CreateMemberGroup(nodeB)
	group1 := &structure.ArrayList{}
	l1 := &structure.ArrayList{}
	l1.Add("D")
	l1.Add("E")
	l1.Add("F")
	l1.Add("G")
	group1.Add(l1)
	l2 := &structure.ArrayList{}
	l2.Add("D")
	l2.Add("G")
	group1.Add(l2)
	memberCollection1.SetMemberGroup(group1)

	nodeC := donet.CreateDNodeByParamsWithOutTime("C", "127.0.0.1", "5555", "127.0.0.1:5555", 0)
	memberCollection2 := donet.CreateMemberGroup(nodeC)
	group2 := &structure.ArrayList{}
	l21 := &structure.ArrayList{}
	l21.Add(donet.CreateDNodeByParamsWithOutTime("D", "127.0.0.1", "5555", "127.0.0.1:5555", 0))
	l21.Add(donet.CreateDNodeByParamsWithOutTime("E", "127.0.0.1", "5555", "127.0.0.1:5555", 0))
	l21.Add(donet.CreateDNodeByParamsWithOutTime("F", "127.0.0.1", "5555", "127.0.0.1:5555", 0))
	l21.Add(donet.CreateDNodeByParamsWithOutTime("G", "127.0.0.1", "5555", "127.0.0.1:5555", 0))
	group2.Add(l21)
	l22 := &structure.ArrayList{}
	l22.Add(donet.CreateDNodeByParamsWithOutTime("D", "127.0.0.1", "5555", "127.0.0.1:5555", 0))
	l22.Add(donet.CreateDNodeByParamsWithOutTime("G", "127.0.0.1", "5555", "127.0.0.1:5555", 0))
	group2.Add(l22)
	memberCollection2.SetMemberGroup(group2)

	dnet.MemberMap["B"] = memberCollection1
	dnet.MemberMap["C"] = memberCollection2

	fmt.Println("--------创建完成--------")
	fmt.Println(dnet.RootNode)
	//测试获取直属节点
	fmt.Println(dnet.MemberLeaders())
	collection := dnet.ConvertToMemberGroup()
	fmt.Println("OK1")
	fmt.Println(collection)
	fmt.Println("OK2")

}

func main() {
	//	MySelf
	//	self := CreateMySelf1()
	self := CreateMySelf2()
	self.Print()
	//	Come
	//come := CreateCome1()
	fmt.Println("-----------------")
	come := CreateCome3()
	come.Print()
	fmt.Println("-----------------")
	self.MergePerceptionDNETSync(come)
	self.Print()
}

func CreateMySelf1() *donet.PerceptionDNET {
	root := donet.CreateDNodeByParamsWithOutTime("A", "127.0.0.1", "5555", "A", 0)
	result := donet.CreatePerceptionDNET(root)

	//B C
	mySelfDNodeList1 := &structure.ArrayList{}
	nodeB := donet.CreateDNodeByParamsWithOutTime("B", "127.0.0.1", "5555", "B", 0)
	nodeC := donet.CreateDNodeByParamsWithOutTime("C", "127.0.0.1", "5555", "C", 0)
	mySelfDNodeList1.Add(nodeB)
	mySelfDNodeList1.Add(nodeC)
	result.AddDNodesArrayListForDNet(mySelfDNodeList1)

	//B D
	mySelfDNodeList2 := &structure.ArrayList{}
	nodeD := donet.CreateDNodeByParamsWithOutTime("D", "127.0.0.1", "5555", "D", 0)
	mySelfDNodeList2.Add(nodeB)
	mySelfDNodeList2.Add(nodeD)
	result.AddDNodesArrayListForDNet(mySelfDNodeList2)

	return result
}
func CreateMySelf2() *donet.PerceptionDNET {
	root := donet.CreateDNodeByParamsWithOutTime("A", "127.0.0.1", "5555", "A", 0)
	result := donet.CreatePerceptionDNET(root)

	//B C D
	mySelfDNodeList1 := &structure.ArrayList{}
	nodeB := donet.CreateDNodeByParamsWithOutTime("B", "127.0.0.1", "5555", "B", 0)
	nodeC := donet.CreateDNodeByParamsWithOutTime("C", "127.0.0.1", "5555", "C", 0)
	nodeD := donet.CreateDNodeByParamsWithOutTime("D", "127.0.0.1", "5555", "D", 0)
	mySelfDNodeList1.Add(nodeB)
	mySelfDNodeList1.Add(nodeC)
	mySelfDNodeList1.Add(nodeD)
	result.AddDNodesArrayListForDNet(mySelfDNodeList1)

	//B  D
	mySelfDNodeList2 := &structure.ArrayList{}
	mySelfDNodeList2.Add(nodeB)
	mySelfDNodeList2.Add(nodeD)
	result.AddDNodesArrayListForDNet(mySelfDNodeList2)

	//E F
	mySelfDNodeList3 := &structure.ArrayList{}
	nodeE := donet.CreateDNodeByParamsWithOutTime("E", "127.0.0.1", "5555", "E", 0)
	nodeF := donet.CreateDNodeByParamsWithOutTime("F", "127.0.0.1", "5555", "F", 0)
	mySelfDNodeList3.Add(nodeE)
	mySelfDNodeList3.Add(nodeF)

	result.AddDNodesArrayListForDNet(mySelfDNodeList3)

	return result
}

func CreateCome1() *donet.PerceptionDNET {
	root := donet.CreateDNodeByParamsWithOutTime("B", "127.0.0.1", "5555", "B", 0)
	result := donet.CreatePerceptionDNET(root)

	//B C
	mySelfDNodeList1 := &structure.ArrayList{}
	nodeC := donet.CreateDNodeByParamsWithOutTime("C", "127.0.0.1", "5555", "C", 0)
	mySelfDNodeList1.Add(nodeC)
	result.AddDNodesArrayListForDNet(mySelfDNodeList1)

	//B E
	mySelfDNodeList2 := &structure.ArrayList{}
	nodeD := donet.CreateDNodeByParamsWithOutTime("E", "127.0.0.1", "5555", "E", 0)
	mySelfDNodeList2.Add(nodeD)
	result.AddDNodesArrayListForDNet(mySelfDNodeList2)

	return result
}
func CreateCome2() *donet.PerceptionDNET {
	root := donet.CreateDNodeByParamsWithOutTime("B", "127.0.0.1", "5555", "B", 0)
	result := donet.CreatePerceptionDNET(root)

	//B C
	mySelfDNodeList1 := &structure.ArrayList{}
	nodeC := donet.CreateDNodeByParamsWithOutTime("C", "127.0.0.1", "5555", "C", 0)
	mySelfDNodeList1.Add(nodeC)
	result.AddDNodesArrayListForDNet(mySelfDNodeList1)

	//B E
	mySelfDNodeList2 := &structure.ArrayList{}
	nodeD := donet.CreateDNodeByParamsWithOutTime("E", "127.0.0.1", "5555", "E", 0)
	mySelfDNodeList2.Add(nodeD)
	result.AddDNodesArrayListForDNet(mySelfDNodeList2)

	//B C D
	mySelfDNodeList3 := &structure.ArrayList{}
	nodeA := donet.CreateDNodeByParamsWithOutTime("A", "127.0.0.1", "5555", "A", 0)
	nodeE := donet.CreateDNodeByParamsWithOutTime("E", "127.0.0.1", "5555", "E", 0)
	nodeF := donet.CreateDNodeByParamsWithOutTime("F", "127.0.0.1", "5555", "F", 0)
	mySelfDNodeList3.Add(nodeA)
	mySelfDNodeList3.Add(nodeE)
	mySelfDNodeList3.Add(nodeF)
	result.AddDNodesArrayListForDNet(mySelfDNodeList3)

	return result
}
func CreateCome3() *donet.PerceptionDNET {
	nodeA := donet.CreateDNodeByParamsWithOutTime("A", "127.0.0.1", "5555", "A", 0)
	root := donet.CreateDNodeByParamsWithOutTime("B", "127.0.0.1", "5555", "B", 0)
	nodeC := donet.CreateDNodeByParamsWithOutTime("C", "127.0.0.1", "5555", "C", 0)
	nodeD := donet.CreateDNodeByParamsWithOutTime("D", "127.0.0.1", "5555", "D", 0)
	nodeE := donet.CreateDNodeByParamsWithOutTime("E", "127.0.0.1", "5555", "E", 0)
	nodeF := donet.CreateDNodeByParamsWithOutTime("F", "127.0.0.1", "5555", "F", 0)
	nodeG := donet.CreateDNodeByParamsWithOutTime("G", "127.0.0.1", "5555", "G", 0)

	result := donet.CreatePerceptionDNET(root)

	mySelfDNodeList1 := &structure.ArrayList{}
	mySelfDNodeList1.Add(nodeC)
	mySelfDNodeList1.Add(nodeD)
	mySelfDNodeList1.Add(nodeE)
	mySelfDNodeList1.Add(nodeF)
	result.AddDNodesArrayListForDNet(mySelfDNodeList1)

	mySelfDNodeList2 := &structure.ArrayList{}
	mySelfDNodeList2.Add(nodeD)
	mySelfDNodeList2.Add(nodeG)
	result.AddDNodesArrayListForDNet(mySelfDNodeList2)

	mySelfDNodeList3 := &structure.ArrayList{}
	mySelfDNodeList3.Add(nodeA)
	mySelfDNodeList3.Add(nodeE)
	mySelfDNodeList3.Add(nodeF)
	result.AddDNodesArrayListForDNet(mySelfDNodeList3)

	mySelfDNodeList4 := &structure.ArrayList{}
	mySelfDNodeList4.Add(nodeD)
	mySelfDNodeList4.Add(nodeG)
	mySelfDNodeList4.Add(nodeF)
	result.AddDNodesArrayListForDNet(mySelfDNodeList4)

	return result
}
