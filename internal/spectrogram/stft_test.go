package spectrogram

import (
	"math"
	"testing"
)

func TestGenerateSpectrogram_Silence(t *testing.T) {
	samples := make([]float64, 4096)
	spec := GenerateSpectrogram(samples)

	if len(spec) == 0 {
		t.Fatal("expected non-empty spectrogram")
	}

	for i, frame := range spec {
		for j, mag := range frame {
			if mag != 0.0 {
				t.Errorf("frame %d, bin %d: expected 0.0 for silence, got %f", i, j, mag)
			}
		}
	}
}

func TestGenerateSpectrogram_SineWave(t *testing.T) {
	sampleRate := 44100.0
	freq := 440.0
	n := 8192
	samples := make([]float64, n)
	for i := range samples {
		samples[i] = math.Sin(2 * math.Pi * freq * float64(i) / sampleRate)
	}

	spec := GenerateSpectrogram(samples)
	if len(spec) == 0 {
		t.Fatal("expected non-empty spectrogram")
	}

	totalEnergy := 0.0
	for _, frame := range spec {
		for _, mag := range frame {
			totalEnergy += mag
		}
	}
	if totalEnergy == 0 {
		t.Error("expected non-zero energy for sine wave input")
	}

	peakBin := -1
	peakMag := -1.0
	for i, mag := range spec[len(spec)/2] {
		if mag > peakMag {
			peakMag = mag
			peakBin = i
		}
	}
	if peakBin < 0 {
		t.Error("expected a peak frequency bin")
	}
}

func TestGenerateSpectrogram_FrameCount(t *testing.T) {
	samples := make([]float64, 4096)
	spec := GenerateSpectrogram(samples)

	expected := 4096 / 512
	if len(spec) != 5 {
		t.Errorf("expected %d frames for 4096 samples, got %d", expected, len(spec))
	}
}

func TestGenerateSpectrogram_WindowSize(t *testing.T) {
	samples := make([]float64, 2047)
	spec := GenerateSpectrogram(samples)
	if len(spec) != 0 {
		t.Errorf("expected 0 frames for 2047 samples (less than window size), got %d", len(spec))
	}
}

func TestGenerateSpectrogram_BinCount(t *testing.T) {
	samples := make([]float64, 4096)
	spec := GenerateSpectrogram(samples)
	if len(spec) == 0 {
		t.Fatal("expected non-empty spectrogram")
	}
	if len(spec[0]) != 512 {
		t.Errorf("expected 512 magnitude bins per frame, got %d", len(spec[0]))
	}
}

func BenchmarkGenerateSpectrogram(b *testing.B) {
	n := 44100 * 2
	samples := make([]float64, n)
	for i := range samples {
		samples[i] = math.Sin(2 * math.Pi * 440 * float64(i) / 44100)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GenerateSpectrogram(samples)
	}
}
