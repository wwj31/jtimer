package jtimer

import (
	"fmt"
	"github.com/satori/go.uuid"
	"runtime/debug"
)

type TimerMgr struct {
	timers   Heap
	id2timer map[string]*Timer
}

func NewTimerMgr() TimerMgr {
	timer_mgr := TimerMgr{}
	timer_mgr.id2timer = make(map[string]*Timer)
	timer_mgr.timers = NewQueue(nil, MIN_HEAP, QUAD) // 计时器统一用小顶堆

	return timer_mgr
}

func (s *TimerMgr) Reset() {
	s.timers = NewQueue(nil, MIN_HEAP, QUAD)
	s.id2timer = make(map[string]*Timer)
}

func (s *TimerMgr) Empty() bool {
	return len(s.id2timer) == 0
}

// AddTimer
func (s *TimerMgr) AddTimer(timer *Timer) string {
	if timer.timeid == "" {
		timer.timeid = uuid.NewV4().String()
	}

	s.timers.Push(timer)
	s.id2timer[timer.timeid] = timer
	return timer.timeid
}

func (s *TimerMgr) UpdateTimer(key string, endAt int64) error {
	oldTimer, ok := s.id2timer[key]
	if !ok {
		return ErrorUpdateTimer
	}
	newTimer, _ := NewTimer(oldTimer.startAt, endAt, oldTimer.count, oldTimer.cb, oldTimer.timeid)
	newTimer.heapIdx = oldTimer.heapIdx
	return s.timers.Update(newTimer)
}

// CancelTimer  if del == true, the complexity is O(log n) else is O(1)
func (s *TimerMgr) CancelTimer(timeid string, del ...bool) {
	timer, ok := s.id2timer[timeid]
	if !ok {
		return
	}

	if len(del) == 0 || !del[0] {
		timer.disabled = true
		return
	}

	s.timers.Delete(timer)
	delete(s.id2timer, timeid)
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
		head := s.timers.Peek()
		if head == nil {
			break
		}
		timer := head.(*Timer)
		del := timer.disabled
		if !del && timer.interval > 0 {
			delayTime := now - timer.endAt
			if delayTime < 0 {
				break
			}

			overtimes := (delayTime + timer.interval) / timer.interval
			for i := 0; i < int(overtimes); i++ {
				if timer.count == 0 {
					break
				}
				try(func() {
					if timer.count == 1 {
						// 最后一次，间隔+delay
						timer.cb(timer.interval + delayTime)
					} else {
						// 不是最后一次，按照固定间隔执行
						timer.cb(timer.interval)
					}
				})

				if timer.count > 0 {
					timer.count--
				}
			}

			// 还有次数,继续加入优先队列
			if timer.count != 0 && !timer.disabled {
				timer.endAt += timer.interval * overtimes
				s.timers.Topdown()
			} else {
				del = true
			}
		}
		if del {
			s.timers.Pop()
			delete(s.id2timer, timer.timeid)
			timerPool.Put(timer)
		}
	}
}
