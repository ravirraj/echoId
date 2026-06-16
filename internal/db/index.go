package db

import (
	"encoding/gob"
	"fmt"
	"os"

	"github.com/ravirraj/echoid/internal/fingerprint"
)

// Entry stores a song identifier and the anchor time of a matched
// fingerprint.
type Entry struct {
	SongID     string
	AnchorTime int
}

// Index maps hashed fingerprint values to the list of song entries
// that produced them.
type Index struct {
	Data map[uint64][]Entry
}

// NewIndex creates and returns a new empty Index.
func NewIndex() *Index {
	return &Index{
		Data: make(map[uint64][]Entry),
	}
}

// Add inserts fingerprints for a given song ID into the index.
func (idx *Index) Add(songID string, fps []fingerprint.Fingerprint) {
	for _, fp := range fps {
		hash := fingerprint.HashFingerprint(fp.Freq1, fp.Freq2, fp.DeltaTime)
		idx.Data[hash] = append(idx.Data[hash], Entry{
			SongID:     songID,
			AnchorTime: fp.AnchorTime,
		})
	}
}

// Save writes the index data to a file using gob encoding. It syncs
// the file to disk before closing.
func (idx *Index) Save(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create index file: %w", err)
	}
	defer file.Close()

	if err := file.Sync(); err != nil {
		return fmt.Errorf("sync index file: %w", err)
	}

	enc := gob.NewEncoder(file)
	if err := enc.Encode(idx.Data); err != nil {
		return fmt.Errorf("encode index: %w", err)
	}

	return nil
}

// LoadIndex reads an Index from a gob-encoded file.
func LoadIndex(path string) (*Index, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open index file: %w", err)
	}
	defer file.Close()

	idx := NewIndex()
	dec := gob.NewDecoder(file)
	if err := dec.Decode(&idx.Data); err != nil {
		return nil, fmt.Errorf("decode index: %w", err)
	}

	return idx, nil
}
