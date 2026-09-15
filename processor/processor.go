package processor

import (
	"context"
	"fmt"
	"os"
	"sort"
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
	var cpuList, memList []models.ProcStat
	procChan := make(chan models.ProcStat, len(processes))

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
					procChan <- models.ProcStat{
						PID:         proc.Pid,
						Name:        name,
						CPU:         cpuPercent,
						Memory:      memInfo.RSS,
						RamPercent:  ramPercent,
						RunningTime: runningTime,
					}

				}

			}
		}(p)
	}

	go func() {
		wg.Wait()
		close(procChan)
	}()

	for stat := range procChan {
		if stat.CPU > 1 {
			cpuList = append(cpuList, stat)
		}

		if stat.RamPercent > 1 {
			memList = append(memList, stat)
		}
	}

	sort.Slice(cpuList, func(i, j int) bool {
		return cpuList[i].CPU > cpuList[j].CPU
	})

	sort.Slice(memList, func(i, j int) bool {
		return memList[i].RamPercent > memList[j].RamPercent
	})

	output := "=== Top 5 Cpu dung nhieu nhat \n"
	for i := 0; i < len(cpuList) && i < 5; i++ {
		output += fmt.Sprintf("%d. [%d] %s - CPU: %.2f%% - RAM: %.2f MB (%.2f%%) - Running: %s \n",
			i+1,
			cpuList[i].PID,
			cpuList[i].Name,
			cpuList[i].CPU,
			float64(cpuList[i].Memory)/1024.0/1024.0,
			cpuList[i].RamPercent,
			cpuList[i].RunningTime,
		)
	}

	output += "=== Top 5 RAM dung nhieu nhat \n"
	for i := 0; i < len(memList) && i < 5; i++ {
		output += fmt.Sprintf("%d. [%d] %s - CPU: %.2f%% - RAM: %.2f MB (%.2f%%) - Running: %s \n",
			i+1,
			memList[i].PID,
			memList[i].Name,
			memList[i].CPU,
			float64(memList[i].Memory)/1024.0/1024.0,
			memList[i].RamPercent,
			memList[i].RunningTime,
		)
	}

	ExportToCSV(cpuList, memList)

	return output
}

func ExportToCSV(cpuList, memoryList []models.ProcStat) {
	file, err := os.OpenFile("process_stats.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644) // 644: chủ file được đọc và ghi, người khác chỉ đọc file mà ko dc sửa
	if err != nil {
		fmt.Println("[Export To CSV] ghi vao CSV ko thanh cong: ", err)
		return
	}
	defer file.Close()

	if stat, err := file.Stat(); err == nil && stat.Size() == 0 { // header
		file.WriteString("Timestamp,PID,Name,CPU (%), RAM (MB), RAM (%), Running Time \n")
	}

	timeStamp := time.Now().Format(time.RFC3339)

	for i := 0; i < len(cpuList) && i < 5; i++ {
		line := fmt.Sprintf("%s,%d,%s,%.2f,%.2f,%.2f,%s\n",
			timeStamp,
			cpuList[i].PID,
			cpuList[i].Name,
			cpuList[i].CPU,
			float64(cpuList[i].Memory)/1024.0/1024.0,
			cpuList[i].RamPercent,
			cpuList[i].RunningTime,
		)
		file.WriteString(line)
	}

	for i := 0; i < len(memoryList) && i < 5; i++ {
		line := fmt.Sprintf("%s,%d,%s,%.2f,%.2f,%.2f,%s\n",
			timeStamp,
			memoryList[i].PID,
			memoryList[i].Name,
			memoryList[i].CPU,
			float64(memoryList[i].Memory)/1024.0/1024.0,
			memoryList[i].RamPercent,
			memoryList[i].RunningTime,
		)
		file.WriteString(line)
	}
}
