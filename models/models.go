package models

import (
	"context"
	"sync"
)

type SystemStat struct {
	Name  string
	Value string
}

type IMonitor interface {
	GetName() string
	Check(ctx context.Context) string
}

var (
	Mtx   sync.Mutex
	Stats = map[string]SystemStat{}
)
