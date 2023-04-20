package main

import (
	arena "esimov/memarena"
	"fmt"
	"testing"
)

type Struct[T any] struct {
	len  int
	data []T
}

var testCases = []int{100, 10_000, 100_000}

func BenchmarkSimpleStruct_NoArena(b *testing.B) {
	for _, n := range testCases {
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			for x := 0; x < b.N; x++ {
				_ = Struct[int]{}
			}
		})
	}
}

func BenchmarkSimpleStruct_Arena(b *testing.B) {
	for _, n := range testCases {
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			mem := arena.New()
			defer arena.Free(mem)

			for x := 0; x < b.N; x++ {
				_ = arena.NewAlloc[Struct[int]](mem)
			}
		})
	}
}

func BenchmarkComplexStruct_NoArena(b *testing.B) {
	for _, n := range testCases {
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			obj := Struct[int]{
				data: make([]int, 0, n),
				len:  n,
			}
			for i := 0; i < b.N; i++ {
				for i := 0; i < n; i++ {
					obj.data = append(obj.data, i)
				}
			}
		})
	}
}

func BenchmarkComplexStruct_Arena(b *testing.B) {
	for _, n := range testCases {
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			mem := arena.New()
			defer arena.Free(mem)

			obj := arena.NewAlloc[Struct[int]](mem)
			obj.len = n
			obj.data = arena.MakeSlice[int](mem, 0, n)

			for i := 0; i < b.N; i++ {
				for i := 0; i < n; i++ {
					obj.data = arena.Append(mem, obj.data, i)
				}
			}
		})
	}
}

func BenchmarkComplexStruct_IterNoArena(b *testing.B) {
	for _, n := range testCases {
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			for x := 0; x < 1000; x++ {
				obj := Struct[int]{
					len:  n,
					data: make([]int, 0, n),
				}

				for i := 0; i < b.N; i++ {
					for i := 0; i < n; i++ {
						obj.data = append(obj.data, i)
					}
				}
			}
		})
	}
}

func BenchmarkComplexStruct_IterArena(b *testing.B) {
	for _, n := range testCases {
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			mem := arena.New()
			defer arena.Free(mem)

			for x := 0; x < 1000; x++ {
				obj := arena.NewAlloc[Struct[int]](mem)
				obj.len = n
				obj.data = arena.MakeSlice[int](mem, 0, n)

				for i := 0; i < b.N; i++ {
					for i := 0; i < n; i++ {
						obj.data = arena.Append(mem, obj.data, i)
					}
				}
			}
		})
	}
}
