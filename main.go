package main

import (
	// "fmt"
	// "log"
	// "os"
	"qap_solver/internal/qap"
	"qap_solver/internal/solvers"
	"qap_solver/pkg"
)

// Logger setup
var logger = pkg.NewLogger()

func main() {
	instance, err := qap.ReadInstance("instances/bur26a.dat")
	logger.Printf("Loaded instance: Size = %d\n", instance.Size)
	logger.Println("First row of Flow Matrix:", instance.FlowMatrix[0])
	logger.Println("First row of Distance Matrix:", instance.DistanceMatrix[0])

	if err != nil {
		logger.Fatalf("Failed to read instance: %v", err)
	}

	// Run a random solution
	randomSolution := solvers.RandomSolution(instance.Size)
	fitness := qap.CalculateFitness(instance, randomSolution)
	logger.Printf("Random solution fitness: %d", fitness)

	// Run Multiple Start Local Search
	/*
		mslsSolution := solvers.MultipleStartLocalSearch(instance, 10)
		fitness = qap.CalculateFitness(instance, mslsSolution)
		logger.Printf("MSLS solution fitness: %d", fitness) */
}
