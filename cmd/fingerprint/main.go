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
		// songID := addCmd.String("id", *file, "song identifier")

		addCmd.Parse(os.Args[2:])
		songID := addCmd.String("id", *file, "song identifier")
		addCmd.Parse(os.Args[2:])

		fmt.Println(*file)
		fmt.Println(*songID)

		if *file == "" {
			fmt.Println("please provide -file")
			return
		}

		// if *songID == "" {
		// 	fmt.Println("Please Provide the id for the song")
		// 	// return
		// }

		err := runAdd(*file, *songID)

		if err != nil {
			fmt.Println(err)
			return
		}

	case "match":
		matchCmd := flag.NewFlagSet("match", flag.ExitOnError)

		file := matchCmd.String("file", "", "audio file to match")

		matchCmd.Parse(os.Args[2:])

		if *file == "" {
			fmt.Println("please provide -file")
			return
		}

		err := runMatch(*file)
		if err != nil {
			fmt.Println(err)
			return
		}

	case "listen":
		err := recordAudio(10)
		if err != nil {
			return
		}
	default:
		fmt.Println("unknown command:", os.Args[1])
	}

}

func runAdd(file string, songID string) error {

	index := db.NewIndex()

	samples, err := audio.LoadAudio(file)
	if err != nil {
		fmt.Println(err)
		return err

	}

	spec := spectrogram.GenerateSpectrogram(samples)
	fmt.Println(samples[:20])

	p := peak.DetectPeaks(spec)

	fps := fingerprint.GenerateFingerprints(p)

	index.Add(songID, fps)

	err = index.Save("fingerprints.db")
	if err != nil {
		fmt.Println("save error:", err)
		return err
	}

	fmt.Println("song indexed:", songID)
	return nil

}

func runMatch(file string) error {

	index, err := db.LoadIndex("fingerprints.db")
	if err != nil {
		fmt.Println("load error:", err)
		return err
	}

	samples, err := audio.LoadAudio(file)
	if err != nil {
		fmt.Println("load error:", err)
		return err
	}
	// fmt.Println(samples)

	spec := spectrogram.GenerateSpectrogram(samples)
	// fmt.Println(spec[:30])

	p := peak.DetectPeaks(spec)

	// fmt.Println(p[:3])
	query := fingerprint.GenerateFingerprints(p)
	// fmt.Println(query[:2])
	song, score := matcher.Match(index, query)

	fmt.Println("match:", song, "score:", score)

	return nil
}

func recordAudio(duration int) error {
	filePAth, err := audio.RecordAudio(duration)
	if err != nil {
		return err
	}
	err = runMatch(filePAth)
	if err != nil {
		return err
	}
	return nil
}
