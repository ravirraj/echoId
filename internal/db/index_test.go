package db

import (
	"path/filepath"
	"testing"

	"github.com/ravirraj/echoid/internal/fingerprint"
)

func TestNewIndex(t *testing.T) {
	idx := NewIndex()
	if idx == nil {
		t.Fatal("NewIndex returned nil")
	}
	if idx.Data == nil {
		t.Fatal("Index.Data is nil")
	}
	if len(idx.Data) != 0 {
		t.Errorf("expected empty index, got %d entries", len(idx.Data))
	}
}

func TestAdd(t *testing.T) {
	idx := NewIndex()
	fps := []fingerprint.Fingerprint{
		{Freq1: 100, Freq2: 200, DeltaTime: 10, AnchorTime: 0},
		{Freq1: 300, Freq2: 400, DeltaTime: 20, AnchorTime: 5},
		{Freq1: 100, Freq2: 200, DeltaTime: 10, AnchorTime: 10},
	}
	idx.Add("song1", fps)

	if len(idx.Data) != 2 {
		t.Errorf("expected 2 unique hash keys, got %d", len(idx.Data))
	}

	hash1 := fingerprint.HashFingerprint(100, 200, 10)
	entries := idx.Data[hash1]
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries for hash1, got %d", len(entries))
	}
	for _, e := range entries {
		if e.SongID != "song1" {
			t.Errorf("expected SongID 'song1', got %q", e.SongID)
		}
	}

	hash2 := fingerprint.HashFingerprint(300, 400, 20)
	entries2 := idx.Data[hash2]
	if len(entries2) != 1 {
		t.Fatalf("expected 1 entry for hash2, got %d", len(entries2))
	}
	if entries2[0].AnchorTime != 5 {
		t.Errorf("expected AnchorTime 5, got %d", entries2[0].AnchorTime)
	}
}

func TestAddMultipleSongs(t *testing.T) {
	idx := NewIndex()
	idx.Add("song1", []fingerprint.Fingerprint{
		{Freq1: 100, Freq2: 200, DeltaTime: 10, AnchorTime: 0},
	})
	idx.Add("song2", []fingerprint.Fingerprint{
		{Freq1: 100, Freq2: 200, DeltaTime: 10, AnchorTime: 5},
	})

	hash := fingerprint.HashFingerprint(100, 200, 10)
	entries := idx.Data[hash]
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries for shared hash, got %d", len(entries))
	}

	songs := map[string]int{}
	for _, e := range entries {
		songs[e.SongID]++
	}
	if songs["song1"] != 1 {
		t.Errorf("expected song1 to have 1 entry, got %d", songs["song1"])
	}
	if songs["song2"] != 1 {
		t.Errorf("expected song2 to have 1 entry, got %d", songs["song2"])
	}
}

func TestSaveAndLoadIndex(t *testing.T) {
	orig := NewIndex()
	orig.Add("song1", []fingerprint.Fingerprint{
		{Freq1: 100, Freq2: 200, DeltaTime: 10, AnchorTime: 0},
		{Freq1: 300, Freq2: 400, DeltaTime: 20, AnchorTime: 5},
	})
	orig.Add("song2", []fingerprint.Fingerprint{
		{Freq1: 100, Freq2: 200, DeltaTime: 10, AnchorTime: 10},
	})

	path := filepath.Join(t.TempDir(), "index.gob")
	if err := orig.Save(path); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := LoadIndex(path)
	if err != nil {
		t.Fatalf("LoadIndex: %v", err)
	}

	if len(loaded.Data) != len(orig.Data) {
		t.Errorf("key count: got %d, want %d", len(loaded.Data), len(orig.Data))
	}

	for hash, origEntries := range orig.Data {
		loadedEntries, ok := loaded.Data[hash]
		if !ok {
			t.Errorf("hash %d missing in loaded index", hash)
			continue
		}
		if len(loadedEntries) != len(origEntries) {
			t.Errorf("hash %d: got %d entries, want %d", hash, len(loadedEntries), len(origEntries))
			continue
		}
		for i := range origEntries {
			if loadedEntries[i] != origEntries[i] {
				t.Errorf("hash %d entry %d: got %+v, want %+v",
					hash, i, loadedEntries[i], origEntries[i])
			}
		}
	}
}

func BenchmarkSaveLoad(b *testing.B) {
	idx := NewIndex()
	for songID := 0; songID < 100; songID++ {
		fps := make([]fingerprint.Fingerprint, 100)
		for i := range fps {
			fps[i] = fingerprint.Fingerprint{
				Freq1:      (i * 7) % 512,
				Freq2:      (i * 13) % 512,
				DeltaTime:  (i * 3) % 60,
				AnchorTime: i,
			}
		}
		idx.Add("song"+string(rune('A'+songID%26)), fps)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		path := filepath.Join(b.TempDir(), "index.gob")
		if err := idx.Save(path); err != nil {
			b.Fatal(err)
		}
		if _, err := LoadIndex(path); err != nil {
			b.Fatal(err)
		}
	}
}
