package backoff

import (
	"math"

	"github.com/valyala/fastrand"
)

const (
	limit64 = (1 << 63) - 1
)

func randomInt64() (result int64) {
	result |= (int64(fastrand.Uint32()) << 32) & limit64
	result |= int64(fastrand.Uint32())
	return
}

func saturatedMultiply(left int64, right float64) int64 {
	if tmp := float64(left) * right; tmp < math.MaxInt64 {
		return int64(tmp)
	}
	return math.MaxInt64
}

func nextRandomInt64(bound int64) (result int64) {
	if bound <= 0 {
		return bound
	}

	mask := bound - 1
	result = randomInt64()

	if bound&mask == 0 {
		result &= mask
	} else {
		u := result >> 1
		for {
			if result = u % bound; u < result-mask {
				u = randomInt64() >> 1
			} else {
				break
			}
		}
	}

	if result == 0 {
		result = 1
	}

	return
}
