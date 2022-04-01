package slru_test

import (
	"fmt"

	slru "github.com/floatdrop/slru"
)

func ExampleSLRU() {
	cache := slru.New[string, int](256)

	cache.Set("Hello", 5)

	if e := cache.Get("Hello"); e != nil {
		fmt.Println(*e)
		// Output: 5
	}
}
