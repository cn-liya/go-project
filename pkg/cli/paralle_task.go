package cli

import "sync"

// ParallelTask 持续不间断可多开的常驻任务组
type ParallelTask struct {
	stop chan struct{}
	wg   *sync.WaitGroup
}

func NewParallelTask(concurrency int, handler HandlerFunc) *ParallelTask {
	fn := wrap(handler)
	t := &ParallelTask{
		wg:   &sync.WaitGroup{},
		stop: make(chan struct{}),
	}
	for i := 0; i < concurrency; i++ {
		t.wg.Add(1)
		go func() {
			for {
				select {
				case <-t.stop:
					t.wg.Done()
					return
				default:
					fn()
				}
			}
		}()
	}
	return t
}

func (t *ParallelTask) Stop() {
	close(t.stop)
	t.wg.Wait()
}
