package main

import (
	"fmt"
	mapset "github.com/deckarep/golang-set"
	"myagent/src/core/structure"
)

func main() {
	list := structure.ArrayList{}

	list.Add("1")
	list.Add("2")
	list.Add("3")
	list.Add("4")
	list.Add("5")
	list.AppendTo(1, "0")
	fmt.Println(list)
}
func maint1() {
	set := mapset.NewSet()
	set.Add(1)
	set.Add(2)
	set.Add(3)
	fmt.Println(set)
	set.Add(3)
	fmt.Println(set)

	set1 := mapset.NewSet()

	set1.Add(1)
	set1.Add(2)
	set1.Add(10)
	set1.Add(101)
	set1.Add(102)

	set = set1
	fmt.Println(set)
}
