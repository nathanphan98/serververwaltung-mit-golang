package processor

import (
	"context"
	"sync"
	"time"

	"nathanphan.com/serververwaltung-mit-golang/models"
	"nathanphan.com/serververwaltung-mit-golang/monitors"
)

func RunMonitor(ctx context.Context, wg *sync.WaitGroup, cn chan<- models.SystemStat) {
	defer wg.Done()

	ticker := time.NewTicker(2 * time.Second) // 2s tick 1 cái
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			cpuMonitor := monitors.CPUMonitor{}
			memoryMonitor := monitors.MemoryMonitor{}

			cn <- models.SystemStat{
				Name:  cpuMonitor.GetName(),
				Value: cpuMonitor.Check(ctx),
			}

			cn <- models.SystemStat{
				Name:  memoryMonitor.GetName(),
				Value: memoryMonitor.Check(ctx),
			}
		}
	}

}
