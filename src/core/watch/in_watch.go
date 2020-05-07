package watch

//watch 工厂
type IWatchFactory interface {
	//初始化
	Init()
	//更新观察值
	UpdateValue(topic string, value interface{}) error
	//添加被观察者
	AddObserver(topic string, observer *IObserver) error
	//移除观察者
	//RemoveObserver(topic string, observer *IObserver) error


}

//被观察者
type ObserverStruct struct {}
type IObserver interface {
	//通知代理
	WatchHandler(v interface{}) error
}
