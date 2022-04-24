package jtimer

import (
	"container/heap"
	"math/rand"
	"testing"
	"time"
)

func verify(t *testing.T, th timerHeap) {
	for index, val := range th {
		if index != val.index {
			t.Errorf("verify error index != time.index index:[%v], timer:%+v", index, *val)
			return
		}
	}
}

func TestHeapIndex(t *testing.T) {
	h := make(timerHeap, 0)
	heap.Init(&h)

	for i := 0; i < 10000; i++ {
		v := rand.Int63n(10000)
		heap.Push(&h, &timer{endAt: time.Unix(0, v)})
		verify(t, h)
	}

	for i := 0; i < 10000; i++ {
		index := rand.Intn(h.Len())
		heap.Remove(&h, index)
		verify(t, h)
	}
}

func TestTimer(t *testing.T) {
	timerMgr := New()

	newTimer := func() {
		startAt := time.Now()
		delayTime := time.Duration(rand.Int63n(10000)) * time.Millisecond
		timerMgr.Add(
			startAt.Add(delayTime),
			func(dt time.Duration) {
				duration := time.Now().Sub(startAt)
				if duration/time.Millisecond != delayTime/time.Millisecond {
					t.Errorf("timer error delayTime:%v, duration:%v", delayTime, duration)
				} else {
					println("OK")
				}
			},
			1,
		)
	}

	go func() {
		for {
			newTimer()
			time.Sleep(time.Millisecond * 100)
		}
	}()
	go func() {
		tick := time.NewTicker(time.Millisecond)
		for range tick.C {
			timerMgr.Update(time.Now())
		}
	}()
	time.Sleep(time.Hour)
}
