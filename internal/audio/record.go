package audio

import (
	"fmt"
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
		return "", fmt.Errorf("getwd: %w", err)
	}
	tempDir := filepath.Join(dir, "temp")
	cmd.Dir = tempDir

	_, err = cmd.Output()
	if err != nil {
		return "", fmt.Errorf("ffmpeg: %w", err)
	}

	return filepath.Join(tempDir, "record.wav"), nil
}
