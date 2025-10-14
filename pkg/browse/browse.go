package browse

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"wallplayer/pkg/video"
)

const (
	DefaultVideosDir = "videos" // Chemin relatif par défaut
)

var (
	ErrInvalidPath = errors.New("invalid path: must be within Videos directory")
	ErrNoVideosDir = errors.New("VIDEOS_DIR environment variable not set")
	BaseDir        string
)

func Init() error {
	// Get videos directory from environment
	videosDir := os.Getenv("VIDEOS_DIR")
	if videosDir == "" {
		// Fallback to default if not set
		workDir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get working directory: %w", err)
		}
		videosDir = filepath.Join(workDir, DefaultVideosDir)

		// Create directory if it doesn't exist
		if err := os.MkdirAll(videosDir, 0755); err != nil {
			return fmt.Errorf("failed to create videos directory: %w", err)
		}
	}

	// Make sure the directory exists
	if _, err := os.Stat(videosDir); os.IsNotExist(err) {
		return fmt.Errorf("videos directory does not exist: %s", videosDir)
	}

	// Get absolute path
	absPath, err := filepath.Abs(videosDir)
	if err != nil {
		return fmt.Errorf("failed to get absolute path for videos directory: %w", err)
	}

	BaseDir = absPath
	return nil
}

type Item struct {
	Name      string  `json:"name"`
	Path      string  `json:"path"` // Path relative to BaseDir
	FullPath  string  `json:"-"`    // Full filesystem path (not exposed in API)
	Type      string  `json:"type"` // "directory" or "video"
	Size      int64   `json:"size,omitempty"`
	Duration  float64 `json:"duration,omitempty"`
	UpdatedAt string  `json:"updatedAt,omitempty"`
}

func sanitizePath(path string) (string, error) {
	log.Printf("sanitizePath: input path: %q", path)

	// If path is empty or root, use base dir
	if path == "" || path == "/" {
		log.Printf("sanitizePath: using BaseDir: %q", BaseDir)
		return BaseDir, nil
	}

	// Clean the path to resolve any ".." or "." components
	cleanPath := filepath.Join(BaseDir, filepath.Clean(path))
	log.Printf("sanitizePath: cleaned path: %q", cleanPath)

	// Get absolute path
	absPath, err := filepath.Abs(cleanPath)
	if err != nil {
		return "", ErrInvalidPath
	}

	// Check if the path is within BaseDir
	if !strings.HasPrefix(absPath, BaseDir) {
		return "", ErrInvalidPath
	}

	return absPath, nil
}

func List(requestedPath string) ([]Item, error) {
	type workItem struct {
		fullPath string
		info     os.FileInfo
		name     string
		relPath  string
	}

	log.Printf("List: requested path: %q", requestedPath)

	// Sanitize and validate the path
	path, err := sanitizePath(requestedPath)
	if err != nil {
		log.Printf("List: sanitizePath error: %v", err)
		return nil, err
	}

	log.Printf("List: sanitized path: %q", path)

	// Check if path is a directory
	info, err := os.Stat(path)
	if err != nil {
		log.Printf("List: stat error: %v", err)
		return nil, err
	}

	if !info.IsDir() {
		log.Printf("List: path is not a directory: %q", path)
		return nil, fmt.Errorf("path is not a directory: %s", path)
	}

	// Read directory contents
	entries, err := os.ReadDir(path)
	if err != nil {
		log.Printf("List: ReadDir error: %v", err)
		return nil, err
	}

	items := make([]Item, 0)
	var videoItems []workItem
	var wg sync.WaitGroup

	// Premier passage : collecter les dossiers et préparer les vidéos
	for _, entry := range entries {
		// Skip hidden files
		if entry.Name()[0] == '.' {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		fullPath := filepath.Join(path, entry.Name())
		// Create path relative to BaseDir for API response
		relPath, err := filepath.Rel(BaseDir, fullPath)
		if err != nil {
			continue
		}

		if entry.IsDir() {
			items = append(items, Item{
				Name:      entry.Name(),
				Path:      relPath,
				FullPath:  fullPath,
				Type:      "directory",
				UpdatedAt: info.ModTime().Format(time.RFC3339),
			})
			continue
		}

		// Collecter les vidéos pour traitement asynchrone
		if video.IsVideo(entry.Name()) {
			videoItems = append(videoItems, workItem{
				fullPath: fullPath,
				info:     info,
				name:     entry.Name(),
				relPath:  relPath,
			})
		}
	}

	// Traiter les vidéos en parallèle
	resultChan := make(chan Item, len(videoItems))
	for _, item := range videoItems {
		wg.Add(1)
		go func(wi workItem) {
			defer wg.Done()

			// Récupérer les infos de la vidéo
			videoInfo, _ := video.GetInfo(wi.fullPath)
			duration := float64(0)
			if videoInfo != nil {
				duration = videoInfo.Duration
			}

			resultChan <- Item{
				Name:      wi.name,
				Path:      wi.relPath,
				FullPath:  wi.fullPath,
				Type:      "video",
				Size:      wi.info.Size(),
				Duration:  duration,
				UpdatedAt: wi.info.ModTime().Format(time.RFC3339),
			}
		}(item)
	}

	// Attendre la fin des goroutines
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Collecter les résultats
	for item := range resultChan {
		items = append(items, item)
	}

	// Sort items by type (directories first) then by name
	sort.Slice(items, func(i, j int) bool {
		// If types are different, directories come first
		if items[i].Type != items[j].Type {
			return items[i].Type == "directory"
		}
		// If types are the same, sort by name
		return strings.ToLower(items[i].Name) < strings.ToLower(items[j].Name)
	})

	return items, nil
}

//videoExtensions déplacé dans le package player
