# Development Guide  

This guide covers development practices, testing, and contribution guidelines.

## Development Setup

### Prerequisites

- Go 1.23.4 or later
- Git
- Make (optional)

### Quick Start

```bash
git clone https://github.com/khanbekov/go-bitget.git
cd go-bitget
go mod download
go test ./...
go build .
```

## Testing

```bash
go test ./...
go test -cover ./...
go test -race ./...
```

For detailed development information, see CLAUDE.md.
