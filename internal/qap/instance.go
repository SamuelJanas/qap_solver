package qap

import (
	"os"
	"strconv"
	"strings"
)

type QAPInstance struct {
	Size           int
	FlowMatrix     [][]int
	DistanceMatrix [][]int
}


func ReadInstance(filename string) (*QAPInstance, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	size, _ := strconv.Atoi(strings.TrimSpace(lines[0]))

	flowMatrix := make([][]int, size)
	distMatrix := make([][]int, size)

	// Initialize matrices
	for i := 0; i < size; i++ {
		flowMatrix[i] = make([]int, size)
		distMatrix[i] = make([]int, size)
	}

	// Parse flow matrix (starts at line 2)
	lineIndex := 2
	for i := 0; i < size; i++ {
		flowMatrix[i] = parseLine(lines[lineIndex+i])
	}

	// Parse distance matrix (starts after flow matrix)
	lineIndex = 3 + size
	for i := 0; i < size; i++ {
		distMatrix[i] = parseLine(lines[lineIndex+i])
	}

	return &QAPInstance{
		Size:           size,
		FlowMatrix:     flowMatrix,
		DistanceMatrix: distMatrix,
	}, nil
}

func parseLine(line string) []int {
	parts := strings.Fields(line)
	result := make([]int, len(parts))
	for i, v := range parts {
		result[i], _ = strconv.Atoi(v)
	}
	return result
}
