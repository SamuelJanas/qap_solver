package solvers

import (
	"qap_solver/internal/qap"
)

type SolverResult struct {
	Solution []int
	Fitness  int
}

// Solver interface defines the contract that all solvers must implement
type Solver interface {
	// Name returns the name of the solver
	Name() string

	// Solve performs the solution process and returns the best solution found
	Solve(instance *qap.QAPInstance) SolverResult

	// Description returns a description of the solver
	Description() string
}
