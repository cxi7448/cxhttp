package xwork

import (
	"github.com/cxi7448/cxhttp/clUtil/clJson"
	"github.com/cxi7448/cxhttp/clUtil/processbar"
	"sync"
)

type XWork struct {
	Num     int
	locker  sync.RWMutex
	Queue   []XQueue // 队列
	process *processbar.ProcessBar
}
type XQueue struct {
	Param    clJson.M
	Callback func(param clJson.M)
}

func New(num int) *XWork {
	xwork := &XWork{
		Num:   num,
		Queue: []XQueue{},
	}
	return xwork
}

func (this *XWork) AddWork(param clJson.M, _func func(_param clJson.M)) {
	this.locker.Lock()
	defer this.locker.Unlock()
	this.Queue = append(this.Queue, XQueue{
		Param:    param,
		Callback: _func,
	})
}

func (this *XWork) pop() *XQueue {
	this.locker.Lock()
	defer this.locker.Unlock()
	if len(this.Queue) == 0 {
		return nil
	}
	queue := this.Queue[0]
	newSlice := this.Queue[1:]
	this.Queue = newSlice
	return &queue
}

func (this *XWork) doWork() {
	ch := make(chan string)
	for i := 0; i < this.Num; i++ {
		go func(_ch chan string) {
			for {
				queue := this.pop()
				if queue == nil {
					break
				}
				queue.Callback(queue.Param)
				if this.process != nil {
					this.process.Add()
				}
			}
			ch <- "success"
		}(ch)
	}
	for i := 0; i < this.Num; i++ {
		<-ch
	}
}
func (this *XWork) Wait() {
	this.process = processbar.New(float32(len(this.Queue)))
	this.doWork()
}

func (this *XWork) Run() {
	go this.doWork()
}
