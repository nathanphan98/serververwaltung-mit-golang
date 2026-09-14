package processor

import (
	"context"
	"sync"
	"time"

	"nathanphan.com/serververwaltung-mit-golang/models"
)

func RunMonitor(ctx context.Context, wg *sync.WaitGroup, cn chan<- models.SystemStat, m models.IMonitor) {
	defer wg.Done()

	ticker := time.NewTicker(2 * time.Second) // 2s tick 1 cái
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			cn <- models.SystemStat{
				Name:  m.GetName(),
				Value: m.Check(ctx),
			}
		}
	}
}
