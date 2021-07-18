package jtimer

import (
	"fmt"
	"runtime/debug"
)

type TimerMgr struct {
	timers     *Heap
	id2timer   map[int64]*Timer
	cur_timeid int64 // 自增id
}

func NewTimerMgr() *TimerMgr {
	timer_mgr := TimerMgr{}
	timer_mgr.cur_timeid = 0
	timer_mgr.id2timer = make(map[int64]*Timer)
	timer_mgr.timers = NewQueue(nil, MIN_HEAP, QUAD) // 计时器统一用小顶堆

	return &timer_mgr
}

func (s *TimerMgr) Reset() {
	s.timers = NewQueue(nil, MIN_HEAP, QUAD)
	s.id2timer = make(map[int64]*Timer)
	s.cur_timeid = 0
}

// AddTimer
func (s *TimerMgr) AddTimer(timer *Timer) int64 {
	s.cur_timeid += 1
	timer.timeid = s.cur_timeid

	s.timers.Push(timer)
	s.id2timer[timer.timeid] = timer
	return timer.timeid
}

// CancelTimer
func (s *TimerMgr) CancelTimer(timeid int64) {
	if timer, ok := s.id2timer[timeid]; ok {
		timer.disabled = true
	}
}

func try(fn func()) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("[%v] panic recover %v", r, string(debug.Stack()))
		}
	}()
	fn()
}

func (s *TimerMgr) Update(now int64) {
	for {
		intf := s.timers.Peek()
		if intf == nil {
			break
		}
		timer := intf.(*Timer)
		del := timer.disabled
		if !del && timer.interval > 0 {
			// 检查执行时间是否到了
			delayTime := now - timer.next_triggertime
			if delayTime < 0 {
				break
			}
			//fmt.Printf("now:%v  - abtime.next_triggertime:%v = delayTime:%v\n", now, abtime.next_triggertime, delayTime)
			overtimes := (delayTime + timer.interval) / timer.interval
			for i := 0; i < int(overtimes); i++ {
				if timer.trigger_times == 0 {
					break
				}
				try(func() {
					if timer.trigger_times == 1 {
						// 最后一次，间隔+delay
						timer.func_callback(timer.interval + delayTime)
					} else {
						// 不是最后一次，按照固定间隔执行
						timer.func_callback(timer.interval)
					}
				})

				if timer.trigger_times > 0 {
					timer.trigger_times--
				}
			}

			// 还有次数,继续加入优先队列
			if timer.trigger_times != 0 && !timer.disabled {
				timer.next_triggertime += timer.interval * overtimes
				s.timers.Topdown()
			} else {
				del = true
			}
		}
		if del {
			s.timers.Pop()
			delete(s.id2timer, timer.timeid)
		}
	}
}
