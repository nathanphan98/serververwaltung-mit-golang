package monitors

import (
	"context"
	"fmt"

	"github.com/shirou/gopsutil/v4/disk"
)

type DiskMonitor struct {
}

func (cpuM *DiskMonitor) GetName() string {
	return "Disk"
}

func (cpuM *DiskMonitor) Check(ctx context.Context) string {
	path := "/"
	diskStat, err := disk.UsageWithContext(ctx, path)

	if err != nil {
		return "ko có"
	}

	return fmt.Sprintf(" %.2f%% used", diskStat.UsedPercent )
}
