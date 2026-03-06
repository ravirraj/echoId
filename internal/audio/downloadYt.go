package audio

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/kkdai/youtube/v2"
)

func DownloadAudio(url string) (string, string, error) {

	client := youtube.Client{}

	video, err := client.GetVideo(url)
	if err != nil {
		return "", "", err
	}

	title := video.Title

	formats := video.Formats.WithAudioChannels()

	if len(formats) == 0 {
		return "", "", fmt.Errorf("no audio formats found")
	}

	best := &formats[0]
	for i := range formats {
		if formats[i].Bitrate > best.Bitrate {
			best = &formats[i]
		}
	}

	stream, _, err := client.GetStream(video, best)
	if err != nil {
		return "", "", err
	}

	currentDir, _ := os.Getwd()

	filePath := filepath.Join(currentDir, "temp","downloaded", title+".m4a")

	file, err := os.Create(filePath)
	if err != nil {
		return "", "", err
	}

	defer file.Close()

	_, err = io.Copy(file, stream)
	if err != nil {
		return "", "", err
	}

	return title, filePath, nil
}
