package solvers

import (
	"math"
	"math/rand"
	"qap_solver/internal/metrics"
	"qap_solver/internal/qap"
	"time"
)

type SimulatedAnnealingSolver struct {
	Alpha          float64
	P              int
	AcceptanceProb float64
}

func NewSimulatedAnnealingSolver(alpha float64, p int, acceptanceProb float64) *SimulatedAnnealingSolver {
	return &SimulatedAnnealingSolver{
		Alpha:          alpha,
		P:              p,
		AcceptanceProb: acceptanceProb,
	}
}

func (s *SimulatedAnnealingSolver) Name() string {
	return "SimulatedAnnealing"
}

func (s *SimulatedAnnealingSolver) Description() string {
	return "Simulated Annealing with adaptive initial temperature and cooling schedule"
}

func (s *SimulatedAnnealingSolver) Solve(instance *qap.QAPInstance) SolverResult {
	n := instance.Size
	Lk := n * (n - 1) / 2

	current := RandomSolution(n)
	best := make([]int, n)
	copy(best, current)

	currentFitness := qap.CalculateFitness(instance, current)
	bestFitness := currentFitness

	// Estimate average delta for worse moves to set initial temperature
	T := s.estimateInitialTemperature(instance, current, currentFitness)

	minTemp := -1.0 / math.Log(s.AcceptanceProb)
	noImprovementCounter := 0
	maxNoImprovement := s.P * Lk

	for T > minTemp || noImprovementCounter < maxNoImprovement {
		i1, i2 := rand.Intn(n), 1+rand.Intn(n-2)
		i1 = (i1 + i2) % n

		neighbor := make([]int, n)
		copy(neighbor, current)
		neighbor[i1], neighbor[i2] = neighbor[i2], neighbor[i1]

		newFitness := qap.CalculateFitness(instance, neighbor)
		delta := float64(newFitness - currentFitness)

		if delta < 0 || (rand.Float64() < math.Exp(-delta/T) && delta != 0) {
			copy(current, neighbor)
			currentFitness = newFitness

			if currentFitness < bestFitness {
				copy(best, current)
				bestFitness = currentFitness
				noImprovementCounter = 0
			}
		} else {
			noImprovementCounter += 1
		}
		T *= s.Alpha
	}

	return SolverResult{
		Solution: best,
		Fitness:  bestFitness,
	}
}

func (s *SimulatedAnnealingSolver) SolveWithMetrics(
	instance *qap.QAPInstance,
	metricsCollector *metrics.MetricsCollector,
	instanceName string,
	runNumber int,
) SolverResult {
	startTime := time.Now()

	n := instance.Size
	Lk := n * (n - 1) / 2

	current := RandomSolution(n)
	best := make([]int, n)
	copy(best, current)

	currentFitness := qap.CalculateFitness(instance, current)
	bestFitness := currentFitness

	initialSolution := make([]int, n)
	copy(initialSolution, current)
	initialFitness := currentFitness

	T := s.estimateInitialTemperature(instance, current, currentFitness)
	minTemp := -1.0 / math.Log(s.AcceptanceProb)

	noImprovementCounter := 0
	maxNoImprovement := s.P * Lk

	totalSteps := 0
	totalEvaluations := 0
	totalSolutionsChecked := 0

	for T > minTemp || noImprovementCounter < maxNoImprovement {
		i1, i2 := rand.Intn(n), 1+rand.Intn(n-2)
		i1 = (i1 + i2) % n

		neighbor := make([]int, n)
		copy(neighbor, current)
		neighbor[i1], neighbor[i2] = neighbor[i2], neighbor[i1]

		newFitness := qap.CalculateFitness(instance, neighbor)
		totalEvaluations++
		totalSolutionsChecked++

		delta := float64(newFitness - currentFitness)

		if delta < 0 || (rand.Float64() < math.Exp(-delta/T) && delta != 0) {
			totalSteps++
			copy(current, neighbor)
			currentFitness = newFitness

			if currentFitness < bestFitness {
				copy(best, current)
				bestFitness = currentFitness
				noImprovementCounter = 0
			}
		} else {
			noImprovementCounter += 1
		}

		T *= s.Alpha
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
			Solution:         best,
		})
	}

	return SolverResult{
		Solution: best,
		Fitness:  bestFitness,
	}
}

func (s *SimulatedAnnealingSolver) estimateInitialTemperature(instance *qap.QAPInstance, sol []int, fitness int) float64 {
	n := instance.Size
	numSamples := 100
	var totalDelta float64
	count := 0

	for i := 0; i < numSamples; i++ {
		i1, i2 := rand.Intn(n), 1+rand.Intn(n-2)
		i1 = (i1 + i2) % n

		neighbor := make([]int, n)
		copy(neighbor, sol)
		neighbor[i1], neighbor[i2] = neighbor[i2], neighbor[i1]
		newFitness := qap.CalculateFitness(instance, neighbor)
		delta := float64(newFitness - fitness)
		if delta > 0 {
			totalDelta += delta
			count++
		}
	}
	if count == 0 {
		return 69420.0
	}
	avgDelta := totalDelta / float64(count)
	return -avgDelta / math.Log(0.95) // for 95% acceptance
}
