package models

import "context"

type SystemStat struct {
	Name  string
	Value string
}

type IMonitor interface {
	GetName() string
	Check(ctx context.Context) string
}
