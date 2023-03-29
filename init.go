package memarena

import (
	"arena"
	"fmt"
	"runtime"
)

// New allocates a new memory arena.
func New() *arena.Arena {
	return arena.NewArena()
}

// NewAlloc creates a new T value in the allocated memory arena.
func NewAlloc[T any](a *arena.Arena) *T {
	return arena.New[T](a)
}

// Free frees the memory arena without the garbage collection overhead.
func Free(a *arena.Arena) {
	if a == nil {
		return
	}
	(*arena.Arena)(a).Free()
}

// Clone returns a shallow copy of the allocated object in the memory arena.
// After freeing up the arena the cloned object will be moved into the heap,
// which means that it can be referenced after the arena is cleaned up.
// This out-lived object will be cleaned up when the GC cycle will pick it up.
func Clone[T any](s T) T {
	return arena.Clone(s)
}

// MakeSlice creates a new slice and puts it into the arena.
func MakeSlice[T any](a *arena.Arena, l, c int) []T {
	if a == nil {
		return make([]T, l, c)
	}
	return arena.MakeSlice[T](a, l, c)
}

// Append is a helper method to populate the arena allocated slice.
func Append[T any](a *arena.Arena, data []T, v T) []T {
	if a == nil {
		return append(data, v)
	}
	c := 2 * len(data)
	// Increase the slice to the double of its initial capacity.
	if len(data) >= cap(data) {
		if c == 0 {
			c = 1
		}
		slice := arena.MakeSlice[T](a, len(data)+1, c)
		copy(slice, data)
		data = slice
		data[len(data)-1] = v
	} else {
		data = append(data, v)
	}
	return data
}

// PrintStats prints various information related to the memory allocations.
func PrintStats() {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	fmt.Printf("Alloc = %v MiB", bToMb(mem.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(mem.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(mem.Sys))
	fmt.Printf("\tNumGC = %v\n", mem.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
