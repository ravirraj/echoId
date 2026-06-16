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
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: echoid [add|match|listen|youtube]")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "add":
		addCmd := flag.NewFlagSet("add", flag.ExitOnError)
		file := addCmd.String("file", "", "audio file to index")
		songID := addCmd.String("id", *file, "song identifier")
		addCmd.Parse(os.Args[2:])

		if *file == "" {
			fmt.Println("error: -file is required")
			return
		}

		if err := runAdd(*file, *songID); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
		}

	case "match":
		matchCmd := flag.NewFlagSet("match", flag.ExitOnError)
		file := matchCmd.String("file", "", "audio file to match")
		matchCmd.Parse(os.Args[2:])

		if *file == "" {
			fmt.Println("error: -file is required")
			return
		}

		if err := runMatch(*file); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
		}

	case "listen":
		fmt.Println("  Recording...")
		if err := recordAudio(10); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
		}

	case "youtube":
		matchCmd := flag.NewFlagSet("youtube", flag.ExitOnError)
		url := matchCmd.String("url", "", "YouTube URL")
		matchCmd.Parse(os.Args[2:])

		if *url == "" {
			fmt.Println("error: -url is required")
			return
		}

		fmt.Println("  Fetching...")
		title, filePath, err := audio.DownloadAudio(*url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			return
		}
		fmt.Println("  Downloaded")

		if err := runAdd(filePath, title); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			return
		}

		os.Remove(filePath)

	default:
		fmt.Printf("unknown command: %s\n", os.Args[1])
	}
}

func runAdd(file string, songID string) error {
	var index *db.Index

	if _, err := os.Stat("fingerprints.db"); err == nil {
		var err error
		index, err = db.LoadIndex("fingerprints.db")
		if err != nil {
			return err
		}
	} else {
		index = db.NewIndex()
	}

	fmt.Println("  Loading audio...")
	samples, err := audio.LoadAudio(file)
	if err != nil {
		return err
	}

	fmt.Println("  Generating fingerprints...")
	spec := spectrogram.GenerateSpectrogram(samples)
	p := peak.DetectPeaks(spec)
	fps := fingerprint.GenerateFingerprints(p)

	fmt.Printf("  Peaks: %d | Fingerprints: %d\n", len(p), len(fps))

	index.Add(songID, fps)
	if err := index.Save("fingerprints.db"); err != nil {
		return err
	}

	fmt.Printf("  Indexed: %s\n", songID)
	return nil
}

func runMatch(file string) error {
	index, err := db.LoadIndex("fingerprints.db")
	if err != nil {
		return err
	}

	fmt.Println("  Loading audio...")
	samples, err := audio.LoadAudio(file)
	if err != nil {
		return err
	}

	fmt.Println("  Matching...")
	spec := spectrogram.GenerateSpectrogram(samples)
	p := peak.DetectPeaks(spec)
	query := fingerprint.GenerateFingerprints(p)

	song, score := matcher.Match(index, query)
	if song == "" {
		fmt.Println("  No match found")
	} else {
		fmt.Printf("  Match: %s (score: %d)\n", song, score)
	}

	return nil
}

func recordAudio(duration int) error {
	filePath, err := audio.RecordAudio(duration)
	if err != nil {
		return err
	}
	return runMatch(filePath)
}
