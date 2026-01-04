# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

A Go-based CLI application that converts BibTeX entries and URLs to APA 6 format, organized by projects and stored in SQLite.

## Architecture

### Core Components
1. **cmd/**: CLI command definitions using cobra or similar
2. **internal/bibtex/**: BibTeX parser and processor
3. **internal/apa/**: APA 6 formatter
4. **internal/url/**: URL metadata extractor
5. **internal/db/**: SQLite database layer
6. **internal/project/**: Project management logic

### Database Schema
- **projects**: id, name, created_at, updated_at
- **references**: id, project_id, bibtex_entry, apa_format, source_type, created_at
- **metadata**: id, reference_id, key, value (for storing parsed metadata)

## Commands

```bash
# Build the application
go build -o bibapa ./cmd/bibapa

# Run tests
go test ./...

# Run with race detector
go test -race ./...

# Format code
go fmt ./...

# Vet code
go vet ./...

# Install dependencies
go mod tidy

# Run the CLI
./bibapa [command]
```

## CLI Commands Structure

```bash
bibapa project create <name>          # Create new project
bibapa project list                   # List all projects
bibapa project select <name>          # Select active project
bibapa add                           # Add reference (prompts for input)
bibapa list                          # List references in current project
bibapa export                        # Export current project references
bibapa delete <id>                   # Delete a reference
```

## Key Dependencies

```go
// go.mod dependencies to consider
github.com/spf13/cobra              // CLI framework
github.com/mattn/go-sqlite3         // SQLite driver
github.com/nickng/bibtex            // BibTeX parser (or similar)
github.com/PuerkitoBio/goquery      // HTML parsing for URL metadata
```

## Implementation Notes

- Use cobra for CLI structure with subcommands
- Implement interactive prompts for paste input using bufio
- Store both original BibTeX and formatted APA for each reference
- Sort references alphabetically by author surname when displaying
- Handle URL metadata extraction with timeout and error handling
- Use database migrations for schema management