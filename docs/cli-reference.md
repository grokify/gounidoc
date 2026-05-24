# CLI Reference

GoUniDoc provides a unified CLI with 16 commands for document conversion and processing.

## Overview

```bash
gounidoc [command] [flags]
```

Running without a command starts the MCP server (default behavior).

## Document Conversion Commands

### pdf2docx

Convert a PDF file to DOCX format with optional table, image, and layout extraction.

```bash
gounidoc pdf2docx -i input.pdf -o output.docx
```

**Flags:**

| Flag | Short | Description |
|------|-------|-------------|
| `--input` | `-i` | Input PDF file path (required) |
| `--output` | `-o` | Output DOCX file path (optional, defaults to input name with .docx) |
| `--basic` | `-b` | Basic mode: text only, no tables/images/layout |
| `--no-tables` | | Disable table extraction |
| `--no-images` | | Disable image detection |
| `--no-layout` | | Disable layout/heading detection |

**Examples:**

```bash
# Full conversion with all features
gounidoc pdf2docx -i document.pdf -o document.docx

# Basic mode (fastest, text only)
gounidoc pdf2docx -i document.pdf -b

# Disable specific features
gounidoc pdf2docx -i document.pdf --no-tables --no-images
```

**Conversion Features:**

| Feature | Description |
|---------|-------------|
| Tables | Auto-detects tables and converts to DOCX tables with borders |
| Images | Detects images and notes their dimensions (placeholder text) |
| Layout | Detects headings (numbered sections, ALL CAPS, short lines) and applies bold formatting |

---

### docx2pdf

Convert a DOCX file to PDF format.

```bash
gounidoc docx2pdf -i document.docx -o output.pdf
```

**Flags:**

| Flag | Short | Description |
|------|-------|-------------|
| `--input` | `-i` | Input DOCX file path (required) |
| `--output` | `-o` | Output PDF file path (optional, defaults to input name with .pdf) |

---

### pdf2png

Convert a PDF page to PNG image format.

```bash
gounidoc pdf2png -i document.pdf -o page1.png --page 1 --width 1200
```

**Flags:**

| Flag | Short | Description |
|------|-------|-------------|
| `--input` | `-i` | Input PDF file path (required) |
| `--output` | `-o` | Output PNG file path (required) |
| `--page` | | Page number to convert (1-indexed, default: 1) |
| `--width` | | Output width in pixels (default: 1200) |

---

### batch-pdf2docx

Convert multiple PDF files to DOCX format with progress reporting.

```bash
gounidoc batch-pdf2docx -i "*.pdf" -d ./output/
```

**Flags:**

| Flag | Short | Description |
|------|-------|-------------|
| `--input` | `-i` | Input PDF files (can be specified multiple times, supports glob patterns) |
| `--dir` | `-d` | Output directory (optional, defaults to same directory as input) |
| `--basic` | `-b` | Basic mode: text only |
| `--no-tables` | | Disable table extraction |
| `--no-images` | | Disable image detection |
| `--no-layout` | | Disable layout/heading detection |

**Examples:**

```bash
# Convert all PDFs using glob pattern
gounidoc batch-pdf2docx -i "*.pdf" -d ./converted/

# Convert specific files
gounidoc batch-pdf2docx -i file1.pdf -i file2.pdf -i file3.pdf -d ./output/

# Basic mode for faster conversion
gounidoc batch-pdf2docx -i "*.pdf" -b -d ./output/
```

---

### batch-docx2pdf

Convert multiple DOCX files to PDF format with progress reporting.

```bash
gounidoc batch-docx2pdf -i "*.docx" -d ./output/
```

**Flags:**

| Flag | Short | Description |
|------|-------|-------------|
| `--input` | `-i` | Input DOCX files (can be specified multiple times, supports glob patterns) |
| `--dir` | `-d` | Output directory (optional, defaults to same directory as input) |

---

## PDF Operations Commands

### merge

Merge multiple PDF files into a single output file.

```bash
gounidoc merge -i file1.pdf -i file2.pdf -i file3.pdf -o merged.pdf
```

**Flags:**

| Flag | Short | Description |
|------|-------|-------------|
| `--input` | `-i` | Input PDF files (can be specified multiple times) |
| `--output` | `-o` | Output merged PDF file (required) |

Files are merged in the order specified.

---

### split

Split a PDF file into individual pages.

```bash
gounidoc split -i document.pdf -d ./pages/
```

**Flags:**

| Flag | Short | Description |
|------|-------|-------------|
| `--input` | `-i` | Input PDF file path (required) |
| `--dir` | `-d` | Output directory (optional, defaults to input file directory) |

Output files are named `{original}_page_001.pdf`, `{original}_page_002.pdf`, etc.

---

### extract

Extract specific pages from a PDF file into a new PDF.

```bash
gounidoc extract -i document.pdf -o extracted.pdf --pages 1,3,5
```

**Flags:**

| Flag | Short | Description |
|------|-------|-------------|
| `--input` | `-i` | Input PDF file path (required) |
| `--output` | `-o` | Output PDF file path (required) |
| `--pages` | | Page numbers to extract (1-indexed, comma-separated) |

**Examples:**

```bash
# Extract first 3 pages
gounidoc extract -i document.pdf -o first3.pdf --pages 1,2,3

# Extract specific pages
gounidoc extract -i document.pdf -o selected.pdf --pages 1,5,10,15
```

---

### pagecount

Get the number of pages in a PDF file.

```bash
gounidoc pagecount -i document.pdf
```

**Flags:**

| Flag | Short | Description |
|------|-------|-------------|
| `--input` | `-i` | Input PDF file path (required) |

**Output:**

```json
{"file_path":"document.pdf","page_count":42}
```

---

## Document Information Commands

### metadata

Read metadata from an Office document (DOCX, XLSX, PPTX).

```bash
gounidoc metadata -f document.docx
```

**Flags:**

| Flag | Short | Description |
|------|-------|-------------|
| `--file` | `-f` | File path (required) |

**Supported formats:** `.docx`, `.xlsx`, `.pptx`

---

### xlsx-info

Get information about an XLSX workbook including metadata and sheet names.

```bash
gounidoc xlsx-info -f workbook.xlsx
```

**Flags:**

| Flag | Short | Description |
|------|-------|-------------|
| `--file` | `-f` | XLSX file path (required) |

**Output:**

```json
{
  "file_path": "workbook.xlsx",
  "metadata": {
    "author": "John Doe",
    "title": "Sales Report"
  },
  "sheet_count": 3,
  "sheet_names": ["Summary", "Q1 Data", "Q2 Data"]
}
```

---

### pptx-info

Get information about a PPTX presentation including metadata and slide count.

```bash
gounidoc pptx-info -f presentation.pptx
```

**Flags:**

| Flag | Short | Description |
|------|-------|-------------|
| `--file` | `-f` | PPTX file path (required) |

**Output:**

```json
{
  "file_path": "presentation.pptx",
  "metadata": {
    "author": "Jane Smith",
    "title": "Quarterly Review"
  },
  "slide_count": 25
}
```

---

## Server Commands

### serve

Start the MCP server using stdio transport.

```bash
gounidoc serve
```

This is the default behavior when running `gounidoc` without a command.

---

### version

Print version information.

```bash
gounidoc version
```

**Output:**

```
mcp-gounidoc v0.1.0
```

---

## Global Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--help` | `-h` | Show help for any command |
| `--output` | `-o` | Output format: `json` or `pretty` (default: json) |

**Note:** The global `--output` flag controls the output format for commands that return JSON data (like `metadata`, `pagecount`, `xlsx-info`, `pptx-info`).
