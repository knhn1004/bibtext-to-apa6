# BibTeX to APA 6 Converter

A command-line tool for converting BibTeX entries and URLs to APA 6 format, with project organization, in-text citation generation, and rich text clipboard support.

## Installation

```bash
go install github.com/knhn1004/bibtext-to-apa6/cmd/bibapa@latest
```

Or build from source:

```bash
git clone https://github.com/knhn1004/bibtext-to-apa6.git
cd bibtext-to-apa6
make build
```

## Usage

### Project Management

Create a new project:
```bash
bibapa project create research-paper
```

List all projects:
```bash
bibapa project list
```

Select active project:
```bash
bibapa project select research-paper
```

Delete a project and all its references:
```bash
bibapa project delete research-paper
```

### Adding References

Add a reference (BibTeX or URL):
```bash
bibapa add
# Paste your BibTeX entry or URL when prompted
# Press Enter twice to submit
```

Example BibTeX input:
```bibtex
@article{smith2023,
  author = {Smith, John and Doe, Jane},
  title = {An Example Article},
  journal = {Journal of Examples},
  year = {2023},
  volume = {10},
  number = {2},
  pages = {123-145},
  doi = {10.1234/example}
}
```

Example URL input:
```
https://example.com/article
```

### Managing References

List all references in current project (with stable reference numbers):
```bash
bibapa list
```

Export references with rich text formatting (copied to clipboard):
```bash
bibapa export
```

Delete a reference by its number:
```bash
bibapa delete 3
```

### In-Text Citations

Generate APA 6 in-text citations and copy to clipboard:

```bash
# Single citation
bibapa cite 1

# Multiple citations
bibapa cite 1 2 3

# Range of citations
bibapa cite 1-5
```

Output examples:
- Single author: `(Smith, 2023)`
- Two authors: `(Smith & Jones, 2023)`
- Three+ authors: `(Smith et al., 2023)`
- Multiple sources: `(Jones, 2022; Smith, 2023)`

## Features

- **BibTeX Conversion**: Converts BibTeX entries to APA 6 format
- **URL Metadata Extraction**: Extracts metadata from web pages and formats as APA 6
  - Uses HTTP with browser-like headers for standard pages
  - Falls back to Playwright headless browser for JavaScript-heavy sites
- **Project Organization**: Organize references into separate projects
- **Rich Text Clipboard**: Exports preserve italic formatting when pasting into Word, Google Docs, etc.
- **In-Text Citations**: Generate properly formatted in-text citations
- **Duplicate Detection**: Prevents adding the same reference twice
- **Stable Reference Numbers**: References maintain consistent numbering
- **Cross-Platform**: Works on macOS, Windows, and Linux
- **SQLite Storage**: References stored locally in a SQLite database

## Examples

### BibTeX to APA

Input:
```bibtex
@article{jones2024,
  author = {Jones, Alice B. and Smith, Robert C.},
  title = {Climate Change Impacts on Urban Planning},
  journal = {Environmental Studies Quarterly},
  year = {2024},
  volume = {15},
  number = {3},
  pages = {234-256}
}
```

Output:
```
Jones, A. B., & Smith, R. C. (2024). Climate Change Impacts on Urban Planning. Environmental Studies Quarterly, 15(3), 234-256.
```

### URL to APA

Input:
```
https://www.nature.com/articles/example-2024
```

Output:
```
Nature Publishing Group. (2024). Example Article Title. Nature. Retrieved January 4, 2026, from https://www.nature.com/articles/example-2024
```

## Supported BibTeX Entry Types

- `@article` - Journal articles
- `@book` - Books
- `@inproceedings` - Conference papers
- `@incollection` - Book chapters
- `@misc` - Web pages and other sources
- `@phdthesis` / `@mastersthesis` - Theses and dissertations
- `@techreport` - Technical reports

## Requirements

- Go 1.21+
- For URL metadata extraction from JavaScript-heavy sites: Playwright browsers (`npx playwright install chromium`)
