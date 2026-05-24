# MCP Server

GoUniDoc includes an MCP (Model Context Protocol) server that exposes document operations as tools for AI agents.

## Overview

The MCP server provides 12 tools for document conversion and processing, enabling AI assistants like Claude to work with PDF and Office documents.

## Starting the Server

```bash
# Default behavior (start MCP server)
gounidoc

# Or explicitly
gounidoc serve
```

The server uses stdio transport for communication with MCP clients.

## Available Tools

### Document Conversion Tools

#### pdf2docx

Convert a PDF file to DOCX format with optional table, image, and layout extraction.

**Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `input_path` | string | Yes | Path to the input PDF file |
| `output_path` | string | No | Path for the output DOCX file |
| `extract_tables` | boolean | No | Extract tables (default: true) |
| `extract_images` | boolean | No | Detect images (default: true) |
| `detect_layout` | boolean | No | Detect headings (default: true) |

---

#### docx2pdf

Convert a DOCX file to PDF format.

**Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `input_path` | string | Yes | Path to the input DOCX file |
| `output_path` | string | No | Path for the output PDF file |

---

#### pdf2png

Convert a PDF page to PNG image format.

**Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `input_path` | string | Yes | Path to the input PDF file |
| `output_path` | string | Yes | Path for the output PNG file |
| `page_number` | integer | No | Page number (1-indexed, default: 1) |
| `width` | integer | No | Output width in pixels (default: 1200) |

---

#### batch_pdf2docx

Convert multiple PDF files to DOCX format in batch.

**Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `input_paths` | array | Yes | List of input PDF file paths |
| `output_dir` | string | No | Output directory for converted files |
| `extract_tables` | boolean | No | Extract tables (default: true) |
| `extract_images` | boolean | No | Detect images (default: true) |
| `detect_layout` | boolean | No | Detect headings (default: true) |

---

#### batch_docx2pdf

Convert multiple DOCX files to PDF format in batch.

**Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `input_paths` | array | Yes | List of input DOCX file paths |
| `output_dir` | string | No | Output directory for converted files |

---

### PDF Operations Tools

#### pdf_merge

Merge multiple PDF files into a single output file.

**Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `input_paths` | array | Yes | List of input PDF file paths (in order) |
| `output_path` | string | Yes | Path for the merged output PDF |

---

#### pdf_split

Split a PDF file into individual pages.

**Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `input_path` | string | Yes | Path to the input PDF file |
| `output_dir` | string | No | Directory for output files |

---

#### pdf_extract_pages

Extract specific pages from a PDF file into a new PDF.

**Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `input_path` | string | Yes | Path to the input PDF file |
| `output_path` | string | Yes | Path for the output PDF file |
| `pages` | array | Yes | List of page numbers to extract (1-indexed) |

---

#### pdf_page_count

Get the number of pages in a PDF file.

**Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `input_path` | string | Yes | Path to the PDF file |

**Returns:**

```json
{
  "success": true,
  "input_path": "/path/to/document.pdf",
  "page_count": 42
}
```

---

### Document Information Tools

#### metadata

Read metadata from an Office document (DOCX, XLSX, PPTX).

**Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `file_path` | string | Yes | Path to the Office document |

---

#### xlsx_info

Get information about an XLSX workbook including metadata and sheet names.

**Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `file_path` | string | Yes | Path to the XLSX file |

**Returns:**

```json
{
  "success": true,
  "file_path": "/path/to/workbook.xlsx",
  "metadata": {
    "author": "John Doe",
    "title": "Sales Report"
  },
  "sheet_count": 3,
  "sheet_names": ["Summary", "Q1 Data", "Q2 Data"]
}
```

---

#### pptx_info

Get information about a PPTX presentation including metadata and slide count.

**Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `file_path` | string | Yes | Path to the PPTX file |

**Returns:**

```json
{
  "success": true,
  "file_path": "/path/to/presentation.pptx",
  "metadata": {
    "author": "Jane Smith",
    "title": "Quarterly Review"
  },
  "slide_count": 25
}
```

---

## Client Integration

### Claude Desktop

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

!!! tip "Finding the Binary Path"
    After installing with `go install`, find the binary path:
    ```bash
    which gounidoc
    # or
    go env GOPATH
    # Binary is at $GOPATH/bin/gounidoc
    ```

### Claude Code

Add to your Claude Code MCP configuration:

```json
{
  "mcpServers": {
    "gounidoc": {
      "command": "gounidoc",
      "args": ["serve"],
      "env": {
        "UNIDOC_KEY": "your-metered-api-key"
      }
    }
  }
}
```

### Custom MCP Client

The server uses stdio transport, so any MCP client that supports stdio can connect:

```bash
# The server reads JSON-RPC from stdin and writes to stdout
gounidoc serve
```

## Omniskill Integration

GoUniDoc implements the [omniskill](https://github.com/plexusone/omniskill) Skill interface, making it easy to compose with other MCP skills:

```go
package main

import (
    "context"
    "github.com/grokify/gounidoc"
    "github.com/grokify/gounidoc/skills/unidoc"
    "github.com/modelcontextprotocol/go-sdk/mcp"
    runtime "github.com/plexusone/omniskill/mcp/server"
)

func main() {
    ctx := context.Background()

    // Initialize UniDoc
    gounidoc.SetMeteredKeyEnv()

    // Create runtime
    rt := runtime.New(&mcp.Implementation{
        Name:    "my-document-server",
        Version: "v1.0.0",
    }, nil)

    // Create and register UniDoc skill
    unidocSkill := unidoc.New()
    unidocSkill.Init(ctx)
    defer unidocSkill.Close()
    rt.RegisterSkill(unidocSkill)

    // Add other skills here...
    // rt.RegisterSkill(otherSkill)

    // Serve
    rt.ServeStdio(ctx)
}
```

## Logging

The skill supports structured logging via `slog`. To use a custom logger:

```go
import (
    "log/slog"
    "os"
    "github.com/grokify/gounidoc/skills/unidoc"
)

// Create skill with JSON logger
logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
    Level: slog.LevelDebug,
}))
skill := unidoc.NewWithLogger(logger)
```

Log output includes:

- Operation start/completion
- Input/output paths
- Timing information
- Error details
