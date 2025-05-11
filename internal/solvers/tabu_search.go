package solvers

import (
	"math/rand"
	"qap_solver/internal/qap"
	"sort"
	// "time"
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

type move struct {
	i, j       int
	newFitness int
	isTabu     bool
	aspiration bool
}

func (s *TabuSearchSolver) Solve(instance *qap.QAPInstance) SolverResult {
	n := instance.Size
	maxNoImprovement := s.P * n
	tabuTenure := n / 2
	tabuList := make([][]int, n)
	for i := range tabuList {
		tabuList[i] = make([]int, n)
	}

	current := RandomSolution(n)
	currentFitness := qap.CalculateFitness(instance, current)

	best := make([]int, n)
	copy(best, current)
	bestFitness := currentFitness

	noImprovementCounter := 0
	iteration := 0

	for noImprovementCounter < maxNoImprovement {
		iteration++
		var candidateMoves []move

		possibleSwaps := allSwaps(n)
		sampleSize := len(possibleSwaps) / 5
		rand.Shuffle(len(possibleSwaps), func(i, j int) {
			possibleSwaps[i], possibleSwaps[j] = possibleSwaps[j], possibleSwaps[i]
		})
		sampledSwaps := possibleSwaps[:sampleSize]

		for _, sw := range sampledSwaps {
			i, j := sw[0], sw[1]

			newSolution := make([]int, n)
			copy(newSolution, current)
			newSolution[i], newSolution[j] = newSolution[j], newSolution[i]

			newFitness := qap.CalculateFitness(instance, newSolution)

			isTabu := tabuList[i][current[j]] > iteration || tabuList[j][current[i]] > iteration
			aspiration := newFitness < bestFitness

			candidateMoves = append(candidateMoves, move{i, j, newFitness, isTabu, aspiration})
		}

		// Sort candidate moves by newFitness ascending (better first)
		sort.Slice(candidateMoves, func(i, j int) bool {
			return candidateMoves[i].newFitness < candidateMoves[j].newFitness
		})

		// Pick top 20% of candidates
		topSize := len(candidateMoves) / 5
		if topSize == 0 {
			topSize = 1
		}
		candidateMoves = candidateMoves[:topSize]

		// Choose the best allowed move (aspiration or non-tabu)
		var chosen move
		for _, m := range candidateMoves {
			if !m.isTabu || m.aspiration {
				chosen = m
				break
			}
		}
		// If no non-tabu or aspirational move, fallback to least tabu
		if chosen == (move{}) && len(candidateMoves) > 0 {
			chosen = candidateMoves[0]
		}

		// Apply the move
		i, j := chosen.i, chosen.j
		current[i], current[j] = current[j], current[i]
		currentFitness = chosen.newFitness

		// Update tabu list
		tabuList[i][current[i]] = iteration + tabuTenure
		tabuList[j][current[j]] = iteration + tabuTenure

		// Update best solution if needed
		if currentFitness < bestFitness {
			copy(best, current)
			bestFitness = currentFitness
			noImprovementCounter = 0
		} else {
			noImprovementCounter++
		}
	}

	return SolverResult{
		Solution: best,
		Fitness:  bestFitness,
	}
}

// allSwaps returns all unique i < j pairs
func allSwaps(n int) [][2]int {
	var swaps [][2]int
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			swaps = append(swaps, [2]int{i, j})
		}
	}
	return swaps
}
