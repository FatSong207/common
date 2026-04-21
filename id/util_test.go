package id

import (
	"testing"
)

func TestInt64ToA(t *testing.T) {
	tests := []struct {
		num    int64
		expect string
	}{
		{0, "0"},
		{1, "1"},
		{10, "A"},
		{61, "z"},
		{62, "10"},
		{100, "1C"},
		{9223372036854775807, "AzL8n0Y58d7"}, // max int64
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := Int64ToA(tt.num)
			if got != tt.expect {
				t.Errorf("Int64ToA(%d) = %s, want %s", tt.num, got, tt.expect)
			}
		})
	}
}

func TestAToInt64(t *testing.T) {
	tests := []struct {
		s      string
		expect int64
	}{
		{"0", 0},
		{"1", 1},
		{"A", 10},
		{"z", 61},
		{"10", 62},
		{"1C", 100},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := AToInt64(tt.s)
			if got != tt.expect {
				t.Errorf("AToInt64(%s) = %d, want %d", tt.s, got, tt.expect)
			}
		})
	}
}

func TestRoundTrip(t *testing.T) {
	testNums := []int64{0, 1, 100, 1000, 999999, 9223372036854775800}

	for _, num := range testNums {
		t.Run("", func(t *testing.T) {
			str := Int64ToA(num)
			back := AToInt64(str)
			if back != num {
				t.Errorf("round trip failed for %d: got %d, str=%s", num, back, str)
			}
		})
	}
}

func BenchmarkInt64ToA(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Int64ToA(int64(i))
	}
}

func BenchmarkAToInt64(b *testing.B) {
	strs := []string{"0", "1A", "1Az", "1Azz", "1Azzz"}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		AToInt64(strs[i%len(strs)])
	}
}
