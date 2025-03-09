## How to use:

1. List available solvers:
```sh
go run main.go -list
```

2. Run with specific solvers: NOTE the use of `;`, `:`, `,` and the lack of spaces.
```sh
go run main.go -solvers="random:iterations=2000;localsearch:maxIter=5000,maxNonImproving=500,restarts=10"
```

3. Run with specific instance:
```sh
go run main.go -instance="instances/bur26a.dat" -solvers="localsearch"
```

4. Run in experiment mode:
```sh
go run main.go -experiment -instances=instances -runs=10 -solvers="random:iterations=2000"
```
This mode creates a .csv file inside of results/ directory (by default) with run details. It can be further analysed with "TODO.py".


## Add new solvers:

1. Implement the `Solver` interface. See `internal/solvers/random.go` for specifics.
2. Add a creator function in `internal/solvers/solver_factory.go`.
3. Register the `Solver` in `NewSolverFactory`.
4. Append the new solver to `ListAvailable`.
