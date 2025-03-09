package metrics

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// RunMetrics stores metrics for a single algorithm run
type RunMetrics struct {
	InstanceName     string
	SolverName       string
	Run              int
	InitialFitness   int
	FinalFitness     int
	TimeElapsed      time.Duration
	StepsCount       int
	EvaluationsCount int
	SolutionsChecked int
	Solution         []int
}

// ExperimentMetrics collects metrics from multiple runs
type ExperimentMetrics struct {
	InstanceName string
	SolverName   string
	Runs         []RunMetrics
}

// MetricsCollector manages metrics for multiple experiments
type MetricsCollector struct {
	Experiments map[string]map[string]*ExperimentMetrics // Map[InstanceName][SolverName]
	OutputDir   string
}

// NewMetricsCollector creates a new metrics collector
func NewMetricsCollector(outputDir string) *MetricsCollector {
	// Create the output directory if it doesn't exist
	os.MkdirAll(outputDir, 0755)

	return &MetricsCollector{
		Experiments: make(map[string]map[string]*ExperimentMetrics),
		OutputDir:   outputDir,
	}
}

// AddRunMetrics adds a run's metrics to the collector
func (c *MetricsCollector) AddRunMetrics(metrics RunMetrics) {
	// Ensure we have a map for this instance
	if _, exists := c.Experiments[metrics.InstanceName]; !exists {
		c.Experiments[metrics.InstanceName] = make(map[string]*ExperimentMetrics)
	}

	// Ensure we have an experiment for this solver and instance
	instanceSolvers := c.Experiments[metrics.InstanceName]
	if _, exists := instanceSolvers[metrics.SolverName]; !exists {
		instanceSolvers[metrics.SolverName] = &ExperimentMetrics{
			InstanceName: metrics.InstanceName,
			SolverName:   metrics.SolverName,
			Runs:         make([]RunMetrics, 0),
		}
	}

	// Add the run metrics
	experiment := instanceSolvers[metrics.SolverName]
	experiment.Runs = append(experiment.Runs, metrics)
}

func (c *MetricsCollector) SaveToCSV() error {
	// Create a single results file
	dateStr := time.Now().Format("2006-01-02T15_04")
	resultsPath := filepath.Join(c.OutputDir, fmt.Sprintf("results_%s.csv", dateStr))
	resultsFile, err := os.Create(resultsPath)
	if err != nil {
		return err
	}
	defer resultsFile.Close()

	resultsWriter := csv.NewWriter(resultsFile)
	defer resultsWriter.Flush()

	// Write header (No Optimum, No Aggregated Stats)
	header := []string{
		"Instance", "Solver", "Run",
		"InitialFitness", "FinalFitness",
		"TimeMs", "Steps", "Evaluations", "SolutionsChecked",
		"Solution",
	}
	resultsWriter.Write(header)

	// Process each experiment
	for instanceName, solvers := range c.Experiments {
		for solverName, experiment := range solvers {
			// Write each run's details
			for i, run := range experiment.Runs {
				resultsWriter.Write([]string{
					instanceName, solverName, strconv.Itoa(i + 1),
					strconv.Itoa(run.InitialFitness),
					strconv.Itoa(run.FinalFitness),
					strconv.FormatFloat(float64(run.TimeElapsed.Milliseconds()), 'f', 2, 64),
					strconv.Itoa(run.StepsCount),
					strconv.Itoa(run.EvaluationsCount),
					strconv.Itoa(run.SolutionsChecked),
					fmt.Sprintf("%v", run.Solution),
				})
			}
		}
	}

	return nil
}
