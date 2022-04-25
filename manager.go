package jtimer

import (
	"container/heap"
	"time"

	"github.com/rs/xid"
)

type CallbackFn func(dt time.Duration)

type Manager struct {
	timers map[Id]*timer
	heap   timerHeap
}

func New() *Manager {
	h := make(timerHeap, 0)
	heap.Init(&h)

	return &Manager{
		timers: make(map[Id]*timer),
		heap:   h,
	}
}

func (m *Manager) Add(now, endAt time.Time, callback CallbackFn, execCount int, id ...Id) Id {
	var timerId Id
	if len(id) > 0 {
		timerId = id[0]
	} else {
		timerId = xid.New().String()
	}

	if endAt.Before(now) {
		endAt = now
	}

	if oldTimer, exist := m.timers[timerId]; exist {
		oldTimer.startAt = now
		oldTimer.endAt = endAt
		oldTimer.callback = callback
		oldTimer.execCount = execCount
		heap.Fix(&m.heap, oldTimer.index)
		return oldTimer.id
	}

	newTimer := &timer{
		id:        timerId,
		startAt:   now,
		endAt:     endAt,
		callback:  callback,
		execCount: execCount,
	}

	heap.Push(&m.heap, newTimer)
	m.timers[timerId] = newTimer
	return newTimer.id
}

func (m *Manager) Remove(id Id, softRemove ...bool) {
	var softrm bool
	if len(softRemove) > 0 && softRemove[0] {
		softrm = true
	}
	m.remove(id, softrm)
}

func (m *Manager) Len() int {
	return len(m.heap)
}

func (m *Manager) NextUpdateAt() (at time.Time) {
	headTimer := m.heap.peek()
	if headTimer == nil {
		return time.Unix(0, 0)
	}

	return headTimer.endAt
}

func (m *Manager) Update(now time.Time) {
	headTimer := m.heap.peek()

	for headTimer != nil {
		if now.Before(headTimer.endAt) {
			return
		}

		if headTimer.remove {
			m.remove(headTimer.id, false)
			headTimer = m.heap.peek()
			continue
		}

		timerDuration := headTimer.endAt.Sub(headTimer.startAt)
		totalElapsed := now.Sub(headTimer.startAt)

		if headTimer.spareCount() {
			count := totalElapsed / timerDuration
			elapsedDuration := count * timerDuration
			if headTimer.callback != nil {
				headTimer.callback(elapsedDuration)
			}

			headTimer.consumeCount(int(count))
			headTimer.startAt = headTimer.startAt.Add(elapsedDuration)
			headTimer.endAt = headTimer.startAt.Add(timerDuration)
		}

		if !headTimer.spareCount() {
			m.remove(headTimer.id, false)
		} else {
			heap.Fix(&m.heap, headTimer.index)
		}

		headTimer = m.heap.peek()
	}
}

func (m *Manager) remove(id Id, softRemove bool) {
	_timer, found := m.timers[id]
	if !found {
		return
	}

	if softRemove {
		_timer.remove = true
		return
	}

	heap.Remove(&m.heap, _timer.index)
	delete(m.timers, id)
}
