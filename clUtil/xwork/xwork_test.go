package xwork

import (
	"github.com/cxi7448/cxhttp/clUtil/clJson"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	work := New(10)
	var queue = []clJson.M{}
	for i := 0; i < 1000; i++ {
		work.AddWork(clJson.M{"i": i}, func(m clJson.M) {
			//fmt.Println(m["i"])
			queue = append(queue, m)
			<-time.After(time.Millisecond * 100)
		})
	}
	work.Wait()
	<-time.After(time.Second * 5)
}
