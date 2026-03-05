package audio

import (
	"fmt"
	"os/exec"
	"strconv"
)

func RecordAudio(duration int) (string, error) {
	strDuration := strconv.Itoa(duration)
	cmd := exec.Command(
		"ffmpeg",
		"-y",
		"-f", "pulse",
		"-i", "default",
		"-t", strDuration,
		"record.wav",
	)
	err := cmd.Run()
	fmt.Println("Recording...")
	if err != nil {
		return "", err
	}

	return "record.wav", nil
}
