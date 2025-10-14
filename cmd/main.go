package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"wallplayer/pkg/browse"
	"wallplayer/pkg/config"
)

// Les handlers sont dans handlers.go

func main() {
	port := fmt.Sprintf(":%d", config.Port)
	// Initialize videos directory
	if err := browse.Init(); err != nil {
		log.Fatalf("Failed to initialize videos directory: %v", err)
	}
	log.Printf("Videos directory: %s", browse.BaseDir)

	// Create generated directories
	if err := config.EnsureDirectories(); err != nil {
		log.Fatalf("Failed to create required directories: %v", err)
	}

	// Dev mode detection
	devMode := os.Getenv("DEV") == "1"

	// Configure static file serving based on mode
	setupStaticHandlers(devMode)

	// Handle generated directories
	http.Handle("/thumbnails/", http.StripPrefix("/thumbnails/", http.FileServer(http.Dir(config.ThumbnailsDir))))
	http.Handle("/subtitles/", http.StripPrefix("/subtitles/", http.FileServer(http.Dir(config.SubtitlesDir))))

	// API routes
	http.HandleFunc("/api/browse", handleBrowseAPI)
	http.HandleFunc("/api/browse/html", handleBrowseHTML)
	http.HandleFunc("/api/video", handleVideoAPI)
	http.HandleFunc("/api/video/stream", handleVideoStream)
	http.HandleFunc("/api/video/thumbnail", handleVideoThumbnail)
	http.HandleFunc("/api/video/subtitle", handleVideoSubtitle)

	log.Println("Starting server on " + port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
