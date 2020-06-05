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

func (list *ArrayList) Add(values ...interface{}) {
	if list.size+len(values) > len(list.elements) {
		//if list.size+len(values) >= len(list.elements)-1 {
		arrLen := len(values)
		newElements := make([]interface{}, list.size+arrLen+arrLen*2)
		copy(newElements, list.elements)
		list.elements = newElements
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
	copy(result.elements,list.elements)
	result.size = list.Size()

	return result
}

func (list *ArrayList) Destroy()  {
	if nil == list {
		return
	}

	for i:=0;i<list.size;i++ {
		list.Remove(i)
	}
	list = nil
}

func (list *ArrayList) Print() {
	fmt.Println(list.elements)
}
