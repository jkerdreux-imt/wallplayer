# Video Library Presentation Interface Specification

## Context & Use Case

- 75" touch screen display
- Presenter stands on the left side
- Used for video presentations/demonstrations
- Experimental videos organized in directories
- Live presentation context, no room for errors

## User Experience Requirements

- Quick video selection during presentation
- No interruption in presenter's flow
- Clear visibility from standing position
- Easy touch interaction while talking
- Video continues playing while browsing

## Interface Layout

### Left Panel (Navigation)

- Width: Adjustable (25-50% of screen)
- Content: Directory tree + video files
- Height: Full screen, optimized for standing interaction
- Must support both:
  - List view (compact, efficient)
  - Thumbnail view (visual recognition)

### Right Panel (Player)

- Width: Remaining screen space
- Content: Video player only
- Clean interface, no distractions
- Maintains playback during navigation

## Technical Stack

### Frontend

- HTMX for interactivity
  - Minimal JavaScript
  - Partial page updates
  - Server-side driven
- Pico CSS framework
  - Challenges:
    - Default font sizes too small for 75" display
    - Need to adjust spacing for touch
    - Custom CSS required for thumbnail view
  - Benefits:
    - Clean, modern defaults
    - Dark/light themes
    - Semantic HTML
- Material Icons
  - folder: for directories
  - videocam: for video files

### Backend (Go)

- File system navigation
- Video streaming (with embedded subtitles)
- Thumbnail generation
- Directory scanning

## Key Challenges

### Performance

- Thumbnail generation and caching
- Video streaming optimization
- Directory scanning with large files
- Quick interface response on touch

### UI/UX

- Touch target sizes (minimum 48x48px)
- Standing user ergonomics
- Panel resize handling
- Visibility at presentation distance
- Clear current video indication

### Technical
- Video playback continuation during navigation
- Thumbnail generation for various formats
- Deep directory structure handling
- File system performance
- Native HTML5 video subtitle track support

## Views & Modes

### List View

- Default view
- Compact representation
- Shows:
  - Directory/file icons
  - Names
  - Basic metadata
- Optimized for known content

### Thumbnail View

- Optional view
- 2-3 thumbnails per row
- Shows:
  - Video preview image
  - File name
  - Duration if available
- Better for content discovery

## Interaction Design

### Touch Interactions

- Large touch targets
- Clear visual feedback
- No complex gestures
- Simple tap operations

### Navigation

- Clear current position indication
- Easy return to parent directory
- Current video highlighted
- Next video easily selectable

### Video Control

- Basic playback controls
- Volume adjustment
- Clear current status
- Simple seek functionality

## Visual Design

### Layout Proportions

- Left panel: 25% (list) to 50% (thumbnails)
- Right panel: 75% to 50% respectively
- Min touch target: 48x48px
- Comfortable text size for standing distance

### Typography (Pico CSS Override)

- Base size: 20px (up from Pico's default)
- Directory items: 24px
- Current path: 22px
- Video title: 26px
- Touch targets: minimum 48px height

### Colors

- Based on Pico CSS dark/light themes
- High contrast for standing visibility
- Non-distracting during presentation
- Clear active states

### Spacing

- Vertical rhythm: 24px base
- Touch targets: 12px minimum separation
- List view: 48px row height
- Thumbnail view: 180px tile height

## Future Considerations

- Keyboard support for admin
- Remote control integration
- Presentation mode shortcuts
- Multi-screen support
- Network share browsing

## Development Phases

1. Basic Structure
   - Layout implementation
   - File system navigation
   - Basic video playback

2. Enhanced Features
   - Thumbnail generation
   - View switching
   - Panel resizing

3. UI Refinement
   - Touch optimization
   - Visual feedback
   - Performance optimization

4. Testing & Adjustment
   - Real environment testing
   - User feedback
   - Performance tuning

