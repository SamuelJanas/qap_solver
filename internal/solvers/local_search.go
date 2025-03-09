package solvers

import (
	"fmt"
	"qap_solver/internal/metrics"
	"qap_solver/internal/qap"
	"time"
)

type LocalSearchSolver struct {
	MaxIterations  int
	MaxNonImproving int
	RandomRestarts int
}

func NewLocalSearchSolver(maxIterations, maxNonImproving, randomRestarts int) *LocalSearchSolver {
	return &LocalSearchSolver{
		MaxIterations:   maxIterations,
		MaxNonImproving: maxNonImproving,
		RandomRestarts:  randomRestarts,
	}
}

func (s *LocalSearchSolver) Name() string {
	return "LocalSearch"
}

func (s *LocalSearchSolver) Description() string {
	return fmt.Sprintf("Local search with swap neighborhood (Max iterations: %d, Non-improving limit: %d, Random restarts: %d)",
		s.MaxIterations, s.MaxNonImproving, s.RandomRestarts)
}


func (s *LocalSearchSolver) Solve(instance *qap.QAPInstance) SolverResult {
	bestSolution := make([]int, instance.Size)
	bestFitness := -1

	for restart := 0; restart < s.RandomRestarts; restart++ {
		// Start with a random solution
		currentSolution := RandomSolution(instance.Size)
		currentFitness := qap.CalculateFitness(instance, currentSolution)

		nonImprovingCount := 0

		for iter := 0; iter < s.MaxIterations && nonImprovingCount < s.MaxNonImproving; iter++ {
			improved := false

			// Try all possible swaps to find improvement
			for i := 0; i < instance.Size-1; i++ {
				for j := i + 1; j < instance.Size; j++ {
					// Create a new solution by swapping positions i and j
					newSolution := make([]int, instance.Size)
					copy(newSolution, currentSolution)
					newSolution[i], newSolution[j] = newSolution[j], newSolution[i]

					// Calculate fitness of new solution
					newFitness := qap.CalculateFitness(instance, newSolution)

					// If it's better, accept it
					if newFitness < currentFitness {
						copy(currentSolution, newSolution)
						currentFitness = newFitness
						improved = true
						// Break out of inner loop to start again with the new solution
						break
					}
				}
				if improved {
					break
				}
			}

			// If no improvement was found in this iteration
			if !improved {
				nonImprovingCount++
			} else {
				nonImprovingCount = 0
			}
		}

		// Update best solution if this restart found a better one
		if bestFitness == -1 || currentFitness < bestFitness {
			copy(bestSolution, currentSolution)
			bestFitness = currentFitness
		}
	}

	return SolverResult{
		Solution: bestSolution,
		Fitness:  bestFitness,
	}
}


func (s *LocalSearchSolver) SolveWithMetrics(
    instance *qap.QAPInstance,
    metricsCollector *metrics.MetricsCollector,
    instanceName string,
    runNumber int,
) SolverResult {
    startTime := time.Now()

    bestSolution := make([]int, instance.Size)
    bestFitness := -1

    totalSteps := 0
    totalEvaluations := 0
    totalSolutionsChecked := 0

    var initialSolution []int
    var initialFitness int

    for restart := 0; restart < s.RandomRestarts; restart++ {
        currentSolution := RandomSolution(instance.Size)
        currentFitness := qap.CalculateFitness(instance, currentSolution)

        if restart == 0 {
            initialSolution = make([]int, len(currentSolution))
            copy(initialSolution, currentSolution)
            initialFitness = currentFitness
        }

        nonImprovingCount := 0

        for iter := 0; iter < s.MaxIterations && nonImprovingCount < s.MaxNonImproving; iter++ {
            improved := false

            for i := 0; i < instance.Size-1; i++ {
                for j := i + 1; j < instance.Size; j++ {
                    newSolution := make([]int, instance.Size)
                    copy(newSolution, currentSolution)
                    newSolution[i], newSolution[j] = newSolution[j], newSolution[i]

                    newFitness := qap.CalculateFitness(instance, newSolution)
                    totalEvaluations++
                    totalSolutionsChecked++

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
            if !improved {
                nonImprovingCount++
            } else {
                nonImprovingCount = 0
            }
        }

        if bestFitness == -1 || currentFitness < bestFitness {
            copy(bestSolution, currentSolution)
            bestFitness = currentFitness
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

