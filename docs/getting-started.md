# Getting Started

This guide will help you install and configure GoUniDoc for document conversion and processing.

## Installation

### Using Go Install

The easiest way to install GoUniDoc:

```bash
go install github.com/grokify/gounidoc/cmd/gounidoc@latest
```

### Building from Source

Clone and build from source:

```bash
git clone https://github.com/grokify/gounidoc.git
cd gounidoc
go build ./cmd/gounidoc/
```

### Verify Installation

```bash
gounidoc version
# Output: mcp-gounidoc v0.1.0
```

## Configuration

### API Key Setup

GoUniDoc requires a UniDoc API key for operation. Get a free trial key at [unidoc.io](https://unidoc.io/).

Set the environment variable:

=== "Bash/Zsh"

    ```bash
    export UNIDOC_KEY="your-metered-api-key"
    ```

=== "Fish"

    ```fish
    set -gx UNIDOC_KEY "your-metered-api-key"
    ```

=== "PowerShell"

    ```powershell
    $env:UNIDOC_KEY = "your-metered-api-key"
    ```

### Environment Variables

| Variable | Description |
|----------|-------------|
| `UNIDOC_KEY` | API key used for both UniPDF and UniOffice |
| `UNIDOC_KEY_PDF` | API key specifically for UniPDF (PDF operations) |
| `UNIDOC_KEY_OFFICE` | API key specifically for UniOffice (Office operations) |

!!! tip "API Key Priority"
    If `UNIDOC_KEY_PDF` or `UNIDOC_KEY_OFFICE` are set, they take priority over `UNIDOC_KEY` for their respective operations.

## Quick Start Examples

### Convert PDF to DOCX

```bash
gounidoc pdf2docx -i document.pdf -o document.docx
```

### Convert DOCX to PDF

```bash
gounidoc docx2pdf -i document.docx -o document.pdf
```

### Merge Multiple PDFs

```bash
gounidoc merge -i file1.pdf -i file2.pdf -i file3.pdf -o merged.pdf
```

### Batch Convert Files

```bash
# Convert all PDFs in current directory to DOCX
gounidoc batch-pdf2docx -i "*.pdf" -d ./output/
```

### Start MCP Server

```bash
gounidoc serve
```

## Using as a Library

### Basic Usage

```go
package main

import (
    "github.com/grokify/gounidoc"
    "github.com/grokify/gounidoc/unipdfutil"
)

func main() {
    // Initialize UniDoc with API key from environment
    if err := gounidoc.SetMeteredKeyEnv(); err != nil {
        panic(err)
    }

    // Convert PDF to DOCX
    err := unipdfutil.ConvertPDFFileToDocxFile("input.pdf", "output.docx")
    if err != nil {
        panic(err)
    }
}
```

### Using the Skill Interface

```go
package main

import (
    "context"
    "github.com/grokify/gounidoc"
    "github.com/grokify/gounidoc/skills/unidoc"
    "github.com/grokify/gounidoc/unipdfutil"
)

func main() {
    ctx := context.Background()

    // Initialize UniDoc
    gounidoc.SetMeteredKeyEnv()

    // Create and initialize skill
    skill := unidoc.New()
    skill.Init(ctx)
    defer skill.Close()

    // Use typed methods
    opts := unipdfutil.DefaultConversionOptions()
    output, err := skill.ConvertPDFToDocx(ctx, "input.pdf", "output.docx", opts)
    if err != nil {
        panic(err)
    }
    println("Created:", output)
}
```

## Next Steps

- [CLI Reference](cli-reference.md) - Learn all 16 CLI commands
- [MCP Server](mcp-server.md) - Set up AI agent integration
- [Library Reference](library-reference.md) - Go API documentation
- [Examples](examples.md) - More usage examples
