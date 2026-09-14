package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"nathanphan.com/serververwaltung-mit-golang/models"
	"nathanphan.com/serververwaltung-mit-golang/monitors"
	"nathanphan.com/serververwaltung-mit-golang/processor"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cn := make(chan models.SystemStat)

	monitors := []models.IMonitor{ // 1 slice chứa 2 struct kiểu interface Monitor
		&monitors.CPUMonitor{}, // các receiver của struct này yêu cầu * nên giá trị của kiểu interface models.Monitor phải là kiểu con trỏ
		&monitors.MemoryMonitor{},
		&monitors.NetMonitor{},
		&monitors.DiskMonitor{},
	}

	var wg sync.WaitGroup
	for _, m := range monitors {
		wg.Add(1)
		go processor.RunMonitor(ctx, &wg, cn, m)
		// dùng goroutine để tránh main chỉ để dùng làm 1 việc
	}

	go func() {
		wg.Wait() // khi nào Wait mở block thì mới chạy xuống close()
		close(cn)
	}()

	for i := range cn {
		fmt.Println(i)
	}
}
