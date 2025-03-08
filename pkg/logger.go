package pkg

import (
    "log"
    "os"
)

func NewLogger() *log.Logger {
    return log.New(os.Stdout, "[QAP Solver] ", log.LstdFlags)
}
