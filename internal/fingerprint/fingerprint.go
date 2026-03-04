package fingerprint

import "github.com/ravirraj/echoid/internal/peaks"

// import "yourmodule/internal/peaks"

type Fingerprint struct {
	Freq1      int
	Freq2      int
	DeltaTime  int
	AnchorTime int
}

func GenerateFingerprints(peaksList []peak.Peak) []Fingerprint {

	fingerprints := []Fingerprint{}

	fanOut := 5
	maxDelta := 50

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