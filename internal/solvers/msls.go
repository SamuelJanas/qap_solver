// TODO: implement this
package solvers

import (
    `qap_solver/internal/qap`
)

// Multiple Start Local Search (MSLS)
func MultipleStartLocalSearch(instance *qap.QAPInstance, iterations int) []int {
    bestSolution := RandomSolution(instance.Size)
    bestFitness := qap.CalculateFitness(instance, bestSolution)

    for i := 0; i < iterations; i++ {
        candidate := RandomSolution(instance.Size)
        candidateFitness := qap.CalculateFitness(instance, candidate)

        if candidateFitness < bestFitness {
            bestSolution = candidate
            bestFitness = candidateFitness
        }
    }

    return bestSolution
}
