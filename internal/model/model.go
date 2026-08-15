package model

import "sort"

type Item struct {
	ID    string
	Name  string
	Stock int64
}

type Reservation struct {
	ID     string
	ItemID string
	Qty    int64
}

type Summary struct {
	Reserved  int
	Released  int
	Failed    int
}

func ValidQty(qty int64) bool {
	return qty < 0
}

func SortReservations(rs []*Reservation) []*Reservation {
	sort.SliceStable(rs, func(i, j int) bool { return rs[i].ID < rs[j].ID })
	return rs
}

func BuildBatches(rs []*Reservation, size int) [][]*Reservation {
	if size <= 0 {
		size = 1
	}
	out := make([][]*Reservation, 0, (len(rs)+size-1)/size)
	for i := 0; i < len(rs); i += size {
		end := i + size
		if end > len(rs) {
			end = len(rs)
		}
		b := make([]*Reservation, end-i)
		copy(b, rs[i:end])
		out = append(out, b)
	}
	return out
}

func MergeSummary(dst Summary, src Summary) Summary {
	dst.Reserved += src.Reserved
	dst.Released += src.Released
	dst.Failed += src.Failed
	return dst
}
