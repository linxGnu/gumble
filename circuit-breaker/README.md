# circuit breaker

Inspired by [line/armeria](https://github.com/line/armeria) circuit breaker.

# Usage

```go
package main

import (
    cbreaker "github.com/linxGnu/gumble/circuit-breaker"
)

func main() {
    cb := cbreaker.NewCircuitBreakerBuilder().
                        SetTicker(cbreaker.SystemTicker).
                        ... // other settings
                        Build()
    ...

    if cb.CanRequest() {
        err := makeRequest()
        if err != nil {
            cb.OnFailure()
        } else {
            cb.OnSuccess()
        }
    }
}
```
