# longadder

Thread-safe, high performance, contention-aware `LongAdder` and `DoubleAdder` for Go, inspired by OpenJDK9.
Beside JDK-based `LongAdder` and `DoubleAdder`, library includes other adders for various use.

# Usage

## JDKAdder (recommended)

```go
package main

import (
	"fmt"
	"time"

	ga "github.com/linxGnu/gumble/adder"
)

func main() {
	// or ga.DefaultAdder() which uses jdk adder as default
	adder := ga.NewLongAdder(ga.JDKAdderType) 

	for i := 0; i < 100; i++ {
		go func() {
			adder.Add(123)
		}()
	}

	time.Sleep(3 * time.Second)

	// get total added value
	fmt.Println(adder.Sum()) 
}
```

## RandomCellAdder

* A `LongAdder` with simple strategy of preallocating atomic cell and select random cell for update.
* Slower than JDK LongAdder but 1.5-2x faster than AtomicAdder on contention.
* Consume ~1KB to store cells.

```
adder := ga.NewLongAdder(ga.RandomCellAdderType)
```

## AtomicAdder

* A `LongAdder` based on atomic variable. All routines share this variable.

```go
adder := ga.NewLongAdder(ga.AtomicAdderType)
```

## MutexAdder

* A `LongAdder` based on mutex. All routines share same value and mutex.

```go
adder := ga.NewLongAdder(ga.MutexAdderType)
```

# Benchmark

* Hardware: MacBookPro14,3
* OS: macOS 10.14.6 (18G103)

```
GO111MODULE=""
GOARCH="amd64"
GOBIN=""
GOEXE=""
GOFLAGS=""
GOHOSTARCH="amd64"
GOHOSTOS="darwin"
GONOPROXY=""
GONOSUMDB=""
GOOS="darwin"
GOPATH="/Users/JP22782/workspace/go"
GOPRIVATE=""
GOPROXY="https://proxy.golang.org,direct"
GOROOT="/usr/local/Cellar/go/1.13.3/libexec"
GOSUMDB="sum.golang.org"
GOTMPDIR=""
GCCGO="gccgo"
AR="ar"
CC="clang"
CXX="clang++"
CGO_ENABLED="1"
go.mod"
CGO_CFLAGS="-g -O2"
CGO_CPPFLAGS=""
CGO_CXXFLAGS="-g -O2"
CGO_FFLAGS="-g -O2"
CGO_LDFLAGS="-g -O2"
PKG_CONFIG="pkg-config"
GOGCCFLAGS="-fPIC -m64 -pthread -fno-caret-diagnostics -Qunused-arguments -fmessage-length=0 -fdebug-prefix-map=/var/folders/9m/_p6tsqzn1_d1d0rsrbwb_hb80000gp/T/go-build564245481=/tmp/go-build -gno-record-gcc-switches -fno-common"
```
```
BenchmarkAtomicF64AdderSingleRoutine-8                93          11898692 ns/op               0 B/op          0 allocs/op
BenchmarkJDKF64AdderSingleRoutine-8                   99          11614899 ns/op               0 B/op          0 allocs/op
BenchmarkAtomicF64AdderMultiRoutine-8                  4         352340390 ns/op            5968 B/op         16 allocs/op
BenchmarkJDKF64AdderMultiRoutine-8                    26          43095301 ns/op            2028 B/op          6 allocs/op
BenchmarkAtomicF64AdderMultiRoutineMix-8               3         347352986 ns/op              16 B/op          1 allocs/op
BenchmarkJDKF64AdderMultiRoutineMix-8                 22          51581986 ns/op              86 B/op          1 allocs/op
BenchmarkMutexAdderSingleRoutine-8                    40          28532922 ns/op               0 B/op          0 allocs/op
BenchmarkAtomicAdderSingleRoutine-8                  168           7281029 ns/op               0 B/op          0 allocs/op
BenchmarkRandomCellAdderSingleRoutine-8               40          29172354 ns/op              25 B/op          0 allocs/op
BenchmarkJDKAdderSingleRoutine-8                     100          10555204 ns/op               0 B/op          0 allocs/op
BenchmarkMutexAdderMultiRoutine-8                      4         269586062 ns/op             476 B/op          6 allocs/op
BenchmarkAtomicAdderMultiRoutine-8                    16          70678513 ns/op              16 B/op          1 allocs/op
BenchmarkRandomCellAdderMultiRoutine-8                38          30580824 ns/op              52 B/op          1 allocs/op
BenchmarkJDKAdderMultiRoutine-8                       26          41366055 ns/op              82 B/op          1 allocs/op
BenchmarkMutexAdderMultiRoutineMix-8                   4         277455294 ns/op             352 B/op          4 allocs/op
BenchmarkAtomicAdderMultiRoutineMix-8                 15          71483437 ns/op              16 B/op          1 allocs/op
BenchmarkRandomCellAdderMultiRoutineMix-8             31          34591180 ns/op              49 B/op          1 allocs/op
BenchmarkJDKAdderMultiRoutineMix-8                    22          47938895 ns/op             110 B/op          1 allocs/op
```
