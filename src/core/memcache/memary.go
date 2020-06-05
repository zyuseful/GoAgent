package memcache

import (
	"github.com/patrickmn/go-cache"
)


var memcache *cache.Cache

func init() {
	memcache := cache.New(0, 0)
	memcache.Set("k1","HelloV1",cache.NoExpiration)
}


func GetMemcache() *cache.Cache {
	return memcache
}

