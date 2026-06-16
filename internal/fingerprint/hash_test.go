package fingerprint

import "testing"

func TestHashFingerprint_Consistency(t *testing.T) {
	tests := []struct {
		freq1, freq2, delta int
	}{
		{0, 0, 0},
		{1, 2, 3},
		{65535, 65535, 65535},
		{100, 200, 10},
	}
	for _, tt := range tests {
		h1 := HashFingerprint(tt.freq1, tt.freq2, tt.delta)
		h2 := HashFingerprint(tt.freq1, tt.freq2, tt.delta)
		if h1 != h2 {
			t.Errorf("HashFingerprint(%d,%d,%d) not consistent: %d vs %d",
				tt.freq1, tt.freq2, tt.delta, h1, h2)
		}
	}
}

func TestHashFingerprint_Uniqueness(t *testing.T) {
	inputs := []struct {
		freq1, freq2, delta int
	}{
		{0, 0, 0},
		{1, 0, 0},
		{0, 1, 0},
		{0, 0, 1},
		{100, 200, 10},
		{200, 100, 10},
		{100, 200, 20},
	}
	seen := map[uint64]bool{}
	for _, in := range inputs {
		h := HashFingerprint(in.freq1, in.freq2, in.delta)
		if seen[h] {
			t.Errorf("duplicate hash %d for inputs (%d,%d,%d) and a previous input",
				h, in.freq1, in.freq2, in.delta)
		}
		seen[h] = true
	}
}

func TestHashFingerprint_BitPacking(t *testing.T) {
	tests := []struct {
		freq1, freq2, delta int
		want                uint64
	}{
		{0x1234, 0x5678, 0x9ABC, (uint64(0x1234) << 32) | (uint64(0x5678) << 16) | uint64(0x9ABC)},
		{0xFFFF, 0xFFFF, 0xFFFF, (uint64(0xFFFF) << 32) | (uint64(0xFFFF) << 16) | uint64(0xFFFF)},
		{0, 0, 0, 0},
		{1, 2, 3, (uint64(1) << 32) | (uint64(2) << 16) | uint64(3)},
	}
	for _, tt := range tests {
		got := HashFingerprint(tt.freq1, tt.freq2, tt.delta)
		if got != tt.want {
			t.Errorf("HashFingerprint(%#x,%#x,%#x) = %#x, want %#x",
				tt.freq1, tt.freq2, tt.delta, got, tt.want)
		}
	}
}

func TestHashFingerprint_Truncation(t *testing.T) {
	h := HashFingerprint(0x1FFFF, 0x2FFFF, 0x3FFFF)
	upper := (h >> 32) & 0xFFFF
	middle := (h >> 16) & 0xFFFF
	lower := h & 0xFFFF
	if upper != 0xFFFF {
		t.Errorf("freq1 truncated: upper=0x%X, want 0xFFFF", upper)
	}
	if middle != 0xFFFF {
		t.Errorf("freq2 truncated: middle=0x%X, want 0xFFFF", middle)
	}
	if lower != 0xFFFF {
		t.Errorf("delta truncated: lower=0x%X, want 0xFFFF", lower)
	}
}
