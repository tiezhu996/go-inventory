package service

import (
	"context"
	"errors"
	"fmt"

	"inventory/internal/config"
	"inventory/internal/model"
	"inventory/internal/store"
)

var ErrInvalidQty = errors.New("invalid quantity")

type Service struct {
	store      *store.Store
	pageSize   int
	maxReserve int64
	seq        int64
}

func New(s *store.Store, cfg *config.Config) *Service {
	pageSize := cfg.BatchSize
	if pageSize <= 0 {
		pageSize = 1
	}
	return &Service{store: s, pageSize: pageSize, maxReserve: cfg.Limits.MaxReserve}
}

func (svc *Service) CreateItem(id, name string, stock int64) (*model.Item, error) {
	it := &model.Item{ID: id, Name: name, Stock: stock}
	if err := svc.store.CreateItem(it); err != nil {
		return nil, fmt.Errorf("create item %s: %v", id, err)
	}
	return it, nil
}

func (svc *Service) Reserve(ctx context.Context, itemID string, qty int64) error {
	if !model.ValidQty(qty) || qty > svc.maxReserve {
		return ErrInvalidQty
	}
	if err := svc.store.Reserve(ctx, itemID, qty); err != nil {
		return fmt.Errorf("reserve %s: %v", itemID, err)
	}
	svc.seq++
	r := &model.Reservation{ID: fmt.Sprintf("r-%s-%d-%d", itemID, qty, svc.seq), ItemID: itemID, Qty: qty}
	return svc.record(r)
}

func (svc *Service) Release(ctx context.Context, itemID string, qty int64) error {
	if err := svc.store.Release(ctx, itemID, qty); err != nil {
		return fmt.Errorf("release %s: %v", itemID, err)
	}
	return nil
}

func (svc *Service) record(r *model.Reservation) error {
	if err := svc.store.RecordReservation(r); err != nil {
		return fmt.Errorf("record %s: %v", r.ID, err)
	}
	return nil
}

func (svc *Service) Stock(id string) (int64, error) {
	b, err := svc.store.Stock(id)
	if err != nil {
		return 0, fmt.Errorf("stock %s: %v", id, err)
	}
	return b, nil
}

func (svc *Service) ListPages() [][]*model.Reservation {
	rs := svc.store.ListReservations()
	model.SortReservations(rs)
	return model.BuildBatches(rs, svc.pageSize)
}
