package main

import (
	arena "esimov/memarena"
	"fmt"
	"testing"
)

func BenchmarkSimpleStruct_NoArena(b *testing.B) {
	for _, n := range testCases {
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			for x := 0; x < b.N; x++ {
				_ = &T{data: make([]int, 0, n), len: n}
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
				_ = arena.NewAlloc[T](mem)
			}
		})
	}
}

func BenchmarkComplexStruct_NoArena(b *testing.B) {
	for _, n := range testCases {
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			obj := &T{
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

			obj := arena.NewAlloc[T](mem)
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

func BenchmarkComplexStruct_ArenaOP(b *testing.B) {
	for _, n := range testCases {
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			mem := arena.New()
			defer arena.Free(mem)

			obj := arena.NewAlloc[T](mem)
			obj.len = n
			obj.data = arena.MakeSlice[int](mem, 0, n)

			for i := 0; i < b.N; i++ {
				for i := 0; i < n; i++ {
					if len(obj.data) < cap(obj.data) {
						obj.data = append(obj.data, i)
					} else {
						obj.data = arena.Append(mem, obj.data, i)
					}
				}
			}
		})
	}
}

func BenchmarkComplexStruct_IterNoArena(b *testing.B) {
	for _, n := range testCases {
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			mem := arena.New()
			defer arena.Free(mem)

			for x := 0; x < 50; x++ {
				obj := &T{
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

			for x := 0; x < 50; x++ {
				obj := arena.NewAlloc[T](mem)
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