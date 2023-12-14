package xwork

import (
	"fmt"
	"github.com/cxi7448/cxhttp/clUtil/clJson"
	"testing"
)

func TestNew(t *testing.T) {
	work := New(10)
	for i := 0; i < 1000; i++ {
		work.AddWork(clJson.M{"i": i}, func(m clJson.M) {
			fmt.Println(m["i"])
		})
	}
	work.Wait()
}
