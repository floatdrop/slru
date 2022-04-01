package slru

import (
	"github.com/floatdrop/lru"
)

const (
	DefaultProbationRatio = 0.2
)

type SLRU[K comparable, V any] struct {
	probation *lru.LRU[K, V]
	protected *lru.LRU[K, V]
}

// Evicted holds key/value pair that was evicted from cache.
type Evicted[K comparable, V any] struct {
	Key   K
	Value V
}

func (S *SLRU[K, V]) Get(key K) *V {
	if e := S.protected.Get(key); e != nil {
		return e
	}

	if e := S.probation.Get(key); e != nil {
		S.probation.Remove(key)
		if ev := S.protected.Set(key, *e); ev != nil {
			S.probation.Set(ev.Key, ev.Value)
		}
		return e
	}

	return nil
}

func fromLruEvicted[K comparable, V any](e *lru.Evicted[K, V]) *Evicted[K, V] {
	if e == nil {
		return nil
	}

	return &Evicted[K, V]{
		e.Key,
		e.Value,
	}
}

func (S *SLRU[K, V]) Set(key K, value V) *Evicted[K, V] {
	if e := S.protected.Peek(key); e != nil {
		return fromLruEvicted(S.protected.Set(key, value))
	}

	if e := S.probation.Peek(key); e != nil {
		S.probation.Remove(key)
		if ev := S.protected.Set(key, *e); ev != nil {
			S.probation.Set(ev.Key, ev.Value)
		}
		return nil
	}

	return fromLruEvicted(S.probation.Set(key, value))
}

func (S *SLRU[K, V]) Victim(key K) *K {
	if e := S.probation.Peek(key); e != nil {
		return S.protected.Victim()
	}

	return S.probation.Victim()
}

func (S *SLRU[K, V]) Len() int {
	return S.probation.Len() + S.protected.Len()
}

func (S *SLRU[K, V]) Peek(key K) *V {
	if e := S.probation.Peek(key); e != nil {
		return e
	}

	return S.protected.Peek(key)
}

func (S *SLRU[K, V]) Remove(key K) *V {
	if e := S.protected.Remove(key); e != nil {
		return e
	}

	return S.probation.Remove(key)
}

func NewParams[K comparable, V any](probationSize int, protectedSize int) *SLRU[K, V] {
	return &SLRU[K, V]{
		probation: lru.New[K, V](probationSize),
		protected: lru.New[K, V](protectedSize),
	}
}

func New[K comparable, V any](size int) *SLRU[K, V] {
	probationSize := int(DefaultProbationRatio * float64(size))
	protectedSize := size - probationSize
	return NewParams[K, V](
		probationSize,
		protectedSize,
	)
}
