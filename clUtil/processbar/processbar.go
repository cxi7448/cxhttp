package processbar

import (
	"fmt"
	"sync"
)

//go get -u github.com/schollz/progressbar/v2

type ProcessBar struct {
	width  int
	total  float32
	cur    int
	locker sync.RWMutex
}

func New(total float32) *ProcessBar {
	p := &ProcessBar{
		total: total,
		width: 100,
	}
	return p
}

func (this *ProcessBar) Add() {
	this.locker.Lock()
	defer this.locker.Unlock()
	if float32(this.cur) > this.total {
		return
	}
	this.cur += 1
	per := float32(this.cur) / float32(this.total) * 100
	cur_process := int(per)
	var process = ""
	for j := 0; j < cur_process; j++ {
		process += "="
	}
	process += fmt.Sprintf("%0.2f%%", per)
	for j := cur_process; j < this.width; j++ {
		process += " "
	}
	fmt.Printf("\r[%v]%v/%v", process, this.cur, this.total)

	if float32(this.cur) >= this.total {
		fmt.Printf("\n")
	}
}
