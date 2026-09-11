package types

import (
	"sync/atomic"
	"time"
)

type ServerStats struct {
	Requests  atomic.Int64
	Responses atomic.Int64
	Errors    atomic.Int64
	StartTime time.Time
	StopTime  time.Time
}
