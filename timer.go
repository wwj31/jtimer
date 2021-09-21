package jtimer

import (
	"errors"
)

type callback func(dt int64)
type Timer struct {
	timeid   string   // 计时器id
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
	return &Timer{
		timeid:   id,
		startAt:  now,
		interval: endAt - now,
		endAt:    endAt,
		count:    count,
		cb:       callback,
		disabled: false,
	}, nil
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
