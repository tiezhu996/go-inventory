package store

import (
	"context"
	"testing"

	"inventory/internal/model"
)

func TestItemReserveRelease(t *testing.T) {
	s := New()
	if err := s.CreateItem(&model.Item{ID: "a", Stock: 100}); err != nil {
		t.Fatal(err)
	}
	if err := s.Reserve(context.Background(), "a", 30); err != nil {
		t.Fatal(err)
	}
	if err := s.Reserve(context.Background(), "a", 80); err != ErrInsufficient {
		t.Fatalf("err=%v", err)
	}
	if err := s.Release(context.Background(), "a", 30); err != nil {
		t.Fatal(err)
	}
	if st, _ := s.Stock("a"); st != 100 {
		t.Fatalf("stock=%d", st)
	}
}

func TestReservationRecordOrderFresh(t *testing.T) {
	s := New()
	_ = s.RecordReservation(&model.Reservation{ID: "r1"})
	_ = s.RecordReservation(&model.Reservation{ID: "r2"})
	if err := s.RecordReservation(&model.Reservation{ID: "r1"}); err != ErrReservationExists {
		t.Fatalf("err=%v", err)
	}
	ids := s.OrderIDs()
	ids[0] = "x"
	if s.OrderIDs()[0] != "r1" {
		t.Fatal("OrderIDs aliased")
	}
	if len(s.ListReservations()) != 2 {
		t.Fatal("ListReservations len")
	}
}
