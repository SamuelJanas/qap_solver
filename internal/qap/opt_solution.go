package qap

import (
	"bufio"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// OptimalSolutions maps instance names to their optimal fitness values
type OptimalSolutions map[string]int

func LoadOptimalSolutions(instancesDir string) (OptimalSolutions, error) {
	solutions := make(OptimalSolutions)

	entries, err := os.ReadDir(instancesDir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sln") {
			filePath := filepath.Join(instancesDir, entry.Name())

			file, err := os.Open(filePath)
			if err != nil {
				continue // Skip files that can't be opened
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			if scanner.Scan() { // Read the first line
				parts := strings.Fields(scanner.Text())
				if len(parts) >= 2 {
					_, valueStr := parts[0], parts[1] // Ignore size, extract optimal value
					value, err := strconv.Atoi(valueStr)
					if err == nil {
						instanceName := strings.TrimSuffix(entry.Name(), ".sln")
						solutions[instanceName] = value
					}
				}
			}
		}
	}

	return solutions, nil
}

func (o OptimalSolutions) GetOptimalSolution(instanceName string) int {
	// Extract base instance name without path and extension
	base := instanceName
	if idx := strings.LastIndex(base, "/"); idx >= 0 {
		base = base[idx+1:]
	}
	if idx := strings.LastIndex(base, "."); idx >= 0 {
		base = base[:idx]
	}

	return o[base]
}
