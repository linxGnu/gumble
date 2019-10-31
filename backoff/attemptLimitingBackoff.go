package backoff

import (
	"fmt"
)

// AttemptLimitingBackoff a backoff which limits the number of attempts up to the specified value.
type AttemptLimitingBackoff struct {
	delegate Backoff
	limit    int
}

// NewAttemptLimitingBackoff create new AttemptLimitingBackoff
func NewAttemptLimitingBackoff(delegate Backoff, limit int) (b *AttemptLimitingBackoff, err error) {
	if delegate == nil {
		err = fmt.Errorf("Delegate must be not nil")
	} else if limit <= 0 {
		err = fmt.Errorf("maxAttempts: %d (expected: > 0)", limit)
	} else {
		b = &AttemptLimitingBackoff{delegate: delegate, limit: limit}
	}
	return
}

// NextDelayMillis return the number of milliseconds to wait for before attempting a retry.
func (f *AttemptLimitingBackoff) NextDelayMillis(numAttemptsSoFar int) int64 {
	if numAttemptsSoFar >= f.limit {
		return -1
	}
	return f.delegate.NextDelayMillis(numAttemptsSoFar)
}
