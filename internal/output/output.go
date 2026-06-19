package output

import (
	"fmt"
	"os"
	"strings"
)

const (
	clearLine = "\r\033[2K"
	reset     = "\033[0m"
	bold      = "\033[1m"
	dim       = "\033[2m"
	green     = "\033[32m"
	cyan      = "\033[36m"
	yellow    = "\033[33m"
	red       = "\033[31m"
	white     = "\033[97m"
)

func init() {
	if os.Getenv("NO_COLOR") != "" {
		_ = reset
	}
}

func step(msg string) {
	fmt.Fprintf(os.Stderr, "%s  %s%s\n", dim, msg, reset)
}

func Stepf(format string, args ...interface{}) {
	step(fmt.Sprintf(format, args...))
}

func Listening() {
	step(fmt.Sprintf("%sListening...%s", white, reset))
}

func Generating() {
	step(fmt.Sprintf("%sGenerating fingerprints...%s", white, reset))
}

func Matching() {
	step(fmt.Sprintf("%sMatching...%s", white, reset))
}

func PrintError(msg string) {
	fmt.Fprintf(os.Stderr, "%s  ✗ %s%s\n", red, msg, reset)
}

func PrintNoMatch() {
	fmt.Fprintf(os.Stderr, "\n%s  No match found.%s\n\n", dim, reset)
}

func PrintSongFound() {
	fmt.Fprintf(os.Stderr, "\n%s%s ✓ Song Found!%s\n\n", bold, green, reset)
}

func PrintResult(title, artist, album string, duration float64) {
	PrintSongFound()

	printField("Title", title, 12)
	printField("Artist", artist, 12)
	printField("Album", album, 12)
	printField("Time", fmt.Sprintf("%.2fs", duration), 12)
	fmt.Println()
}

func printField(key, value string, pad int) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		trimmed = "N/A"
	}
	fmt.Fprintf(os.Stderr, "  %s%s%s : %s\n", bold, key, reset, trimmed)
}

func PrintIndexed(songID string) {
	fmt.Fprintf(os.Stderr, "\n%s%s ✓ Indexed:%s %s\n\n", bold, green, reset, songID)
}

func PrintFingerprintStats(peaks, fingerprints int) {
	step(fmt.Sprintf("Peaks: %d | Fingerprints: %d", peaks, fingerprints))
}
