package qap

func CalculateFitness(instance *QAPInstance, solution []int) int {
	size := instance.Size
	totalCost := 0

	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			totalCost += instance.FlowMatrix[i][j] * instance.DistanceMatrix[solution[i]][solution[j]]
		}
	}

	return totalCost
}
