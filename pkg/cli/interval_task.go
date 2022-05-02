package cli

import "time"

// IntervalTask 间隔时间执行的任务
type IntervalTask struct {
	stop chan struct{}
	done chan struct{}
}

func NewIntervalTask(duration time.Duration, handler HandlerFunc) *IntervalTask {
	fn := wrap(handler)
	t := &IntervalTask{
		stop: make(chan struct{}),
		done: make(chan struct{}),
	}
	go func() {
		tick := time.NewTicker(duration)
		for {
			select {
			case <-t.stop:
				tick.Stop()
				close(t.done)
				return
			case <-tick.C:
				fn()
			}
		}
	}()
	return t
}

func (t *IntervalTask) Stop() {
	close(t.stop)
	<-t.done
}
