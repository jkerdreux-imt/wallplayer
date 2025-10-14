package video

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"wallplayer/pkg/config"
)

// GenerateThumbnail generates a thumbnail for a video
// Returns the path to the generated thumbnail or path to no-preview image if generation fails
func GenerateThumbnail(videoPath string) (string, error) {
	// Check if video file exists
	if _, err := os.Stat(videoPath); os.IsNotExist(err) {
		return "/static/img/no-preview.jpg", nil
	}

	// Use generated thumbnails directory
	thumbPath := filepath.Join(config.ThumbnailsDir, filepath.Base(videoPath)+".jpg")

	// Generate thumbnail using ffmpeg with the following arguments:
	cmd := exec.Command("ffmpeg",
		"-v", "error", // Only show errors in output
		"-ss", "10", // Seek to 1o second (avoid black frames at start)
		"-i", videoPath, // Input file
		"-frames:v", "1", // Extract exactly one frame
		"-q:v", "2", // Quality factor (2-31, lower is better quality)
		"-vf", "scale=320:-1", // Scale width to 320px, height auto (-1)
		"-y",      // Overwrite output file if exists
		thumbPath, // Output file path
	)

	if err := cmd.Run(); err != nil {
		log.Printf("Error generating thumbnail for %s: %v", videoPath, err)
		return "/static/img/no-preview.jpg", nil
	}

	return "/thumbnails/" + filepath.Base(thumbPath), nil
}
