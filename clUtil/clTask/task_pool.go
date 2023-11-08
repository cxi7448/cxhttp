package clTask

import (
	"github.com/cxi7448/cxhttp/clUtil/clLog"
	"github.com/cxi7448/cxhttp/clUtil/clRedis"
	"github.com/cxi7448/cxhttp/clUtil/clTime"
	"time"
)

var redisPtr *clRedis.RedisObject

// 建立一个新的任务池
func NewPool(_redis *clRedis.RedisObject) *TaskPool {
	redisPtr = _redis
	return &TaskPool{
		pool: make([]*TaskInfo, 0),
	}
}

func GetRedis() *clRedis.RedisObject {
	return redisPtr
}

// 开始检查任务池中的任务是否得以执行
func (this *TaskPool) Start() {

	for {

		nowTime, _ := clTime.NewDate("")
		todayBegin, _ := nowTime.GetDayBetween(0)
		for _, val := range this.pool {
			go val.Run(nowTime.TimeStamp, (nowTime.TimeStamp-todayBegin)%86400, uint32(1<<(nowTime.Week)), this.recoverCallback)
		}

		<-time.After(1 * time.Second)
	}

}

// 添加新的计划任务
func (this *TaskPool) AddNew(task *TaskInfo) {
	clLog.Info("定时任务添加: %+v", task)
	this.pool = append(this.pool, task)
}

// 设置异常捕获回调
func (this *TaskPool) SetRecoverCallback(_f func(string, string)) {
	this.recoverCallback = _f
}
