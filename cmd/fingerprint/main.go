package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ravirraj/echoid/internal/audio"
	"github.com/ravirraj/echoid/internal/db"
	"github.com/ravirraj/echoid/internal/fingerprint"
	"github.com/ravirraj/echoid/internal/matcher"
	peak "github.com/ravirraj/echoid/internal/peaks"
	spectrogram "github.com/ravirraj/echoid/internal/spectogram"
	// "github.com/ravirraj/echoid/internal/spectrogram"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("usage: echoid [add|match] <file>")
		os.Exit(1)
	}

	switch os.Args[1] {

	case "add":
		addCmd := flag.NewFlagSet("add", flag.ExitOnError)

		file := addCmd.String("file", "", "audio file to index")
		songID := addCmd.String("id", "song1", "song identifier")

		addCmd.Parse(os.Args[2:])

		if *file == "" {
			fmt.Println("please provide -file")
			return
		}

		runAdd(*file, *songID)

	case "match":
		matchCmd := flag.NewFlagSet("match", flag.ExitOnError)

		file := matchCmd.String("file", "", "audio file to match")

		matchCmd.Parse(os.Args[2:])

		if *file == "" {
			fmt.Println("please provide -file")
			return
		}

		runMatch(*file)

	default:
		fmt.Println("unknown command:", os.Args[1])
	}

}

func runAdd(file string, songID string) {

	index := db.NewIndex()

	samples, _, err := audio.LoadWav(file)

	spec := spectrogram.GenerateSpectrogram(samples)

	p := peak.DetectPeaks(spec)

	fps := fingerprint.GenerateFingerprints(p)

	index.Add(songID, fps)

	err = index.Save("fingerprints.db")
	if err != nil {
		fmt.Println("save error:", err)
		return
	}

	fmt.Println("song indexed:", songID)

}

func runMatch(file string) {

	index, err := db.LoadIndex("fingerprints.db")
	if err != nil {
		fmt.Println("load error:", err)
		return
	}

	samples, _, err := audio.LoadWav(file)
	if err != nil {
		fmt.Println("load error:", err)
		return
	}

	spec := spectrogram.GenerateSpectrogram(samples)

	p := peak.DetectPeaks(spec)

	query := fingerprint.GenerateFingerprints(p)

	song, score := matcher.Match(index, query)

	fmt.Println("match:", song, "score:", score)

}
