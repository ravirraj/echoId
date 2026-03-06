package db

import (
	"encoding/gob"
	"os"

	"github.com/ravirraj/echoid/internal/fingerprint"
)

type Entry struct {
	SongID     string
	AnchorTime int
}

type Index struct {
	Data map[uint64][]Entry
}

func NewIndex() *Index {
	return &Index{
		Data: make(map[uint64][]Entry),
	}
}

func (idx *Index) Add(songID string, fps []fingerprint.Fingerprint) {

	for _, fp := range fps {

		hash := fingerprint.HashFingerprint(fp.Freq1, fp.Freq2, fp.DeltaTime)

		// fmt.Println("Adding:", songID)

		idx.Data[hash] = append(idx.Data[hash], Entry{
			SongID:     songID,
			AnchorTime: fp.AnchorTime,
		})
	}
}

func (idx *Index) Save(path string) error {

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	enc := gob.NewEncoder(file)

	return enc.Encode(idx.Data)
}

func LoadIndex(path string) (*Index, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	idx := NewIndex()

	dec := gob.NewDecoder(file)

	err = dec.Decode(&idx.Data)

	return idx, err
}
