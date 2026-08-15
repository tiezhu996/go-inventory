package config

import "testing"

func TestLoad(t *testing.T) {
	t.Setenv("INV_WORKERS", "")
	t.Setenv("INV_BATCH_SIZE", "")
	c := Load()
	if c.Workers != 2 || c.BatchSize != 2 || c.Limits == nil || c.Limits.MaxReserve != 1000 {
		t.Fatalf("%+v", c)
	}
}
