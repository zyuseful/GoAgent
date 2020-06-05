package main

import (
	"fmt"
	"time"
)

func main() {
	// 1.创建一个管道
	myCh1 := make(chan int, 5)
	myCh2 := make(chan int, 5)
	exitCh := make(chan bool)

	// 2.开启一个协程生产数据
	go func() {
		//time.Sleep(time.Second * 5) 如果存在这行数据会打印超时了,不存在则会正常消费
		for i := 0; i < 10; i++ {
			myCh1 <- i
			fmt.Println("生产者1生产了", i)
		}
		close(myCh1)
		exitCh <- true
	}()

	go func() {
		time.Sleep(time.Second * 5)
		for i := 0; i < 10; i++ {
			myCh2 <- i
			fmt.Println("生产者2生产了", i)
		}
		close(myCh2)
	}()

	for {
		select {
		case num1 := <-myCh1:
			fmt.Println("------消费者消费了myCh1", num1)
		case <-time.After(3):
			fmt.Println("超时了")
			return
		case <-exitCh:
			close(exitCh)
			goto ss
		}
		time.Sleep(time.Millisecond)
	}
ss:
	fmt.Println("程序结束了")
}
