package slru

import (
	"math/rand"
	"testing"
)

func BenchmarkSLRU_Rand(b *testing.B) {
	l := New[int64, int64](8192)

	trace := make([]int64, b.N*2)
	for i := 0; i < b.N*2; i++ {
		trace[i] = rand.Int63() % 32768
	}

	b.ResetTimer()

	var hit, miss int
	for i := 0; i < 2*b.N; i++ {
		if i%2 == 0 {
			l.Set(trace[i], trace[i])
		} else {
			if l.Get(trace[i]) == nil {
				miss++
			} else {
				hit++
			}
		}
	}
	b.Logf("hit: %d miss: %d ratio: %f", hit, miss, float64(hit)/float64(miss))
}

func BenchmarkSLRU_Freq(b *testing.B) {
	l := New[int64, int64](8192)

	trace := make([]int64, b.N*2)
	for i := 0; i < b.N*2; i++ {
		if i%2 == 0 {
			trace[i] = rand.Int63() % 16384
		} else {
			trace[i] = rand.Int63() % 32768
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		l.Set(trace[i], trace[i])
	}
	var hit, miss int
	for i := 0; i < b.N; i++ {
		if l.Get(trace[i]) == nil {
			miss++
		} else {
			hit++
		}
	}
	b.Logf("hit: %d miss: %d ratio: %f", hit, miss, float64(hit)/float64(miss))
}

func TestSLRU_zero(t *testing.T) {
	l := New[int, int](0)
	i := 5

	if l.Len() != 0 {
		t.Errorf("should have 0 length")
	}

	if l.Victim(i) != nil {
		t.Errorf("should have no victims in zero cache")
	}

	if e := l.Set(i, i); e == nil || e.Value != i {
		t.Fatalf("value should be evicted")
	}

	if e := l.Remove(i); e != nil {
		t.Fatalf("value should not be removed")
	}
}

func TestSLRU(t *testing.T) {
	l := NewParams[int, int](2, 3)

	l.Set(1, 1)
	l.Set(2, 2)
	if e := l.Set(3, 3); e == nil || e.Key != 1 {
		t.Fatalf("value should be removed from probation cache")
	}

	l.Get(2) // Promote 2 to protected with Get
	if l.Set(4, 4) != nil {
		t.Fatalf("value should not be removed from probation cache")
	}
	if e := l.Set(5, 5); e == nil || e.Key != 3 {
		t.Fatalf("value should be removed from probation cache")
	}

	if e := l.Get(2); e == nil || *e != 2 {
		t.Fatalf("value should stay in protected cache")
	}

	if e := l.Set(2, 22); e == nil || e.Value != 2 {
		t.Fatalf("value should be updatable in protected cache: %+v", e)
	}

	if l.Set(4, 4) != nil { // Promote 4 to protected with Set
		t.Fatalf("value should be promoted to protected without eviction")
	}
}
