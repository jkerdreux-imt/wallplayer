package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"wallplayer/pkg/browse"
	"wallplayer/pkg/player"
	"wallplayer/pkg/video"
	"wallplayer/web"
)

// setContentType sets the Content-Type header based on file extension
func setContentType(w http.ResponseWriter, path string) {
	switch {
	case strings.HasSuffix(path, ".css"):
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
	case strings.HasSuffix(path, ".js"):
		w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	case strings.HasSuffix(path, ".html"):
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
	case strings.HasSuffix(path, ".jpg"), strings.HasSuffix(path, ".jpeg"):
		w.Header().Set("Content-Type", "image/jpeg")
	case strings.HasSuffix(path, ".png"):
		w.Header().Set("Content-Type", "image/png")
	case strings.HasSuffix(path, ".ico"):
		w.Header().Set("Content-Type", "image/x-icon")
	case strings.HasSuffix(path, ".webp"):
		w.Header().Set("Content-Type", "image/webp")
	}
}

// redirectRoot redirects "/" to "/static/index.html" and returns 404 for other paths
func redirectRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.Redirect(w, r, "/static/", http.StatusFound)
		return
	}
	http.NotFound(w, r)
}

// setupStaticHandlers configures static file serving for dev and prod modes
func setupStaticHandlers(devMode bool) {
	// Handle root path
	http.HandleFunc("/", redirectRoot)

	if devMode {
		log.Println("Running in development mode")
		fs := http.FileServer(http.Dir("web/static"))
		http.Handle("/static/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			setContentType(w, r.URL.Path)
			http.StripPrefix("/static/", fs).ServeHTTP(w, r)
		}))
	} else {
		log.Println("Running in production mode")
		// Handle static files with special handling for index.html
		http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
			// Strip /static/ prefix from path
			path := strings.TrimPrefix(r.URL.Path, "/static/")

			// be default fs.embedded doesn't return index.html for directory
			if path == "" {
				path = "index.html"
			}

			// Read file from embedded filesystem
			data, err := web.StaticFiles.ReadFile("static/" + path)
			if err != nil {
				http.Error(w, "File not found", http.StatusNotFound)
				return
			}

			// Set content type based on file extension
			setContentType(w, path)
			w.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
			w.Write(data)
		})
	}
}

// formatName removes file extension and replaces _ and - with spaces
func formatName(name string) string {
	ext := filepath.Ext(name)
	nameWithoutExt := name[:len(name)-len(ext)]
	formatted := strings.ReplaceAll(nameWithoutExt, "_", " ")
	formatted = strings.ReplaceAll(formatted, "-", " ")
	return formatted
}

func handleBrowseHTML(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	if path == "" {
		path = "/"
	}
	log.Printf("handleBrowseHTML: requested path: %q", path)

	if filepath.Ext(path) != "" {
		http.Error(w, "Cannot list a file", http.StatusBadRequest)
		return
	}

	items, err := browse.List(path)
	if err != nil {
		if err == browse.ErrInvalidPath {
			http.Error(w, "Invalid path", http.StatusBadRequest)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	html := `<ul class="file-list">`
	if path != "/" {
		parentPath := filepath.Dir(path)
		if parentPath == "." {
			parentPath = "/"
		}
		html += fmt.Sprintf(`
			<li hx-get="/api/browse/html?path=%s" hx-trigger="click" hx-target="#path-browser">
				<span class="material-symbols-rounded">folder</span>
				<span>..</span>
			</li>`, parentPath)
	}

	for _, item := range items {
		if item.Type == "directory" {
			html += fmt.Sprintf(`
				<li hx-get="/api/browse/html?path=%s" hx-trigger="click" hx-target="#path-browser">
					<span class="material-symbols-rounded">folder</span>
					<span>%s</span>
				</li>`, item.Path, formatName(item.Name))
		} else {
			durationStr := "⋯"
			if item.Duration > 0 {
				minutes := int(item.Duration / 60)
				seconds := int(item.Duration) % 60
				durationStr = fmt.Sprintf("%02d:%02d", minutes, seconds)
			}
			html += fmt.Sprintf(`
				<li onclick="playVideo('%s')">
					<span class="material-symbols-rounded video" data-hide-in-expanded="true">movie_info</span>
					<img class="thumbnail" src="/api/video/thumbnail?path=%s" loading="lazy" alt="">
					<span class="name">%s</span>
					<span class="duration">%s</span>
				</li>`, item.Path, item.Path, formatName(item.Name), durationStr)
		}
	}
	html += "</ul>"

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

func handleBrowseAPI(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	if path == "" {
		path = "/"
	}
	fullPath := filepath.Join(browse.BaseDir, path)
	items, err := browse.List(fullPath)
	if err != nil {
		if err == browse.ErrInvalidPath {
			http.Error(w, "Invalid path", http.StatusBadRequest)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}
	response := struct {
		Path  string        `json:"path"`
		Items []browse.Item `json:"items"`
	}{
		Path:  path,
		Items: items,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleVideoAPI(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	if path == "" {
		http.Error(w, "path parameter required", http.StatusBadRequest)
		return
	}
	fullPath := filepath.Join(browse.BaseDir, path)
	info, err := video.GetInfo(fullPath)
	if err != nil {
		log.Printf("Error getting video info: %v", err)
		http.Error(w, "Error reading video info", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Path string           `json:"path"`
		Type string           `json:"type"`
		Info *video.VideoInfo `json:"info"`
	}{
		Path: path,
		Type: "video",
		Info: info,
	})
}

func handleVideoSubtitle(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	lang := r.URL.Query().Get("lang")
	if path == "" || lang == "" {
		http.Error(w, "path and lang parameters required", http.StatusBadRequest)
		return
	}
	fullPath := filepath.Join(browse.BaseDir, path)
	subtitlePath, err := video.EnsureSubtitle(fullPath, lang)
	if err != nil {
		log.Printf("Error handling subtitle: %v", err)
		switch {
		case strings.Contains(err.Error(), "no subtitle found for language"):
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, "Error handling subtitle", http.StatusInternalServerError)
		}
		return
	}
	http.Redirect(w, r, "/subtitles/"+filepath.Base(subtitlePath), http.StatusFound)
}

func handleVideoThumbnail(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	if path == "" {
		http.Error(w, "path parameter required", http.StatusBadRequest)
		return
	}
	fullPath := filepath.Join(browse.BaseDir, path)
	thumbPath, err := video.GenerateThumbnail(fullPath)
	if err != nil {
		log.Printf("Error generating thumbnail: %v", err)
		http.Error(w, "Error generating thumbnail", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, thumbPath, http.StatusFound)
}

func handleVideoStream(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	if path == "" {
		http.Error(w, "path parameter required", http.StatusBadRequest)
		return
	}
	err := player.Stream(w, r, path)
	if err != nil {
		if err == player.ErrInvalidPath {
			http.Error(w, "Invalid path", http.StatusBadRequest)
		} else {
			log.Printf("Streaming error: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}
}
