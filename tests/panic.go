package main

// import (
// 	arena "esimov/memarena"
// )

// type T struct {
// 	val int
// }

// func main() {
// 	// Calling this function with the `-asan` option will panic.
// 	memAllocArenaPanic()
// }

// func memAllocArenaPanic() *T {
// 	mem := arena.New()

// 	obj := arena.NewAlloc[T](mem)
// 	arena.Free(mem)
// 	// Accessing a variable after the allocated memory has been released will panic.
// 	obj.val = 1

// 	return obj
// }
