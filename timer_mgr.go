package jtimer

import (
	"fmt"
	"./priorityqueue"
	"runtime/debug"
	"time"
)

var (
	/*
		为了保证时间轮简单，暂时固定参数设置：
			轮动间隔:100ms
			计时范围:1天
			所需槽数:864000=DAY_MS/ROTATE_INTERVAL
	*/
	wheel_dayTime  = 24 * time.Hour         // 转动间隔时间(毫秒)
	wheel_interval = 100 * time.Millisecond // 一天(毫秒)
)

type FuncCallback func(dt int64)

type TimerMgr struct {
	timers     *priorityqueue.Heap
	wheel      *wheel
	id2timer   map[int64]*Timer
	cur_timeid int64 // 自增id
}

// 创建一个timer管理器
func NewTimerMgr() *TimerMgr {
	timer_mgr := TimerMgr{}
	timer_mgr.cur_timeid = 0
	timer_mgr.id2timer = make(map[int64]*Timer)
	timer_mgr.timers = priorityqueue.NewQueue(nil, priorityqueue.MIN_HEAP) // 计时器统一用小顶堆

	return &timer_mgr
}

func (s *TimerMgr) newWheel() *wheel {
	return newWheel(int(wheel_dayTime/wheel_interval), int64(wheel_interval))
}
func (s *TimerMgr) AvailWheel() error {
	s.wheel = s.newWheel()
	now := time.Now().UnixNano()
	wheel_timer, err := NewTimer(now, now+int64(wheel_interval), -1, func(dt int64) {
		s.wheel.update()
	})
	if err != nil {
		return err
	}
	s.AddTimer(wheel_timer, true)
	return nil
}

// 重置 会丢弃所有timer
func (s *TimerMgr) Reset() {
	s.timers = priorityqueue.NewQueue(nil, priorityqueue.MIN_HEAP)
	if s.wheel != nil {
		s.wheel = s.newWheel()
	}
	s.id2timer = make(map[int64]*Timer)
	s.cur_timeid = 0
}

// 增加一个timer
func (s *TimerMgr) AddTimer(timer *Timer, forceHeap bool) int64 {
	s.cur_timeid += 1
	timer.timeid = s.cur_timeid

	// 时间范围小于时间轮计时范围，优先丢给时间轮处理，否则丢给时间堆处理
	if s.wheel != nil && timer.interval < s.wheel.limit() && !forceHeap {
		s.wheel.addTime(timer)
	} else {
		s.timers.Push(timer)
	}
	s.id2timer[timer.timeid] = timer
	return timer.timeid
}

// 注销一个timer
func (s *TimerMgr) CancelTimer(timeid int64) {
	if timer, ok := s.id2timer[timeid]; ok {
		timer.disabled = true // 到期判断有效性
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

		if !timer.disabled && timer.interval > 0{
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
			if timer.trigger_times != 0 && !timer.disabled{
				timer.next_triggertime += timer.interval * overtimes
				s.timers.Push(timer)
			}
		}
		s.timers.Pop()
	}
}