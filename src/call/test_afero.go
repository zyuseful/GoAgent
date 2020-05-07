package main

import (
	"fmt"
	"github.com/spf13/afero"
	MyTaskPools "github.com/wazsmwazsm/mortar"
	"myagent/src/common"
	MyWatch "myagent/src/core/watch"
	MyOS "myagent/src/myos"
	"strconv"
	"strings"
	"sync"
)

func main1() {
	/*
		var AppFs = afero.NewOsFs()
		open, err := AppFs.Open("/Users/zys")
		if (err != nil) {
			fmt.Println(err.Error())
			return
		}
		readdir, err := open.Readdir(1000)
		for _,d := range readdir {
			fmt.Println(d.Name())
		}
		defer open.Close()
	*/

	var AppFs = afero.NewOsFs()
	//ss := MyOS.SystemFileResult{}
	var err error
	//err = MyOS.CopyFile(AppFs, "/Users/zys/asdf.log", "/tmp/asdf.log")
	//if err != nil {
	//	fmt.Println(err.Error())
	//}

	err = MyOS.Copy(AppFs, "/Users/zys/testDir", "/tmp/zz/testDir")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("OK")

}

func main2() {
	var AppFs = afero.NewMemMapFs()
	filename := "/tmp/zz/hello.txt"
	info := "这是一个测试文件\nThis is a test file\nHelloWorld"

	err := afero.WriteFile(AppFs, filename, []byte(info), 755)
	if err != nil {
		fmt.Println(err.Error())
	}

	file, err := afero.ReadFile(AppFs, filename)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(string(file[:]))
}

func main4()  {
	queue := common.Queue{}
	queue.Push("Hello")
	queue.Push("World")
	pop := queue.Pop()
	fmt.Println(pop)
	pop = queue.Pop()
	fmt.Println(pop)
}

func main5()  {
	list := common.ArrayList{}
	list.Add("h1","h2","h3")
	list.Print()
	list.Add("asdfasdf","asdfasdfd")
	list.Print()
	fmt.Println(list.Contains("h3"))
}

func main6()  {
	// 创建容量为 10 的任务池
	pool, err := MyTaskPools.NewPool(10)
	if err != nil {
		panic(err)
	}

	wg := new(sync.WaitGroup)
	// 创建任务
	task := &MyTaskPools.Task{
		Handler: func(v ...interface{}) {
			wg.Done()
			fmt.Println(v)
		},
	}

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		// 添加任务函数的参数
		task.Params = []interface{}{i, i * 2, "hello"}
		// 将任务放入任务池
		pool.Put(task)
	}

	wg.Add(1)
	// 再创建一个任务
	pool.Put(&MyTaskPools.Task{
		Handler: func(v ...interface{}) {
			wg.Done()
			fmt.Println(v)
		},
		Params: []interface{}{"hi!"}, // 也可以在创建任务时设置参数
	})

	wg.Wait()

	// 安全关闭任务池（保证已加入池中的任务被消费完）
	pool.Close()
	// 如果任务池已经关闭, Put() 方法会返回 ErrPoolAlreadyClosed 错误
	err = pool.Put(&MyTaskPools.Task{
		Handler: func(v ...interface{}) {},
	})
	if err != nil {
		fmt.Println(err) // print: pool already closed
	}
}

func main()  {
	c := make(chan bool,1)

	factory := MyWatch.MyWatchFactory{}
	factory.Init()

	/*
	var w1 MyWatch.IObserver
	var w2 MyWatch.IObserver
	w1 = &s1{"s1"}
	w2 = &s1{"s2"}
	factory.AddObserver("S1",w1)
	factory.AddObserver("S1",w2)
	factory.AddObserver("S2",w1)
	factory.UpdateValue("S1","Hello")
	factory.UpdateValue("S1","World")
	factory.UpdateValue("S2","World")
	*/

	for i:=0;i<5;i++ {
		var topic,name strings.Builder
		topic.WriteString("tp")
		topic.WriteString(strconv.Itoa(i))
		for j:=0;j<10;j++ {
			name.WriteString("task")
			name.WriteString(strconv.Itoa(i))
			name.WriteString("-")
			name.WriteString(strconv.Itoa(j))
			testThread(&factory,name.String(),topic.String())
		}
	}

	for i:=0;i<100;i++ {
		var topic,val strings.Builder
		topic.WriteString("tp")
		topic.WriteString(strconv.Itoa(i))
		val.WriteString("val")
		val.WriteString(strconv.Itoa(i))

		go testY(&factory,topic.String(),val.String())
	}
	<-c
}

func testThread(factory *MyWatch.MyWatchFactory,name string,topic string) {
	var w1 MyWatch.IObserver
	w1= &s1{name}
	factory.AddObserver(topic,w1)
}
func testY(factory *MyWatch.MyWatchFactory,topic string,val string) {
	factory.UpdateValue(topic,val)
}

type s1 struct {
	Name string
}
func (this *s1) WatchHandler(v interface{}) error {
	fmt.Println("I am ",this.Name,"  do val:",v)
	return nil
}