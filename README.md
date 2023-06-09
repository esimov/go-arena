# go-arena

Testing and benchmarking the new experimental memory management system called **arenas**, introduced in Go 1.20.

## Benchmarks

You can analyze the performance of the new memory management system by running the following command: 

```bash
$ GOEXPERIMENT=arenas go test -asan -v ./benchmarks/... -bench="." -cover -benchmem
```
It's important to use the `GOEXPERIMENT=arenas` environment variable. The `-asan` flag means to use the program with the memory address sanitizer option enabled. 

### Results:

```bash
goos: linux
goarch: amd64
pkg: esimov/memarena/benchmarks
cpu: 11th Gen Intel(R) Core(TM) i7-11850H @ 2.50GHz
BenchmarkSlice_NoArena
BenchmarkSlice_NoArena/n=100
BenchmarkSlice_NoArena/n=100-16                           332364              4674 ns/op            4523 B/op          0 allocs/op
BenchmarkSlice_NoArena/n=10000
BenchmarkSlice_NoArena/n=10000-16                           3516            607361 ns/op          455522 B/op          0 allocs/op
BenchmarkSlice_NoArena/n=100000
BenchmarkSlice_NoArena/n=100000-16                          325           3838836 ns/op         4157606 B/op          0 allocs/op
BenchmarkSlice_Arena
BenchmarkSlice_Arena/n=100
BenchmarkSlice_Arena/n=100-16                             402762              2760 ns/op            2095 B/op          0 allocs/op
BenchmarkSlice_Arena/n=10000
BenchmarkSlice_Arena/n=10000-16                             4900            423092 ns/op          268697 B/op          0 allocs/op
BenchmarkSlice_Arena/n=100000
BenchmarkSlice_Arena/n=100000-16                             387           2653274 ns/op         2130343 B/op          0 allocs/op
BenchmarkSimpleStruct_NoArena
BenchmarkSimpleStruct_NoArena/n=100
BenchmarkSimpleStruct_NoArena/n=100-16                  128466096                9.292 ns/op           0 B/op          0 allocs/op
BenchmarkSimpleStruct_NoArena/n=10000
BenchmarkSimpleStruct_NoArena/n=10000-16                131011968                9.477 ns/op           0 B/op          0 allocs/op
BenchmarkSimpleStruct_NoArena/n=100000
BenchmarkSimpleStruct_NoArena/n=100000-16               132393255                9.086 ns/op           0 B/op          0 allocs/op
BenchmarkSimpleStruct_Arena
BenchmarkSimpleStruct_Arena/n=100
BenchmarkSimpleStruct_Arena/n=100-16                     6948849               168.9 ns/op            64 B/op          1 allocs/op
BenchmarkSimpleStruct_Arena/n=10000
BenchmarkSimpleStruct_Arena/n=10000-16                   6737854               167.7 ns/op            64 B/op          1 allocs/op
BenchmarkSimpleStruct_Arena/n=100000
BenchmarkSimpleStruct_Arena/n=100000-16                  7434373               164.8 ns/op            64 B/op          1 allocs/op
BenchmarkComplexStruct_NoArena
BenchmarkComplexStruct_NoArena/n=100
BenchmarkComplexStruct_NoArena/n=100-16                   317937              3506 ns/op            4728 B/op          0 allocs/op
BenchmarkComplexStruct_NoArena/n=10000
BenchmarkComplexStruct_NoArena/n=10000-16                   3271            320152 ns/op          489641 B/op          0 allocs/op
BenchmarkComplexStruct_NoArena/n=100000
BenchmarkComplexStruct_NoArena/n=100000-16                   387           3372102 ns/op         4366697 B/op          0 allocs/op
BenchmarkComplexStruct_Arena
BenchmarkComplexStruct_Arena/n=100
BenchmarkComplexStruct_Arena/n=100-16                     215119              5077 ns/op            1973 B/op          0 allocs/op
BenchmarkComplexStruct_Arena/n=10000
BenchmarkComplexStruct_Arena/n=10000-16                     2253            501648 ns/op          293497 B/op          0 allocs/op
BenchmarkComplexStruct_Arena/n=100000
BenchmarkComplexStruct_Arena/n=100000-16                     223           4943917 ns/op         1860246 B/op          0 allocs/op
BenchmarkComplexStruct_IterNoArena
BenchmarkComplexStruct_IterNoArena/n=100
BenchmarkComplexStruct_IterNoArena/n=100-16                  432           2763199 ns/op         3687444 B/op         39 allocs/op
BenchmarkComplexStruct_IterNoArena/n=10000
BenchmarkComplexStruct_IterNoArena/n=10000-16                  4         262156029 ns/op        366595616 B/op      1778 allocs/op
BenchmarkComplexStruct_IterNoArena/n=100000
BenchmarkComplexStruct_IterNoArena/n=100000-16                 1        2158265388 ns/op        802818304 B/op      1018 allocs/op
BenchmarkComplexStruct_IterArena
BenchmarkComplexStruct_IterArena/n=100
BenchmarkComplexStruct_IterArena/n=100-16                    228           4985802 ns/op         1802962 B/op          4 allocs/op
BenchmarkComplexStruct_IterArena/n=10000
BenchmarkComplexStruct_IterArena/n=10000-16                    3         511031599 ns/op        192949549 B/op       359 allocs/op
BenchmarkComplexStruct_IterArena/n=100000
BenchmarkComplexStruct_IterArena/n=100000-16                   1        4703096115 ns/op        838895496 B/op      1110 allocs/op
PASS
```

As we have intuited the arena outperformed significantly the standard memory allocation both in terms of `allocations/operations` and `bytes/operations`. Though one interesting result got my attention, specifically when we tested the arena allocation with a slice with a capacity of 100000. Somehow surprisingly the results are a little worse than the results obtained from the standard allocations. The explication is coming from checking the arena source code. In the runtime/arena package we have the following on line [351](https://github.com/golang/go/blob/d9c29ec6a54f929f4b0736db6b7598a4c2305e5e/src/runtime/arena.go#L351):

> Alloc reserves space in the current chunk or calls refill and reserves space in a new chunk.

What does this means? It means that the arena pre-allocates memory space in chunks. If the data does not fits in the current chunk then it reserves space in a new chunk. Then further down under the `refill()` method , which is invoked within the `alloc()`, we get the following:

> Refill inserts the current arena chunk onto the full list and obtains a new one, either from the partial list or allocating a new one, both from mheap.

I assume that this is the reason why on the last benchmark with a slice having a capacity of 100000 this kind of re-allocation is happening more frequently.

## GC test
The `test` folder contains a few simple examples of how you can use the memory arena together with some analysis about the number of GC calls and heap allocations.

```bash
$ GOEXPERIMENT=arenas go run -asan tests/gc.go

Alloc = 0 MiB   TotalAlloc = 800 MiB    Sys = 26 MiB    NumGC = 1
Alloc = 3 MiB   TotalAlloc = 1676 MiB   Sys = 50 MiB    NumGC = 9
Alloc = 0 MiB   TotalAlloc = 1753 MiB   Sys = 50 MiB    NumGC = 31
```

For a detailed explanation I wrote an article here: https://medium.com/@esimov/go-memory-arenas-1ba930bf79c1

## Author
* Endre Simo ([@simo_endre](https://twitter.com/simo_endre))
