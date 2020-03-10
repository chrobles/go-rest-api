package cmcapiclient

import (
	"math"
	"testing"
)

func TestConfigure(t *testing.T) {
	got := math.Abs(-1)
	if got != 1 {
		t.Errorf("Abs(-1) = %f; want 1", got)
	}
}
