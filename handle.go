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
	m  []interface{}
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

	go func() {
		defer close(m.ch)
		m.wg.Wait()
	}()

	for d := range m.ch {
		m.m = append(m.m, d)
	}
}

func (m *multiTaskHandler) Visit(f func(v interface{})) {
	for _, data := range m.m {
		f(data)
	}
}

func (m *multiTaskHandler) exec(ctx context.Context, t Tasker) {
	defer m.wg.Done()

	if t.timeout() > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Duration(t.timeout())*time.Second)
		defer cancel()
	}

	select {
	case <-ctx.Done():
		fmt.Println(t.getName(), "timeout")
		return

	case m.ch <- t.exec():

	}

}
