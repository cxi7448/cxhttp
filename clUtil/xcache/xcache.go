package xcache

import (
	"sync"
	"time"
)

// 常驻内容的任务

type CacheList struct {
	Data   map[string]Cache
	locker sync.RWMutex
}
type Cache struct {
	Key    string
	Value  interface{}
	Expire int64
}

var sLocker sync.RWMutex
var pLen = 100
var cachePools = make(map[int]*CacheList)

func init() {
	for i := 0; i < pLen; i++ {
		cachePools[i] = &CacheList{
			Data: map[string]Cache{},
		}
	}
	go func() {
		for {
			clearCache()
			<-time.After(time.Second * 60)
		}
	}()
}
func getSliceKey(key string) int {
	total := 0
	for _, a := range key {
		total += int(a)
	}
	return total % pLen
}

func Get(key string) interface{} {
	index := getSliceKey(key)
	sLocker.RLock()
	c := cachePools[index]
	sLocker.RUnlock()
	if c == nil {
		return nil
	}
	c.locker.RLock()
	defer c.locker.RUnlock()
	data, ok := c.Data[key]
	if !ok {
		return nil
	}
	if data.Expire < time.Now().Unix() {
		return nil
	}
	return data.Value
}

func Set(key string, value interface{}, expire int64) {
	index := getSliceKey(key)
	sLocker.RLock()
	c := cachePools[index]
	sLocker.RUnlock()
	if c != nil {
		c.locker.RLock()
		defer c.locker.RUnlock()
		c.Data[key] = Cache{
			Key:    key,
			Value:  value,
			Expire: time.Now().Unix() + expire,
		}
	}
}

func (this *CacheList) clearCache() {
	this.locker.Lock()
	defer this.locker.Unlock()
	now := time.Now().Unix()
	for key := range this.Data {
		if this.Data[key].Expire < now {
			delete(this.Data, key)
		}
	}
}

func clearCache() {
	sLocker.RLock()
	defer sLocker.RUnlock()
	for i := range cachePools {
		cachePools[i].clearCache()
	}
}
