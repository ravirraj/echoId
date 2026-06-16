package audio

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

// RecordAudio captures audio from the default PulseAudio source for the
// given duration (in seconds) and returns the absolute path to the
// resulting WAV file.
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
