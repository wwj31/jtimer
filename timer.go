package jtimer

import (
	"errors"
)

type FuncCallback func(dt int64)
type Timer struct {
	timeid           string       // 计时器id
	interval         int64        // 时间间隔
	next_triggertime int64        // 下一次执行的时间戳
	trigger_times    int32        // 执行次数
	func_callback    FuncCallback // 到时回调
	disabled         bool         // 失效标记
}

//创建一个timer
//now:               创建的时刻
//next_triggertime:  触发时刻
//trigger_times:     触发次数 -1 表示无限次,如果trigger_times==0 返回error
//callback:          回调函数
func NewTimer(now, next_triggertime int64, trigger_times int32, callback FuncCallback, timeId ...string) (*Timer, error) {
	if trigger_times == 0 {
		return nil, errors.New("trigger_times == 0")
	}
	if next_triggertime <= now {
		next_triggertime = now + 1
	}
	id := ""
	if len(timeId) > 0 {
		id = timeId[0]
	}
	return &Timer{
		timeid:           id,
		interval:         next_triggertime - now,
		next_triggertime: next_triggertime,
		trigger_times:    trigger_times,
		func_callback:    callback,
		disabled:         false,
	}, nil
}

func (s *Timer) Priority() int64 {
	return s.next_triggertime
}
