# dependency-linearization
An experiment in finding a fast dependency linearization algorithm for go.

Not sure if this will ever develop this concept into a usable public library
(I might at some point). But regardless, the code is here in case it might help
someone.

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

The first implementation I picked (somewhat randomly) for zoom demonstrated that the
concept could work, and that it made the code easier to manage. But it had a huge impact
on performance, so I went on the hunt for something faster. This code is a series of
benchmarks and tests against different dependency linearization implementations. In
the end, I found an implementation that was up to 15-20x faster when linearizing
10 phases.

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