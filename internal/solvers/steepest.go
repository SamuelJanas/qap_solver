package solvers

import (
	"fmt"
	"qap_solver/internal/metrics"
	"qap_solver/internal/qap"
	"time"
)

type SteepestSolver struct {
	MaxIterations  int
	RandomRestarts int
}

func NewSteepestSolver(maxIterations int) *SteepestSolver {
	return &SteepestSolver{
		MaxIterations: maxIterations,
	}
}

func (s *SteepestSolver) Name() string {
	return "Steepest"
}

func (s *SteepestSolver) Description() string {
	return fmt.Sprintf("Steepest search")
}

func (s *SteepestSolver) Solve(instance *qap.QAPInstance) SolverResult {

	currentSolution := RandomSolution(instance.Size)
	currentFitness := qap.CalculateFitness(instance, currentSolution)

	for {
		bestNeighbor := make([]int, instance.Size)
		copy(bestNeighbor, currentSolution)
		bestNeighborFitness := currentFitness

		for i := 0; i < instance.Size-1; i++ {
			for j := i + 1; j < instance.Size; j++ {
				newSolution := make([]int, instance.Size)
				copy(newSolution, currentSolution)
				newSolution[i], newSolution[j] = newSolution[j], newSolution[i]
				newFitness := qap.CalculateFitness(instance, newSolution)

				if newFitness < bestNeighborFitness {
					copy(bestNeighbor, newSolution)
					bestNeighborFitness = newFitness
				}
			}
		}
		if bestNeighborFitness < currentFitness {
			copy(currentSolution, bestNeighbor)
			currentFitness = bestNeighborFitness
		} else {
			break
		}
	}
	return SolverResult{Solution: currentSolution, Fitness: currentFitness}
}

func (s *SteepestSolver) SolveWithMetrics(
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

	// Start the steepest descent iterations
	for {
		bestNeighbor := make([]int, instance.Size)
		copy(bestNeighbor, currentSolution)
		bestNeighborFitness := currentFitness

		// Check all possible neighbors
		for i := 0; i < instance.Size-1; i++ {
			for j := i + 1; j < instance.Size; j++ {
				newSolution := make([]int, instance.Size)
				copy(newSolution, currentSolution)
				newSolution[i], newSolution[j] = newSolution[j], newSolution[i]
				newFitness := qap.CalculateFitness(instance, newSolution)

				totalEvaluations++
				totalSolutionsChecked++

				// Update the best neighbor if a better fitness is found
				if newFitness < bestNeighborFitness {
					copy(bestNeighbor, newSolution)
					bestNeighborFitness = newFitness
				}
			}
		}

		totalSteps++

		// If a better solution was found, accept it
		if bestNeighborFitness < currentFitness {
			copy(currentSolution, bestNeighbor)
			currentFitness = bestNeighborFitness
		} else {
			// If no improvement is found, exit the loop
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
