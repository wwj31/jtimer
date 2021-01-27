package abtime

import (
	"fmt"
	"testing"
	"time"
)

func TestTimerMgr(t *testing.T) {
	tt := time.Millisecond * 100
	tt1 := time.Now().UnixNano() / int64(time.Millisecond)
	fmt.Println("begin -------------------------", tt1)
	timer := time.NewTimer(tt)
	timer2 := time.NewTimer(time.Second * 1)

	cout := 0
	timerMgr := NewTimerMgr()
	timerMgr.AddTimer(&Timer{
		interval:         1000,
		next_triggertime: time.Now().UnixNano()/int64(time.Millisecond) + 1000,
		trigger_times:    10,
		func_callback: func(dt int64) {
			cout++
			fmt.Printf("dt:%v  count:%v \n", dt, cout)
		},
	})
	for {
		select {
		case <-timer.C:
			timer.Reset(tt)
			timerMgr.Update(time.Now().UnixNano() / int64(time.Millisecond))
		case <-timer2.C:
			timer2.Reset(time.Second * 1)
			tt2 := time.Now().UnixNano() / int64(time.Millisecond)
			fmt.Println("end -------------------------", tt2, tt2-tt1)
			tt1 = tt2
		}
	}

}
