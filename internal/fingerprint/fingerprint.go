package fingerprint

import (
	"sort"

	peak "github.com/ravirraj/echoid/internal/peaks"
)

// Fingerprint represents a hashed audio landmark derived from a pair
// of peaks, storing their frequency bins, time delta, and anchor time.
type Fingerprint struct {
	Freq1      int
	Freq2      int
	DeltaTime  int
	AnchorTime int
}

// GenerateFingerprints creates a set of Fingerprints from a list of
// spectral peaks by pairing each anchor peak with nearby target peaks
// within a constrained time window.
func GenerateFingerprints(peaksList []peak.Peak) []Fingerprint {

	if len(peaksList) == 0 {
		return nil
	}

	sort.Slice(peaksList, func(i, j int) bool {
		return peaksList[i].TimeIndex < peaksList[j].TimeIndex
	})

	const (
		fanOut      = 10
		minDelta    = 1
		maxDelta    = 60
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
