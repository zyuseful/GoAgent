package structure

import (
	"fmt"
)

type ArrayList struct {
	elements []interface{}
	size     int
}

//func New(values ...interface{}) *ArrayList {
//	list := &ArrayList{}
//	list.elements = make([]interface{}, 10)
//	if len(values) > 0 {
//		list.Add(values...)
//	}
//	return list
//}

//扩容 num <= 0 默认扩容原长度的1.5倍 num != 0 则扩容 num长度
func (list *ArrayList) expand(num int) {
	var expandSize int = list.size
	if num <= 0 {
		expandSize += expandSize / 2
	} else {
		expandSize = num
	}

	newElements := make([]interface{}, expandSize)
	copy(newElements, list.elements)
	list.elements = newElements
}

func (list *ArrayList) Add(values ...interface{}) {
	if list.size+len(values) > len(list.elements) {
		//if list.size+len(values) >= len(list.elements)-1 {
		arrLen := len(values)
		//newElements := make([]interface{}, list.size+arrLen+arrLen*2)
		list.expand(list.size + arrLen + (list.size / 2))
	}

	for _, value := range values {
		list.elements[list.size] = value
		list.size++
	}

}

func (list *ArrayList) Remove(index int) interface{} {
	if index < 0 || index >= list.size {
		return nil
	}

	curEle := list.elements[index]
	list.elements[index] = nil
	copy(list.elements[index:], list.elements[index+1:list.size])
	list.size--
	return curEle
}

func (list *ArrayList) Get(index int) interface{} {
	if index < 0 || index >= list.size {
		return nil
	}
	return list.elements[index]
}

func (list *ArrayList) IsEmpty() bool {
	return list.size == 0
}

func (list *ArrayList) Size() int {
	return list.size
}
func (list *ArrayList) Contains(value interface{}) (int, bool) {
	for index, curValue := range list.elements {
		if curValue == value {
			return index, true
		}
	}
	return -1, false
}
func (list *ArrayList) CopyArrary() *ArrayList {
	if list.Size() <= 0 {
		//return nil,fmt.Errorf("原始数组为空")
		return nil
	}

	result := &ArrayList{}
	result.elements = make([]interface{}, list.size)
	copy(result.elements, list.elements)
	result.size = list.Size()

	return result
}

func (list *ArrayList) Destroy() {
	if nil == list {
		return
	}

	for i := 0; i < list.size; i++ {
		list.Remove(i)
	}
	list = nil
}

func (list *ArrayList) Print() {
	fmt.Println(list.elements)
}

func (list *ArrayList) AppendTo(index int, values ...interface{}) error {
	//如果空间不足先进行扩容
	if list.size+len(values) > len(list.elements) {
		arrLen := len(values)
		list.expand(list.size + arrLen + (list.size / 2))
	}

	//计算追加元素以后的偏移量 , 该值 > 0说明是在源数组中间进行追加 =0 说明是在末尾追加 <0 说明数组越界
	afterIndexNum := list.size - index - 1

	if afterIndexNum < 0 {
		var err error
		fmt.Errorf("Fatal error config file: %s \n", err)
		return err
	} else if afterIndexNum > 0 {
		//中间追加，先将原index后的元素向后平移（这里不做 size++）
		var i int
		for i = 1; i <= afterIndexNum; i++ {
			list.elements[list.size+afterIndexNum+i-1] = list.elements[index+i]
		}
		//将追加内容向空位进行插入 同时 size++
		for i = 0; i < len(values); i++ {
			list.elements[index+i+1] = values[i]
			list.size++
		}
	} else {
		//与原append模式一样处理
		for _, value := range values {
			list.elements[list.size] = value
			list.size++
		}
	}
	return nil
}
