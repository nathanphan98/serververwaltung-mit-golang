package processor

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/process"
	"nathanphan.com/serververwaltung-mit-golang/models"
)

func RunMonitor(ctx context.Context, wg *sync.WaitGroup, cn chan<- models.SystemStat, m models.IMonitor) {
	defer wg.Done()

	ticker := time.NewTicker(2 * time.Second) // 2s tick 1 cái
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			cn <- models.SystemStat{
				Name:  m.GetName(),
				Value: m.Check(ctx),
			}
		}
	}
}

func GetTopProcesses(ctx context.Context) string {
	vmStat, err := mem.VirtualMemoryWithContext(ctx)
	if err != nil {
		return fmt.Sprintf("[Get Top Processes] khong nhan duoc Memory info: %v \n", err)
	}

	totalMem := vmStat.Total

	processes, err := process.ProcessesWithContext(ctx)
	if err != nil {
		return fmt.Sprintf("[Get Top Processes] khong nhan duoc danh sach process: %v \n", err)
	}

	var wg sync.WaitGroup
	for _, p := range processes {
		wg.Add(1)
		go func(proc *process.Process) {
			defer wg.Done()
			select {
			case <-ctx.Done():
				return
			default:
				name, err := proc.NameWithContext(ctx)
				if err != nil {
					return
				}

				cpuPercent, err := proc.CPUPercentWithContext(ctx)
				if err != nil {
					return
				}

				// RSS: Resident set size => Ram thực sự mà tiến trình sử dụng
				// VMS: Vitual Memory Size => Tổng bộ nhớ ảo mà hdh cấp phát cho tiến trình
				memInfo, err := proc.MemoryInfoExWithContext(ctx)
				if err != nil {
					return
				}

				ramPercent := (float64(memInfo.RSS) / float64(totalMem)) * 100

				createTime, err := proc.CreateTimeWithContext(ctx) // => milisecond
				if err != nil {
					return
				}

				// Time Unix là tổng số giây từ mốc thời gian 1970-0-0 00:00:00 UTC đến hiện tại
				// Vì createTime là milisecond nên phải chia 1000 để ra second
				runningTime := time.Since(time.Unix(createTime/1000, 0))

				if cpuPercent > 1 || ramPercent > 1 {
					procStat := models.ProcStat{
						PID: proc.Pid,
						Name: name,
						CPU: cpuPercent,
						Memory: memInfo.RSS,
						RamPercent: ramPercent,
						RunningTime: runningTime,
					}

					fmt.Printf("%v \n", procStat)
				}


			}
		}(p)
	}
	wg.Wait()
	return ""
}
