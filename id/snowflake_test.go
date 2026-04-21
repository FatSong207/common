package id

import (
	"testing"
	"time"
)

func TestNewSnowflake(t *testing.T) {
	tests := []struct {
		name       string
		datacenter int64
		machine    int64
		wantErr    bool
	}{
		{"valid", 0, 0, false},
		{"valid max", 31, 31, false},
		{"invalid datacenter negative", -1, 0, true},
		{"invalid datacenter overflow", 32, 0, true},
		{"invalid machine negative", 0, -1, true},
		{"invalid machine overflow", 0, 32, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewSnowflake(tt.datacenter, tt.machine)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSnowflake() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGenerate(t *testing.T) {
	sf, _ := NewSnowflake(1, 1)

	id1 := sf.Generate()
	id2 := sf.Generate()

	if id1 >= id2 {
		t.Errorf("id1 should be less than id2, got id1=%d, id2=%d", id1, id2)
	}

	if id1 <= 0 || id2 <= 0 {
		t.Error("id should be positive")
	}
}

func TestGenerateUnique(t *testing.T) {
	sf, _ := NewSnowflake(1, 1)

	ids := make(map[int64]bool)
	count := 100000

	for i := 0; i < count; i++ {
		id := sf.Generate()
		if ids[id] {
			t.Fatalf("duplicate id found: %d", id)
		}
		ids[id] = true
	}
}

func TestGenerateString(t *testing.T) {
	sf, _ := NewSnowflake(1, 1)

	str1 := sf.GenerateString()
	str2 := sf.GenerateString()

	if str1 == "" || str2 == "" {
		t.Error("string id should not be empty")
	}

	if str1 == str2 {
		t.Error("string ids should be unique")
	}

	if AToInt64(str1) >= AToInt64(str2) {
		t.Errorf("str1 should be less than str2, got str1=%s(%d), str2=%s(%d)",
			str1, AToInt64(str1), str2, AToInt64(str2))
	}
}

func TestConcurrency(t *testing.T) {
	sf, _ := NewSnowflake(1, 1)

	ids := make(map[int64]bool)
	ch := make(chan int64, 1000)
	done := make(chan bool, 100)

	// 100個 goroutine，每個生成100個ID
	for i := 0; i < 100; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				ch <- sf.Generate()
			}
			done <- true
		}()
	}

	// 等待所有 goroutine 完成
	for i := 0; i < 100; i++ {
		<-done
	}

	// 檢查所有ID的唯一性
	for i := 0; i < 10000; i++ {
		id := <-ch
		if ids[id] {
			t.Fatalf("concurrent: duplicate id found: %d", id)
		}
		ids[id] = true
	}
}

func BenchmarkGenerate(b *testing.B) {
	sf, _ := NewSnowflake(1, 1)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sf.Generate()
	}
}

func BenchmarkGenerateString(b *testing.B) {
	sf, _ := NewSnowflake(1, 1)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sf.GenerateString()
	}
}

func TestTimestampExtraction(t *testing.T) {
	sf, _ := NewSnowflake(1, 2)
	before := time.Now().UnixMilli()
	id := sf.Generate()
	after := time.Now().UnixMilli()

	// 提取時間戳部分 (最高41位)
	extractedTs := id >> timestampShift
	realTs := extractedTs + epoch

	if realTs < before || realTs > after {
		t.Errorf("extracted timestamp out of range: got %d, expected between %d and %d",
			realTs, before, after)
	}
}
