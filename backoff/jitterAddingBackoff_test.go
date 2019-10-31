package backoff

import (
	"testing"
)

func TestJitterAddingBackoff(t *testing.T) {
	if _, err := NewJitterAddingBackoff(nil, 0, 0); err == nil {
		t.FailNow()
	}

	if _, err := NewJitterAddingBackoff(NoDelayBackoff, -1.1, 1); err == nil {
		t.FailNow()
	}

	if _, err := NewJitterAddingBackoff(NoDelayBackoff, 1.1, 1); err == nil {
		t.FailNow()
	}

	if _, err := NewJitterAddingBackoff(NoDelayBackoff, 0.9, -1.1); err == nil {
		t.FailNow()
	}

	if _, err := NewJitterAddingBackoff(NoDelayBackoff, 0.9, 1.1); err == nil {
		t.FailNow()
	}

	if _, err := NewJitterAddingBackoff(NoDelayBackoff, 0.9, 0.5); err == nil {
		t.FailNow()
	}

	// fake backoff
	if b, err := NewJitterAddingBackoff(&FixedBackoff{delayMillis: -1}, 0.5, 0.9); err != nil || b == nil {
		t.FailNow()
	} else if b.NextDelayMillis(2) != 0 {
		t.FailNow()
	}

	// real backoff
	if b, err := NewJitterAddingBackoff(&ExponentialBackoff{initialDelayMillis: 100, maxDelayMillis: 1200, multiplier: 1.2},
		0.7, 0.97); err != nil || b == nil {
		t.FailNow()
	} else {
		// fake call
		for i := 0; i < 10000; i++ {
			b.NextDelayMillis(i)
		}
	}
}
