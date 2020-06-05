package perception

import (
	"myagent/src/core/structure"
	"time"
)

//感知节点
type PerceptionNode struct {
	//节点名称
	Name string
	IP   string
	PORT string
	//IP:PORT
	Addr string
	//是否活跃
	Active bool
	//最后一次活跃时间
	LastActiveTime time.Time
}

type PerceptionAgent struct {
	//当前Agent
	MySelf            *PerceptionNode

	//计算后的关系 ArrayList<ArrayList<PerceptionNode>
	RelationshipLines *structure.ArrayList

	
}
