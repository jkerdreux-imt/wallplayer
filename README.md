# Wallplayer

Wallplayer is a modern, web-based media player designed for large touch screens and collective spaces. It allows you to browse and play videos from a local or network-mounted directory, with a simple, robust, and touch-friendly interface for big screens.

## Overview

![Wallplayer Screenshot](screenshots/1.png)

## Features

- **Single binary, zero install**: Everything is included in one statically built executable—no dependencies, no Python, no Node, no database, nothing to install. Just run the binary and you’re ready.
- **Video thumbnails**: Automatic generation and display of video thumbnails for quick visual navigation. [See screenshot](screenshots/2.png)
- **Seamless browsing**: Browse folders and select new videos while a video is playing, without interrupting playback.
- **Subtitle support**: Display and select subtitles (if available) for your videos.
- **Touch-friendly UI**: Optimized for large touch screens and public/shared environments.
- **Instant playback**: Play videos directly in the browser with no extra plugins.
- **Modern stack**: Built with HTMX, Pico CSS, and Go for speed and simplicity.

## Requirements

Wallplayer requires FFmpeg to be installed on your system for video thumbnail generation and subtitle extraction:

```bash
# Ubuntu/Debian
sudo apt install ffmpeg

# macOS
brew install ffmpeg

# Windows (using Chocolatey)
choco install ffmpeg
```

## Configuration

### Videos Directory

By default, Wallplayer creates and uses a `videos` directory in the current working directory. You can change this by setting the `VIDEOS_DIR` environment variable:

```bash
# Set videos directory
export VIDEOS_DIR=/path/to/your/videos

# Run with custom videos directory
VIDEOS_DIR=/path/to/your/videos ./wallplayer
```

## License

This project is licensed under the GNU General Public License v3.0 (GPLv3).
See the [LICENSE](LICENSE) file for details.

## Commands

### Development

To run in development mode:

```bash
go run cmd/main.go
```

### Production

To build the production executable (recommended idiomatic Go command):

```bash
go build -o wallplayer ./cmd
```

To run the executable:

```bash
./wallplayer
```
