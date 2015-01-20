# dependency-linearization
An experiment in finding a fast dependency linearization algorithm for go.

This code is a series of benchmarks and tests against different dependency
linearization implementations. In the end, I found an implementation that was
up to 15-20x faster when linearizing 10 phases. The code is here in case someone
might find it helpful or interesting.

### The Goal

I needed to find a faster way to linearize dependencies in [zoom](https://github.com/albrow/zoom),
which is a fast, lightweight ORM for go backed by redis. The problem is as follows:

1. There is a "transaction" object which consists of any number of "phases". (Each
phase gets translated into a redis-style transaction using MULTI and EXEC).
2. Each phase consists of some number of "commands". (Redis instructions to be executed
line by line inside of a MULTI/EXEC statement).
2. Some phases depend on the results from other phases in the transaction, and hence
some phases need to be executed before others.
3. I needed a way to "linearize" the dependent phases, i.e. sort them into a list of
phases which is safe to execute sequentially.
4. Whatever method I choose needs a way to detect cycles and report an error if there
is one. This would indicate there is some cyclical dependencies which cannot be
linearized.

The first implementation I picked (somewhat randomly) for zoom demonstrated that the
concept could work, and that it made the code easier to manage. But it had a huge impact
on performance, so I went on the hunt for something faster.

### How to Run the Tests

`test/correctness_test.go` tests each implementation for correctness. You can run
these tests with the following command:

```
go test ./...
```

`test/bench_test.go` contains a few different benchmarks with different graph
topologies for each implementation. You can run the benchmarks with the following:

```
go test ./... -run NONE -bench .
```

The `-run NONE` part is optional. It tells go to skip the tests and only run the
benchmarks, since the pattern "NONE" does not appear in the name of any test functions.

These were the results on my laptop:

```
BenchmarkLinear1Goraph  1000000      2465 ns/op
BenchmarkLinear1Unix       1000   1725134 ns/op
BenchmarkLinear1Graph    500000      3740 ns/op
BenchmarkLinear3Goraph   200000     10602 ns/op
BenchmarkLinear3Unix       1000   1838640 ns/op
BenchmarkLinear3Graph    200000      9324 ns/op
BenchmarkLinear10Goraph   30000     42080 ns/op
BenchmarkLinear10Unix      1000   1772527 ns/op
BenchmarkLinear10Graph    50000     28542 ns/op
BenchmarkTree1Goraph     300000      5171 ns/op
BenchmarkTree1Unix         1000   1780040 ns/op
BenchmarkTree1Graph      200000      6545 ns/op
BenchmarkTree3Goraph      50000     35202 ns/op
BenchmarkTree3Unix         1000   1812251 ns/op
BenchmarkTree3Graph      200000     11133 ns/op
BenchmarkTree10Goraph      5000    362511 ns/op
BenchmarkTree10Unix        1000   1853921 ns/op
BenchmarkTree10Graph      50000     34455 ns/op
```
