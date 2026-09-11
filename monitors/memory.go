package monitors

import (
	"context"
	"fmt"

	"github.com/shirou/gopsutil/v4/mem"
)

type MemoryMonitor struct {
}

func (memoryM *MemoryMonitor) GetName() string {
	return "Memory"
}

func (memoryM *MemoryMonitor) Check(ctx context.Context) string {
	virtualMemoryStat, err := mem.VirtualMemoryWithContext(ctx)

	if err != nil {
		return "ko có"
	}

	return fmt.Sprintf("%.2f %%", virtualMemoryStat.UsedPercent)
}
