package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ravirraj/echoid/internal/audio"
	"github.com/ravirraj/echoid/internal/config"
	"github.com/ravirraj/echoid/internal/db"
	"github.com/ravirraj/echoid/internal/fingerprint"
	"github.com/ravirraj/echoid/internal/matcher"
	"github.com/ravirraj/echoid/internal/output"
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
		artist := addCmd.String("artist", "", "artist name")
		album := addCmd.String("album", "", "album name")
		addCmd.Usage = func() {
			fmt.Fprintf(os.Stderr, "Usage: echoid add -file <path> [-id <name>] [-artist <name>] [-album <name>]\n")
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

		meta := db.SongMeta{
			Title:  id,
			Artist: *artist,
			Album:  *album,
		}

		if err := runAdd(*file, id, meta); err != nil {
			output.PrintError(err.Error())
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
			output.PrintError(err.Error())
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

		output.Listening()
		if err := recordAudio(*seconds); err != nil {
			output.PrintError(err.Error())
			os.Exit(1)
		}

	case "youtube":
		ytCmd := flag.NewFlagSet("youtube", flag.ExitOnError)
		url := ytCmd.String("url", "", "YouTube URL")
		artist := ytCmd.String("artist", "", "artist name (overrides metadata)")
		album := ytCmd.String("album", "", "album name (overrides metadata)")
		ytCmd.Usage = func() {
			fmt.Fprintf(os.Stderr, "Usage: echoid youtube -url <youtube-url> [-artist <name>] [-album <name>]\n")
			ytCmd.PrintDefaults()
		}
		ytCmd.Parse(os.Args[2:])

		if *url == "" {
			ytCmd.Usage()
			os.Exit(1)
		}

		output.Stepf("Fetching...")
		meta, filePath, err := audio.DownloadAudio(*url)
		if err != nil {
			output.PrintError(err.Error())
			os.Exit(1)
		}
		output.Stepf("Downloaded")

		if *artist != "" {
			meta.Artist = *artist
		}
		if *album != "" {
			meta.Album = *album
		}

		songID := sanitizeFilename(meta.Title)
		dbMeta := db.SongMeta{Title: meta.Title, Artist: meta.Artist, Album: meta.Album}
		if err := runAdd(filePath, songID, dbMeta); err != nil {
			os.Remove(filePath)
			output.PrintError(err.Error())
			os.Exit(1)
		}

		if err := os.Remove(filePath); err != nil {
			output.Stepf("warning: could not remove temp file: %v", err)
		}

	case "-h", "--help", "help":
		usage()

	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", os.Args[1])
		usage()
		os.Exit(1)
	}
}

func runAdd(file string, songID string, meta db.SongMeta) error {
	var index *db.Index

	if _, err := os.Stat(config.DBPath); err == nil {
		index, err = db.LoadIndex(config.DBPath)
		if err != nil {
			return fmt.Errorf("loading index: %w", err)
		}
	} else {
		index = db.NewIndex()
	}

	output.Stepf("Loading audio...")
	samples, err := audio.LoadAudio(file)
	if err != nil {
		return fmt.Errorf("loading audio: %w", err)
	}

	output.Generating()
	spec := spectrogram.GenerateSpectrogram(samples)
	p := peak.DetectPeaks(spec)
	fps := fingerprint.GenerateFingerprints(p)

	output.PrintFingerprintStats(len(p), len(fps))

	index.Add(songID, meta, fps)
	if err := index.Save(config.DBPath); err != nil {
		return fmt.Errorf("saving index: %w", err)
	}

	output.PrintIndexed(songID)
	return nil
}

func runMatch(file string) error {
	index, err := db.LoadIndex(config.DBPath)
	if err != nil {
		return fmt.Errorf("loading index: %w", err)
	}

	output.Stepf("Loading audio...")
	samples, err := audio.LoadAudio(file)
	if err != nil {
		return fmt.Errorf("loading audio: %w", err)
	}

	output.Matching()
	start := time.Now()
	spec := spectrogram.GenerateSpectrogram(samples)
	p := peak.DetectPeaks(spec)
	query := fingerprint.GenerateFingerprints(p)

	song, score := matcher.Match(index, query)
	elapsed := time.Since(start).Seconds()
	_ = score

	if song == "" {
		output.PrintNoMatch()
	} else {
		meta := index.Metadata[song]
		title := meta.Title
		if title == "" {
			title = song
		}

		output.PrintResult(title, meta.Artist, meta.Album, elapsed)
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

func sanitizeFilename(name string) string {
	replacer := strings.NewReplacer(
		"/", "_", "\\", "_", ":", "_",
		"*", "_", "?", "_", "\"", "_",
		"<", "_", ">", "_", "|", "_",
		" ", "_",
	)
	return replacer.Replace(name)
}
