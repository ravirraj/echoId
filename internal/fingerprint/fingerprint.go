package fingerprint

import (
	"sort"

	peak "github.com/ravirraj/echoid/internal/peaks"
)

// import "yourmodule/internal/peaks"

type Fingerprint struct {
	Freq1      int
	Freq2      int
	DeltaTime  int
	AnchorTime int
}

func GenerateFingerprints(peaksList []peak.Peak) []Fingerprint {

	fingerprints := []Fingerprint{}

	fanOut := 3
	maxDelta := 30

	sort.Slice(peaksList, func(i, j int) bool {
		return peaksList[i].TimeIndex < peaksList[j].TimeIndex
	})
	for i := 0; i < len(peaksList); i++ {

		for j := i + 1; j < len(peaksList) && j < i+fanOut; j++ {

			deltaTime := peaksList[j].TimeIndex - peaksList[i].TimeIndex

			if deltaTime > maxDelta {
				break
			}

			fingerprints = append(fingerprints, Fingerprint{
				Freq1:      peaksList[i].FreqIndex,
				Freq2:      peaksList[j].FreqIndex,
				DeltaTime:  deltaTime,
				AnchorTime: peaksList[i].TimeIndex,
			})
		}
	}

	return fingerprints
}
