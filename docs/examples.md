# Examples

This page provides practical examples for common use cases.

## CLI Examples

### Document Conversion Workflow

Convert a PDF report to editable DOCX, make changes, then convert back to PDF:

```bash
# Convert PDF to DOCX
gounidoc pdf2docx -i report.pdf -o report.docx

# Edit report.docx in your word processor...

# Convert back to PDF
gounidoc docx2pdf -i report.docx -o report_final.pdf
```

### Batch Processing

Convert an entire folder of documents:

```bash
# Convert all PDFs in a directory
gounidoc batch-pdf2docx -i "./documents/*.pdf" -d ./converted/

# Convert with basic mode for speed
gounidoc batch-pdf2docx -i "*.pdf" -b -d ./output/

# Convert all DOCX files to PDF
gounidoc batch-docx2pdf -i "./drafts/*.docx" -d ./finals/
```

### PDF Manipulation

Combine multiple PDFs and extract specific pages:

```bash
# Merge chapter PDFs into a book
gounidoc merge -i chapter1.pdf -i chapter2.pdf -i chapter3.pdf -o book.pdf

# Split a PDF into individual pages
gounidoc split -i book.pdf -d ./pages/

# Extract the table of contents and first chapter
gounidoc extract -i book.pdf -o preview.pdf --pages 1,2,3,4,5

# Check how many pages before processing
gounidoc pagecount -i large_document.pdf
```

### Creating Thumbnails

Generate PNG thumbnails for PDF documents:

```bash
# Create thumbnail of first page
gounidoc pdf2png -i document.pdf -o thumbnail.png --page 1 --width 400

# Create high-res preview
gounidoc pdf2png -i document.pdf -o preview.png --page 1 --width 1920
```

### Document Inventory

List information about Office documents:

```bash
# Get metadata from a Word document
gounidoc metadata -f report.docx

# List all sheets in an Excel workbook
gounidoc xlsx-info -f data.xlsx

# Check slide count in a presentation
gounidoc pptx-info -f slides.pptx
```

## Library Examples

### Basic Conversion

```go
package main

import (
    "fmt"
    "github.com/grokify/gounidoc"
    "github.com/grokify/gounidoc/unipdfutil"
    "github.com/grokify/gounidoc/uniofficeutil"
)

func main() {
    // Initialize API key
    if err := gounidoc.SetMeteredKeyEnv(); err != nil {
        panic(err)
    }

    // Convert PDF to DOCX
    err := unipdfutil.ConvertPDFFileToDocxFile("input.pdf", "output.docx")
    if err != nil {
        panic(err)
    }
    fmt.Println("Converted PDF to DOCX")

    // Convert DOCX to PDF
    err = uniofficeutil.ConvertDOCXFileToPDFFile("document.docx", "document.pdf")
    if err != nil {
        panic(err)
    }
    fmt.Println("Converted DOCX to PDF")
}
```

### PDF Merge Tool

```go
package main

import (
    "fmt"
    "os"
    "github.com/grokify/gounidoc"
    "github.com/grokify/gounidoc/unipdfutil"
)

func main() {
    if len(os.Args) < 3 {
        fmt.Println("Usage: merge output.pdf input1.pdf input2.pdf ...")
        os.Exit(1)
    }

    gounidoc.SetMeteredKeyEnv()

    outputPath := os.Args[1]
    inputPaths := os.Args[2:]

    if err := unipdfutil.MergePDFFiles(inputPaths, outputPath); err != nil {
        fmt.Printf("Error: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("Merged %d files into %s\n", len(inputPaths), outputPath)
}
```

### Batch Converter with Progress

```go
package main

import (
    "fmt"
    "os"
    "path/filepath"
    "github.com/grokify/gounidoc"
    "github.com/grokify/gounidoc/unipdfutil"
    "github.com/grokify/mogo/fmt/progress"
)

func main() {
    gounidoc.SetMeteredKeyEnv()

    // Find all PDFs
    inputPaths, _ := unipdfutil.ExpandGlobPatterns([]string{"*.pdf"})
    inputPaths = unipdfutil.FilterByExtension(inputPaths, ".pdf")

    if len(inputPaths) == 0 {
        fmt.Println("No PDF files found")
        return
    }

    // Create progress renderer
    renderer := progress.NewSingleStageRenderer(os.Stdout).
        WithBarWidth(30).
        WithTextWidth(40)

    // Convert with progress
    opts := unipdfutil.DefaultConversionOptions()
    progressFunc := func(current, total int, inputPath string) {
        renderer.Update(current, total, filepath.Base(inputPath))
    }

    result := unipdfutil.BatchConvertPDFToDocxWithProgress(
        inputPaths, "./output", opts, progressFunc)

    renderer.Done("")
    fmt.Printf("Converted: %d, Failed: %d\n", result.Succeeded, result.Failed)
}
```

### Document Inspector

```go
package main

import (
    "encoding/json"
    "fmt"
    "os"
    "path/filepath"
    "strings"
    "github.com/grokify/gounidoc"
    "github.com/grokify/gounidoc/uniofficeutil"
    "github.com/grokify/gounidoc/unipdfutil"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: inspector <file>")
        os.Exit(1)
    }

    gounidoc.SetMeteredKeyEnv()

    filePath := os.Args[1]
    ext := strings.ToLower(filepath.Ext(filePath))

    var info any
    var err error

    switch ext {
    case ".pdf":
        count, e := unipdfutil.GetPDFPageCount(filePath)
        info = map[string]any{"file": filePath, "type": "pdf", "pages": count}
        err = e
    case ".xlsx":
        info, err = uniofficeutil.GetXlsxInfo(filePath)
    case ".pptx":
        info, err = uniofficeutil.GetPptxInfo(filePath)
    case ".docx":
        info, err = uniofficeutil.ReadFileMetadata(filePath)
    default:
        fmt.Printf("Unsupported file type: %s\n", ext)
        os.Exit(1)
    }

    if err != nil {
        fmt.Printf("Error: %v\n", err)
        os.Exit(1)
    }

    output, _ := json.MarshalIndent(info, "", "  ")
    fmt.Println(string(output))
}
```

### Custom MCP Server

```go
package main

import (
    "context"
    "log/slog"
    "os"
    "github.com/grokify/gounidoc"
    "github.com/grokify/gounidoc/skills/unidoc"
    "github.com/modelcontextprotocol/go-sdk/mcp"
    runtime "github.com/plexusone/omniskill/mcp/server"
)

func main() {
    ctx := context.Background()

    // Initialize UniDoc
    if err := gounidoc.SetMeteredKeyEnv(); err != nil {
        slog.Error("Failed to initialize UniDoc", "error", err)
        os.Exit(1)
    }

    // Create runtime
    rt := runtime.New(&mcp.Implementation{
        Name:    "my-document-server",
        Version: "v1.0.0",
    }, nil)

    // Create skill with custom logger
    logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
        Level: slog.LevelDebug,
    }))
    skill := unidoc.NewWithLogger(logger)

    if err := skill.Init(ctx); err != nil {
        slog.Error("Failed to initialize skill", "error", err)
        os.Exit(1)
    }
    defer skill.Close()

    // Register skill
    rt.RegisterSkill(skill)

    // Serve
    slog.Info("Starting MCP server")
    if err := rt.ServeStdio(ctx); err != nil {
        slog.Error("Server error", "error", err)
        os.Exit(1)
    }
}
```

### Using Skill Typed Methods

```go
package main

import (
    "context"
    "fmt"
    "github.com/grokify/gounidoc"
    "github.com/grokify/gounidoc/skills/unidoc"
    "github.com/grokify/gounidoc/unipdfutil"
)

func main() {
    ctx := context.Background()
    gounidoc.SetMeteredKeyEnv()

    // Create and initialize skill
    skill := unidoc.New()
    if err := skill.Init(ctx); err != nil {
        panic(err)
    }
    defer skill.Close()

    // Get PDF info
    count, err := skill.GetPDFPageCount(ctx, "document.pdf")
    if err != nil {
        panic(err)
    }
    fmt.Printf("Document has %d pages\n", count)

    // Convert with options
    opts := unipdfutil.ConversionOptions{
        ExtractTables: true,
        ExtractImages: false, // Skip images for faster conversion
        DetectLayout:  true,
    }
    output, err := skill.ConvertPDFToDocx(ctx, "document.pdf", "", opts)
    if err != nil {
        panic(err)
    }
    fmt.Printf("Converted to: %s\n", output)

    // Get XLSX info
    xlsxInfo, err := skill.GetXlsxInfo(ctx, "data.xlsx")
    if err != nil {
        panic(err)
    }
    fmt.Printf("Workbook has %d sheets: %v\n", xlsxInfo.SheetCount, xlsxInfo.SheetNames)
}
```

## Shell Script Examples

### Batch Convert Directory

```bash
#!/bin/bash
# convert_all.sh - Convert all PDFs in a directory

INPUT_DIR="${1:-.}"
OUTPUT_DIR="${2:-./converted}"

mkdir -p "$OUTPUT_DIR"

for pdf in "$INPUT_DIR"/*.pdf; do
    if [ -f "$pdf" ]; then
        name=$(basename "$pdf" .pdf)
        echo "Converting: $name"
        gounidoc pdf2docx -i "$pdf" -o "$OUTPUT_DIR/${name}.docx"
    fi
done

echo "Done!"
```

### Create PDF Book from Chapters

```bash
#!/bin/bash
# make_book.sh - Merge chapter PDFs into a book

OUTPUT="${1:-book.pdf}"
shift

if [ $# -eq 0 ]; then
    echo "Usage: make_book.sh output.pdf chapter1.pdf chapter2.pdf ..."
    exit 1
fi

# Build the merge command
CMD="gounidoc merge -o $OUTPUT"
for chapter in "$@"; do
    CMD="$CMD -i $chapter"
done

echo "Merging into $OUTPUT..."
eval $CMD

# Show result
gounidoc pagecount -i "$OUTPUT"
```

### Generate Thumbnails

```bash
#!/bin/bash
# thumbnails.sh - Generate thumbnails for all PDFs

WIDTH="${1:-400}"
OUTPUT_DIR="${2:-./thumbnails}"

mkdir -p "$OUTPUT_DIR"

for pdf in *.pdf; do
    if [ -f "$pdf" ]; then
        name=$(basename "$pdf" .pdf)
        echo "Creating thumbnail: $name"
        gounidoc pdf2png -i "$pdf" -o "$OUTPUT_DIR/${name}.png" --page 1 --width "$WIDTH"
    fi
done

echo "Thumbnails created in $OUTPUT_DIR"
```
