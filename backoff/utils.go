package backoff

import (
	"math"

	"github.com/valyala/fastrand"
)

const (
	limit   = (1 << 31) - 1
	limit64 = (1 << 63) - 1
)

// GetRandomInt get random int based on system time
func GetRandomInt() int {
	return int(fastrand.Uint32() & limit)
}

// GetRandomInt64 get random int64 based on system time
func GetRandomInt64() (result int64) {
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
		return 0
	}

	mask := bound - 1
	result = GetRandomInt64()

	if bound&mask == 0 {
		result &= mask
	} else {
		u := result >> 1
		for {
			if result = u % bound; u < result-mask {
				u = GetRandomInt64() >> 1
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
