# BibTeX to APA 6 Converter

A command-line tool for converting BibTeX entries and URLs to APA 6 format, organized by projects.

## Installation

```bash
go install github.com/knhn1004/bibtext-to-apa6/cmd/bibapa@latest
```

Or build from source:

```bash
git clone https://github.com/knhn1004/bibtext-to-apa6.git
cd bibtext-to-apa6
go build -o bibapa ./cmd/bibapa
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

List all references in current project:
```bash
bibapa list
```

Export references (formatted for copy/paste):
```bash
bibapa export
```

Delete a reference by ID:
```bash
bibapa delete 123
```

## Features

- Converts BibTeX entries to APA 6 format
- Extracts metadata from URLs and formats as APA 6
- Organizes references by projects
- Stores references in SQLite database
- Alphabetically sorts references by author
- Supports common BibTeX entry types (article, book, inproceedings, etc.)

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

Output (example):
```
Nature Publishing Group. (2024). Example Article Title. Nature. Retrieved December 29, 2024, from https://www.nature.com/articles/example-2024
```