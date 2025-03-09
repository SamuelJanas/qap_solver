package solvers

import (
	"fmt"
	"qap_solver/internal/metrics"
	"qap_solver/internal/qap"
	// "qap_solver/pkg"
	"time"
)

// SteepestLocalSearchSolver implements a steepest descent local search
type SteepestLocalSearchSolver struct {
	MaxIterations   int
	MaxNonImproving int
	RandomRestarts  int
}

// NewSteepestLocalSearchSolver creates a new steepest local search solver
func NewSteepestLocalSearchSolver(maxIterations, maxNonImproving, randomRestarts int) *SteepestLocalSearchSolver {
	return &SteepestLocalSearchSolver{
		MaxIterations:   maxIterations,
		MaxNonImproving: maxNonImproving,
		RandomRestarts:  randomRestarts,
	}
}

func (s *SteepestLocalSearchSolver) Name() string {
	return "SteepestLocalSearch"
}

func (s *SteepestLocalSearchSolver) Description() string {
	return fmt.Sprintf("Steepest local search with swap neighborhood (Max iterations: %d, Non-improving limit: %d, Random restarts: %d)",
		s.MaxIterations, s.MaxNonImproving, s.RandomRestarts)
}

// Solve runs the steepest local search algorithm and collects metrics
func (s *SteepestLocalSearchSolver) Solve(instance *qap.QAPInstance) SolverResult {
	return s.SolveWithMetrics(instance, nil, "", 0, 0)
}

// SolveWithMetrics runs the steepest local search algorithm with detailed metrics collection
func (s *SteepestLocalSearchSolver) SolveWithMetrics(
	instance *qap.QAPInstance,
	metricsCollector *metrics.MetricsCollector,
	instanceName string,
	runNumber int,
	optimalFitness int,
) SolverResult {
	startTime := time.Now()

	bestSolution := make([]int, instance.Size)
	bestFitness := -1

	totalSteps := 0
	totalEvaluations := 0
	totalSolutionsChecked := 0

	// Track initial solution quality for reporting
	var initialSolution []int
	var initialFitness int

	for restart := 0; restart < s.RandomRestarts; restart++ {
		// Start with a random solution
		currentSolution := RandomSolution(instance.Size)
		currentFitness := qap.CalculateFitness(instance, currentSolution)

		// For the first restart, record the initial solution
		if restart == 0 {
			initialSolution = make([]int, len(currentSolution))
			copy(initialSolution, currentSolution)
			initialFitness = currentFitness
		}

		steps := 0
		evaluations := 0
		solutionsChecked := 0

		for iter := 0; iter < s.MaxIterations; iter++ {
			// Steepest descent - evaluate all neighbors and choose the best
			bestSwapI := -1
			bestSwapJ := -1
			bestSwapFitness := currentFitness

			// Try all possible swaps
			for i := 0; i < instance.Size-1; i++ {
				for j := i + 1; j < instance.Size; j++ {
					// We're checking a new potential solution
					solutionsChecked++

					// Evaluate the swap
					swapFitness := evaluateSwap(instance, currentSolution, currentFitness, i, j)
					evaluations++

					// If it's better, remember it
					if swapFitness < bestSwapFitness {
						bestSwapI = i
						bestSwapJ = j
						bestSwapFitness = swapFitness
					}
				}
			}

			// If we found a better solution, make the swap
			if bestSwapI != -1 {
				currentSolution[bestSwapI], currentSolution[bestSwapJ] = currentSolution[bestSwapJ], currentSolution[bestSwapI]
				currentFitness = bestSwapFitness
				steps++
			} else {
				// No improvement found, we're at a local optimum
				break
			}
		}

		// Update totals
		totalSteps += steps
		totalEvaluations += evaluations
		totalSolutionsChecked += solutionsChecked

		// Update best solution if this restart found a better one
		if bestFitness == -1 || currentFitness < bestFitness {
			copy(bestSolution, currentSolution)
			bestFitness = currentFitness
		}
	}

	// Calculate elapsed time
	elapsedTime := time.Since(startTime)

	// Collect metrics if a collector was provided
	if metricsCollector != nil {
		metricsCollector.AddRunMetrics(metrics.RunMetrics{
			InstanceName:     instanceName,
			SolverName:       s.Name(),
			Run:              runNumber,
			InitialFitness:   initialFitness,
			FinalFitness:     bestFitness,
			OptimalFitness:   optimalFitness,
			TimeElapsed:      elapsedTime,
			StepsCount:       totalSteps,
			EvaluationsCount: totalEvaluations,
			SolutionsChecked: totalSolutionsChecked,
			Solution:         bestSolution,
		})
	}

	return SolverResult{
		Solution: bestSolution,
		Fitness:  bestFitness,
	}
}

// evaluateSwap calculates the fitness after swapping facilities i and j
// This is a more efficient implementation that doesn't create a new solution
func evaluateSwap(instance *qap.QAPInstance, solution []int, currentFitness int, i, j int) int {
	diff := 0

	// For each pair of locations
	for k := 0; k < instance.Size; k++ {
		if k != i && k != j {
			// Calculate the change in flow*distance due to the swap
			diff += (instance.FlowMatrix[solution[i]][solution[k]] - instance.FlowMatrix[solution[j]][solution[k]]) *
				   (instance.DistanceMatrix[i][k] - instance.DistanceMatrix[j][k])

			diff += (instance.FlowMatrix[solution[k]][solution[i]] - instance.FlowMatrix[solution[k]][solution[j]]) *
				   (instance.DistanceMatrix[k][i] - instance.DistanceMatrix[k][j])
		}
	}

	// Add the contribution of the direct flow between i and j
	diff += (instance.FlowMatrix[solution[i]][solution[j]] + instance.FlowMatrix[solution[j]][solution[i]]) *
		   (instance.DistanceMatrix[i][j] - instance.DistanceMatrix[j][i])

	return currentFitness + diff
}
