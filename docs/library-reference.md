# Library Reference

GoUniDoc can be used as a Go library for PDF and Office document operations.

## Package Overview

| Package | Description |
|---------|-------------|
| `gounidoc` | Root package for API key initialization |
| `unipdfutil` | PDF operations (conversion, merge, split, etc.) |
| `uniofficeutil` | Office document operations (DOCX, XLSX, PPTX) |
| `skills/unidoc` | Omniskill interface for MCP integration |

## Initialization

### Setting API Key

```go
import "github.com/grokify/gounidoc"

// Load from environment variables (UNIDOC_KEY, UNIDOC_KEY_PDF, UNIDOC_KEY_OFFICE)
err := gounidoc.SetMeteredKeyEnv()

// Or set directly
err := gounidoc.SetMeteredKey("office-api-key", "pdf-api-key")
```

## unipdfutil Package

### PDF to DOCX Conversion

```go
import "github.com/grokify/gounidoc/unipdfutil"

// Simple conversion with default options
err := unipdfutil.ConvertPDFFileToDocxFile("input.pdf", "output.docx")

// Conversion with custom options
opts := unipdfutil.ConversionOptions{
    ExtractTables: true,   // Extract and convert tables
    ExtractImages: true,   // Detect images (placeholder text)
    DetectLayout:  true,   // Detect headings and structure
}
err := unipdfutil.ConvertPDFFileToDocxFileWithOptions("input.pdf", "output.docx", opts)

// Get default options
opts := unipdfutil.DefaultConversionOptions()
```

### PDF to PNG Conversion

```go
// Convert page 1 to PNG with width of 1200 pixels
err := unipdfutil.ConvertPDFFilePageToPNGFile("input.pdf", "page1.png", 1, 1200)
```

### PDF Merge

```go
// Merge multiple PDFs into one
err := unipdfutil.MergePDFFiles([]string{"file1.pdf", "file2.pdf", "file3.pdf"}, "merged.pdf")
```

### PDF Split

```go
// Split PDF into individual pages
// Returns list of created file paths
outputPaths, err := unipdfutil.SplitPDFFile("input.pdf", "./pages/")

// Output files: input_page_001.pdf, input_page_002.pdf, etc.
```

### PDF Page Extraction

```go
// Extract specific pages
err := unipdfutil.ExtractPDFPages("input.pdf", "extracted.pdf", []int{1, 3, 5, 7})
```

### PDF Page Count

```go
count, err := unipdfutil.GetPDFPageCount("document.pdf")
fmt.Printf("Pages: %d\n", count)
```

### Batch Conversion

```go
// Batch convert PDFs to DOCX
opts := unipdfutil.DefaultConversionOptions()
result := unipdfutil.BatchConvertPDFToDocx(inputPaths, outputDir, opts)

fmt.Printf("Total: %d, Succeeded: %d, Failed: %d\n",
    result.TotalFiles, result.Succeeded, result.Failed)

// With progress callback
progress := func(current, total int, inputPath string) {
    fmt.Printf("Processing %d/%d: %s\n", current, total, inputPath)
}
result := unipdfutil.BatchConvertPDFToDocxWithProgress(inputPaths, outputDir, opts, progress)
```

### Glob Pattern Expansion

```go
// Expand glob patterns to file list
files, err := unipdfutil.ExpandGlobPatterns([]string{"*.pdf", "docs/*.pdf"})

// Filter by extension
pdfFiles := unipdfutil.FilterByExtension(files, ".pdf")
```

## uniofficeutil Package

### DOCX to PDF Conversion

```go
import "github.com/grokify/gounidoc/uniofficeutil"

err := uniofficeutil.ConvertDOCXFileToPDFFile("input.docx", "output.pdf")
```

### Reading Metadata

```go
// Read metadata from any Office document (DOCX, XLSX, PPTX)
metadata, err := uniofficeutil.ReadFileMetadata("document.docx")

fmt.Printf("Title: %s\n", metadata.Title)
fmt.Printf("Author: %s\n", metadata.Author)
fmt.Printf("Created: %s\n", metadata.Created)
fmt.Printf("Modified: %s\n", metadata.Modified)
```

### Metadata Structure

```go
type Metadata struct {
    Author         string    `json:"author,omitempty"`
    Category       string    `json:"category,omitempty"`
    ContentStatus  string    `json:"content_status,omitempty"`
    Created        time.Time `json:"created,omitempty"`
    Description    string    `json:"description,omitempty"`
    LastModifiedBy string    `json:"last_modified_by,omitempty"`
    Modified       time.Time `json:"modified,omitempty"`
    Title          string    `json:"title,omitempty"`
    FileType       string    `json:"file_type,omitempty"`
}
```

### XLSX Information

```go
info, err := uniofficeutil.GetXlsxInfo("workbook.xlsx")

fmt.Printf("Sheet count: %d\n", info.SheetCount)
fmt.Printf("Sheets: %v\n", info.SheetNames)
fmt.Printf("Author: %s\n", info.Metadata.Author)
```

### PPTX Information

```go
info, err := uniofficeutil.GetPptxInfo("presentation.pptx")

fmt.Printf("Slide count: %d\n", info.SlideCount)
fmt.Printf("Title: %s\n", info.Metadata.Title)
```

## skills/unidoc Package

The skill package provides typed methods with logging support.

### Creating a Skill

```go
import "github.com/grokify/gounidoc/skills/unidoc"

// Default logger
skill := unidoc.New()

// Custom logger
logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
skill := unidoc.NewWithLogger(logger)

// Initialize
err := skill.Init(ctx)
defer skill.Close()
```

### Typed Methods

```go
// PDF to DOCX
opts := unipdfutil.DefaultConversionOptions()
outputPath, err := skill.ConvertPDFToDocx(ctx, "input.pdf", "output.docx", opts)

// DOCX to PDF
outputPath, err := skill.ConvertDocxToPDF(ctx, "input.docx", "output.pdf")

// PDF to PNG
err := skill.ConvertPDFPageToPNG(ctx, "input.pdf", "page1.png", 1, 1200)

// Read metadata
metadata, err := skill.ReadMetadata(ctx, "document.docx")

// Merge PDFs
err := skill.MergePDFs(ctx, []string{"file1.pdf", "file2.pdf"}, "merged.pdf")

// Split PDF
outputPaths, err := skill.SplitPDF(ctx, "input.pdf", "./pages/")

// Extract pages
err := skill.ExtractPDFPages(ctx, "input.pdf", "extracted.pdf", []int{1, 3, 5})

// Get page count
count, err := skill.GetPDFPageCount(ctx, "document.pdf")

// Batch conversion with progress
progress := func(current, total int, inputPath string) {
    fmt.Printf("Processing %d/%d\n", current, total)
}
result := skill.BatchConvertPDFToDocxWithProgress(ctx, inputPaths, outputDir, opts, progress)

// XLSX info
info, err := skill.GetXlsxInfo(ctx, "workbook.xlsx")

// PPTX info
info, err := skill.GetPptxInfo(ctx, "presentation.pptx")
```

### Available Tools

Access MCP tools programmatically:

```go
tools := skill.Tools()
for _, tool := range tools {
    fmt.Printf("Tool: %s - %s\n", tool.Name(), tool.Description())
}
```

### Skill Interface

The skill implements the omniskill `Skill` interface:

```go
type Skill interface {
    Name() string
    Description() string
    Tools() []Tool
    Init(ctx context.Context) error
    Close() error
}
```

## Error Handling

All functions return errors that should be checked:

```go
output, err := skill.ConvertPDFToDocx(ctx, input, output, opts)
if err != nil {
    // Handle error
    log.Printf("Conversion failed: %v", err)
    return err
}
```

Batch operations return results that include per-file errors:

```go
result := skill.BatchConvertPDFToDocx(ctx, inputs, outputDir, opts)
for _, r := range result.Results {
    if !r.Success {
        log.Printf("Failed: %s - %s", r.InputPath, r.Error)
    }
}
```

## GoDoc

Full API documentation is available at:

- [pkg.go.dev/github.com/grokify/gounidoc](https://pkg.go.dev/github.com/grokify/gounidoc)
