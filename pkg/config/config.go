package config

import (
	"os"
	"path/filepath"
	"strconv"
)

const (
	// Default values
	DefaultPort         = 9999
	DefaultVideosDir    = "videos"
	DefaultGeneratedDir = "data"
)

var (
	// Paths
	ThumbnailsDir = filepath.Join(DefaultGeneratedDir, "thumbnails")
	SubtitlesDir  = filepath.Join(DefaultGeneratedDir, "subtitles")

	// Runtime configuration
	Port = getPort()
)

// getPort returns the port number from environment variable or default
func getPort() int {
	if envPort := os.Getenv("PORT"); envPort != "" {
		if port, err := strconv.Atoi(envPort); err == nil {
			return port
		}
	}
	return DefaultPort
}

// EnsureDirectories creates necessary directories if they don't exist
func EnsureDirectories() error {
	dirs := []string{
		ThumbnailsDir,
		SubtitlesDir,
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	return nil
}
