package cache

import (
	"time"

	"github.com/patrickmn/go-cache"
)

var Cache cache.Cache

func LoadCache() *cache.Cache {

	c := cache.New(5*time.Minute, 7*time.Minute)

	return c
}
