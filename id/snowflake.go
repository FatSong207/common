package id

import (
	"sync"
	"time"
)

const (
	epoch        = int64(1640966400000)    // 2022-01-01 00:00:00 UTC（毫秒）
	timestampLen = 41
	datacenterLen = 5
	machineLen    = 5
	seqLen        = 12

	maxTimestamp  = int64((1 << timestampLen) - 1)
	maxDatacenter = int64((1 << datacenterLen) - 1)
	maxMachine    = int64((1 << machineLen) - 1)
	maxSeq        = int64((1 << seqLen) - 1)

	timestampShift  = datacenterLen + machineLen + seqLen
	datacenterShift = machineLen + seqLen
	machineShift    = seqLen
)

// Snowflake 雪花算法ID生成器
type Snowflake struct {
	mu         sync.Mutex
	lastTs     int64
	lastSeq    int64
	datacenter int64
	machine    int64
}

// NewSnowflake 建立雪花算法生成器
// datacenter: 數據中心ID (0-31)
// machine: 工作機器ID (0-31)
func NewSnowflake(datacenter, machine int64) (*Snowflake, error) {
	if datacenter < 0 || datacenter > maxDatacenter {
		return nil, ErrInvalidDatacenter
	}
	if machine < 0 || machine > maxMachine {
		return nil, ErrInvalidMachine
	}

	return &Snowflake{
		lastTs:     0,
		lastSeq:    0,
		datacenter: datacenter,
		machine:    machine,
	}, nil
}

// Generate 生成唯一ID
func (s *Snowflake) Generate() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UnixMilli()
	ts := now - epoch

	if ts > maxTimestamp {
		panic("timestamp overflow")
	}

	if ts == s.lastTs {
		s.lastSeq = (s.lastSeq + 1) & maxSeq
		if s.lastSeq == 0 {
			// 序列號溢位，等待下一毫秒
			for ts <= s.lastTs {
				ts = time.Now().UnixMilli() - epoch
			}
		}
	} else if ts > s.lastTs {
		s.lastSeq = 0
	}

	s.lastTs = ts

	id := (ts << timestampShift) |
		(s.datacenter << datacenterShift) |
		(s.machine << machineShift) |
		s.lastSeq

	return id
}

// GenerateString 生成字符串格式ID
func (s *Snowflake) GenerateString() string {
	return Int64ToA(s.Generate())
}
