package backoff

import "testing"

func TestAttemptLimitingBackoff(t *testing.T) {
	if _, err := NewAttemptLimitingBackoff(nil, 1); err == nil {
		t.FailNow()
	}

	fixedBackoff, _ := NewFixedBackoff(123)
	if _, err := NewAttemptLimitingBackoff(fixedBackoff, 0); err == nil {
		t.FailNow()
	}

	b, err := NewAttemptLimitingBackoff(fixedBackoff, 3)
	if err != nil || b == nil {
		t.FailNow()
	}

	for i := 0; i < 5; i++ {
		if next := b.NextDelayMillis(i); !(i < 3 && next == 123) && !(i >= 3 && next == -1) {
			t.FailNow()
		}
	}
}
