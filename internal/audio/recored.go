package audio

import (
	"os"
	"os/exec"
	"path/filepath"
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

	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	filePath := filepath.Join(dir, "/temp")
	cmd.Dir = filePath

	_, err = cmd.Output()
	if err != nil {
		return "", err
	}

	// fmt.Println(string(output))


	return "temp/record.wav", nil
}
