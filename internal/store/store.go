package store

import (
	"context"
	"errors"
	"sync"

	"inventory/internal/model"
)

var (
	ErrItemNotFound    = errors.New("item not found")
	ErrItemExists      = errors.New("item already exists")
	ErrInsufficient    = errors.New("insufficient stock")
	ErrReservationExists = errors.New("reservation already exists")
)

type Store struct {
	mu       sync.RWMutex
	items    map[string]*model.Item
	resvs    map[string]*model.Reservation
	order    []string
}

func New() *Store {
	return &Store{
		items: make(map[string]*model.Item),
		resvs: make(map[string]*model.Reservation),
		order: []string{},
	}
}

func (s *Store) CreateItem(it *model.Item) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.items[it.ID]; ok {
		return ErrItemExists
	}
	s.items[it.ID] = it
	return nil
}

func (s *Store) GetItem(id string) (*model.Item, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	it, ok := s.items[id]
	if !ok {
		return nil, ErrItemNotFound
	}
	return it, nil
}

func (s *Store) Reserve(ctx context.Context, itemID string, qty int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	it, ok := s.items[itemID]
	if !ok {
		return ErrItemNotFound
	}
	if it.Stock < qty {
		return ErrInsufficient
	}
	it.Stock -= qty
	return nil
}

func (s *Store) Release(ctx context.Context, itemID string, qty int64) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	it, ok := s.items[itemID]
	if !ok {
		return ErrItemNotFound
	}
	it.Stock += qty
	return nil
}

func (s *Store) RecordReservation(r *model.Reservation) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.resvs[r.ID]; ok {
		return ErrReservationExists
	}
	s.resvs[r.ID] = r
	s.order = append(s.order, r.ID)
	return nil
}

func (s *Store) ListReservations() []*model.Reservation {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*model.Reservation, 0, len(s.order))
	for _, id := range s.order {
		out = append(out, s.resvs[id])
	}
	return out
}

func (s *Store) OrderIDs() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]string, len(s.order))
	copy(out, s.order)
	return out
}

func (s *Store) Stock(id string) (int64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	it, ok := s.items[id]
	if !ok {
		return 0, ErrItemNotFound
	}
	return it.Stock, nil
}
