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

func (memoryM *MemoryMonitor) Check(ctx context.Context) (string , bool) {
	virtualMemoryStat, err := mem.VirtualMemoryWithContext(ctx)

	if err != nil {
		return "ko có", false
	}

	return fmt.Sprintf("%.2f %%", virtualMemoryStat.UsedPercent), virtualMemoryStat.UsedPercent > 60
}
