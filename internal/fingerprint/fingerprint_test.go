package fingerprint

import (
	"testing"

	peak "github.com/ravirraj/echoid/internal/peaks"
)

func TestGenerateFingerprints_Empty(t *testing.T) {
	got := GenerateFingerprints(nil)
	if got != nil {
		t.Errorf("expected nil, got %v", got)
	}

	got = GenerateFingerprints([]peak.Peak{})
	if got != nil {
		t.Errorf("expected nil for empty slice, got %v", got)
	}
}

func TestGenerateFingerprints_TwoPeaks(t *testing.T) {
	peaks := []peak.Peak{
		{TimeIndex: 0, FreqIndex: 100, Magnitude: 1.0},
		{TimeIndex: 10, FreqIndex: 200, Magnitude: 1.0},
	}

	fps := GenerateFingerprints(peaks)
	if len(fps) != 1 {
		t.Fatalf("expected 1 fingerprint, got %d: %+v", len(fps), fps)
	}

	expected := Fingerprint{
		Freq1:      (100 / 4) * 4,
		Freq2:      (200 / 4) * 4,
		DeltaTime:  (10 / 2) * 2,
		AnchorTime: 0,
	}
	if fps[0] != expected {
		t.Errorf("fingerprint: got %+v, want %+v", fps[0], expected)
	}
}

func TestGenerateFingerprints_SinglePeak(t *testing.T) {
	peaks := []peak.Peak{
		{TimeIndex: 0, FreqIndex: 100, Magnitude: 1.0},
	}

	fps := GenerateFingerprints(peaks)
	if len(fps) != 0 {
		t.Errorf("expected 0 fingerprints from single peak, got %d", len(fps))
	}
}

func TestGenerateFingerprints_SortsByTime(t *testing.T) {
	peaks := []peak.Peak{
		{TimeIndex: 20, FreqIndex: 200, Magnitude: 1.0},
		{TimeIndex: 10, FreqIndex: 100, Magnitude: 1.0},
	}

	fps := GenerateFingerprints(peaks)
	if len(fps) != 1 {
		t.Fatalf("expected 1 fingerprint, got %d", len(fps))
	}
	if fps[0].AnchorTime != 10 {
		t.Errorf("expected anchor time 10 (earliest peak), got %d", fps[0].AnchorTime)
	}
}

func TestGenerateFingerprints_FanOutLimit(t *testing.T) {
	peaks := make([]peak.Peak, 15)
	for i := range peaks {
		peaks[i] = peak.Peak{
			TimeIndex: i * 4,
			FreqIndex: 100 + i*10,
			Magnitude: 1.0,
		}
	}

	fps := GenerateFingerprints(peaks)

	for _, p := range peaks {
		anchorCount := 0
		for _, fp := range fps {
			if fp.AnchorTime == p.TimeIndex {
				anchorCount++
			}
		}
		if anchorCount > 10 {
			t.Errorf("anchor at time %d produced %d pairs, expected ≤10", p.TimeIndex, anchorCount)
		}
	}
}

func TestGenerateFingerprints_MinDeltaSkip(t *testing.T) {
	peaks := []peak.Peak{
		{TimeIndex: 0, FreqIndex: 100, Magnitude: 1.0},
		{TimeIndex: 0, FreqIndex: 200, Magnitude: 1.0},
	}

	fps := GenerateFingerprints(peaks)
	for _, fp := range fps {
		if fp.DeltaTime == 0 {
			t.Errorf("expected no fingerprints with DeltaTime=0, got %+v", fp)
		}
	}
}

func TestGenerateFingerprints_MaxDeltaBreak(t *testing.T) {
	peaks := []peak.Peak{
		{TimeIndex: 0, FreqIndex: 100, Magnitude: 1.0},
		{TimeIndex: 10, FreqIndex: 110, Magnitude: 1.0},
		{TimeIndex: 70, FreqIndex: 120, Magnitude: 1.0},
	}

	fps := GenerateFingerprints(peaks)
	for _, fp := range fps {
		if fp.DeltaTime > 60 {
			t.Errorf("expected no fingerprints with DeltaTime > 60, got %+v", fp)
		}
	}
}

func BenchmarkGenerateFingerprints(b *testing.B) {
	peaks := make([]peak.Peak, 500)
	for i := range peaks {
		peaks[i] = peak.Peak{
			TimeIndex: i * 2,
			FreqIndex: (i * 37) % 256,
			Magnitude: float64(i % 100),
		}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GenerateFingerprints(peaks)
	}
}
