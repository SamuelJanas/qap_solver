package solvers

import (
	"fmt"
	"math"
	"qap_solver/internal/metrics"
	"qap_solver/internal/qap"
	"qap_solver/pkg"
	"time"
)

type NearestNeighborSolver struct {
	RandomStarts int
}

// NewNearestNeighborSolver creates a new nearest neighbor solver with specified random starts
func NewNearestNeighborSolver(randomStarts int) *NearestNeighborSolver {
	return &NearestNeighborSolver{
		RandomStarts: randomStarts,
	}
}

func (s *NearestNeighborSolver) Name() string {
	return "NearestNeighbor"
}

func (s *NearestNeighborSolver) Description() string {
	return fmt.Sprintf("Nearest Neighbor heuristic solver (%d random starts)", s.RandomStarts)
}

func (s *NearestNeighborSolver) Solve(instance *qap.QAPInstance) SolverResult {
	bestSolution := make([]int, instance.Size)
	bestFitness := -1

	for i := 0; i < s.RandomStarts; i++ {
		solution := NearestNeighborSolution(instance)
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

func (s *NearestNeighborSolver) SolveWithMetrics(
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

	for i := 0; i < s.RandomStarts; i++ {
		solution := NearestNeighborSolution(instance)
		fitness := qap.CalculateFitness(instance, solution)

		if i == 0 {
			// record initial solution
			initialSolution = make([]int, len(solution))
			copy(initialSolution, solution)
			initialFitness = fitness
		}

		totalSteps += instance.Size // Each solution construction takes n steps
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

// NearestNeighborSolution generates a solution using the nearest neighbor heuristic,
// starting from a random facility and location
func NearestNeighborSolution(instance *qap.QAPInstance) []int {
	size := instance.Size
	solution := make([]int, size)
	for i := range solution {
		solution[i] = -1 // Initialize with -1 to mark as unassigned
	}

	// Use a separate slice to track which locations have been assigned
	assignedLocations := make([]bool, size)

	// Start with a random facility
	currentFacility := pkg.RandomInt(0, size-1)

	// Assign the first facility to a random location
	firstLocation := pkg.RandomInt(0, size-1)
	solution[currentFacility] = firstLocation
	assignedLocations[firstLocation] = true

	// Assign remaining facilities
	for i := 1; i < size; i++ {
		// Find next unassigned facility
		nextFacility := -1
		for j := 0; j < size; j++ {
			if solution[j] == -1 { // Unassigned facility
				if nextFacility == -1 ||
					getFlowCost(instance, j, currentFacility) > getFlowCost(instance, nextFacility, currentFacility) {
					nextFacility = j
				}
			}
		}

		// Find best location for this facility
		bestLocation := -1
		bestCost := math.MaxInt32

		for loc := 0; loc < size; loc++ {
			if !assignedLocations[loc] {
				// Calculate partial cost for this assignment
				cost := calculatePartialCost(instance, solution, nextFacility, loc)
				if cost < bestCost {
					bestCost = cost
					bestLocation = loc
				}
			}
		}

		// Assign the facility to the best location
		solution[nextFacility] = bestLocation
		assignedLocations[bestLocation] = true
		currentFacility = nextFacility
	}

	return solution
}

// getFlowCost returns the flow cost between two facilities
func getFlowCost(instance *qap.QAPInstance, facility1, facility2 int) int {
	return instance.FlowMatrix[facility1][facility2]
}

// calculatePartialCost calculates the cost of assigning a facility to a location,
// considering only the interactions with already assigned facilities
func calculatePartialCost(instance *qap.QAPInstance, partialSolution []int, facility, location int) int {
	cost := 0
	for i := 0; i < instance.Size; i++ {
		if partialSolution[i] != -1 { // Only consider assigned facilities
			flow := instance.FlowMatrix[facility][i]
			distance := instance.DistanceMatrix[location][partialSolution[i]]
			cost += flow * distance

			// Add reverse flow (symmetric)
			flow = instance.FlowMatrix[i][facility]
			distance = instance.DistanceMatrix[partialSolution[i]][location]
			cost += flow * distance
		}
	}
	return cost
}
