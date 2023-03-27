package main

import (
	arena "esimov/memarena"
	"runtime"
)

type Struct[T any] struct {
	len  int
	data []T
}

var (
	size = 1000
	iter = 10_000_00
)

func main() {
	for i := 0; i < iter; i++ {
		_ = memAllocArena[int]()
	}
	arena.PrintStats()

	runtime.GC()
	// We can safely use the cloned slice here, because this is
	// going into the heap before freeing up the memory arena.
	// For this reason the number of GC calls should be
	// with one more then the previous method invocation.
	_ = memAllocArenaClone[int]()
	arena.PrintStats()

	runtime.GC()
	for i := 0; i < iter; i++ {
		_ = memAllocClassic[int]()
	}
	arena.PrintStats()
}

func memAllocClassic[T any]() []T {
	s := &Struct[T]{len: size, data: make([]T, 0, size)}

	for i := 0; i < size; i++ {
		s.data = append(s.data, (any(i)).(T))
	}
	return s.data
}

func memAllocArena[T any]() []T {
	mem := arena.New()
	defer arena.Free(mem)

	obj := arena.NewAlloc[Struct[T]](mem)
	obj.len = size
	obj.data = arena.MakeSlice[T](mem, 0, size)

	for i := 0; i < size; i++ {
		obj.data = arena.Append(mem, obj.data, (any(i)).(T))

	}
	return obj.data
}

func memAllocArenaClone[T any]() []T {
	mem := arena.New()
	defer arena.Free(mem)

	obj := arena.NewAlloc[Struct[T]](mem)
	obj.len = size
	obj.data = arena.MakeSlice[T](mem, 0, size)

	for i := 0; i < size; i++ {
		obj.data = arena.Append(mem, obj.data, (any(i)).(T))
	}

	return arena.Clone(obj.data)
}
