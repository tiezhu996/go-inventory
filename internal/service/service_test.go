package service

import (
	"context"
	"errors"
	"testing"

	"inventory/internal/config"
	"inventory/internal/store"
)

func newSvc(t *testing.T) (*store.Store, *Service) {
	s := store.New()
	cfg := config.Load()
	return s, New(s, cfg)
}

func TestReserveFlow(t *testing.T) {
	_, svc := newSvc(t)
	if _, err := svc.CreateItem("a", "apple", 100); err != nil {
		t.Fatal(err)
	}
	if err := svc.Reserve(context.Background(), "a", 30); err != nil {
		t.Fatal(err)
	}
	if err := svc.Reserve(context.Background(), "a", 2000); !errors.Is(err, ErrInvalidQty) {
		t.Fatalf("over limit err=%v", err)
	}
	st, _ := svc.Stock("a")
	if st != 70 {
		t.Fatalf("stock=%d", st)
	}
}

func TestWrapping(t *testing.T) {
	_, svc := newSvc(t)
	if err := svc.Reserve(context.Background(), "nope", 1); !errors.Is(err, store.ErrItemNotFound) {
		t.Fatalf("errors.Is=false err=%v", err)
	}
	if _, err := svc.Stock("nope"); !errors.Is(err, store.ErrItemNotFound) {
		t.Fatalf("errors.Is=false err=%v", err)
	}
}

func TestListPagesOrder(t *testing.T) {
	_, svc := newSvc(t)
	_, _ = svc.CreateItem("a", "apple", 100)
	_ = svc.Reserve(context.Background(), "a", 10)
	_ = svc.Reserve(context.Background(), "a", 20)
	pages := svc.ListPages()
	if len(pages) == 0 || pages[0][0].ID != "r-a-10-1" {
		t.Fatalf("pages=%v", pages)
	}
}
