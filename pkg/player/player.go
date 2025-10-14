package player

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"wallplayer/pkg/browse"
	"wallplayer/pkg/video"
)

var (
	ErrInvalidPath = errors.New("invalid path: must be within Videos directory")
)

func Stream(w http.ResponseWriter, r *http.Request, path string) error {
	// Validate path
	fullPath, err := validatePath(path)
	if err != nil {
		return err
	}

	// Open the video file
	file, err := os.Open(fullPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Get file info
	info, err := file.Stat()
	if err != nil {
		return err
	}
	fileSize := info.Size()

	// Check if the file is a video
	if !isVideoFile(fullPath) {
		return errors.New("not a video file")
	}

	// Parse range header
	rangeHeader := r.Header.Get("Range")
	if rangeHeader == "" {
		// No range requested, serve full file
		w.Header().Set("Content-Length", strconv.FormatInt(fileSize, 10))
		w.Header().Set("Content-Type", getContentType(fullPath))
		_, err = io.Copy(w, file)
		return err
	}

	// Parse the range header
	start, end, err := parseRange(rangeHeader, fileSize)
	if err != nil {
		return err
	}

	// Seek to start position
	if _, err := file.Seek(start, 0); err != nil {
		return err
	}

	// Set headers for partial content
	w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, fileSize))
	w.Header().Set("Accept-Ranges", "bytes")
	w.Header().Set("Content-Length", strconv.FormatInt(end-start+1, 10))
	w.Header().Set("Content-Type", getContentType(fullPath))
	// Ne pas écrire le status si c'est une requête HEAD
	if r.Method != "HEAD" {
		w.WriteHeader(http.StatusPartialContent)
	}

	// Stream the content
	_, err = io.CopyN(w, file, end-start+1)
	return err
}

func validatePath(path string) (string, error) {
	// Clean the path to resolve any ".." or "." components
	cleanPath := filepath.Clean(path)

	// Join with base directory
	fullPath := filepath.Join(browse.BaseDir, cleanPath)

	// Get absolute path
	absPath, err := filepath.Abs(fullPath)
	if err != nil {
		return "", ErrInvalidPath
	}

	// Check if the path is within BaseDir
	if !strings.HasPrefix(absPath, browse.BaseDir) {
		return "", ErrInvalidPath
	}

	return absPath, nil
}

func parseRange(rangeHeader string, fileSize int64) (start int64, end int64, err error) {
	// Expected format: "bytes=0-1023"
	parts := strings.Split(strings.TrimPrefix(rangeHeader, "bytes="), "-")
	if len(parts) != 2 {
		return 0, 0, errors.New("invalid range format")
	}

	// Parse start
	if parts[0] == "" {
		start = 0
	} else {
		start, err = strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			return 0, 0, err
		}
	}

	// Parse end
	if parts[1] == "" {
		end = fileSize - 1
	} else {
		end, err = strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return 0, 0, err
		}
	}

	// Validate range
	if start < 0 || end >= fileSize || start > end {
		return 0, 0, errors.New("invalid range values")
	}

	return start, end, nil
}

func getContentType(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".mp4":
		return "video/mp4"
	case ".webm":
		return "video/webm"
	case ".mkv":
		return "video/x-matroska"
	case ".avi":
		return "video/x-msvideo"
	case ".mov":
		return "video/quicktime"
	default:
		return "application/octet-stream"
	}
}

func isVideoFile(path string) bool {
	return video.IsVideo(path)
}
