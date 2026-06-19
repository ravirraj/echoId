package matcher

import (
	"github.com/ravirraj/echoid/internal/db"
	"github.com/ravirraj/echoid/internal/fingerprint"
)

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

			// bin nearby offsets so slight timing differences still agree
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

	minThreshold := max(len(query)/25, 15) // need at least some fraction of query to match

	if bestScore < minThreshold {
		return "", 0
	}

	// if top two results are too close, it's ambiguous
	if secondBest > 0 {
		confidence := float64(bestScore) / float64(secondBest)
		if confidence < 1.3 {
			return "", 0
		}
	}
	return bestSong, bestScore
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
