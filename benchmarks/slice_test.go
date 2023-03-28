package main

import (
	arena "esimov/memarena"
	"fmt"
	"testing"
)

func BenchmarkSlice_NoArena(b *testing.B) {
	for _, n := range testCases {
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			s := make([]*int, 0, n) // make this pointer type
			for i := 0; i < b.N; i++ {
				s = s[:0] // reset slice
				for i := 0; i < n; i++ {
					s = append(s, &i)
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
				s = s[:0]
				for i := 0; i < n; i++ {
					s = arena.Append(mem, s, i)
				}
			}
		})
	}
}
