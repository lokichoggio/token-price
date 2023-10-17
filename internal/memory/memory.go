package memory

import (
	"time"

	cache "github.com/patrickmn/go-cache"
)

var memoryCache *cache.Cache

func InitCache() {
	memoryCache = cache.New(time.Second*10, time.Second*20)
}

func Set(k string, v interface{}, d time.Duration) {
	memoryCache.Set(k, v, d)
}

func Get(k string) (interface{}, bool) {
	return memoryCache.Get(k)
}
