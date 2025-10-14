# Agent Guidelines for WallPlayer

## Stack
- Frontend: HTMX + Pico CSS + Material Icons
- Backend: Go

## Build & Test Commands
- Development: `go run cmd/main.go`
- Build: `go build -o wallplayer cmd/main.go`
- Run: `./wallplayer`
- Test all: `go test ./...`
- Single test: `go test ./path/to/package -run TestName`
- Lint: `go vet ./...`
- Format: `go fmt ./...`

## Code Style Guidelines

### Frontend
- Use semantic HTML5 elements
- HTMX attributes for interactivity - minimize JavaScript
- Follow size conventions: Base: 20px, Directory: 24px, Title: 26px
- Touch targets: min 48x48px, Vertical rhythm: 24px

#### CSS Guidelines
- Use CSS nesting for related styles
- Follow BEM-like naming for components
- Use CSS variables for themes and colors
- Prefer flexbox/grid for layouts
- Ensure consistent button sizes to prevent layout shifts
- Handle both light and dark themes
- Maintain smooth transitions for state changes

#### JavaScript Guidelines
- Use camelCase for function names
- Use localStorage for persistence
- Use event listeners for DOM manipulation
- Minimize global state

### Backend (Go)
- **Imports**: Group as std lib, external, internal
- **Error handling**: Explicit with meaningful messages, use %w for wrapping
- **Naming**: camelCase for unexported, PascalCase for exported functions/structs
- **Constants**: Use for default values and magic numbers
- **Error variables**: Define common errors (e.g., ErrInvalidPath)
- **Concurrency**: Use goroutines/channels, sync.WaitGroup for coordination
- **Logging**: Use log.Printf with context
- **Paths**: Use filepath.Clean/Rel, validate with filepath.Abs
- **Timestamps**: Use time.RFC3339 format
- **Interfaces**: Use for testability
- **Caching**: Implement for expensive operations
- **HTTP responses**: 400 for invalid input, 500 for internal errors
- **Configuration**: Use environment variables (VIDEOS_DIR, PORT, etc.)

## Safety Guidelines
- Never modify AGENTS.md without explicit user approval
- Never modify files without user's explicit approval
- Never commit changes without user's explicit request
- All code and comments must be in English only
