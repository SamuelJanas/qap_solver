package pkg

import (
    "log"
    "time"
)

func TimeTrack(start time.Time, name string, logger *log.Logger) {
    duration := time.Since(start)
    logger.Printf("%s took %s", name, duration)
}
