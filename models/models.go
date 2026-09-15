package models

import (
	"context"
	"sync"
	"time"
)

type SystemStat struct {
	Name  string
	Value string
}

type IMonitor interface {
	GetName() string
	Check(ctx context.Context) (string, bool)
}

type ProcStat struct {
	PID         int32
	Name        string
	CPU         float64
	Memory      uint64
	RamPercent  float64
	RunningTime time.Duration
}

var (
	Mtx   sync.Mutex
	Stats = map[string]SystemStat{}
)
