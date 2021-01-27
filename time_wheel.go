package abtime

import (
	"root/pkg/log"
	"root/pkg/tools"
)

/*
	一个简单的单轮盘时间轮，外部依赖timeMgr对delay的分割处理
*/

type (
	wheel struct {
		index    int        // 更新槽位坐标
		slots    []slotInfo // 轮槽位
		interval int64      // 间隔时间
		lasttime int64      // 最后一次更新的时间戳
	}

	slotInfo []*Timer
)

func newWheel(size int, inter int64) *wheel {
	w := &wheel{
		index:    0,
		interval: inter,
		slots:    make([]slotInfo, size, size),
		lasttime: Now().UnixNano(),
	}
	return w
}

func (s *wheel) limit() int64 {
	// 槽数-1，防止出现在update中调用addTime时，刚好添加到index所在槽位引起的问题，故此减少一个槽位的时间
	return s.interval * int64(len(s.slots)-1)
}

func (s *wheel) addTime(timer *Timer) {
	if timer.interval >= s.limit() {
		log.KVs(log.Fields{"interval": timer.interval, "limit": s.limit()}).ErrorStack(3, "timer interval out of the range")
		return
	}

	diff := timer.next_triggertime - s.lasttime
	inc := int(diff / s.interval)
	fitPos := s.index + inc // 算出需要安置的槽位
	if fitPos >= len(s.slots) {
		fitPos = fitPos - len(s.slots)
	}
	if s.slots[fitPos] == nil {
		s.slots[fitPos] = make(slotInfo, 0)
	}
	s.slots[fitPos] = append(s.slots[fitPos], timer)
}

func (s *wheel) update() {
	times := s.slots[s.index]
	defer func() {
		s.index++
		if s.index >= len(s.slots) {
			s.index = 0
		}
		s.lasttime += s.interval
	}()

	if len(times) == 0 {
		return
	}

	for _, time := range times {
		if !time.disabled {
			tools.Try(func() {
				time.func_callback(s.interval)
			}, nil)

			if time.trigger_times > 0 {
				time.trigger_times--
			}
			if time.trigger_times > 0 || time.trigger_times == -1 {
				time.next_triggertime += time.interval
				s.addTime(time)
			}
		}
	}
	s.slots[s.index] = times[0:0]
}
