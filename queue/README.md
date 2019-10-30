# queue

High performance and thread-safe queue implementation in Go:
* JDKLinkedQueue: a non-blocking, thread-safe linked list queue ported from OpenJDK ConcurrentLinkedQueue.
* MutexLinkedQueue: thread-safe linked list queue based on mutex.

# Usage

```
package main

import (
    queue "github.com/linxGnu/gumble/queue"
)

func main() {
    q := DefaultQueue() // default using jdk linked queue

    q.Offer(struct{}{}) // push

    polled := q.Poll() // remove and return head queue

    head := q.Peak() // return head queue but not remove
}
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
goos: darwin
goarch: amd64
pkg: github.com/linxGnu/gumble/queue
Benchmark_MutexLinkedQueue_50P50C-8            4         250396229 ns/op        32814560 B/op    1000050 allocs/op
Benchmark_JDKLinkedQueue_50P50C-8             10         102264315 ns/op        32804635 B/op    1500015 allocs/op
Benchmark_LFQueue_50P50C-8                     6         174677691 ns/op        72807312 B/op    1500039 allocs/op
Benchmark_MutexLinkedQueue_50P10C-8            5         230457374 ns/op        32803932 B/op    1000015 allocs/op
Benchmark_JDKLinkedQueue_50P10C-8             10         101427604 ns/op        32804708 B/op    1500023 allocs/op
Benchmark_LFQueue_50P10C-8                     7         147318500 ns/op        72803601 B/op    1500018 allocs/op
Benchmark_MutexLinkedQueue_10P50C-8            9         125443296 ns/op         7203642 B/op     200012 allocs/op
Benchmark_JDKLinkedQueue_10P50C-8             62          19857979 ns/op         7203588 B/op     300010 allocs/op
Benchmark_LFQueue_10P50C-8                    33          33143888 ns/op        15203620 B/op     300016 allocs/op
Benchmark_MutexLinkedQueue_100P-8              5         233750795 ns/op        64001536 B/op    2000018 allocs/op
Benchmark_JDKLinkedQueue_100P-8                8         138373969 ns/op        64000696 B/op    3000008 allocs/op
Benchmark_LFQueue_100P-8                       5         246579063 ns/op        144000449 B/op   3000010 allocs/op
Benchmark_MutexLinkedQueue_100C-8              6         187289416 ns/op        64000320 B/op    2000004 allocs/op
Benchmark_JDKLinkedQueue_100C-8                9         117081262 ns/op        64000016 B/op    3000001 allocs/op
```
