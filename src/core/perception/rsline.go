package perception

import (
	"myagent/src/core/structure"
)

type RsLine struct {
	/**key节点 OR 上游节点
	(A)
	*/
	keyNode *PNode

	/** 维护关系
	--------------------------------------------------------------------------------
	| *structure.ArrayList<*structure.ArrayList>	|  *structure.ArrayList<*PNode>	|
	--------------------------------------------------------------------------------
	| 第一纵列		 								|		第二纵列					|
	| [B.ADDR]										|       (F)						|
	|												|       (C,D,E)					|
	|    											|       (D,E)					|
	|												|       (E)						|
	--------------------------------------------------------------------------------
	| [E.ADDR]    									|		()						|
	| [F.ADDR]    									|		()						|
	--------------------------------------------------------------------------------

	*/
	linealMap map[string]*structure.ArrayList
}

//---------------------------------构造创建---------------------------------------
//创建 RsLine
func CreateRsLine() *RsLine {
	var result *RsLine
	result = &RsLine{}
	result.keyNode = CreatePNode()
	result.linealMap = CreateLinealMap()
	return result
}

//创建 RsLine.linealMap
func CreateLinealMap() map[string]*structure.ArrayList {
	return make(map[string]*structure.ArrayList)
}

//---------------------------------基本 GetSet---------------------------------------
//获取 RsLine.keyNode
func (this *RsLine) GetKeyNode() *PNode {
	return this.keyNode
}

//设置 RsLine.keyNode
func (this *RsLine) SetKeyNode(n *PNode) {
	this.keyNode = n
}

//获取 RsLine.linealMap
func (this *RsLine) GetLinealMap() map[string]*structure.ArrayList {
	return this.linealMap
}

//设置 RsLine.linealMap
func (this *RsLine) SetLinealMap(lp map[string]*structure.ArrayList) {
	this.linealMap = lp
}

//---------------------------------基本 linealMap 操作---------------------------------------
//linealMap 中是否包含 节点  -- true:存在/false:不存在
func (this *RsLine) ExistFromLinealMap(n *PNode) bool {
	return this.linealMap[n.GetAddr()] != nil
}

//linealMap 中是否包含 节点  -- true:存在/false:不存在
func (this *RsLine) ExistFromLinealMapByAddr(addr string) bool {
	return this.linealMap[addr] != nil
}


//---------------------------------复合 linealMap 操作---------------------------------------
//通过 addrKey + 第一纵列index  获取  实际元素ArrayList 指针
func (this *RsLine) FindSecondArrayListFromLinealMap(keyAddr string,index int) *structure.ArrayList {
	if nil == this.linealMap[keyAddr] {
		return nil
	}
	return this.linealMap[keyAddr].Get(index).(*structure.ArrayList)
}
/**
	从 RsLine.linealMap 查询是否存在该节点
	返回结果：structure.ArrayList<*FindNodeIndex>
*/
func (this *RsLine) FindNodeIndexsFromLinealMap(addr string) *structure.ArrayList {
	keyArr := &structure.ArrayList{}
	//0 不做处理 , 1 精确查询myself ,
	for kv,_:=range this.GetLinealMap() {
		if len(kv) > 0 {
			keyArr.Add(kv)
		}
	}

	result := &structure.ArrayList{}

	for ki:=0;ki<keyArr.Size();ki++ {
		//
		key := keyArr.Get(ki).(string)
		//第一纵列
		nodeArrsList := this.linealMap[key]
		if nil != nodeArrsList && !nodeArrsList.IsEmpty() {
			for ni:=0;ni<nodeArrsList.Size();ni++ {
				//实际节点列
				nodeArr := nodeArrsList.Get(ni).(*structure.ArrayList)
				if nil != nodeArr && !nodeArr.IsEmpty() {
					for sni:=0;sni<nodeArr.Size();sni++ {
						node := nodeArr.Get(sni).(*PNode)
						if node.GetAddr() == addr {
							result.Add(CreateFindNodeIndex(key,ni,sni,node))
						}
					}
				}

			}
		}
	}
	return result
}


//---------------------------------FindNodeIndex -- RsLine.linealMap 结果struct---------------------------------------
//查询 RsLine.linealMap 的结果
type FindNodeIndex struct {
	rkey string
	columIndex int
	rowIndex int
	node *PNode
}
func CreateFindNodeIndex(rkey string,columIndex int,rowIndex int,node *PNode) *FindNodeIndex{
	result := &FindNodeIndex{
		rkey, columIndex, rowIndex, node,
	}
	return result
}
func (this *FindNodeIndex) GetRKey() string {
	return this.rkey
}
func (this *FindNodeIndex) GetColumIndex() int {
	return this.columIndex
}
func (this *FindNodeIndex) GetRowIndex() int {
	return this.rowIndex
}
func (this *FindNodeIndex) GetPNode() *PNode {
	return this.node
}
//---------------------------------FindNodeIndex -- RsLine.linealMap 结果struct---------------------------------------