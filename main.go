package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"nathanphan.com/serververwaltung-mit-golang/models"
	"nathanphan.com/serververwaltung-mit-golang/processor"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cn := make(chan models.SystemStat)

	var wg sync.WaitGroup
	wg.Add(1)
	go processor.RunMonitor(ctx, &wg, cn)
	// dùng goroutine để tránh main chỉ để dùng làm 1 việc

	go func() {
		wg.Wait() // khi nào Wait mở block thì mới chạy xuống close()
		close(cn)
	}()

	for i := range cn {
		fmt.Println(i)
	}
}
