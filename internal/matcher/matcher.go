// package matcher

// import (
// 	"github.com/ravirraj/echoid/internal/db"
// 	"github.com/ravirraj/echoid/internal/fingerprint"
// )

// func Match(index *db.Index, query []fingerprint.Fingerprint) (string, int) {

// 	offsetVotes := make(map[string]map[int]int)

// 	for _, fp := range query {

// 		hash := fingerprint.HashFingerprint(fp.Freq1, fp.Freq2, fp.DeltaTime)

// 		matches, ok := index.Data[hash]
// 		if !ok {
// 			continue
// 		}

// 		for _, match := range matches {

// 			offset := match.AnchorTime - fp.AnchorTime

// 			if _, ok := offsetVotes[match.SongID]; !ok {
// 				offsetVotes[match.SongID] = make(map[int]int)
// 			}

// 			offsetVotes[match.SongID][offset]++
// 		}
// 	}

// 	bestSong := ""
// 	bestScore := 0

// 	for songID, offsets := range offsetVotes {

// 		for _, count := range offsets {

// 			if count > bestScore {
// 				bestScore = count
// 				bestSong = songID
// 			}
// 		}
// 	}

// 	return bestSong, bestScore
// }

package matcher

import (
	"github.com/ravirraj/echoid/internal/db"
	"github.com/ravirraj/echoid/internal/fingerprint"
)

func Match(index *db.Index, query []fingerprint.Fingerprint) (string, int) {

	offsetVotes := make(map[string]map[int]int)

	// if len(query) > 500 {
	// 	query = query[:500]
	// }
	for _, fp := range query {

		hash := fingerprint.HashFingerprint(fp.Freq1, fp.Freq2, fp.DeltaTime)

		matches, ok := index.Data[hash]
		if !ok {
			continue
		}
		// println("hash matched")

		for _, match := range matches {

			offset := match.AnchorTime - fp.AnchorTime
			// fmt.Println(offset)

			if _, ok := offsetVotes[match.SongID]; !ok {
				offsetVotes[match.SongID] = make(map[int]int)
			}

			offsetVotes[match.SongID][offset]++
			// fmt.Println("hash matches:", len(matches))

		}
	}

	bestSong := ""
	bestScore := 0

	for songID, offsets := range offsetVotes {

		maxOffsetVotes := 0

		for _, count := range offsets {
			if count > maxOffsetVotes {
				maxOffsetVotes = count
			}
		}

		if maxOffsetVotes > bestScore {
			bestScore = maxOffsetVotes
			bestSong = songID
		}
	}

	// Dynamic threshold: require at least 5% of query fingerprints to match
	// or minimum of 50 matches, whichever is higher
	minThreshold := len(query) * 5 / 100
	if minThreshold < 50 {
		minThreshold = 50
	}

	if bestScore < minThreshold {
		return "", 0
	}

	return bestSong, bestScore

	// fmt.Println(bestSong)
	// return bestSong, bestScore
}
