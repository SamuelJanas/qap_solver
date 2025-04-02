package solvers

import (
	"qap_solver/internal/metrics"
	"qap_solver/internal/qap"
	"sort"
	"time"
)

type GreedyConstructionSolver struct{}

// NewGreedyConstructionSolver creates a new instance of the greedy heuristic solver
func NewGreedyConstructionSolver() *GreedyConstructionSolver {
	return &GreedyConstructionSolver{}
}

func (s *GreedyConstructionSolver) Name() string {
	return "Heuristic"
}

func (s *GreedyConstructionSolver) Description() string {
	return "Greedy heuristic for Quadratic Assignment Problem (QAP)"
}

func (s *GreedyConstructionSolver) Solve(instance *qap.QAPInstance) SolverResult {
	solution := greedyConstruction(instance, nil)
	fitness := qap.CalculateFitness(instance, solution)
	return SolverResult{Solution: solution, Fitness: fitness}
}

func (s *GreedyConstructionSolver) SolveWithMetrics(
	instance *qap.QAPInstance,
	metricsCollector *metrics.MetricsCollector,
	instanceName string,
	runNumber int,
) SolverResult {
	startTime := time.Now()
	totalSteps := 0
	totalEvaluations := 0
	solution := greedyConstruction(instance, &totalSteps)
	fitness := qap.CalculateFitness(instance, solution)
	totalEvaluations++

	elapsedTime := time.Since(startTime)

	if metricsCollector != nil {
		metricsCollector.AddRunMetrics(metrics.RunMetrics{
			InstanceName:     instanceName,
			SolverName:       s.Name(),
			Run:              runNumber,
			InitialFitness:   fitness,
			FinalFitness:     fitness,
			TimeElapsed:      elapsedTime,
			StepsCount:       totalSteps,
			EvaluationsCount: totalEvaluations,
			SolutionsChecked: totalSteps,
			Solution:         solution,
		})
	}

	return SolverResult{Solution: solution, Fitness: fitness}
}

func greedyConstruction(instance *qap.QAPInstance, stepsCounter *int) []int {
	size := instance.Size
	unassignedFacilities := make([]int, size)
	unassignedLocations := make([]int, size)
	assigned := make([][2]int, 0, size)

	for i := 0; i < size; i++ {
		unassignedFacilities[i] = i
		unassignedLocations[i] = i
	}

	sort.Slice(unassignedFacilities, func(i, j int) bool {
		return facilityFlowSum(instance, unassignedFacilities[i]) > facilityFlowSum(instance, unassignedFacilities[j])
	})

	for len(unassignedFacilities) > 0 {
		facility := unassignedFacilities[len(unassignedFacilities)-1]
		unassignedFacilities = unassignedFacilities[:len(unassignedFacilities)-1]

		sort.Slice(unassignedLocations, func(i, j int) bool {
			return calculateIncrementalCost(instance, facility, unassignedLocations[i], assigned) <
				calculateIncrementalCost(instance, facility, unassignedLocations[j], assigned)
		})

		location := unassignedLocations[len(unassignedLocations)-1]
		unassignedLocations = unassignedLocations[:len(unassignedLocations)-1]

		assigned = append(assigned, [2]int{facility, location})

		if stepsCounter != nil {
			*stepsCounter++
		}
	}

	solution := make([]int, size)
	for _, pair := range assigned {
		solution[pair[1]] = pair[0]
	}

	return solution
}

func facilityFlowSum(instance *qap.QAPInstance, facility int) int {
	sum := 0
	for i := 0; i < instance.Size; i++ {
		sum += instance.FlowMatrix[facility][i]
	}
	return sum
}

func calculateIncrementalCost(instance *qap.QAPInstance, facility, location int, assigned [][2]int) int {
	cost := 0
	for _, pair := range assigned {
		f, l := pair[0], pair[1]
		cost += instance.FlowMatrix[facility][f] * instance.DistanceMatrix[location][l]
		cost += instance.FlowMatrix[f][facility] * instance.DistanceMatrix[l][location]
	}
	return cost
}
