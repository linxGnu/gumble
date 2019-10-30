package retry

import (
	"fmt"
	"strconv"
	"strings"
)

// Backoff control back off between attempts in a single retry operation.
type Backoff interface {
	// NextDelayMillis return the number of milliseconds to wait for before attempting a retry.
	NextDelayMillis(numAttemptsSoFar int) int64
}

// BackoffBuilder backoff builder
type BackoffBuilder struct {
	layer []interface{}
	base  Backoff
	spec  string
}

type withLimit struct {
	limit int
}

type withJitter struct {
	minJitterRate float64
	maxJitterRate float64
}

// NewBackoffBuilder create new backoff builder.
func NewBackoffBuilder() *BackoffBuilder {
	return &BackoffBuilder{
		layer: make([]interface{}, 0, 4),
	}
}

// BaseBackoffSpec set specification for building base backoff. Then base backoff could be chaining
// WithJitter and WithLimit number of attempts.
//
// This is the format for the specification:
//   // "exponential=initialDelayMillis:maxDelayMillis:multiplier" is for ExponentialBackoff.
//   // multiplier will be 2.0 if it's omitted.
//   // initialDelayMillis will be 200 if its omitted.
//   // maxDelayMillis will be 10000 if its omitted.
//   //
//   // "fixed=delayMillis" is for FixedBackoff. delayMillis will be 200 if its omitted
//   //
//   // "random=minDelayMillis:maxDelayMillis" is for RandomBackoff.
//   // minDelayMillis will be 0 if its omitted.
//   // maxDelayMillis will be 200 if its omitted.
//
// To omit a value, just make it blank but keep separation ':'.
// For example: "exponential=12::3" means initialDelayMillis = 12, maxDelayMillis is default = 10000 and multiplier = 3
func (b *BackoffBuilder) BaseBackoffSpec(spec string) *BackoffBuilder {
	b.spec = spec
	return b
}

// BaseBackoff set base backoff. Base backoff could be chaining
// WithJitter and WithLimit number of attempts.
func (b *BackoffBuilder) BaseBackoff(base Backoff) *BackoffBuilder {
	b.base = base
	return b
}

// WithLimit wrap base backoff with limiting the number of attempts up to the specified value.
func (b *BackoffBuilder) WithLimit(limit int) *BackoffBuilder {
	b.layer = append(b.layer, &withLimit{limit})
	return b
}

// WithJitter wrap a base backoff, adds a random jitter value to the original delay using full jitter strategy.
// ThejitterRate is used to calculate the lower and upper bound of the ultimate delay.
//
// The lower bound will be ((1 - jitterRate) * originalDelay) and the upper bound will be
// ((1 + jitterRate) * originalDelay).
//
// For example, if the delay returned by
// exponentialBackoff(long, long) is 1000 milliseconds and the provided jitter value is 0.3,
// the ultimate backoff delay will be chosen between 1000 * (1 - 0.3) and 1000 * (1 + 0.3)
// by randomer. The rate value should be between 0.0 and 1.0.
func (b *BackoffBuilder) WithJitter(jitterRate float64) *BackoffBuilder {
	b.layer = append(b.layer, &withJitter{minJitterRate: -jitterRate, maxJitterRate: jitterRate})
	return b
}

// WithJitterBound similar to WithJitter but with specific min-maxJitterRate
func (b *BackoffBuilder) WithJitterBound(minJitterRate, maxJitterRate float64) *BackoffBuilder {
	b.layer = append(b.layer, &withJitter{minJitterRate: minJitterRate, maxJitterRate: maxJitterRate})
	return b
}

// Build backoff
func (b *BackoffBuilder) Build() (r Backoff, err error) {
	if r = b.base; r == nil {
		// try to parse base from spec
		if b.spec == "" {
			err = fmt.Errorf("Base Backoff is required. Please provide it by")
			return
		} else if r, err = parseFromSpec(b.spec); err != nil {
			return
		}
	}

	for _, layer := range b.layer {
		switch l := layer.(type) {
		case *withLimit:
			if r, err = NewAttemptLimitingBackoff(r, l.limit); err != nil {
				return
			}
		case *withJitter:
			if r, err = NewJitterAddingBackoff(r, l.minJitterRate, l.maxJitterRate); err != nil {
				return
			}
		}
	}

	return
}

func parseFromSpec(spec string) (r Backoff, err error) {
	index := strings.Index(spec, "=")
	if index < 0 {
		err = ErrInvalidSpecFormat
		return
	}

	// get key and values
	key, values := spec[:index], spec[index+1:]
	switch key {
	case "exponential": // exponential=initialDelayMillis:maxDelayMillis:multiplier
		splited := strings.Split(values, ":")
		if len(splited) != 3 {
			err = ErrInvalidSpecFormat
			return
		}

		initialDelayMillis, maxDelayMillis, multiplier := DefaultInitialDelayMillis, DefaultMaxDelayMillis, DefaultMultiplier
		if splited[0] != "" {
			if initialDelayMillis, err = strconv.ParseInt(splited[0], 10, 64); err != nil {
				return
			}
		}
		if splited[1] != "" {
			if maxDelayMillis, err = strconv.ParseInt(splited[1], 10, 64); err != nil {
				return
			}
		}
		if splited[2] != "" {
			if multiplier, err = strconv.ParseFloat(splited[2], 64); err != nil {
				return
			}
		}

		r, err = NewExponentialBackoff(initialDelayMillis, maxDelayMillis, multiplier)
		return
	case "fixed": // fixed=delayMillis
		delayMillis := DefaultDelayMillis

		if values != "" {
			if delayMillis, err = strconv.ParseInt(values, 10, 64); err != nil {
				return
			}
		}

		r, err = NewFixedBackoff(delayMillis)
		return
	case "random": // random=minDelayMillis:maxDelayMillis
		splited := strings.Split(values, ":")
		if len(splited) != 2 {
			err = ErrInvalidSpecFormat
			return
		}

		minDelayMillis, maxDelayMillis := DefaultMinDelayMillis, DefaultMaxDelayMillis
		if splited[0] != "" {
			if minDelayMillis, err = strconv.ParseInt(splited[0], 10, 64); err != nil {
				return
			}
		}
		if splited[1] != "" {
			if maxDelayMillis, err = strconv.ParseInt(splited[1], 10, 64); err != nil {
				return
			}
		}

		r, err = NewRandomBackoff(minDelayMillis, maxDelayMillis)
		return
	}

	err = ErrInvalidSpecFormat
	return
}
