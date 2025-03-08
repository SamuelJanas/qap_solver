package main

import (
	"flag"
	"strings"
	"time"
	"qap_solver/internal/qap"
	"qap_solver/internal/solvers"
	"qap_solver/pkg"
)

var logger = pkg.NewLogger()

func main() {
	// parse command line
	instanceFile := flag.String("instance", "instances/bur26a.dat", "Path to instance file")
	solverConfigs := flag.String("solvers", "random:iterations=1000", "Comma-separated list of solvers to run")
	listSolvers := flag.Bool("list", false, "List available solvers")
	flag.Parse()

	factory := solvers.NewSolverFactory()

	if *listSolvers {
		for _, line := range factory.ListAvailable() {
			logger.Println(line)
		}
		return
	}

	startTime := time.Now()
	instance, err := qap.ReadInstance(*instanceFile)
	if err != nil {
		logger.Fatalf("Failed to read instance: %v", err)
	}
	pkg.TimeTrack(startTime, "Instance loading", logger)

	logger.Printf("Loaded instance: Size = %d\n", instance.Size)

	// Run all requested solvers
	solverList := strings.Split(*solverConfigs, ";")
	bestOverallSolution := solvers.SolverResult{Fitness: -1}

	for _, config := range solverList {
		solver, err := factory.Create(config)
		if err != nil {
			logger.Printf("Error creating solver from config '%s': %v", config, err)
			continue
		}

		logger.Printf("Running solver: %s (%s)", solver.Name(), solver.Description())
		startTime := time.Now()
		result := solver.Solve(instance)
		pkg.TimeTrack(startTime, solver.Name()+" execution", logger)

		logger.Printf("%s fitness: %d", solver.Name(), result.Fitness)

		if bestOverallSolution.Fitness == -1 || result.Fitness < bestOverallSolution.Fitness {
			bestOverallSolution = result
			logger.Printf("New best solution found by %s", solver.Name())
		}
	}

	logger.Printf("Best overall solution has fitness: %d", bestOverallSolution.Fitness)
	logger.Printf("Solution: %v", bestOverallSolution.Solution)
}
