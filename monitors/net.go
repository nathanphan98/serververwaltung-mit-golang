package monitors

import (
	"context"
	"fmt"

	"github.com/shirou/gopsutil/v4/net"
)

type NetMonitor struct {
}

func (cpuM *NetMonitor) GetName() string {
	return "Net"
}

func (cpuM *NetMonitor) Check(ctx context.Context) string {
	netStat, err := net.IOCountersWithContext(ctx, false)
	// https://www.youtube.com/watch?v=clN5qlfpsWE&list=PLTasIXUHepx0tYvXfFsbEl69VbMCZDH1i&index=94 (demo có hướng dẫn đọc tài liệu)

	if err != nil {
		return "ko có"
	}

	return fmt.Sprintf("Send: %d KB, Recv: %d KB", netStat[0].BytesSent/1024, netStat[0].BytesRecv/1024) // đổi byte ra kb
}
