# GoUniDoc

[![Go CI][go-ci-svg]][go-ci-url]
[![Go Lint][go-lint-svg]][go-lint-url]
[![Go SAST][go-sast-svg]][go-sast-url]
[![Go Report Card][goreport-svg]][goreport-url]
[![Docs][docs-godoc-svg]][docs-godoc-url]
[![Docs][docs-mkdoc-svg]][docs-mkdoc-url]
[![Visualization][viz-svg]][viz-url]
[![License][license-svg]][license-url]

 [go-ci-svg]: https://github.com/grokify/gounidoc/actions/workflows/go-ci.yaml/badge.svg?branch=main
 [go-ci-url]: https://github.com/grokify/gounidoc/actions/workflows/go-ci.yaml
 [go-lint-svg]: https://github.com/grokify/gounidoc/actions/workflows/go-lint.yaml/badge.svg?branch=main
 [go-lint-url]: https://github.com/grokify/gounidoc/actions/workflows/go-lint.yaml
 [go-sast-svg]: https://github.com/grokify/gounidoc/actions/workflows/go-sast-codeql.yaml/badge.svg?branch=main
 [go-sast-url]: https://github.com/grokify/gounidoc/actions/workflows/go-sast-codeql.yaml
 [goreport-svg]: https://goreportcard.com/badge/github.com/grokify/gounidoc
 [goreport-url]: https://goreportcard.com/report/github.com/grokify/gounidoc
 [docs-godoc-svg]: https://pkg.go.dev/badge/github.com/grokify/gounidoc
 [docs-godoc-url]: https://pkg.go.dev/github.com/grokify/gounidoc
 [docs-mkdoc-svg]: https://img.shields.io/badge/Go-dev%20guide-blue.svg
 [docs-mkdoc-url]: https://grokify.github.io/gounidoc
 [viz-svg]: https://img.shields.io/badge/visualization-Go-blue.svg
 [viz-url]: https://mango-dune-07a8b7110.1.azurestaticapps.net/?repo=grokify%2Fgounidoc
 [loc-svg]: https://tokei.rs/b1/github/grokify/gounidoc
 [repo-url]: https://github.com/grokify/gounidoc
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-url]: https://github.com/grokify/gounidoc/blob/main/LICENSE

![](logo_unidoc.png)

`gounidoc` is a Go module providing wrapper utilities for [UniDoc](https://unidoc.io/) services, enabling PDF and Office document operations including format conversion. It includes both a unified CLI tool and an MCP (Model Context Protocol) server for AI agent integration.

## Features

### Document Conversion

- **PDF to DOCX** - Convert PDF to Word with table extraction, image detection, and layout-aware formatting
- **DOCX to PDF** - Convert Word documents to PDF format
- **PDF to PNG** - Render PDF pages as PNG images
- **Batch Conversion** - Convert multiple files with progress reporting

### PDF Operations

- **Merge** - Combine multiple PDF files into one
- **Split** - Split PDF into individual pages
- **Extract** - Extract specific pages from a PDF
- **Page Count** - Get the number of pages in a PDF

### Office Document Support

- **Metadata** - Read metadata from DOCX, XLSX, and PPTX files
- **XLSX Info** - Get workbook information including sheet names
- **PPTX Info** - Get presentation information including slide count

### Integration

- **MCP Server** - Expose all features as MCP tools for AI agents (12 tools)
- **Omniskill** - Composable skill interface for building custom MCP servers
- **Structured Logging** - slog-based logging for all operations

## Installation

```bash
go install github.com/grokify/gounidoc/cmd/gounidoc@latest
```

Or build from source:

```bash
git clone https://github.com/grokify/gounidoc.git
cd gounidoc
go build ./cmd/gounidoc/
```

## Environment Variables

UniDoc requires a metered API key for operation. Get a free trial key at [unidoc.io](https://unidoc.io/).

| Variable | Description |
|----------|-------------|
| `UNIDOC_KEY` | API key used for both UniPDF and UniOffice (if specific keys not set) |
| `UNIDOC_KEY_PDF` | API key specifically for UniPDF (PDF operations) |
| `UNIDOC_KEY_OFFICE` | API key specifically for UniOffice (Office document operations) |

```bash
export UNIDOC_KEY="your-metered-api-key"
```

## CLI Commands

The `gounidoc` CLI provides 16 commands:

### Document Conversion

```bash
# Convert PDF to DOCX
gounidoc pdf2docx -i input.pdf -o output.docx

# Convert with options (disable features for faster conversion)
gounidoc pdf2docx -i input.pdf --no-tables --no-images --no-layout

# Basic mode (text only, fastest)
gounidoc pdf2docx -i input.pdf -b

# Convert DOCX to PDF
gounidoc docx2pdf -i document.docx -o output.pdf

# Convert PDF page to PNG
gounidoc pdf2png -i document.pdf -o page1.png --page 1 --width 1200
```

### Batch Conversion

```bash
# Batch convert PDF files to DOCX (with progress bar)
gounidoc batch-pdf2docx -i "*.pdf" -d ./output/
gounidoc batch-pdf2docx -i file1.pdf -i file2.pdf -d ./output/

# Batch convert DOCX files to PDF (with progress bar)
gounidoc batch-docx2pdf -i "*.docx" -d ./output/
```

### PDF Operations

```bash
# Merge multiple PDFs
gounidoc merge -i file1.pdf -i file2.pdf -i file3.pdf -o merged.pdf

# Split PDF into individual pages
gounidoc split -i document.pdf -d ./pages/

# Extract specific pages
gounidoc extract -i document.pdf -o extracted.pdf --pages 1,3,5

# Get page count
gounidoc pagecount -i document.pdf
```

### Document Information

```bash
# Read Office document metadata (DOCX, XLSX, PPTX)
gounidoc metadata -f document.docx

# Get XLSX workbook info (sheets, metadata)
gounidoc xlsx-info -f workbook.xlsx

# Get PPTX presentation info (slides, metadata)
gounidoc pptx-info -f presentation.pptx
```

### MCP Server

```bash
# Start MCP server (default behavior)
gounidoc

# Or explicitly
gounidoc serve

# Show version
gounidoc version
```

## MCP Server

gounidoc includes an MCP (Model Context Protocol) server for AI agent integration. The server exposes 12 tools:

| Tool | Description |
|------|-------------|
| `pdf2docx` | Convert PDF to DOCX with table, image, and layout extraction |
| `docx2pdf` | Convert DOCX to PDF |
| `pdf2png` | Convert PDF page to PNG image |
| `metadata` | Read Office document metadata (DOCX, XLSX, PPTX) |
| `pdf_merge` | Merge multiple PDF files into one |
| `pdf_split` | Split PDF into individual pages |
| `pdf_extract_pages` | Extract specific pages from a PDF |
| `pdf_page_count` | Get the number of pages in a PDF |
| `batch_pdf2docx` | Batch convert PDF files to DOCX |
| `batch_docx2pdf` | Batch convert DOCX files to PDF |
| `xlsx_info` | Get XLSX workbook info (sheets, metadata) |
| `pptx_info` | Get PPTX presentation info (slides, metadata) |

### Claude Desktop Integration

Add to your `claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "gounidoc": {
      "command": "/path/to/gounidoc",
      "args": ["serve"],
      "env": {
        "UNIDOC_KEY": "your-metered-api-key"
      }
    }
  }
}
```

## Library Usage

### PDF to DOCX

```go
import "github.com/grokify/gounidoc/unipdfutil"

// Default conversion with all features enabled
err := unipdfutil.ConvertPDFFileToDocxFile("input.pdf", "output.docx")

// Custom options
opts := unipdfutil.ConversionOptions{
    ExtractTables: true,  // Extract and convert tables
    ExtractImages: true,  // Detect images (placeholder text)
    DetectLayout:  true,  // Detect headings and structure
}
err := unipdfutil.ConvertPDFFileToDocxFileWithOptions("input.pdf", "output.docx", opts)
```

### PDF Operations

```go
import "github.com/grokify/gounidoc/unipdfutil"

// Merge PDFs
err := unipdfutil.MergePDFFiles([]string{"file1.pdf", "file2.pdf"}, "merged.pdf")

// Split PDF into pages
outputPaths, err := unipdfutil.SplitPDFFile("input.pdf", "./pages/")

// Extract specific pages
err := unipdfutil.ExtractPDFPages("input.pdf", "extracted.pdf", []int{1, 3, 5})

// Get page count
count, err := unipdfutil.GetPDFPageCount("input.pdf")
```

### Batch Conversion

```go
import "github.com/grokify/gounidoc/unipdfutil"

// Batch convert with progress callback
opts := unipdfutil.DefaultConversionOptions()
progress := func(current, total int, inputPath string) {
    fmt.Printf("Processing %d/%d: %s\n", current, total, inputPath)
}
result := unipdfutil.BatchConvertPDFToDocxWithProgress(inputPaths, outputDir, opts, progress)

fmt.Printf("Converted: %d, Failed: %d\n", result.Succeeded, result.Failed)
```

### Office Document Info

```go
import "github.com/grokify/gounidoc/uniofficeutil"

// Read metadata from any Office document
metadata, err := uniofficeutil.ReadFileMetadata("document.docx")

// Get XLSX workbook info
info, err := uniofficeutil.GetXlsxInfo("workbook.xlsx")
fmt.Printf("Sheets: %v\n", info.SheetNames)

// Get PPTX presentation info
info, err := uniofficeutil.GetPptxInfo("presentation.pptx")
fmt.Printf("Slides: %d\n", info.SlideCount)
```

### Using as Omniskill

gounidoc implements the [omniskill](https://github.com/plexusone/omniskill) Skill interface:

```go
import (
    "github.com/grokify/gounidoc/skills/unidoc"
    "github.com/modelcontextprotocol/go-sdk/mcp"
    runtime "github.com/plexusone/omniskill/mcp/server"
)

// Create runtime
rt := runtime.New(&mcp.Implementation{
    Name:    "my-server",
    Version: "v1.0.0",
}, nil)

// Initialize and register skill
skill := unidoc.New()
skill.Init(ctx)
rt.RegisterSkill(skill)

// Serve
rt.ServeStdio(ctx)
```

### Custom Logger

```go
import (
    "log/slog"
    "github.com/grokify/gounidoc/skills/unidoc"
)

// Create skill with custom logger
logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
skill := unidoc.NewWithLogger(logger)

// Or set logger after creation
skill := unidoc.New()
skill.SetLogger(logger)
```

## Documentation

Full documentation is available at:

- [MkDocs Site](https://grokify.github.io/gounidoc/) (coming soon)
- [GoDoc](https://pkg.go.dev/github.com/grokify/gounidoc)

## References

1. [UniOffice Documentation](https://unidoc.io/unioffice/)
1. [UniOffice API Reference](https://apidocs.unidoc.io/unioffice/latest/)
1. [UniPDF Documentation](https://unidoc.io/unipdf/)
1. [UniPDF API Reference](https://apidocs.unidoc.io/unipdf/latest/)
1. [Get UniDoc API Key](https://unidoc.io/)
1. [Model Context Protocol](https://modelcontextprotocol.io/)
1. [Omniskill](https://github.com/plexusone/omniskill)
