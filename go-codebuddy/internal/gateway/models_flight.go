package gateway

import (
	"sync"

	"github.com/wnddd839/codebuddy-proxy/internal/models"
)

// modelsFlight 合并同一 cache key 的并发模型列表请求，不引入 golang.org/x/sync。
type modelsFlight struct {
	mu sync.Mutex
	m  map[string]*modelsCall
}

type modelsCall struct {
	wg sync.WaitGroup

	val models.ListResult
	err error
}

func (g *modelsFlight) Do(key string, fn func() (models.ListResult, error)) (models.ListResult, error, bool) {
	g.mu.Lock()
	if g.m == nil {
		g.m = map[string]*modelsCall{}
	}
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		c.wg.Wait()
		return c.val, c.err, true
	}
	c := &modelsCall{}
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()

	c.val, c.err = fn()
	c.wg.Done()

	g.mu.Lock()
	delete(g.m, key)
	g.mu.Unlock()
	return c.val, c.err, false
}
