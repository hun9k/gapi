package cache

import (
	"fmt"
	"time"

	"github.com/hun9k/gapi/conf"
)

// cacher interface
type Cacher interface {
	Set(key string, value any, expiration time.Duration)
	Get(key string) (any, bool)
	Delete(key string)
	Clear()
	Stats() map[string]any
}

const (
	DEFAULT_KEY     = "default"
	NoExpiration    = 0
	CleanupInterval = 10 * time.Minute
)

var cachePool = map[string]Cacher{}

func Inst(keys ...string) Cacher {
	key := DEFAULT_KEY
	if len(keys) > 0 {
		key = keys[0]
	}

	// create cacher if isn't exists
	if cachePool[key] == nil {
		switch conf.Get[string](fmt.Sprintf("cache.%s.driver", key)) {
		case "redis":
		case "local":
			cachePool[key] = NewLocal(NoExpiration, CleanupInterval)
		}
	}

	return cachePool[key]
}

// func Set(key string, value any, expiration time.Duration) {
// 	// TODO: Implement cache set
// }

// func Get(key string) any {
// 	// TODO: Implement cache get
// 	return ""
// }

// func Delete(key string) {
// 	// TODO: Implement cache delete
// }

// func Clear() {
// 	// TODO: Implement cache clear
// }

// func Exists(key string) bool {

// 	return false
// }

// func Stats() map[string]any {
// 	return map[string]any{}
// }
