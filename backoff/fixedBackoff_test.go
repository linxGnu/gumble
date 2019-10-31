package backoff

import (
	"testing"
)

func TestFixedBackoff(t *testing.T) {
	if backoff, _ := NewFixedBackoff(12); backoff.delayMillis != 12 || backoff.NextDelayMillis(31) != 12 {
		t.FailNow()
	}

	if backoff, err := NewFixedBackoff(-1); err == nil || backoff != nil {
		t.FailNow()
	}
}
