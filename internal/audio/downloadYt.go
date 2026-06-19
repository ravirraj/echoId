package audio

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func findYtDlp() (string, error) {
	if p, err := exec.LookPath("yt-dlp"); err == nil {
		return p, nil
	}
	dirs := []string{
		filepath.Join(os.Getenv("HOME"), ".local", "bin"),
		"/usr/local/bin",
		"/usr/bin",
		"/opt/homebrew/bin",
	}
	for _, d := range dirs {
		p := filepath.Join(d, "yt-dlp")
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}
	return "", fmt.Errorf("yt-dlp not found; install with: pip3 install --break-system-packages yt-dlp")
}

func browserCookieArg() string {
	for _, b := range []string{"firefox", "brave", "chromium", "chrome", "vivaldi"} {
		var cfgDir string
		switch b {
		case "firefox":
			cfgDir = filepath.Join(os.Getenv("HOME"), ".mozilla", "firefox")
		case "chromium":
			cfgDir = filepath.Join(os.Getenv("HOME"), ".config", "chromium")
		case "chrome":
			cfgDir = filepath.Join(os.Getenv("HOME"), ".config", "google-chrome")
		default:
			cfgDir = filepath.Join(os.Getenv("HOME"), ".config", b)
		}
		if _, err := os.Stat(cfgDir); err == nil {
			return fmt.Sprintf("--cookies-from-browser=%s", b)
		}
	}
	return ""
}

func sanitizeFilename(name string) string {
	replacer := strings.NewReplacer(
		"/", "_", "\\", "_", ":", "_", "：", "_",
		"*", "_", "?", "_", "\"", "_",
		"<", "_", ">", "_", "|", "_", "｜", "_",
		" ", "_",
	)
	return replacer.Replace(name)
}

func cleanTitle(raw string) string {
	if idx := strings.IndexAny(raw, "|｜"); idx != -1 {
		return strings.TrimSpace(raw[:idx])
	}
	return strings.TrimSpace(raw)
}

type YtMeta struct {
	Title  string
	Artist string
	Album  string
}

// DownloadAudio downloads audio from a YouTube URL using yt-dlp. It
// returns metadata (title, artist, album), the absolute path to the
// downloaded audio file (in m4a format), and any error.
func DownloadAudio(url string) (YtMeta, string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return YtMeta{}, "", err
	}

	downloadDir := filepath.Join(currentDir, "temp", "downloaded")
	err = os.MkdirAll(downloadDir, os.ModePerm)
	if err != nil {
		return YtMeta{}, "", fmt.Errorf("failed creating download directory: %w", err)
	}

	yt, err := findYtDlp()
	if err != nil {
		return YtMeta{}, "", err
	}

	cookieArg := browserCookieArg()
	denoPath := filepath.Join(os.Getenv("HOME"), ".local", "bin", "deno")
	jsRuntime := ""
	if _, err := os.Stat(denoPath); err == nil {
		jsRuntime = fmt.Sprintf("deno:%s", denoPath)
	}

	baseArgs := []string{"--remote-components", "ejs:github", "--retries", "5", "--quiet", "--no-warnings"}
	if jsRuntime != "" {
		baseArgs = append(baseArgs, "--js-runtimes", jsRuntime)
	}
	if cookieArg != "" {
		baseArgs = append(baseArgs, cookieArg)
	}

	infoArgs := append(baseArgs, "--print", "%(title)s|||%(artist)s|||%(album)s", "--skip-download", url)
	infoCmd := exec.Command(yt, infoArgs...)
	infoOut, err := infoCmd.Output()
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			return YtMeta{}, "", fmt.Errorf("yt-dlp failed: %s", string(ee.Stderr))
		}
		return YtMeta{}, "", fmt.Errorf("failed to get video info: %w", err)
	}

	parts := strings.SplitN(strings.TrimSpace(string(infoOut)), "|||", 3)
	rawTitle := ""
	if len(parts) > 0 {
		rawTitle = parts[0]
	}
	artist := ""
	if len(parts) > 1 {
		artist = strings.TrimSpace(parts[1])
	}
	album := ""
	if len(parts) > 2 {
		album = strings.TrimSpace(parts[2])
	}

	meta := YtMeta{
		Title:  cleanTitle(rawTitle),
		Artist: artist,
		Album:  album,
	}
	sanitized := sanitizeFilename(rawTitle)

	outTmpl := filepath.Join(downloadDir, "%(title)s.%(ext)s")
	dlArgs := append(baseArgs, "-x", "--audio-format", "m4a", "-o", outTmpl, url)
	cmd := exec.Command(yt, dlArgs...)
	err = cmd.Run()
	if err != nil {
		return YtMeta{}, "", fmt.Errorf("download failed: %w", err)
	}

	filePath := filepath.Join(downloadDir, sanitized+".m4a")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		entries, _ := os.ReadDir(downloadDir)
		for _, e := range entries {
			sanitizedName := sanitizeFilename(e.Name())
			if strings.HasPrefix(sanitizedName, sanitized) {
				filePath = filepath.Join(downloadDir, e.Name())
				break
			}
		}
	}

	return meta, filePath, nil
}
