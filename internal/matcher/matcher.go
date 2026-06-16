package matcher

import (
	"github.com/ravirraj/echoid/internal/db"
	"github.com/ravirraj/echoid/internal/fingerprint"
)

// Match compares query fingerprints against a song index and returns
// the best-matching song ID along with its match score.
func Match(index *db.Index, query []fingerprint.Fingerprint) (string, int) {
	matchedHashes := 0

	offsetVotes := make(map[string]map[int]int)

	const offsetBinSize = 5

	for _, fp := range query {

		hash := fingerprint.HashFingerprint(
			fp.Freq1,
			fp.Freq2,
			fp.DeltaTime,
		)

		matches, ok := index.Data[hash]
		if !ok {
			continue
		}
		matchedHashes++

		for _, match := range matches {

			offset := match.AnchorTime - fp.AnchorTime

			// Bin nearby offsets together
			offset = (offset / offsetBinSize) * offsetBinSize

			if _, ok := offsetVotes[match.SongID]; !ok {
				offsetVotes[match.SongID] = make(map[int]int)
			}

			offsetVotes[match.SongID][offset]++
		}
	}

	bestSong := ""
	bestScore := 0
	secondBest := 0

	for songID, offsets := range offsetVotes {

		songScore := 0

		for _, count := range offsets {
			if count > songScore {
				songScore = count
			}
		}

		if songScore > bestScore {
			secondBest = bestScore
			bestScore = songScore
			bestSong = songID
		} else if songScore > secondBest {
			secondBest = songScore
		}
	}
	if bestScore == 0 {
		return "", 0
	}

	// Require minimum votes
	minThreshold := max(len(query)/25, 15)

	if bestScore < minThreshold {
		return "", 0
	}

	// Confidence check
	if secondBest > 0 {

		confidence := float64(bestScore) / float64(secondBest)

		// If top result is too close to second result,
		// treat as ambiguous.
		if confidence < 1.3 {
			return "", 0
		}
	}
	return bestSong, bestScore
}

// max returns the larger of two integers.
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
