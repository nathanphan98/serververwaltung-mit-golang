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

func (cpuM *CPUMonitor) Check(ctx context.Context) (string , bool) {
	cpuStat, err := cpu.PercentWithContext(ctx, 1*time.Second, false)

	if err != nil {
		return "ko có", false
	}

	return fmt.Sprintf("%.2f %%", cpuStat[0]), cpuStat[0] > 60
}
