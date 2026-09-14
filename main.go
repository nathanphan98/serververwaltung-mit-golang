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
		for i := range cn {
			models.Mtx.Lock()
			models.Stats[i.Name] = i
			models.Mtx.Unlock()
		}
	}()

	printTicker := time.NewTicker(4 * time.Second)

	go func() {
		for range printTicker.C { // đợi 4s trước khi có data từ goroutine RunMonitor
			fmt.Println("=== System status ===")
			for _, stat := range models.Stats {
				models.Mtx.Lock()
				fmt.Printf("[%s] %s \n", stat.Name, stat.Value)
				models.Mtx.Unlock()
			}
		}

	}()

	time.Sleep(30 * time.Second)
	cancel()
	wg.Wait() // khi nào Wait mở block thì mới chạy xuống close()
	close(cn)
	printTicker.Stop()
}
