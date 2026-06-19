package matcher

import (
	"testing"

	"github.com/ravirraj/echoid/internal/db"
	"github.com/ravirraj/echoid/internal/fingerprint"
)

func TestMatch_EmptyIndex(t *testing.T) {
	idx := db.NewIndex()
	query := []fingerprint.Fingerprint{
		{Freq1: 100, Freq2: 200, DeltaTime: 10, AnchorTime: 0},
	}
	songID, score := Match(idx, query)
	if songID != "" || score != 0 {
		t.Errorf("expected ('', 0), got (%q, %d)", songID, score)
	}
}

func TestMatch_ExactMatch(t *testing.T) {
	idx := db.NewIndex()
	fps := make([]fingerprint.Fingerprint, 20)
	for i := range fps {
		fps[i] = fingerprint.Fingerprint{
			Freq1: 100, Freq2: 200, DeltaTime: 10, AnchorTime: 0,
		}
	}
	idx.Add("song1", db.SongMeta{Title: "song1"}, fps)

	query := make([]fingerprint.Fingerprint, 20)
	for i := range query {
		query[i] = fingerprint.Fingerprint{
			Freq1: 100, Freq2: 200, DeltaTime: 10, AnchorTime: 0,
		}
	}

	songID, score := Match(idx, query)
	if songID != "song1" {
		t.Errorf("expected song1, got %q", songID)
	}
	if score == 0 {
		t.Error("expected non-zero score")
	}
}

func TestMatch_NoMatch(t *testing.T) {
	idx := db.NewIndex()
	idx.Add("song1", db.SongMeta{Title: "song1"}, []fingerprint.Fingerprint{
		{Freq1: 100, Freq2: 200, DeltaTime: 10, AnchorTime: 0},
	})

	query := []fingerprint.Fingerprint{
		{Freq1: 999, Freq2: 888, DeltaTime: 77, AnchorTime: 0},
	}

	songID, score := Match(idx, query)
	if songID != "" || score != 0 {
		t.Errorf("expected ('', 0) for no match, got (%q, %d)", songID, score)
	}
}

func TestMatch_ScoreBelowThreshold(t *testing.T) {
	idx := db.NewIndex()
	idx.Add("song1", db.SongMeta{Title: "song1"}, []fingerprint.Fingerprint{
		{Freq1: 100, Freq2: 200, DeltaTime: 10, AnchorTime: 0},
	})

	query := make([]fingerprint.Fingerprint, 10)
	for i := range query {
		query[i] = fingerprint.Fingerprint{
			Freq1: 100, Freq2: 200, DeltaTime: 10, AnchorTime: 0,
		}
	}

	songID, score := Match(idx, query)
	if songID != "" || score != 0 {
		t.Errorf("expected ('', 0) for below-threshold, got (%q, %d)", songID, score)
	}
}

func TestMatch_MultipleSongs(t *testing.T) {
	idx := db.NewIndex()

	fps1 := make([]fingerprint.Fingerprint, 10)
	for i := range fps1 {
		fps1[i] = fingerprint.Fingerprint{
			Freq1: 100, Freq2: 200, DeltaTime: 10, AnchorTime: 0,
		}
	}
	idx.Add("song1", db.SongMeta{Title: "song1"}, fps1)

	fps2 := make([]fingerprint.Fingerprint, 10)
	for i := range fps2 {
		fps2[i] = fingerprint.Fingerprint{
			Freq1: 300, Freq2: 400, DeltaTime: 20, AnchorTime: 5,
		}
	}
	idx.Add("song2", db.SongMeta{Title: "song2"}, fps2)

	query := make([]fingerprint.Fingerprint, 10)
	for i := range query {
		query[i] = fingerprint.Fingerprint{
			Freq1: 100, Freq2: 200, DeltaTime: 10, AnchorTime: 0,
		}
	}

	songID, score := Match(idx, query)
	if songID != "song1" {
		t.Errorf("expected song1, got %q", songID)
	}
	if score == 0 {
		t.Error("expected non-zero score")
	}
}

func TestMatch_ConfidenceCheck(t *testing.T) {
	idx := db.NewIndex()

	fps1 := make([]fingerprint.Fingerprint, 20)
	for i := range fps1 {
		fps1[i] = fingerprint.Fingerprint{
			Freq1: 100, Freq2: 200, DeltaTime: 10, AnchorTime: 0,
		}
	}
	idx.Add("song1", db.SongMeta{Title: "song1"}, fps1)

	fps2 := make([]fingerprint.Fingerprint, 16)
	for i := range fps2 {
		fps2[i] = fingerprint.Fingerprint{
			Freq1: 100, Freq2: 200, DeltaTime: 10, AnchorTime: 5,
		}
	}
	idx.Add("song2", db.SongMeta{Title: "song2"}, fps2)

	query := make([]fingerprint.Fingerprint, 20)
	for i := range query {
		query[i] = fingerprint.Fingerprint{
			Freq1: 100, Freq2: 200, DeltaTime: 10, AnchorTime: 0,
		}
	}

	songID, score := Match(idx, query)
	if songID != "" || score != 0 {
		t.Errorf("expected ('', 0) for low confidence, got (%q, %d)", songID, score)
	}
}

func BenchmarkMatch(b *testing.B) {
	idx := db.NewIndex()

	for songID := 0; songID < 50; songID++ {
		fps := make([]fingerprint.Fingerprint, 200)
		for i := range fps {
			fps[i] = fingerprint.Fingerprint{
				Freq1:      (i * 7) % 512,
				Freq2:      (i * 13) % 512,
				DeltaTime:  (i * 3) % 60,
				AnchorTime: i,
			}
		}
		idx.Add("song"+string(rune('A'+songID%26)), db.SongMeta{Title: "song"}, fps)
	}

	query := make([]fingerprint.Fingerprint, 100)
	for i := range query {
		query[i] = fingerprint.Fingerprint{
			Freq1:      (i * 7) % 512,
			Freq2:      (i * 13) % 512,
			DeltaTime:  (i * 3) % 60,
			AnchorTime: i,
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Match(idx, query)
	}
}
