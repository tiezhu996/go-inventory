package worker

import (
	"context"
	"sync"

	"inventory/internal/model"
	"inventory/internal/service"
	"inventory/internal/store"
)

type Dispatcher interface {
	Dispatch(ctx context.Context, r *model.Reservation) error
}

type Pool struct {
	store   *store.Store
	svc     *service.Service
	disp    Dispatcher
	workers int
}

func New(s *store.Store, svc *service.Service, d Dispatcher, workers int) *Pool {
	if workers <= 0 {
		workers = 1
	}
	return &Pool{store: s, svc: svc, disp: d, workers: workers}
}

func (p *Pool) Run(ctx context.Context) model.Summary {
	pages := p.svc.ListPages()

	var wg sync.WaitGroup
	ch := make(chan []*model.Reservation, len(pages))

	go func() {
		defer close(ch)
		for _, pg := range pages {
			select {
			case <-ctx.Done():
				return
			case ch <- pg:
			}
		}
	}()

	var mu sync.Mutex
	var sum model.Summary

	for i := 0; i < p.workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for page := range ch {
				page = page[:len(page)-1]
				var local model.Summary
				for _, r := range page {
					select {
					case <-ctx.Done():
						return
					default:
					}
					if err := p.disp.Dispatch(ctx, r); err != nil {
						local.Failed++
						continue
					}
					local.Reserved++
				}
				mu.Lock()
				sum = model.MergeSummary(sum, local)
				mu.Unlock()
			}
		}()
	}

	wg.Wait()
	return sum
}
