package jtimer

import (
	"errors"
	"sync"
)

var timerPool = sync.Pool{New: func() interface{} { return new(Timer) }}

type callback func(dt int64)
type Timer struct {
	timeId   string   // 计时器id
	startAt  int64    // 开启时间
	interval int64    // 时间间隔
	endAt    int64    // 结束时间
	count    int32    // 执行次数
	cb       callback // 到时回调
	disabled bool     // 失效标记
	heapIdx  int      // 堆索引
}

//创建一个timer
//now:    	创建的时刻
//endAt:  	触发时刻
//count:  	触发次数 -1 表示无限次,如果trigger_times==0 返回error
//callback: 回调函数
func NewTimer(now, endAt int64, count int32, callback callback, timeId ...string) (*Timer, error) {
	if count == 0 {
		return nil, errors.New("count == 0")
	}
	if endAt <= now {
		endAt = now + 1
	}
	id := ""
	if len(timeId) > 0 {
		id = timeId[0]
	}
	timer := timerPool.Get().(*Timer)
	timer.Rest()
	timer.timeId = id
	timer.startAt = now
	timer.interval = endAt - now
	timer.endAt = endAt
	timer.count = count
	timer.cb = callback
	timer.disabled = false
	return timer, nil
}

func (s *Timer) Priority() int64 {
	return s.endAt
}

func (s *Timer) SetIndex(i int) {
	s.heapIdx = i
}

func (s *Timer) GetIndex() int {
	return s.heapIdx
}

func (s *Timer) Rest() {
	s.timeId = ""
	s.startAt = 0
	s.interval = 0
	s.endAt = 0
	s.count = 0
	s.cb = nil
	s.disabled = false
	s.heapIdx = 0
}
