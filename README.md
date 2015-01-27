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

The current winning implementation, called "Presort", takes advantage of the fact that in most
cases, dependencies will be added before the phases that depend on them. This is true
in both zoom and the test cases. It keeps all the phases in a doubly linked list (out
of the standard library). On each call to addDependency, it checks if the dependency
is already satisfied by the current ordering of phases (90% of the time it is!). If
it's not, it reorders the list by either moving the phase directly after the one it
depends on or moving the phase it depends on directly before it. If nether is possible
because of existing dependencies, it returns an error declaring there is a cycle. It's
quite fast, but very much tuned to my particular use case.

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
BenchmarkLinear1Goraph	 1000000	      2001 ns/op
BenchmarkLinear1Unix	    1000	   1750961 ns/op
BenchmarkLinear1Graph	  500000	      3146 ns/op
BenchmarkLinear1Maps	 1000000	      1563 ns/op
BenchmarkLinear1Lists	 1000000	      1207 ns/op
BenchmarkLinear1Presort	 2000000	       757 ns/op
BenchmarkLinear3Goraph	  200000	      8927 ns/op
BenchmarkLinear3Unix	    1000	   1778778 ns/op
BenchmarkLinear3Graph	  200000	      5560 ns/op
BenchmarkLinear3Maps	  500000	      3783 ns/op
BenchmarkLinear3Lists	  300000	      4569 ns/op
BenchmarkLinear3Presort	  500000	      3065 ns/op
BenchmarkLinear10Goraph	   50000	     33329 ns/op
BenchmarkLinear10Unix	    1000	   1765673 ns/op
BenchmarkLinear10Graph	   50000	     25135 ns/op
BenchmarkLinear10Maps	   50000	     25532 ns/op
BenchmarkLinear10Lists	  100000	     19553 ns/op
BenchmarkLinear10Presort	  200000	     10988 ns/op
BenchmarkTree1Goraph	  300000	      4862 ns/op
BenchmarkTree1Unix	    1000	   1757980 ns/op
BenchmarkTree1Graph	  300000	      5590 ns/op
BenchmarkTree1Maps	  500000	      3566 ns/op
BenchmarkTree1Lists	  500000	      2951 ns/op
BenchmarkTree1Presort	 1000000	      1935 ns/op
BenchmarkTree3Goraph	   50000	     35126 ns/op
BenchmarkTree3Unix	    1000	   1780420 ns/op
BenchmarkTree3Graph	  200000	     10869 ns/op
BenchmarkTree3Maps	  200000	      5886 ns/op
BenchmarkTree3Lists	  200000	      6660 ns/op
BenchmarkTree3Presort	  300000	      4021 ns/op
BenchmarkTree10Goraph	    5000	    308559 ns/op
BenchmarkTree10Unix	    1000	   1775422 ns/op
BenchmarkTree10Graph	   50000	     28792 ns/op
BenchmarkTree10Maps	  100000	     21667 ns/op
BenchmarkTree10Lists	  100000	     20339 ns/op
BenchmarkTree10Presort	  100000	     12157 ns/op
```
