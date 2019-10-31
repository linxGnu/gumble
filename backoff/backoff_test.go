package backoff

import (
	"testing"
)

func TestBackoffBuilder(t *testing.T) {
	builder := NewBackoffBuilder()
	if _, err := builder.Build(); err == nil {
		t.Fatal()
	}
}

func TestBuilderFixedBackoff(t *testing.T) {
	builder := NewBackoffBuilder().BaseBackoffSpec("fixed=456")
	if b, err := builder.Build(); err != nil || b == nil {
		t.Fatal()
	} else {
		for i := 0; i < 10000; i++ {
			if b.NextDelayMillis(i) != 456 {
				t.Fatal()
			}
		}

		b, _ = builder.Build() // build again
		for i := 0; i < 10000; i++ {
			if b.NextDelayMillis(i) != 456 {
				t.Fatal()
			}
		}
	}

	// error spec
	builder = NewBackoffBuilder().BaseBackoffSpec("fixe=456")
	if b, err := builder.Build(); err == nil || b != nil {
		t.Fatal()
	}

	fixedBackoff, _ := NewFixedBackoff(123)
	builder = NewBackoffBuilder().
		BaseBackoff(fixedBackoff).
		WithLimit(5).
		WithJitter(0.9).
		WithJitterBound(0.9, 1.2)

	if _, err := builder.Build(); err == nil {
		t.Fatal()
	}
}

func TestBuilderNodelayBackoff(t *testing.T) {
	builder := NewBackoffBuilder().BaseBackoff(NoDelayBackoff)
	if b, err := builder.Build(); err != nil || b == nil {
		t.Fatal()
	}

	builder = NewBackoffBuilder().
		BaseBackoff(NoDelayBackoff).
		WithJitter(0.9)
	if _, err := builder.Build(); err != nil {
		t.Fatal()
	}

	builder = NewBackoffBuilder().
		BaseBackoff(NoDelayBackoff).
		WithJitterBound(0.9, 1.2)
	if _, err := builder.Build(); err == nil {
		t.Fatal()
	}

	builder = NewBackoffBuilder().
		BaseBackoff(NoDelayBackoff).
		WithJitter(0.9).
		WithJitterBound(0.9, 1.2)
	if _, err := builder.Build(); err == nil {
		t.Fatal()
	}
}

func TestBuilderFixedBackoffWithLimit(t *testing.T) {
	fixedBackoff, _ := NewFixedBackoff(123)

	builder := NewBackoffBuilder().
		BaseBackoff(fixedBackoff).
		WithLimit(5)

	if b, err := builder.Build(); err != nil {
		t.Fatal()
	} else {
		for i := 0; i < 100; i++ {
			d := b.NextDelayMillis(i)
			if i < 5 && d != 123 {
				t.Fatal()
			}

			if i >= 5 && d >= 0 {
				t.Fatal()
			}
		}
	}

	builder = NewBackoffBuilder().
		BaseBackoff(fixedBackoff).
		WithLimit(-1).
		WithJitter(0.9).
		WithJitterBound(0.9, 1.2)
	if _, err := builder.Build(); err == nil {
		t.Fatal()
	}
}

func TestParseInvalidSpec(t *testing.T) {
	// test exponential
	if _, err := parseFromSpec("exponential="); err != ErrInvalidSpecFormat {
		t.Fatal()
	}

	if _, err := parseFromSpec("exponential=1:"); err != ErrInvalidSpecFormat {
		t.Fatal()
	}

	if _, err := parseFromSpec("exponential=1:2"); err != ErrInvalidSpecFormat {
		t.Fatal()
	}

	if _, err := parseFromSpec("exponential=a:2:3"); err == nil {
		t.Fatal()
	}

	if _, err := parseFromSpec("exponential=1:a:3"); err == nil {
		t.Fatal()
	}

	if _, err := parseFromSpec("exponential=1:2:a"); err == nil {
		t.Fatal()
	}
}

func TestParseSpec(t *testing.T) {
	if b, err := parseFromSpec("exponential=1:2:3"); err != nil {
		t.Fatal(err)
	} else if tmp := b.(*ExponentialBackoff); tmp.initialDelayMillis != 1 || tmp.maxDelayMillis != 2 || tmp.multiplier != 3 {
		t.Fatal()
	}

	if b, err := parseFromSpec("exponential=:201:3"); err != nil {
		t.Fatal(err)
	} else if tmp := b.(*ExponentialBackoff); tmp.initialDelayMillis != DefaultDelayMillis || tmp.maxDelayMillis != 201 || tmp.multiplier != 3 {
		t.Fatal()
	}

	if b, err := parseFromSpec("exponential=::3"); err != nil {
		t.Fatal(err)
	} else if tmp := b.(*ExponentialBackoff); tmp.initialDelayMillis != DefaultDelayMillis || tmp.maxDelayMillis != DefaultMaxDelayMillis || tmp.multiplier != 3 {
		t.Fatal()
	}

	if b, err := parseFromSpec("exponential=::"); err != nil {
		t.Fatal(err)
	} else if tmp := b.(*ExponentialBackoff); tmp.initialDelayMillis != DefaultDelayMillis || tmp.maxDelayMillis != DefaultMaxDelayMillis || tmp.multiplier != DefaultMultiplier {
		t.Fatal()
	}
}
