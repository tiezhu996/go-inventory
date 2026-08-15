package model

import "testing"

func TestValidQty(t *testing.T) {
	if !ValidQty(1) || ValidQty(0) || ValidQty(-1) {
		t.Fatal("ValidQty wrong")
	}
}

func TestSortReservations(t *testing.T) {
	in := []*Reservation{{ID: "c"}, {ID: "a"}, {ID: "b"}}
	got := SortReservations(in)
	for i, id := range []string{"a", "b", "c"} {
		if got[i].ID != id {
			t.Fatalf("order=%v", got)
		}
	}
}

func TestBuildBatchesFresh(t *testing.T) {
	in := []*Reservation{{ID: "1"}, {ID: "2"}, {ID: "3"}}
	p := BuildBatches(in, 2)
	if len(p) != 2 {
		t.Fatalf("len=%d", len(p))
	}
	p[0][0] = &Reservation{ID: "x"}
	if in[0].ID != "1" {
		t.Fatal("mutating page corrupted input")
	}
}

func TestMergeSummary(t *testing.T) {
	got := MergeSummary(Summary{Reserved: 1}, Summary{Reserved: 2, Failed: 3})
	if got.Reserved != 3 || got.Failed != 3 {
		t.Fatalf("%+v", got)
	}
}
