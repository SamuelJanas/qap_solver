package solvers

import (
	"fmt"
	"strconv"
	"strings"
)

// SolverFactory creates solver instances based on configuration strings
type SolverFactory struct {
	// Registry of available solvers
	solverCreators map[string]func(args []string) (Solver, error)
}

// NewSolverFactory creates a new factory with registered solvers
func NewSolverFactory() *SolverFactory {
	factory := &SolverFactory{
		solverCreators: make(map[string]func(args []string) (Solver, error)),
	}

	// Register the built-in solvers
	factory.Register("random", factory.createRandomSolver)
	factory.Register("localsearch", factory.createLocalSearchSolver)
	factory.Register("greedy", factory.createGreedySolver)
	factory.Register("steepest", factory.createSteepestSolver)
	factory.Register("randomwalk", factory.createRandomWalkSolver)
	factory.Register("heuristic", factory.createHeuristicSolver)

	return factory
}

// Register adds a new solver type to the factory
func (f *SolverFactory) Register(name string, creator func(args []string) (Solver, error)) {
	f.solverCreators[strings.ToLower(name)] = creator
}

// Create instantiates a solver based on a configuration string
// Format: "solverName:param1=value1,param2=value2,..."
func (f *SolverFactory) Create(config string) (Solver, error) {
	parts := strings.SplitN(config, ":", 2)
	solverType := strings.ToLower(parts[0])

	creator, exists := f.solverCreators[solverType]
	if !exists {
		return nil, fmt.Errorf("unknown solver type: %s", solverType)
	}

	// Parse arguments if provided
	var args []string
	if len(parts) > 1 && parts[1] != "" {
		args = strings.Split(parts[1], ",")
	}

	return creator(args)
}

func (f *SolverFactory) ListAvailable() []string {
	var result []string

	result = append(result, "Available solvers:")
	result = append(result, "  random:iterations=1000 - Random solution generator with 1000 iterations")
	result = append(result, "  localsearch:maxIter=10000,maxNonImproving=1000,restarts=5 - Local search with parameters")
	result = append(result, "  greedy:maxIter=10000,restarts=5 - Greedy search with max iterations and restarts")
	result = append(result, "  steepest:maxIter=10000,restarts=5 - Steepest ascent search with max iterations and restarts")
	result = append(result, "  randomwalk:maxIter=10000,restarts=5 - Random walk search with max iterations 10000 and random restarts")
	result = append(result, "  heuristic:maxIter=10000,restarts=5 - Heuristic search with max iterations 1000 and random restarts:")

	return result
}

/*
------------------------------------------
 Helper functions to create specific solvers
------------------------------------------
*/

func (f *SolverFactory) createRandomSolver(args []string) (Solver, error) {
	iterations := 1000 // Default value

	// Process arguments
	for _, arg := range args {
		parts := strings.SplitN(arg, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.ToLower(parts[0])
		value := parts[1]

		if key == "iterations" {
			if i, err := strconv.Atoi(value); err == nil && i > 0 {
				iterations = i
			}
		}
	}

	return NewRandomSolver(iterations), nil
}

func (f *SolverFactory) createLocalSearchSolver(args []string) (Solver, error) {
	// defaults
	maxIterations := 10000
	maxNonImproving := 1000
	randomRestarts := 5

	// Process arguments
	for _, arg := range args {
		parts := strings.SplitN(arg, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.ToLower(parts[0])
		value := parts[1]

		switch key {
		case "maxiter":
			if i, err := strconv.Atoi(value); err == nil && i > 0 {
				maxIterations = i
			}
		case "maxnonimproving":
			if i, err := strconv.Atoi(value); err == nil && i > 0 {
				maxNonImproving = i
			}
		case "restarts":
			if i, err := strconv.Atoi(value); err == nil && i >= 0 {
				randomRestarts = i
			}
		}
	}

	return NewLocalSearchSolver(maxIterations, maxNonImproving, randomRestarts), nil
}

func (f *SolverFactory) createGreedySolver(args []string) (Solver, error) {
	maxIterations := 10000
	randomRestarts := 5

	for _, arg := range args {
		parts := strings.SplitN(arg, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.ToLower(parts[0])
		value := parts[1]
		switch key {
		case "maxiter":
			if i, err := strconv.Atoi(value); err == nil && i > 0 {
				maxIterations = i
			}
		case "restarts":
			if i, err := strconv.Atoi(value); err == nil && i >= 0 {
				randomRestarts = i
			}
		}
	}
	return NewGreedySolver(maxIterations, randomRestarts), nil
}

func (f *SolverFactory) createSteepestSolver(args []string) (Solver, error) {
	maxIterations := 10000
	randomRestarts := 5

	for _, arg := range args {
		parts := strings.SplitN(arg, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.ToLower(parts[0])
		value := parts[1]
		switch key {
		case "maxiter":
			if i, err := strconv.Atoi(value); err == nil && i > 0 {
				maxIterations = i
			}
		case "restarts":
			if i, err := strconv.Atoi(value); err == nil && i >= 0 {
				randomRestarts = i
			}
		}
	}
	return NewSteepestSolver(maxIterations, randomRestarts), nil
}

func (f *SolverFactory) createRandomWalkSolver(args []string) (Solver, error) {
	maxIterations := 10000
	randomRestarts := 5

	for _, arg := range args {
		parts := strings.SplitN(arg, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.ToLower(parts[0])
		value := parts[1]
		switch key {
		case "maxiter":
			if i, err := strconv.Atoi(value); err == nil && i > 0 {
				maxIterations = i
			}
		case "restarts":
			if i, err := strconv.Atoi(value); err == nil && i >= 0 {
				randomRestarts = i
			}
		}
	}
	return NewRandomWalkSolver(maxIterations, randomRestarts), nil
}

func (f *SolverFactory) createHeuristicSolver(args []string) (Solver, error) {
	return NewGreedyConstructionSolver(), nil
}
