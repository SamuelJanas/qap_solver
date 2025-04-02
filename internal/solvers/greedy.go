package solvers

import (
	"fmt"
	"qap_solver/internal/metrics"
	"qap_solver/internal/qap"
	"time"
)

type GreedySolver struct {
	MaxIterations  int
	RandomRestarts int
}

func NewGreedySolver(maxIterations int) *GreedySolver {
	return &GreedySolver{
		MaxIterations: maxIterations,
	}
}

func (s *GreedySolver) Name() string {
	return "Greedy"
}

func (s *GreedySolver) Description() string {
	return fmt.Sprintf("Greedy search")
}

func (s *GreedySolver) Solve(instance *qap.QAPInstance) SolverResult {
	currentSolution := RandomSolution(instance.Size)
	currentFitness := qap.CalculateFitness(instance, currentSolution)

	for {
		improved := false
		for i := 0; i < instance.Size-1; i++ {
			for j := i + 1; j < instance.Size; j++ {
				newSolution := make([]int, instance.Size)
				copy(newSolution, currentSolution)
				newSolution[i], newSolution[j] = newSolution[j], newSolution[i]
				newFitness := qap.CalculateFitness(instance, newSolution)

				if newFitness < currentFitness {
					copy(currentSolution, newSolution)
					currentFitness = newFitness
					improved = true
					break
				}
			}
			if improved {
				break
			}
		}
		if !improved {
			break
		}
	}
	return SolverResult{Solution: currentSolution, Fitness: currentFitness}
}

func (s *GreedySolver) SolveWithMetrics(
	instance *qap.QAPInstance,
	metricsCollector *metrics.MetricsCollector,
	instanceName string,
	runNumber int,
) SolverResult {
	startTime := time.Now()

	// Initial values for solution and fitness
	currentSolution := RandomSolution(instance.Size)
	currentFitness := qap.CalculateFitness(instance, currentSolution)

	// Metrics counters
	totalSteps := 0
	totalEvaluations := 0
	totalSolutionsChecked := 0

	var initialSolution []int
	var initialFitness int

	// Record initial solution
	initialSolution = make([]int, len(currentSolution))
	copy(initialSolution, currentSolution)
	initialFitness = currentFitness

	// Start the greedy search iterations until no improvement
	for iter := 0; iter < s.MaxIterations; iter++ {
		improved := false

		// Try to improve the current solution by checking neighbors
		for i := 0; i < instance.Size-1; i++ {
			for j := i + 1; j < instance.Size; j++ {
				newSolution := make([]int, instance.Size)
				copy(newSolution, currentSolution)
				newSolution[i], newSolution[j] = newSolution[j], newSolution[i]
				newFitness := qap.CalculateFitness(instance, newSolution)

				totalEvaluations++
				totalSolutionsChecked++

				// If a better solution is found, accept it
				if newFitness < currentFitness {
					copy(currentSolution, newSolution)
					currentFitness = newFitness
					improved = true
					break
				}
			}
			if improved {
				break
			}
		}

		totalSteps++

		// If no improvement is found, exit the loop
		if !improved {
			break
		}
	}

	// Calculate elapsed time
	elapsedTime := time.Since(startTime)

	// Record metrics if the collector is provided
	if metricsCollector != nil {
		metricsCollector.AddRunMetrics(metrics.RunMetrics{
			InstanceName:     instanceName,
			SolverName:       s.Name(),
			Run:              runNumber,
			InitialFitness:   initialFitness,
			FinalFitness:     currentFitness,
			TimeElapsed:      elapsedTime,
			StepsCount:       totalSteps,
			EvaluationsCount: totalEvaluations,
			SolutionsChecked: totalSolutionsChecked,
			Solution:         currentSolution,
		})
	}

	// Return the result
	return SolverResult{
		Solution: currentSolution,
		Fitness:  currentFitness,
	}
}
