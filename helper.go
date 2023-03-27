package memarena

import (
	"arena"
	"fmt"
	"runtime"
)

func New() *arena.Arena {
	return arena.NewArena()
}

func NewAlloc[T any](a *arena.Arena) *T {
	return arena.New[T](a)
}

func Free(a *arena.Arena) {
	if a == nil {
		return
	}
	(*arena.Arena)(a).Free()
}

func Clone[T any](s T) T {
	return arena.Clone(s)
}

func MakeSlice[T any](a *arena.Arena, l, c int) []T {
	if a == nil {
		return make([]T, l, c)
	}
	return arena.MakeSlice[T](a, l, c)
}

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
