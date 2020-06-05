package watch

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	MyTaskPools "github.com/wazsmwazsm/mortar"
	"myagent/src/core/structure"
	"reflect"
	"sync"
)

type (
	//Watch 工厂
	MyWatchFactory struct {
		//存储 topic--{value--observer1}
		// 				    |_observer2
		// 				    |_observer3
		saveTopicAndValueAndObserver map[string]*ValueAndObservers
		//主题锁 -- 用于控制主题操作:创建主题、移除主题
		TopicLock sync.Mutex

		//任务信号队列
		TaskChan chan ValAndObservers
		//TaskChan chan interface{}

		//任务管理
		taskPool  *MyTaskPools.Pool
		waitGroup *sync.WaitGroup
	}

	//观察值 + 被观察者们
	ValueAndObservers struct {
		//锁
		Lock sync.Mutex
		//观察值
		Value interface{}
		//被观察代理队列<IObserver>
		Observers *structure.ArrayList
	}

	//主题、值（拷贝传递)、被观察者
	ValAndObservers struct {
		Value     interface{}
		//被观察代理队列<IObserver>
		Observers *structure.ArrayList
	}
)

//实现 初始化
func (this *MyWatchFactory) Init(num uint64) {
	//初始化组件 : 核心存储初始化 + 任务信号队列
	this.initWatchComponents()
	this.initTaskPoolComponents(num)
}

//初始化组件 : 核心存储初始化 + 任务信号队列
func (this *MyWatchFactory) initWatchComponents() {
	//核心存储初始化
	this.saveTopicAndValueAndObserver = make(map[string]*ValueAndObservers)
	//任务信号队列
	this.TaskChan = make(chan ValAndObservers, 4096)
	//this.TaskChan = make(chan interface{}, 4096)
}

//初始化组件 : 任务管理 TaskPool
func (this *MyWatchFactory) initTaskPoolComponents(taskNums uint64) {
	//初始化 任务池
	if this.taskPool == nil {
		pool, err := MyTaskPools.NewPool(taskNums)
		if err != nil {
			panic(err)
		}
		this.taskPool = pool
	}
	//初始化任务组
	this.waitGroup = new(sync.WaitGroup)

	//启动 任务信号队列 监听
	go func() {
		for {
			select {
			case come := <-this.TaskChan:
				observers := ValAndObservers(come)
				this.dispatchPool(&observers)
			}
		}
	}()

}

//实现 更新观察值
func (this *MyWatchFactory) UpdateValue(topic string, value interface{}) error {
	observers := this.saveTopicAndValueAndObserver[topic]
	if observers == nil {
		return fmt.Errorf("不存在 ValueAndObservers 1、确保操作步骤争取；2、确保topic有效.topic=", topic)
	}

	//进行拷贝 value 和 被观察者的拷贝， 这样可以介绍被观察者锁的占用时间
	cpValueAndObservers := ValAndObservers{}
	this.saveTopicAndValueAndObserver[topic].Lock.Lock()
	this.saveTopicAndValueAndObserver[topic].Value = value
	createValAndObserverFromValueAndObservers(&cpValueAndObservers, this.saveTopicAndValueAndObserver[topic])
	this.saveTopicAndValueAndObserver[topic].Lock.Unlock()

	//将结合后的 value+观察者s 加入到chan中
	this.TaskChan <- cpValueAndObservers

	return nil
}

//实现 根据主题添加 被观察者
func (this *MyWatchFactory) AddObserver(topic string, observer IObserver) error {
	//判断主题是否为空
	if this.saveTopicAndValueAndObserver[topic] == nil {
		this.TopicLock.Lock()
		if this.saveTopicAndValueAndObserver[topic] == nil {
			this.saveTopicAndValueAndObserver[topic] = &ValueAndObservers{}
			this.saveTopicAndValueAndObserver[topic].Observers = &structure.ArrayList{}
		}
		this.TopicLock.Unlock()
	}

	//将observer加入
	this.saveTopicAndValueAndObserver[topic].Lock.Lock()
	//排重
	_, b := this.saveTopicAndValueAndObserver[topic].Observers.Contains(observer)
	//不重复则进行添加
	if !b {
		this.saveTopicAndValueAndObserver[topic].Observers.Add(observer)
	}
	this.saveTopicAndValueAndObserver[topic].Lock.Unlock()

	return nil
}
//实现 移除观察者
func (this *MyWatchFactory) RemoveObserver(topic string, observer IObserver) error  {
	if this.saveTopicAndValueAndObserver[topic] == nil {
		return nil
	}
	
	if this.saveTopicAndValueAndObserver[topic] != nil {
		this.TopicLock.Lock()
		defer  this.TopicLock.Unlock()
		if this.saveTopicAndValueAndObserver[topic] != nil {

			valueAndObservers := this.saveTopicAndValueAndObserver[topic]
			if valueAndObservers != nil {
				valueAndObservers.Lock.Lock()
				if valueAndObservers != nil {
					if valueAndObservers.Observers != nil {
						contains, b := valueAndObservers.Observers.Contains(observer)
						if b {
							valueAndObservers.Observers.Remove(contains)
						}
					}
				}
				defer valueAndObservers.Lock.Unlock()
			}
		}
	}
	return nil
}
//实现 移除主题
func (this *MyWatchFactory) RemoveTopic(topic string) error {
	if this.saveTopicAndValueAndObserver[topic] == nil {
		return nil
	}

	if this.saveTopicAndValueAndObserver[topic] != nil {
		this.TopicLock.Lock()
		if this.saveTopicAndValueAndObserver[topic] != nil {
			this.saveTopicAndValueAndObserver[topic].Lock.Lock()
			if this.saveTopicAndValueAndObserver[topic].Observers != nil {
				this.saveTopicAndValueAndObserver[topic].Observers.Destroy()
			}
			this.saveTopicAndValueAndObserver[topic].Lock.Unlock()
		}
		delete(this.saveTopicAndValueAndObserver, topic)
		this.TopicLock.Unlock()
	}
	return nil
}

//根据给出的 ValueAndObservers 创建相应的 ValAndObservers
func createValAndObserverFromValueAndObservers(to *ValAndObservers, from *ValueAndObservers) {
	to.Value = from.Value
	to.Observers = from.Observers.CopyArrary()
}


//任务调度处理
func (this *MyWatchFactory) dispatchPool(observers *ValAndObservers) {
	// 创建一个任务并加入
	this.waitGroup.Add(1)
	task := MyTaskPools.Task{
		Handler: func(v ...interface{}) {
			observer := v[0].(*ValAndObservers)
			if observer.Observers.Size() > 0 {
				for i := 0; i < observer.Observers.Size(); i++ {
					ob := observer.Observers.Get(i)
					if ob != nil {
						get := ob.(IObserver)
						get.WatchHandler(observers.Value)
					}
				}
			}
		},
		Params:  []interface{}{observers},
	}
	this.taskPool.Put(&task)
}

func deepCopy(dst, src interface{}) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)
}
func SimpleCopyProperties(dst, src interface{}) (err error) {
	// 防止意外panic
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprintf("%v", e))
		}
	}()

	dstType, dstValue := reflect.TypeOf(dst), reflect.ValueOf(dst)
	srcType, srcValue := reflect.TypeOf(src), reflect.ValueOf(src)

	// dst必须结构体指针类型
	if dstType.Kind() != reflect.Ptr || dstType.Elem().Kind() != reflect.Struct {
		return errors.New("dst type should be a struct pointer")
	}

	// src必须为结构体或者结构体指针
	if srcType.Kind() == reflect.Ptr {
		srcType, srcValue = srcType.Elem(), srcValue.Elem()
	}
	if srcType.Kind() != reflect.Struct {
		return errors.New("src type should be a struct or a struct pointer")
	}

	// 取具体内容
	dstType, dstValue = dstType.Elem(), dstValue.Elem()

	// 属性个数
	propertyNums := dstType.NumField()

	for i := 0; i < propertyNums; i++ {
		// 属性
		property := dstType.Field(i)
		// 待填充属性值
		propertyValue := srcValue.FieldByName(property.Name)

		// 无效，说明src没有这个属性 || 属性同名但类型不同
		if !propertyValue.IsValid() || property.Type != propertyValue.Type() {
			continue
		}

		if dstValue.Field(i).CanSet() {
			dstValue.Field(i).Set(propertyValue)
		}
	}

	return nil
}

func tt(ss string) {
	fmt.Println("tt ",ss)
}