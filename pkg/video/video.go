package video

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/vansante/go-ffprobe"

	"wallplayer/pkg/config"
)

type cacheEntry struct {
	info      *VideoInfo
	timestamp time.Time
}

var (
	cache     = make(map[string]cacheEntry)
	cacheLock sync.RWMutex
	// Cache for 1 hour
	cacheDuration = 1 * time.Hour
)

var Extensions = map[string]bool{
	".mp4":  true,
	".webm": true,
	".mkv":  true,
	".avi":  true,
	".mov":  true,
	".m4v":  true,
}

func IsVideo(path string) bool {
	ext := filepath.Ext(path)
	return Extensions[strings.ToLower(ext)]
}

type SubtitleInfo struct {
	Language    string `json:"language"`
	Title       string `json:"title,omitempty"`
	StreamIndex int    `json:"streamIndex"`
	Codec       string `json:"codec"`
}

type VideoInfo struct {
	Duration  float64        `json:"duration"` // in seconds
	Width     int            `json:"width"`
	Height    int            `json:"height"`
	Bitrate   int64          `json:"bitrate"`
	Format    string         `json:"format"`
	Subtitles []SubtitleInfo `json:"subtitles,omitempty"`
}

func GetInfo(path string) (*VideoInfo, error) {
	// Check the cache
	cacheLock.RLock()
	if entry, ok := cache[path]; ok {
		if time.Since(entry.timestamp) < cacheDuration {
			cacheLock.RUnlock()
			return entry.info, nil
		}
	}
	cacheLock.RUnlock()

	// If not in cache or expired, load from file
	data, err := ffprobe.GetProbeData(path, 3*time.Second)
	if err != nil {
		return nil, err
	}

	bitrate := int64(0)
	if data.Format.BitRate != "" {
		fmt.Sscanf(data.Format.BitRate, "%d", &bitrate)
	}

	info := &VideoInfo{
		Duration: data.Format.DurationSeconds,
		Format:   data.Format.FormatName,
		Bitrate:  bitrate,
	}

	// Store in cache
	cacheLock.Lock()
	cache[path] = cacheEntry{
		info:      info,
		timestamp: time.Now(),
	}
	cacheLock.Unlock()

	// Get video and subtitle stream info
	for _, stream := range data.Streams {
		switch stream.CodecType {
		case "video":
			info.Width = stream.Width
			info.Height = stream.Height
		case "subtitle":
			lang := "und"
			if stream.Tags.Language != "" {
				lang = stream.Tags.Language
			}

			subtitle := SubtitleInfo{
				StreamIndex: stream.Index,
				Codec:       stream.CodecName,
				Language:    lang,
			}

			info.Subtitles = append(info.Subtitles, subtitle)
		}
	}

	return info, nil
}

// GetSubtitlePath returns the path where a subtitle file for the given video and language should be stored
func GetSubtitlePath(videoPath, language string) string {
	baseName := filepath.Base(videoPath[:len(videoPath)-len(filepath.Ext(videoPath))])
	return filepath.Join(config.SubtitlesDir, fmt.Sprintf("%s_%s.vtt", baseName, language))
}

// EnsureSubtitle ensures subtitle file exists for the given video and language, extracting it if needed.
// Returns the path to the subtitle file or an error.
func EnsureSubtitle(videoPath, language string) (string, error) {
	subtitlePath := GetSubtitlePath(videoPath, language)

	// Check if subtitle already exists
	if _, err := os.Stat(subtitlePath); err == nil {
		return subtitlePath, nil
	}

	// Get video info to find stream index for this language
	info, err := GetInfo(videoPath)
	if err != nil {
		return "", fmt.Errorf("failed to get video info: %w", err)
	}

	// Find subtitle stream with matching language
	var streamIndex int = -1
	for _, sub := range info.Subtitles {
		if sub.Language == language {
			streamIndex = sub.StreamIndex
			break
		}
	}
	if streamIndex == -1 {
		return "", fmt.Errorf("no subtitle found for language %s", language)
	}

	// Extract subtitle
	if err := ExtractSubtitle(videoPath, streamIndex, subtitlePath); err != nil {
		return "", fmt.Errorf("failed to extract subtitle: %w", err)
	}

	return subtitlePath, nil
}

// ExtractSubtitle extracts a subtitle stream from a video file and saves it as WebVTT
func ExtractSubtitle(videoPath string, streamIndex int, outputPath string) error {
	// Ensure the output directory exists
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Build ffmpeg command to extract directly to WebVTT
	args := []string{
		"-i", videoPath,
		"-map", fmt.Sprintf("0:%d", streamIndex),
		"-f", "webvtt",
		"-c:s", "webvtt",
		outputPath,
	}

	// Execute ffmpeg command to extract subtitles
	cmd := exec.Command("ffmpeg", args...)
	// Suppress ffmpeg stderr output
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to extract subtitles: %w", err)
	}

	// Check if extraction produced a non-empty file
	if stat, err := os.Stat(outputPath); os.IsNotExist(err) || stat.Size() == 0 {
		os.Remove(outputPath) // Clean up empty file
		return fmt.Errorf("failed to extract subtitles or empty subtitle stream")
	}

	return nil
}
