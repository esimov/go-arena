package main

import (
	arena "esimov/memarena"
	"fmt"
	"testing"
)

type T struct {
	len  int
	data []int
}

var testCases = []int{1000, 100_000, 1_000_000}

func BenchmarkSlice_NoArena(b *testing.B) {
	for _, n := range testCases {
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			s := make([]int, 0, n)
			for i := 0; i < b.N; i++ {
				for i := 0; i < n; i++ {
					s = append(s, i)
				}
			}
		})
	}
}

func BenchmarkSlice_Arena(b *testing.B) {
	for _, n := range testCases {
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			mem := arena.New()
			s := arena.MakeSlice[int](mem, 0, n)
			defer arena.Free(mem)

			for i := 0; i < b.N; i++ {
				s = s[:0] // reset slice
				for i := 0; i < n; i++ {
					if len(s) < cap(s) {
						s = append(s, i)
					} else {
						s = arena.Append(mem, s, i)
					}
				}
			}
		})
	}
}

func newStructAllocClone() *T {
	mem := arena.New()
	defer arena.Free(mem)

	obj := arena.NewAlloc[T](mem)
	obj.data = []int{1, 2}
	obj.len = 2
	clone := arena.Clone(obj)

	return clone
}
