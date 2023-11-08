package clSuperMap

import "sync"

// 超级Map
type SuperMap struct {
	data map[string] string
	locker sync.RWMutex
}