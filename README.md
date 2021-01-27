##一个高效的、无依赖的、无内部协程的、计时器系统，在需要的协程自行调度
###默认用时间堆，有大量短时计时器需求，可以开启时间轮
## Quick Start
```golang
func main() {
    ticker := time.NewTicker(time.Millisecond * 100)
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
        case <-ticker.C:
            timerMgr.Update(time.Now().UnixNano())
        }   
    }
}
```
