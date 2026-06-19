package fingerprint

import (
	"sort"

	peak "github.com/ravirraj/echoid/internal/peaks"
)

type Fingerprint struct {
	Freq1      int
	Freq2      int
	DeltaTime  int
	AnchorTime int
}

func GenerateFingerprints(peaksList []peak.Peak) []Fingerprint {
	if len(peaksList) == 0 {
		return nil
	}

	sort.Slice(peaksList, func(i, j int) bool {
		return peaksList[i].TimeIndex < peaksList[j].TimeIndex
	})

	const (
		fanOut      = 10 // max pairs per anchor peak
		minDelta    = 1
		maxDelta    = 60 // don't look too far ahead
		freqBinSize = 4
		timeBinSize = 2
	)

	fingerprints := make([]Fingerprint, 0, len(peaksList)*fanOut)

	for i := 0; i < len(peaksList); i++ {
		anchor := peaksList[i]
		pairsCreated := 0

		for j := i + 1; j < len(peaksList) && pairsCreated < fanOut; j++ {
			target := peaksList[j]
			deltaTime := target.TimeIndex - anchor.TimeIndex

			if deltaTime < minDelta {
				continue
			}

			if deltaTime > maxDelta {
				break
			}

			// quantize to reduce sensitivity to small shifts
			freq1 := (anchor.FreqIndex / freqBinSize) * freqBinSize
			freq2 := (target.FreqIndex / freqBinSize) * freqBinSize
			deltaTime = (deltaTime / timeBinSize) * timeBinSize

			fingerprints = append(fingerprints, Fingerprint{
				Freq1:      freq1,
				Freq2:      freq2,
				DeltaTime:  deltaTime,
				AnchorTime: anchor.TimeIndex,
			})

			pairsCreated++
		}
	}

	return fingerprints
}
