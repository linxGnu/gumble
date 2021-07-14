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

```
GO111MODULE=""
GOARCH="amd64"
GOBIN=""
GOCACHE="/home/gnu/.cache/go-build"
GOENV="/home/gnu/.config/go/env"
GOEXE=""
GOFLAGS=""
GOHOSTARCH="amd64"
GOHOSTOS="linux"
GOINSECURE=""
GOMODCACHE="/home/gnu/workspace/go/pkg/mod"
GONOPROXY=""
GONOSUMDB=""
GOOS="linux"
GOPATH="/home/gnu/workspace/go"
GOPRIVATE=""
GOPROXY="https://proxy.golang.org,direct"
GOROOT="/home/gnu/workspace/goroot"
GOSUMDB="sum.golang.org"
GOTMPDIR=""
GOTOOLDIR="/home/gnu/workspace/goroot/pkg/tool/linux_amd64"
GOVCS=""
GOVERSION="go1.16.6"
GCCGO="gccgo"
AR="ar"
CC="gcc"
CXX="g++"
CGO_ENABLED="1"
GOMOD="/home/gnu/workspace/go/src/git.linecorp.com/LINE-DevOps/go-utils.git/go.mod"
CGO_CFLAGS="-g -O2"
CGO_CPPFLAGS=""
CGO_CXXFLAGS="-g -O2"
CGO_FFLAGS="-g -O2"
CGO_LDFLAGS="-g -O2"
PKG_CONFIG="pkg-config"
GOGCCFLAGS="-fPIC -m64 -pthread -fmessage-length=0 -fdebug-prefix-map=/tmp/go-build1824823233=/tmp/go-build -gno-record-gcc-switches"
```
```
goos: linux
goarch: amd64
pkg: git.linecorp.com/LINE-DevOps/go-utils.git/queue
cpu: AMD Ryzen 9 3950X 16-Core Processor            
PASS
benchmark                              iter      time/iter      bytes alloc              allocs
---------                              ----      ---------      -----------              ------
Benchmark_MutexLinkedQueue_50P50C-32      2   557.08 ms/op    32036812 B/op   1000161 allocs/op
Benchmark_JDKLinkedQueue_50P50C-32        4   250.64 ms/op    32029912 B/op   1500112 allocs/op
Benchmark_LFQueue_50P50C-32               3   399.04 ms/op    68013296 B/op   1500068 allocs/op
Benchmark_MutexLinkedQueue_50P10C-32      2   665.07 ms/op    32000572 B/op   1000011 allocs/op
Benchmark_JDKLinkedQueue_50P10C-32        5   225.47 ms/op    32006699 B/op   1500042 allocs/op
Benchmark_LFQueue_50P10C-32               3   386.33 ms/op    68003088 B/op   1500041 allocs/op
Benchmark_MutexLinkedQueue_10P50C-32      7   162.23 ms/op     6400554 B/op    200010 allocs/op
Benchmark_JDKLinkedQueue_10P50C-32       12    84.28 ms/op     6401472 B/op    300018 allocs/op
Benchmark_LFQueue_10P50C-32              14    76.56 ms/op    13601710 B/op    300023 allocs/op
Benchmark_MutexLinkedQueue_100P-32        3   459.11 ms/op    64000096 B/op   2000003 allocs/op
Benchmark_JDKLinkedQueue_100P-32          2   574.58 ms/op    64000256 B/op   3000003 allocs/op
Benchmark_LFQueue_100P-32                 3   377.33 ms/op   136001592 B/op   3000016 allocs/op
Benchmark_MutexLinkedQueue_100C-32        3   390.75 ms/op    64004176 B/op   2000018 allocs/op
Benchmark_JDKLinkedQueue_100C-32          4   298.37 ms/op    64000016 B/op   3000001 allocs/op
```

