package qap

type Solver interface {
    Solve(instance *QAPInstance) []int
}
