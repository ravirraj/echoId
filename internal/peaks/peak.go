package peak

import (
	"sort"
)

type Peak struct {
	TimeIndex int
	FreqIndex int
	Magnitude float64
}

func DetectPeaks(spectrogram [][]float64) []Peak {

	if len(spectrogram) == 0 {
		return nil
	}

	var peaks []Peak

	const (
		freqWindow       = 5
		timeWindow       = 3
		maxPeaksPerFrame = 8
	)

	for t := timeWindow; t < len(spectrogram)-timeWindow; t++ {

		frame := spectrogram[t]

		type candidate struct {
			freq int
			mag  float64
		}

		var candidates []candidate

		var frameMags []float64
		frameMags = append(frameMags, frame...)

		sort.Float64s(frameMags)

		frameThreshold := frameMags[len(frameMags)*90/100]

		for f := freqWindow; f < len(frame)-freqWindow; f++ {

			mag := frame[f]

			if mag < frameThreshold {
				continue
			}

			isMax := true

			for dt := -timeWindow; dt <= timeWindow && isMax; dt++ {

				neighborFrame := spectrogram[t+dt]

				for df := -freqWindow; df <= freqWindow; df++ {

					if dt == 0 && df == 0 {
						continue
					}

					if neighborFrame[f+df] >= mag {
						isMax = false
						break
					}
				}
			}

			if isMax {
				candidates = append(candidates, candidate{
					freq: f,
					mag:  mag,
				})
			}
		}

		sort.Slice(candidates, func(i, j int) bool {
			return candidates[i].mag > candidates[j].mag
		})

		if len(candidates) > maxPeaksPerFrame {
			candidates = candidates[:maxPeaksPerFrame]
		}

		for _, c := range candidates {
			peaks = append(peaks, Peak{
				TimeIndex: t,
				FreqIndex: c.freq,
				Magnitude: c.mag,
			})
		}
	}

	return peaks
}
