package monitors

import (
	"context"
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
)

type CPUMonitor struct {
}

func (cpuM *CPUMonitor) GetName() string {
	return "Cpu"
}

func (cpuM *CPUMonitor) Check(ctx context.Context) string {
	percent, err := cpu.PercentWithContext(ctx, 1*time.Second, false)

	if err != nil {
		return "ko có"
	}

	return fmt.Sprintf("%.2f %%", percent[0])
}
