# slru
[![Go Reference](https://pkg.go.dev/badge/github.com/floatdrop/slru.svg)](https://pkg.go.dev/github.com/floatdrop/slru)
[![CI](https://github.com/floatdrop/slru/actions/workflows/ci.yml/badge.svg)](https://github.com/floatdrop/slru/actions/workflows/ci.yml)
![Coverage](https://img.shields.io/badge/Coverage-44.4%25-yellow)
[![Go Report Card](https://goreportcard.com/badge/github.com/floatdrop/slru)](https://goreportcard.com/report/github.com/floatdrop/slru)

Thread safe GoLang S(2)LRU cache.

## Example

```go
import (
	"fmt"

	slru "github.com/floatdrop/slru"
)

func main() {
	cache := slru.New[string, int](256)

	cache.Set("Hello", 5)

	if e := cache.Get("Hello"); e != nil {
		fmt.Println(*e)
		// Output: 5
	}
}
```

## Benchmarks

```
floatdrop/slru:
	BenchmarkSLRU_Rand-8   	 5600960	       206.9 ns/op	      44 B/op	       3 allocs/op
	BenchmarkSLRU_Freq-8   	 5927858	       201.0 ns/op	      43 B/op	       3 allocs/op
```