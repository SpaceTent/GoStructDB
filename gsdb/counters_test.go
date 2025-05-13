package gsdb

import (
	"context"
	"testing"
)

func TestCounters(t *testing.T) {

	New("", nil, context.Background())
	DB.StartCounter("test")
	DB.IncCounter("test")
	if DB.GetCounter("test") != 1 {
		t.Errorf("Expected: %d, got: %d", 1, DB.GetCounter("test"))
	}
}
