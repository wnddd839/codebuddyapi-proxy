package gateway

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/wnddd839/codebuddy-proxy/internal/models"
)

func TestModelsFlightCoalescesConcurrentCalls(t *testing.T) {
	var g modelsFlight
	var calls atomic.Int64
	var wg sync.WaitGroup
	for range 16 {
		wg.Go(func() {
			_, err, shared := g.Do("k", func() (models.ListResult, error) {
				calls.Add(1)
				time.Sleep(50 * time.Millisecond)
				return models.ListResult{OK: true, Models: []models.Model{{ID: "auto"}}}, nil
			})
			if err != nil {
				t.Errorf("Do: %v", err)
			}
			_ = shared
		})
	}
	wg.Wait()
	if got := calls.Load(); got != 1 {
		t.Fatalf("upstream calls=%d want 1", got)
	}
}
