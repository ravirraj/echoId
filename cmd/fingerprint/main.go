package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ravirraj/echoid/internal/audio"
	"github.com/ravirraj/echoid/internal/config"
	"github.com/ravirraj/echoid/internal/db"
	"github.com/ravirraj/echoid/internal/fingerprint"
	"github.com/ravirraj/echoid/internal/matcher"
	peak "github.com/ravirraj/echoid/internal/peaks"
	"github.com/ravirraj/echoid/internal/spectrogram"
)

func usage() {
	fmt.Fprintf(os.Stderr, `Usage: echoid <command> [options]

Commands:
  add     Index an audio file
  match   Match an audio file against the index
  listen  Record and match audio (default: 10 seconds)
  youtube Download and index a YouTube video

Options:
  -h, --help  Show this help

Use "echoid <command> -h" for command-specific help.
`)
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "add":
		addCmd := flag.NewFlagSet("add", flag.ExitOnError)
		file := addCmd.String("file", "", "audio file to index")
		songID := addCmd.String("id", "", "song identifier (defaults to filename)")
		addCmd.Usage = func() {
			fmt.Fprintf(os.Stderr, "Usage: echoid add -file <path> [-id <name>]\n")
			addCmd.PrintDefaults()
		}
		addCmd.Parse(os.Args[2:])

		if *file == "" {
			addCmd.Usage()
			os.Exit(1)
		}

		id := *songID
		if id == "" {
			id = *file
		}

		if err := runAdd(*file, id); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}

	case "match":
		matchCmd := flag.NewFlagSet("match", flag.ExitOnError)
		file := matchCmd.String("file", "", "audio file to match")
		matchCmd.Usage = func() {
			fmt.Fprintf(os.Stderr, "Usage: echoid match -file <path>\n")
			matchCmd.PrintDefaults()
		}
		matchCmd.Parse(os.Args[2:])

		if *file == "" {
			matchCmd.Usage()
			os.Exit(1)
		}

		if err := runMatch(*file); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}

	case "listen":
		listenCmd := flag.NewFlagSet("listen", flag.ExitOnError)
		seconds := listenCmd.Int("seconds", 10, "recording duration in seconds")
		listenCmd.Usage = func() {
			fmt.Fprintf(os.Stderr, "Usage: echoid listen [-seconds <n>]\n")
			listenCmd.PrintDefaults()
		}
		listenCmd.Parse(os.Args[2:])

		fmt.Println("  Recording...")
		if err := recordAudio(*seconds); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}

	case "youtube":
		ytCmd := flag.NewFlagSet("youtube", flag.ExitOnError)
		url := ytCmd.String("url", "", "YouTube URL")
		ytCmd.Usage = func() {
			fmt.Fprintf(os.Stderr, "Usage: echoid youtube -url <youtube-url>\n")
			ytCmd.PrintDefaults()
		}
		ytCmd.Parse(os.Args[2:])

		if *url == "" {
			ytCmd.Usage()
			os.Exit(1)
		}

		fmt.Println("  Fetching...")
		title, filePath, err := audio.DownloadAudio(*url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("  Downloaded")

		if err := runAdd(filePath, title); err != nil {
			os.Remove(filePath)
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}

		if err := os.Remove(filePath); err != nil {
			fmt.Fprintf(os.Stderr, "warning: could not remove temp file: %v\n", err)
		}

	case "-h", "--help", "help":
		usage()

	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", os.Args[1])
		usage()
		os.Exit(1)
	}
}

func runAdd(file string, songID string) error {
	var index *db.Index

	if _, err := os.Stat(config.DBPath); err == nil {
		index, err = db.LoadIndex(config.DBPath)
		if err != nil {
			return fmt.Errorf("loading index: %w", err)
		}
	} else {
		index = db.NewIndex()
	}

	fmt.Println("  Loading audio...")
	samples, err := audio.LoadAudio(file)
	if err != nil {
		return fmt.Errorf("loading audio: %w", err)
	}

	fmt.Println("  Generating fingerprints...")
	spec := spectrogram.GenerateSpectrogram(samples)
	p := peak.DetectPeaks(spec)
	fps := fingerprint.GenerateFingerprints(p)

	fmt.Printf("  Peaks: %d | Fingerprints: %d\n", len(p), len(fps))

	index.Add(songID, fps)
	if err := index.Save(config.DBPath); err != nil {
		return fmt.Errorf("saving index: %w", err)
	}

	fmt.Printf("  Indexed: %s\n", songID)
	return nil
}

func runMatch(file string) error {
	index, err := db.LoadIndex(config.DBPath)
	if err != nil {
		return fmt.Errorf("loading index: %w", err)
	}

	fmt.Println("  Loading audio...")
	samples, err := audio.LoadAudio(file)
	if err != nil {
		return fmt.Errorf("loading audio: %w", err)
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
		return fmt.Errorf("recording audio: %w", err)
	}
	defer func() {
		if err := os.Remove(filePath); err != nil {
			fmt.Fprintf(os.Stderr, "warning: could not remove temp file: %v\n", err)
		}
	}()
	return runMatch(filePath)
}

func init() {
	flag.Usage = usage
}
