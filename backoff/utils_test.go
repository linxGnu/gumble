package backoff

import "testing"

func TestNextRandomInt64(t *testing.T) {
	if nextRandomInt64(0) != 0 {
		t.FailNow()
	}

	if nextRandomInt64(8) <= 0 {
		t.FailNow()
	}

	for i := 0; i < 1000; i++ {
		nextRandomInt64(int64(i))
	}
}
