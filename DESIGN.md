# WallPlayer Design Document

## Project Structure

```
wallplayer/
├── cmd/
│   └── main.go          # Main entry point and HTTP handlers
├── pkg/
│   ├── browse/          # Directory browsing
│   │   └── browse.go
│   ├── player/          # Video streaming
│   │   └── player.go
│   └── video/           # Video processing
│       ├── thumbnail.go # Thumbnail generation
│       └── video.go     # Video info and metadata
├── web/
│   ├── static/          # Static files (embedded in production)
│   │   ├── css/
│   │   │   └── style.css
│   │   ├── img/
│   │   │   └── no-preview.jpg
│   │   └── index.html
│   └── embed.go         # Static files embedding
├── data/               # Dynamic generated files
│   ├── thumbnails/     # Generated video thumbnails
│   └── subtitles/      # Generated video subtitles
├── go.mod
└── go.sum
```

## API Design

### Directory Browsing

```go
// List directory contents (JSON)
GET /api/browse?path={path}
Response: {
  "path": "string",
  "items": [
    {
      "name": "string",
      "type": "directory|video",
      "path": "string",
      "size": number,      // File size in bytes
      "duration": number,  // Only for videos (seconds)
      "updatedAt": string // Last modified time in RFC3339 format
    }
  ]
}

// List directory contents (HTML fragment)
GET /api/browse/html?path={path}
Response: HTML fragment with directory listing
```

### Video Handling

```go
// Get video details
GET /api/video?path={path}
Response: {
  "path": "string",
  "type": "video",
  "info": {
    "duration": number,   // seconds
    "width": number,
    "height": number,
    "bitrate": number,
    "format": "string",
    "subtitles": [       // Optional subtitle streams
      {
        "language": "string",     // ISO 639-1 language code
        "title": "string",        // Optional title
        "streamIndex": number,    // FFmpeg stream index
        "codec": "string"         // Subtitle codec (e.g., subrip)
      }
    ]
  }
}

// Stream video file
GET /api/video/stream?path={path}
Response: Binary video stream (supports range requests)

// Get video thumbnail
GET /api/video/thumbnail?path={path}
Response: 302 Redirect to static thumbnail image

// Get video subtitle
GET /api/video/subtitle?path={path}&stream={index}
Response: 302 Redirect to WebVTT subtitle file
```

## Data Models

### Browse

```go
type Item struct {
    Name      string    // File or directory name
    Type      string    // "directory" or "video"
    Path      string    // Relative path from videos root
    Size      int64     // File size in bytes
    Duration  float64   // Video duration in seconds (videos only)
    UpdatedAt string    // Last modified time in RFC3339 format
}
```

### Video

```go
type VideoInfo struct {
    Duration  float64        // Duration in seconds
    Width     int           // Video width in pixels
    Height    int           // Video height in pixels
    Bitrate   int64         // Bitrate in bits per second
    Format    string        // Container format (mp4, mkv, etc)
    Subtitles []SubtitleInfo // Available subtitle streams
}

type SubtitleInfo struct {
    Language    string // ISO 639-1 language code
    Title       string // Optional title
    StreamIndex int    // FFmpeg stream index
    Codec       string // Subtitle codec (e.g., subrip)
}
```

## Implementation Details

### Video Processing

- Uses ffprobe to extract video metadata
- Caches video info for 1 hour to avoid repeated probing
- Supports common video formats: mp4, webm, mkv, avi, mov, m4v

### Thumbnail Generation

Thumbnails are generated using ffmpeg with the following settings:
- Extract frame at 10 seconds to avoid black frames at start
- Scale to 320px width maintaining aspect ratio
- JPEG quality factor 2 (high quality)
- Generated thumbnails are stored in data/thumbnails/
- Falls back to no-preview.jpg if generation fails

### Static and Generated Files

#### Static Files
- Source files in web/static/ (HTML, CSS, JS, images)
- Embedded in binary in production mode
- Served directly from disk in development mode (DEV=1)
- Development mode enables hot-reloading of static files

#### Generated Files
- Stored in data/ directory (not embedded)
- Separate routes for different types:
  - /thumbnails/ → data/thumbnails/
  - /subtitles/ → data/subtitles/
- Served via dedicated FileServer handlers
- Support for proper caching and range requests

### Error Handling

- Invalid paths return 400 Bad Request
- File access errors return 500 Internal Server Error
- Missing thumbnails fallback to no-preview.jpg
- All errors are properly logged for debugging

### User Interface Components

#### Controls
- Play/Pause toggle with state-based icon
- Volume control with custom dropdown menu and persistence
  - Values: 0%, 25%, 50%, 75%, 100%
  - Persisted in localStorage
- Theme toggle (light/dark)
- Autoplay toggle with visual feedback
  - Persisted in localStorage
  - Auto-advances to next video

#### Navigation
- Three-state navigation panel:
  - Hidden (0%)
  - Normal (300px)
  - Expanded (600px) with thumbnails
- Smooth transitions between states
- File browser with nested directory support
- Current playing file highlighting
