## 简介
a timer implementation in proprity queue,like event-loop and without extra gorountine.

## Quick Start
```golang
package main

import (
	"fmt"
	"time"

	"github.com/wwj31/jtimer"
)

func main() {
	ticker := time.NewTicker(time.Millisecond * 100)
	timerMgr := jtimer.NewTimerMgr()

	c := 0
	now := time.Now().UnixNano()
	timer,err := jtimer.NewTimer(now,now + int64(2*time.Second),-1, func(dt int64) {
		c++
		fmt.Printf("dt:%v  count:%v \n", dt, c)
	})
	if err != nil {
		return
	}
	timerMgr.AddTimer(timer)
	for {
		select {
		case <-ticker.C:
			timerMgr.Update(time.Now().UnixNano())
		}
	}
}
```
