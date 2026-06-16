package peak

import (
	"testing"
)

func TestDetectPeaks_Empty(t *testing.T) {
	got := DetectPeaks([][]float64{})
	if got != nil {
		t.Errorf("expected nil, got %v", got)
	}
}

func TestDetectPeaks_SinglePeak(t *testing.T) {
	spec := make([][]float64, 20)
	for i := range spec {
		spec[i] = make([]float64, 30)
	}
	spec[10][15] = 100.0

	peaks := DetectPeaks(spec)
	if len(peaks) != 1 {
		t.Fatalf("expected 1 peak, got %d: %+v", len(peaks), peaks)
	}
	if peaks[0].TimeIndex != 10 {
		t.Errorf("TimeIndex: got %d, want 10", peaks[0].TimeIndex)
	}
	if peaks[0].FreqIndex != 15 {
		t.Errorf("FreqIndex: got %d, want 15", peaks[0].FreqIndex)
	}
	if peaks[0].Magnitude != 100.0 {
		t.Errorf("Magnitude: got %f, want 100.0", peaks[0].Magnitude)
	}
}

func TestDetectPeaks_MultiplePeaks(t *testing.T) {
	spec := make([][]float64, 20)
	for i := range spec {
		spec[i] = make([]float64, 30)
	}
	spec[5][10] = 100.0
	spec[14][22] = 200.0

	peaks := DetectPeaks(spec)
	if len(peaks) != 2 {
		t.Fatalf("expected 2 peaks, got %d: %+v", len(peaks), peaks)
	}

	found := map[int]bool{}
	for _, p := range peaks {
		if p.TimeIndex == 5 && p.FreqIndex == 10 && p.Magnitude == 100.0 {
			found[1] = true
		}
		if p.TimeIndex == 14 && p.FreqIndex == 22 && p.Magnitude == 200.0 {
			found[2] = true
		}
	}
	if !found[1] {
		t.Error("missing peak at (time=5, freq=10, mag=100)")
	}
	if !found[2] {
		t.Error("missing peak at (time=14, freq=22, mag=200)")
	}
}

func TestDetectPeaks_Boundary(t *testing.T) {
	spec := make([][]float64, 20)
	for i := range spec {
		spec[i] = make([]float64, 30)
	}
	spec[1][15] = 100.0
	spec[10][2] = 100.0

	peaks := DetectPeaks(spec)
	if len(peaks) != 0 {
		t.Errorf("expected 0 boundary peaks, got %d: %+v", len(peaks), peaks)
	}
}

func TestDetectPeaks_MaxPeaksPerFrame(t *testing.T) {
	spec := make([][]float64, 20)
	for i := range spec {
		spec[i] = make([]float64, 30)
	}
	for f := 6; f <= 25; f++ {
		spec[10][f] = float64(100 - f)
	}

	peaks := DetectPeaks(spec)
	peaksAtTime10 := 0
	for _, p := range peaks {
		if p.TimeIndex == 10 {
			peaksAtTime10++
		}
	}
	if peaksAtTime10 > 8 {
		t.Errorf("at most 8 peaks per frame, got %d", peaksAtTime10)
	}
}

func BenchmarkDetectPeaks(b *testing.B) {
	spec := make([][]float64, 200)
	for i := range spec {
		spec[i] = make([]float64, 256)
		for j := range spec[i] {
			spec[i][j] = float64((i*17 + j*31) % 100)
		}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		DetectPeaks(spec)
	}
}
