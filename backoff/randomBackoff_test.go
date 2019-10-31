package backoff

import "testing"

func TestRandomBackoff(t *testing.T) {
	if _, err := NewRandomBackoff(-1, 2); err == nil {
		t.FailNow()
	}

	if _, err := NewRandomBackoff(3, 2); err == nil {
		t.FailNow()
	}

	if r, err := NewRandomBackoff(1000, 1000); err != nil {
		t.Error(err)
		t.FailNow()
	} else if d := r.NextDelayMillis(1); d != 1000 {
		t.FailNow()
	}

	if r, err := NewRandomBackoff(1000, 1200); err != nil {
		t.Error(err)
		t.FailNow()
	} else {
		for i := 0; i < 1000; i++ {
			if d := r.NextDelayMillis(i); d < 1000 || d > 1200 {
				t.FailNow()
			}
		}
	}
}
