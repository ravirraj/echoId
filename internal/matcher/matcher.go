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

	for _, fp := range query {

		hash := fingerprint.HashFingerprint(fp.Freq1, fp.Freq2, fp.DeltaTime)

		matches, ok := index.Data[hash]
		if !ok {
			continue
		}

		for _, match := range matches {

			offset := match.AnchorTime - fp.AnchorTime

			if _, ok := offsetVotes[match.SongID]; !ok {
				offsetVotes[match.SongID] = make(map[int]int)
			}

			offsetVotes[match.SongID][offset]++
		}
	}

	bestSong := ""
	bestScore := 0

	for songID, offsets := range offsetVotes {

		for _, count := range offsets {

			if count > bestScore {
				bestScore = count
				bestSong = songID
			}
		}
	}

	return bestSong, bestScore
}
