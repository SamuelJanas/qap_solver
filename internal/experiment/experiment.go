package experiment

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"qap_solver/internal/metrics"
	"qap_solver/internal/qap"
	"qap_solver/internal/solvers"
	"strings"
)

// ExperimentConfig holds configuration for running experiments
type ExperimentConfig struct {
	InstanceSample  int
	InstancesDir    string
	OutputDir       string
	Solvers         []solvers.Solver
	RunsPerInstance int
	Logger          *log.Logger
}

// RunAll runs experiments on all instances with all solvers
func RunAll(config ExperimentConfig) error {
	// Create metrics collector
	metricsCollector := metrics.NewMetricsCollector(config.OutputDir)

	// Get list of instance files
	instanceFiles, err := findInstanceFiles(config.InstancesDir)
	if err != nil {
		return fmt.Errorf("error finding instance files: %v", err)
	}

	if len(instanceFiles) == 0 {
		return fmt.Errorf("no instance files found in %s", config.InstancesDir)
	}

	config.Logger.Printf("Found %d instance files", len(instanceFiles))

	if config.InstanceSample > len(instanceFiles) {
		return fmt.Errorf("sample was provided, but sample exceeds the total number of instance files")
	}

	if config.InstanceSample > 0 {
		config.Logger.Printf("Using only first %d instances", config.InstanceSample)
		instanceFiles = instanceFiles[:config.InstanceSample]
	}

	// Process each instance
	for _, instanceFile := range instanceFiles {
		instanceName := filepath.Base(instanceFile)
		config.Logger.Printf("Processing instance: %s", instanceName)

		// Load the instance
		instance, err := qap.ReadInstance(instanceFile)
		if err != nil {
			config.Logger.Printf("Error loading instance %s: %v", instanceName, err)
			continue
		}

		// Run each solver multiple times
		for _, solver := range config.Solvers {
			config.Logger.Printf("Running %s on %s (%d runs)", solver.Name(), instanceName, config.RunsPerInstance)

			// Run solver multiple times
			for run := 1; run <= config.RunsPerInstance; run++ {
				config.Logger.Printf("  Run %d/%d", run, config.RunsPerInstance)

				// Check if the solver supports metrics collection
				if metricsSolver, ok := solver.(MetricsSolver); ok {
					metricsSolver.SolveWithMetrics(instance, metricsCollector, instanceName, run)
				} else {
					// Run standard solver and collect basic metrics
					result := solver.Solve(instance)
					config.Logger.Printf("    Fitness: %d", result.Fitness)
				}
			}
		}
	}

	// Save all metrics to CSV
	err = metricsCollector.SaveToCSV()
	if err != nil {
		return fmt.Errorf("error saving metrics: %v", err)
	}

	config.Logger.Printf("Experiments completed. Results saved to %s", config.OutputDir)
	return nil
}

// MetricsSolver extends the Solver interface with metrics collection
type MetricsSolver interface {
	solvers.Solver
	SolveWithMetrics(instance *qap.QAPInstance, metricsCollector *metrics.MetricsCollector,
		instanceName string, runNumber int) solvers.SolverResult
}

// Helper function to find all instance files in a directory
func findInstanceFiles(dir string) ([]string, error) {
	var files []string

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		// Only include files with .dat extension or other QAP formats
		if strings.HasSuffix(name, ".dat") || strings.HasSuffix(name, ".qap") {
			files = append(files, filepath.Join(dir, name))
		}
	}

	return files, nil
}
