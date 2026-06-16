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
		"/", "_", "\\", "_", ":", "_",
		"*", "_", "?", "_", "\"", "_",
		"<", "_", ">", "_", "|", "_",
	)
	return replacer.Replace(name)
}

// DownloadAudio downloads audio from a YouTube URL using yt-dlp. It
// returns the sanitized video title and the absolute path to the
// downloaded audio file (in m4a format).
func DownloadAudio(url string) (string, string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", "", err
	}

	downloadDir := filepath.Join(currentDir, "temp", "downloaded")
	err = os.MkdirAll(downloadDir, os.ModePerm)
	if err != nil {
		return "", "", fmt.Errorf("failed creating download directory: %w", err)
	}

	yt, err := findYtDlp()
	if err != nil {
		return "", "", err
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

	args := append(baseArgs, "--print", "title", "--skip-download", url)
	infoCmd := exec.Command(yt, args...)
	titleOut, err := infoCmd.Output()
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			return "", "", fmt.Errorf("yt-dlp failed: %s", string(ee.Stderr))
		}
		return "", "", fmt.Errorf("failed to get video info: %w", err)
	}
	title := sanitizeFilename(strings.TrimSpace(string(titleOut)))

	outTmpl := filepath.Join(downloadDir, "%(title)s.%(ext)s")
	args = append(baseArgs, "-x", "--audio-format", "m4a", "-o", outTmpl, url)
	cmd := exec.Command(yt, args...)
	err = cmd.Run()
	if err != nil {
		return "", "", fmt.Errorf("download failed: %w", err)
	}

	filePath := filepath.Join(downloadDir, title+".m4a")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		entries, _ := os.ReadDir(downloadDir)
		for _, e := range entries {
			if strings.HasPrefix(e.Name(), title) {
				filePath = filepath.Join(downloadDir, e.Name())
				break
			}
		}
	}

	return title, filePath, nil
}
