package watch

import (
	"fmt"
	MyCommon "myagent/src/common"
	)

type TopicConsumer struct {
	CanComsumFlag chan bool
	topicQueue *MyCommon.Queue
}

func (this *TopicConsumer) ConsumeTopics() {
	for ;;{
		select {
		case val :=<-this.CanComsumFlag:
			if val == true {
				fmt.Println(val)
			}
		}
	}
}