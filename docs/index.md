# GoUniDoc

GoUniDoc is a Go module providing wrapper utilities for [UniDoc](https://unidoc.io/) services, enabling PDF and Office document operations including format conversion. It includes both a unified CLI tool and an MCP (Model Context Protocol) server for AI agent integration.

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

## Quick Start

```bash
# Install
go install github.com/grokify/gounidoc/cmd/gounidoc@latest

# Set API key
export UNIDOC_KEY="your-metered-api-key"

# Convert PDF to DOCX
gounidoc pdf2docx -i input.pdf -o output.docx

# Start MCP server
gounidoc serve
```

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                         gounidoc                            │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐     │
│  │   CLI       │    │ MCP Server  │    │  Library    │     │
│  │ (16 cmds)   │    │ (12 tools)  │    │   (Go API)  │     │
│  └──────┬──────┘    └──────┬──────┘    └──────┬──────┘     │
│         │                  │                  │             │
│         └────────────┬─────┴──────────────────┘             │
│                      │                                      │
│              ┌───────▼───────┐                              │
│              │  Skill Layer  │                              │
│              │ (Typed Methods)│                              │
│              └───────┬───────┘                              │
│                      │                                      │
│         ┌────────────┴────────────┐                         │
│         │                         │                         │
│  ┌──────▼──────┐          ┌──────▼──────┐                   │
│  │  unipdfutil │          │uniofficeutil│                   │
│  │ (PDF ops)   │          │(Office ops) │                   │
│  └──────┬──────┘          └──────┬──────┘                   │
│         │                        │                          │
└─────────┴────────────────────────┴──────────────────────────┘
          │                        │
   ┌──────▼──────┐          ┌──────▼──────┐
   │   UniPDF    │          │  UniOffice  │
   └─────────────┘          └─────────────┘
```

## Requirements

- Go 1.21 or later
- UniDoc API key (get a free trial at [unidoc.io](https://unidoc.io/))

## License

MIT License - see [LICENSE](https://github.com/grokify/gounidoc/blob/master/LICENSE)
