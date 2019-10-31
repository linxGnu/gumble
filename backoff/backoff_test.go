package backoff

import (
	"testing"
)

func TestBackoffBuilder(t *testing.T) {
	builder := NewBackoffBuilder()
	if _, err := builder.Build(); err == nil {
		t.FailNow()
	}
}

func TestBuilderFixedBackoff(t *testing.T) {
	builder := NewBackoffBuilder().BaseBackoffSpec("fixed=456")
	if b, err := builder.Build(); err != nil || b == nil {
		t.FailNow()
	} else {
		for i := 0; i < 10000; i++ {
			if b.NextDelayMillis(i) != 456 {
				t.FailNow()
			}
		}

		b, _ = builder.Build() // build again
		for i := 0; i < 10000; i++ {
			if b.NextDelayMillis(i) != 456 {
				t.FailNow()
			}
		}
	}

	// error spec
	builder = NewBackoffBuilder().BaseBackoffSpec("fixe=456")
	if b, err := builder.Build(); err == nil || b != nil {
		t.FailNow()
	}

	fixedBackoff, _ := NewFixedBackoff(123)
	builder = NewBackoffBuilder().
		BaseBackoff(fixedBackoff).
		WithLimit(5).
		WithJitter(0.9).
		WithJitterBound(0.9, 1.2)

	if _, err := builder.Build(); err == nil {
		t.FailNow()
	}
}

func TestBuilderNodelayBackoff(t *testing.T) {
	builder := NewBackoffBuilder().BaseBackoff(NoDelayBackoff)
	if b, err := builder.Build(); err != nil || b == nil {
		t.FailNow()
	}

	builder = NewBackoffBuilder().
		BaseBackoff(NoDelayBackoff).
		WithJitter(0.9)
	if _, err := builder.Build(); err != nil {
		t.FailNow()
	}

	builder = NewBackoffBuilder().
		BaseBackoff(NoDelayBackoff).
		WithJitterBound(0.9, 1.2)
	if _, err := builder.Build(); err == nil {
		t.FailNow()
	}

	builder = NewBackoffBuilder().
		BaseBackoff(NoDelayBackoff).
		WithJitter(0.9).
		WithJitterBound(0.9, 1.2)
	if _, err := builder.Build(); err == nil {
		t.FailNow()
	}
}

func TestBuilderFixedBackoffWithLimit(t *testing.T) {
	fixedBackoff, _ := NewFixedBackoff(123)

	builder := NewBackoffBuilder().
		BaseBackoff(fixedBackoff).
		WithLimit(5)

	if b, err := builder.Build(); err != nil {
		t.FailNow()
	} else {
		for i := 0; i < 100; i++ {
			d := b.NextDelayMillis(i)
			if i < 5 && d != 123 {
				t.FailNow()
			}

			if i >= 5 && d >= 0 {
				t.FailNow()
			}
		}
	}

	builder = NewBackoffBuilder().
		BaseBackoff(fixedBackoff).
		WithLimit(-1).
		WithJitter(0.9).
		WithJitterBound(0.9, 1.2)
	if _, err := builder.Build(); err == nil {
		t.FailNow()
	}
}

func TestParseSpec(t *testing.T) {
	// test exponential
	if _, err := parseFromSpec("exponential="); err != ErrInvalidSpecFormat {
		t.FailNow()
	}

	if _, err := parseFromSpec("exponential=1:"); err != ErrInvalidSpecFormat {
		t.FailNow()
	}

	if _, err := parseFromSpec("exponential=1:2"); err != ErrInvalidSpecFormat {
		t.FailNow()
	}

	if _, err := parseFromSpec("exponential=a:2:3"); err == nil {
		t.FailNow()
	}

	if _, err := parseFromSpec("exponential=1:a:3"); err == nil {
		t.FailNow()
	}

	if _, err := parseFromSpec("exponential=1:2:a"); err == nil {
		t.FailNow()
	}

	if b, err := parseFromSpec("exponential=1:2:3"); err != nil {
		t.Error(err)
		t.FailNow()
	} else if tmp := b.(*ExponentialBackoff); tmp.initialDelayMillis != 1 || tmp.maxDelayMillis != 2 || tmp.multiplier != 3 {
		t.FailNow()
	}

	if b, err := parseFromSpec("exponential=:201:3"); err != nil {
		t.Error(err)
		t.FailNow()
	} else if tmp := b.(*ExponentialBackoff); tmp.initialDelayMillis != DefaultDelayMillis || tmp.maxDelayMillis != 201 || tmp.multiplier != 3 {
		t.FailNow()
	}

	if b, err := parseFromSpec("exponential=::3"); err != nil {
		t.Error(err)
		t.FailNow()
	} else if tmp := b.(*ExponentialBackoff); tmp.initialDelayMillis != DefaultDelayMillis || tmp.maxDelayMillis != DefaultMaxDelayMillis || tmp.multiplier != 3 {
		t.FailNow()
	}

	if b, err := parseFromSpec("exponential=::"); err != nil {
		t.Error(err)
		t.FailNow()
	} else if tmp := b.(*ExponentialBackoff); tmp.initialDelayMillis != DefaultDelayMillis || tmp.maxDelayMillis != DefaultMaxDelayMillis || tmp.multiplier != DefaultMultiplier {
		t.FailNow()
	}
}
