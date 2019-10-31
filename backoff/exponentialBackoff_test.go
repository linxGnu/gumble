package backoff

import (
	"math"
	"testing"
)

func TestExponentialBackoff(t *testing.T) {
	backoff, _ := NewExponentialBackoff(100, 2000, 1.3)
	for i := 1; i < 100; i++ {
		backoff.NextDelayMillis(i)
	}

	if _, err := NewExponentialBackoff(-1, 12, 3); err == nil {
		t.FailNow()
	}

	if _, err := NewExponentialBackoff(3, 2, 3); err == nil {
		t.FailNow()
	}

	if _, err := NewExponentialBackoff(3, 12, 0.3); err == nil {
		t.FailNow()
	}

	// fake
	if saturatedMultiply(3, float64(math.MaxInt64)) != math.MaxInt64 {
		t.FailNow()
	}
}
