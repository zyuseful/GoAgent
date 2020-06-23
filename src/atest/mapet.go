package main

import (
	"fmt"
	mapset "github.com/deckarep/golang-set"
)
func main() {
	set := mapset.NewSet()
	set.Add(1)
	set.Add(2)
	set.Add(3)
	fmt.Println(set)
	set.Add(3)
	fmt.Println(set)


}
