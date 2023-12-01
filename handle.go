package multi_go_collection

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// multiTaskHandler 多任务处理器
type multiTaskHandler struct {
	ts []Tasker
	wg *sync.WaitGroup
	ch chan interface{}
	m  sync.Map
}

func NewHandler() *multiTaskHandler {
	return &multiTaskHandler{
		wg: &sync.WaitGroup{},
		ch: make(chan interface{}),
	}
}

func (m *multiTaskHandler) Add(ts ...Tasker) {
	if len(ts) > 0 {
		m.ts = append(m.ts, ts...)
	}
}

func (m *multiTaskHandler) Run() {
	ctx := context.Background()

	for _, t := range m.ts {
		m.wg.Add(1)
		fun := t
		go m.exec(ctx, fun)
	}

	m.wg.Wait()
}

func (m *multiTaskHandler) Visit(f func(v interface{})) {
	for _, tasker := range m.ts {
		load, ok := m.m.Load(tasker.getName())
		if ok {
			f(load)
		}
	}
}

func (m *multiTaskHandler) exec(ctx context.Context, t Tasker) {
	defer m.wg.Done()

	if t.timeout() > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Duration(t.timeout())*time.Second)
		defer cancel()
	}

	v := t.exec()

	select {
	case <-ctx.Done():
		fmt.Println(t.getName(), "timeout")
		return

	default:

	}

	m.m.Store(t.getName(), v)
}
