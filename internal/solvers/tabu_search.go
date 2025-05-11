package solvers

import (
	"qap_solver/internal/qap"
)

type TabuSearchSolver struct {
	P int
}

func NewTabuSearchSolver(p int) *TabuSearchSolver {
	return &TabuSearchSolver{P: p}
}

func (s *TabuSearchSolver) Name() string {
	return "TabuSearch"
}

func (s *TabuSearchSolver) Description() string {
	return "Tabu Search with elite candidate list, aspiration criteria, and fixed tabu tenure"
}

func (s *TabuSearchSolver) Solve(instance *qap.QAPInstance) SolverResult {
	n := instance.Size
	neighbourhoodSize := n * (n - 1) / 2
	maxNoImprovement := s.P * n
	tabuTenure := neighbourhoodSize / 4
	tabuList := make([][]int, n)

	current := RandomSolution(n)
	best := make([]int, n)
	copy(best, current)

	currentFitness := qap.CalculateFitness(instance, current)
	bestFitness := currentFitness

	return SolverResult{
		Solution: best,
		Fitness:  bestFitness,
	}
}
