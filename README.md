## 简介
一个高效的、无依赖的、无内部协程的、计时器系统，自行调度更新。  
默认用时间堆，有大量短时计时器需求，可以开启时间轮。  
## Quick Start
```golang
package main

import (
	"fmt"
	"github.com/Archer-26/jtimer"
	"time"
)

func main() {
	ticker := time.NewTicker(time.Millisecond * 100)
	timerMgr := jtimer.NewTimerMgr()

	//timerMgr.AvailWheel() 开启时间轮

	c := 0
	now := time.Now().UnixNano()
	timer,err := jtimer.NewTimer(now,now + int64(2*time.Second),-1, func(dt int64) {
		c++
		fmt.Printf("dt:%v  count:%v \n", dt, c)
	})
	if err != nil {
		return
	}
	timerMgr.AddTimer(timer,false)
	for {
		select {
		case <-ticker.C:
			timerMgr.Update(time.Now().UnixNano())
		}
	}
}
```
