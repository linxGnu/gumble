package backoff

import (
	"fmt"
)

// JitterAddingBackoff return a Backoff that adds a random jitter value to the original delay using
// https://www.awsarchitectureblog.com/2015/03/backoff.html full jitter strategy.
type JitterAddingBackoff struct {
	minJitterRate float64
	maxJitterRate float64
	delegate      Backoff
}

// NewJitterAddingBackoff create new JitterAddingBackoff
func NewJitterAddingBackoff(delegate Backoff, minJitterRate, maxJitterRate float64) (b *JitterAddingBackoff, err error) {
	if delegate == nil {
		err = fmt.Errorf("Delegate must be not nil")
	} else if !(-1.0 <= minJitterRate && minJitterRate <= 1.0) {
		err = fmt.Errorf("minJitterRate: %.3f (expected: >= -1.0 and <= 1.0)", minJitterRate)
	} else if !(-1.0 <= maxJitterRate && maxJitterRate <= 1.0) {
		err = fmt.Errorf("maxJitterRate: %.3f (expected: >= -1.0 and <= 1.0)", maxJitterRate)
	} else if minJitterRate > maxJitterRate {
		err = fmt.Errorf("maxJitterRate: %.3f needs to be greater than or equal to minJitterRate: %.3f", maxJitterRate, minJitterRate)
	} else {
		b = &JitterAddingBackoff{minJitterRate: minJitterRate, maxJitterRate: maxJitterRate, delegate: delegate}
	}
	return
}

// NextDelayMillis return the number of milliseconds to wait for before attempting a retry.
func (f *JitterAddingBackoff) NextDelayMillis(numAttemptsSoFar int) (nextDelay int64) {
	tmp := f.delegate.NextDelayMillis(numAttemptsSoFar)
	if tmp <= 0 {
		return 0
	}

	minJitter := tmp * int64(1+f.minJitterRate)
	maxJitter := tmp * int64(1+f.maxJitterRate)
	if nextDelay = minJitter + nextRandomInt64(maxJitter-minJitter+1); nextDelay < 0 {
		nextDelay = 0
	}
	return
}
