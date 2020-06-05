package main

import (
	"github.com/patrickmn/go-cache"
	"myagent/src/core/memcache"
)

func main() {
	/*
	c := cache.New(0, 0)
	c.Set("k1","HelloV1",cache.NoExpiration)
	fmt.Println(c.Get("k1"))
	c.Set("k1","HelloV2",cache.NoExpiration)
	fmt.Println(c.Get("k1"))
	*/

	getMemcache := memcache.GetMemcache()
	getMemcache.Set("k1","V1",cache.NoExpiration)
	getMemcache.Get("k1")


}
