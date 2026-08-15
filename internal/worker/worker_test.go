package worker

import (
	"context"
	"fmt"
	"testing"

	"inventory/internal/config"
	"inventory/internal/model"
	"inventory/internal/service"
	"inventory/internal/store"
)

type okDispatcher struct{}

func (okDispatcher) Dispatch(ctx context.Context, r *model.Reservation) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return nil
}

func newPool(t *testing.T) (*store.Store, *service.Service, *Pool) {
	s := store.New()
	svc := service.New(s, config.Load())
	return s, svc, New(s, svc, okDispatcher{}, 4)
}

func TestRunSummary(t *testing.T) {
	_, svc, p := newPool(t)
	_, _ = svc.CreateItem("a", "apple", 100)
	for i := 0; i < 10; i++ {
		if err := svc.Reserve(context.Background(), "a", 1); err != nil {
			t.Fatal(err)
		}
	}
	sum := p.Run(context.Background())
	if sum.Reserved != 10 {
		t.Fatalf("Reserved=%d want 10", sum.Reserved)
	}
}

type failingDispatcher struct{ failID string }

func (f failingDispatcher) Dispatch(ctx context.Context, r *model.Reservation) error {
	if r.ID == f.failID {
		return fmt.Errorf("fail %s", r.ID)
	}
	return nil
}

func TestRunFailed(t *testing.T) {
	s, svc, _ := newPool(t)
	_, _ = svc.CreateItem("a", "apple", 100)
	_ = svc.Reserve(context.Background(), "a", 10)
	_ = svc.Reserve(context.Background(), "a", 20)
	p := New(s, svc, failingDispatcher{"r-a-10-1"}, 2)
	sum := p.Run(context.Background())
	if sum.Failed != 1 {
		t.Fatalf("Failed=%d want 1", sum.Failed)
	}
	if sum.Reserved != 1 {
		t.Fatalf("Reserved=%d want 1", sum.Reserved)
	}
}

func TestRunCancel(t *testing.T) {
	_, svc, p := newPool(t)
	_, _ = svc.CreateItem("a", "apple", 100)
	_ = svc.Reserve(context.Background(), "a", 10)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	sum := p.Run(ctx)
	if sum.Reserved != 0 {
		t.Fatalf("Reserved=%d want 0", sum.Reserved)
	}
}
