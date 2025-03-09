package solvers

import (
	"fmt"
	"qap_solver/internal/metrics"
	"qap_solver/internal/qap"
	"qap_solver/pkg"
	"time"
)

type RandomSolver struct {
	Iterations int
}

// NewRandomSolver creates a new random solver with specified iterations
func NewRandomSolver(iterations int) *RandomSolver {
	return &RandomSolver{
		Iterations: iterations,
	}
}

func (s *RandomSolver) Name() string {
	return "Random"
}

func (s *RandomSolver) Description() string {
	return fmt.Sprintf("Random solution generator (%d iterations)", s.Iterations)
}

func (s *RandomSolver) Solve(instance *qap.QAPInstance) SolverResult {
	bestSolution := make([]int, instance.Size)
	bestFitness := -1

	for i := 0; i < s.Iterations; i++ {
		solution := RandomSolution(instance.Size)
		fitness := qap.CalculateFitness(instance, solution)

		if bestFitness == -1 || fitness < bestFitness {
			copy(bestSolution, solution)
			bestFitness = fitness
		}
	}

	return SolverResult{
		Solution: bestSolution,
		Fitness:  bestFitness,
	}
}

func (s *RandomSolver) SolveWithMetrics(
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

	var initialSolution []int
	var initialFitness int

	for i := 0; i < s.Iterations; i++ {
		solution := RandomSolution(instance.Size)
		fitness := qap.CalculateFitness(instance, solution)

		if i == 0 {
			// record initial solution
			initialSolution = make([]int, len(solution))
			copy(initialSolution, solution)
			initialFitness = fitness
		}

		totalSteps += 1
		totalEvaluations += 1
		totalSolutionsChecked += 1

		if bestFitness == -1 || fitness < bestFitness {
			copy(bestSolution, solution)
			bestFitness = fitness
		}
	}

	elapsedTime := time.Since(startTime)

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

func RandomSolution(size int) []int {
	solution := make([]int, size)
	for i := range solution {
		solution[i] = i
	}
	pkg.ShuffleSlice(solution)
	return solution
}
