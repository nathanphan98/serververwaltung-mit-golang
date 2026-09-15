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

func (cpuM *DiskMonitor) Check(ctx context.Context) (string , bool) {
	path := "/"
	diskStat, err := disk.UsageWithContext(ctx, path)

	if err != nil {
		return "ko có", false
	}

	return fmt.Sprintf(" %.2f%% used", diskStat.UsedPercent ), diskStat.UsedPercent > 60
}
