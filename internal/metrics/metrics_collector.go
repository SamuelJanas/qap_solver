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
	OptimalFitness   int
	TimeElapsed      time.Duration
	StepsCount       int
	EvaluationsCount int
	SolutionsChecked int
	Solution         []int
}

// GetGapFromOptimum returns the gap between the found solution and the optimum as a percentage
func (m *RunMetrics) GetGapFromOptimum() float64 {
	if m.OptimalFitness <= 0 {
		return 0 // No optimum provided
	}
	return float64(m.FinalFitness-m.OptimalFitness) / float64(m.OptimalFitness) * 100
}

// GetEfficiency returns a measure of efficiency (quality over time)
// We use improvement per second as our efficiency metric
func (m *RunMetrics) GetEfficiency() float64 {
	improvement := float64(m.InitialFitness - m.FinalFitness)
	if improvement <= 0 {
		return 0
	}
	seconds := m.TimeElapsed.Seconds()
	if seconds <= 0 {
		return 0
	}
	return improvement / seconds
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

// SaveToCSV saves all metrics to CSV files
func (c *MetricsCollector) SaveToCSV() error {
	// Create summary file
	summaryPath := filepath.Join(c.OutputDir, "summary.csv")
	summaryFile, err := os.Create(summaryPath)
	if err != nil {
		return err
	}
	defer summaryFile.Close()

	summaryWriter := csv.NewWriter(summaryFile)
	defer summaryWriter.Flush()

	// Write summary header
	summaryHeader := []string{
		"Instance", "Solver", "RunCount",
		"AvgInitialFitness", "BestInitialFitness", "StdDevInitialFitness",
		"AvgFinalFitness", "BestFinalFitness", "StdDevFinalFitness",
		"AvgGapFromOptimum", "BestGapFromOptimum",
		"AvgTimeMs", "AvgSteps", "AvgEvaluations", "AvgSolutionsChecked",
		"AvgEfficiency",
	}
	summaryWriter.Write(summaryHeader)

	// For each instance and solver, create detailed run file and add summary data
	for instanceName, solvers := range c.Experiments {
		for solverName, experiment := range solvers {
			// Create detailed runs file
			runsPath := filepath.Join(c.OutputDir, fmt.Sprintf("%s_%s_runs.csv", instanceName, solverName))
			runsFile, err := os.Create(runsPath)
			if err != nil {
				return err
			}

			runsWriter := csv.NewWriter(runsFile)

			// Write detailed header
			runsHeader := []string{
				"Run", "InitialFitness", "FinalFitness", "OptimalFitness",
				"GapFromOptimum", "TimeMs", "Steps", "Evaluations", "SolutionsChecked",
				"Efficiency", "Solution",
			}
			runsWriter.Write(runsHeader)

			// Track statistics for summary
			var sumInitial, sumFinal, sumTime, sumSteps, sumEvals, sumChecked, sumEfficiency float64
			var sumGap float64
			bestInitial := -1
			bestFinal := -1
			bestGap := -1.0

			// For calculating standard deviations
			initialValues := make([]float64, 0, len(experiment.Runs))
			finalValues := make([]float64, 0, len(experiment.Runs))

			// Write each run to the detailed file and collect stats
			for i, run := range experiment.Runs {
				gap := run.GetGapFromOptimum()
				efficiency := run.GetEfficiency()
				timeMs := float64(run.TimeElapsed.Milliseconds())

				// Update stats
				sumInitial += float64(run.InitialFitness)
				sumFinal += float64(run.FinalFitness)
				sumTime += timeMs
				sumSteps += float64(run.StepsCount)
				sumEvals += float64(run.EvaluationsCount)
				sumChecked += float64(run.SolutionsChecked)
				sumGap += gap
				sumEfficiency += efficiency

				initialValues = append(initialValues, float64(run.InitialFitness))
				finalValues = append(finalValues, float64(run.FinalFitness))

				if bestInitial == -1 || run.InitialFitness < bestInitial {
					bestInitial = run.InitialFitness
				}

				if bestFinal == -1 || run.FinalFitness < bestFinal {
					bestFinal = run.FinalFitness
				}

				if bestGap == -1 || gap < bestGap {
					bestGap = gap
				}

				// Format solution as a string
				solutionStr := fmt.Sprintf("%v", run.Solution)

				// Write run to detailed file
				runsWriter.Write([]string{
					strconv.Itoa(i + 1),
					strconv.Itoa(run.InitialFitness),
					strconv.Itoa(run.FinalFitness),
					strconv.Itoa(run.OptimalFitness),
					strconv.FormatFloat(gap, 'f', 4, 64),
					strconv.FormatFloat(timeMs, 'f', 2, 64),
					strconv.Itoa(run.StepsCount),
					strconv.Itoa(run.EvaluationsCount),
					strconv.Itoa(run.SolutionsChecked),
					strconv.FormatFloat(efficiency, 'f', 4, 64),
					solutionStr,
				})
			}

			runsWriter.Flush()
			runsFile.Close()

			// Calculate averages
			numRuns := float64(len(experiment.Runs))
			avgInitial := sumInitial / numRuns
			avgFinal := sumFinal / numRuns
			avgTime := sumTime / numRuns
			avgSteps := sumSteps / numRuns
			avgEvals := sumEvals / numRuns
			avgChecked := sumChecked / numRuns
			avgGap := sumGap / numRuns
			avgEfficiency := sumEfficiency / numRuns

			// Calculate standard deviations
			stdDevInitial := calculateStdDev(initialValues, avgInitial)
			stdDevFinal := calculateStdDev(finalValues, avgFinal)

			// Write summary record
			summaryWriter.Write([]string{
				instanceName,
				solverName,
				strconv.Itoa(len(experiment.Runs)),
				strconv.FormatFloat(avgInitial, 'f', 2, 64),
				strconv.Itoa(bestInitial),
				strconv.FormatFloat(stdDevInitial, 'f', 2, 64),
				strconv.FormatFloat(avgFinal, 'f', 2, 64),
				strconv.Itoa(bestFinal),
				strconv.FormatFloat(stdDevFinal, 'f', 2, 64),
				strconv.FormatFloat(avgGap, 'f', 4, 64),
				strconv.FormatFloat(bestGap, 'f', 4, 64),
				strconv.FormatFloat(avgTime, 'f', 2, 64),
				strconv.FormatFloat(avgSteps, 'f', 2, 64),
				strconv.FormatFloat(avgEvals, 'f', 2, 64),
				strconv.FormatFloat(avgChecked, 'f', 2, 64),
				strconv.FormatFloat(avgEfficiency, 'f', 4, 64),
			})
		}
	}

	return nil
}

// Helper function to calculate standard deviation
func calculateStdDev(values []float64, mean float64) float64 {
	if len(values) <= 1 {
		return 0
	}

	var variance float64
	for _, v := range values {
		variance += (v - mean) * (v - mean)
	}

	variance /= float64(len(values) - 1)
	return float64(float64(int(variance*100)) / 100)
}
