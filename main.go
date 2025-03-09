package main

import (
	"flag"
	"os"
	"path/filepath"
	"qap_solver/internal/experiment"
	"qap_solver/internal/qap"
	"qap_solver/internal/solvers"
	"qap_solver/pkg"
	"strings"
	"time"
)

var logger = pkg.NewLogger()

func main() {
	// Parse command line arguments
	instanceDir := flag.String("instances", "instances", "Directory containing instance files")
	outputDir := flag.String("output", "results", "Directory for output files")
	solverConfigs := flag.String("solvers", "random:iterations=1000", "See README or baseline for more info. "+
		"Separate solvers by ; and arguments with ,. List arguments after :")
	runsPerInstance := flag.Int("runs", 10, "Number of runs per solver per instance")
	experimentMode := flag.Bool("experiment", false, "Run in experiment mode (batch processing)")
	singleInstanceFile := flag.String("instance", "", "Path to a single instance file (ignored in experiment mode)")
	listSolvers := flag.Bool("list", false, "List available solvers")
	flag.Parse()

	// Create solver factory
	factory := solvers.NewSolverFactory()

	// List available solvers if requested
	if *listSolvers {
		for _, line := range factory.ListAvailable() {
			logger.Println(line)
		}
		return
	}

	// Parse solver configurations
	solverList := strings.Split(*solverConfigs, ";")
	solverInstances := make([]solvers.Solver, 0, len(solverList))

	for _, config := range solverList {
		solver, err := factory.Create(config)
		if err != nil {
			logger.Printf("Error creating solver from config '%s': %v", config, err)
			continue
		}
		solverInstances = append(solverInstances, solver)
	}

	if len(solverInstances) == 0 {
		logger.Fatalf("No valid solvers specified")
	}

	// Run in experiment mode or single instance mode
	if !*experimentMode {
		// Run on a single instance
		instanceFile := *singleInstanceFile
		if instanceFile == "" {
			// Find first .dat file in instance directory
			entries, err := os.ReadDir(*instanceDir)
			if err == nil {
				for _, entry := range entries {
					if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".dat") {
						instanceFile = filepath.Join(*instanceDir, entry.Name())
						break
					}
				}
			}

			if instanceFile == "" {
				logger.Fatalf("No instance file specified and none found in instance directory")
			}
		}

		// Load instance
		startTime := time.Now()
		instance, err := qap.ReadInstance(instanceFile)
		if err != nil {
			logger.Fatalf("Failed to read instance: %v", err)
		}
		pkg.TimeTrack(startTime, "Instance loading", logger)

		logger.Printf("Loaded instance: %s (Size = %d)", instanceFile, instance.Size)

		// Run all solvers on the instance
		bestOverallSolution := solvers.SolverResult{Fitness: -1}

		for _, solver := range solverInstances {
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
	} else {
		// Run batch experiment on all instances
		err := experiment.RunAll(experiment.ExperimentConfig{
			InstancesDir:    *instanceDir,
			OutputDir:       *outputDir,
			Solvers:         solverInstances,
			RunsPerInstance: *runsPerInstance,
			Logger:          logger,
		})

		if err != nil {
			logger.Fatalf("Experiment failed: %v", err)
		}
	}
}
