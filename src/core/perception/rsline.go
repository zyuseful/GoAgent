package perception

import (
	"myagent/src/core/structure"
)

type RsLine struct {
	LinealNode *PNode
	LinealMap map[string]*structure.ArrayList
}


func CreateRsLine() *RsLine {
	var result *RsLine
	result = &RsLine{}
	result.LinealNode = CreatePNode()
	result.LinealMap = CreateRsLine_LinealMap()
	return result
}

func CreateRsLine_LinealMap() map[string]*structure.ArrayList{
	//return new(map[string]*structure.ArrayList)
	return make(map[string]*structure.ArrayList)
}

func (this *RsLine) SetRsLine_LinealNode(name string,ip string,port string)  {
	this.LinealNode.SetPNode(name,ip,port)
}

