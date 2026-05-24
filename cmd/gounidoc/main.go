// gounidoc is a CLI tool and MCP server for PDF and Office document operations.
//
// It can be used as:
//   - A command-line tool for converting documents
//   - An MCP server for AI agents
//
// Running without a subcommand starts the MCP server (default behavior).
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/grokify/gounidoc"
	"github.com/grokify/gounidoc/skills/unidoc"
	"github.com/grokify/gounidoc/uniofficeutil"
	"github.com/grokify/gounidoc/unipdfutil"
	"github.com/grokify/mogo/fmt/progress"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	runtime "github.com/plexusone/omniskill/mcp/server"
	"github.com/spf13/cobra"
)

const (
	serverName    = "mcp-gounidoc"
	serverVersion = "v0.1.0"
)

var outputFormat string

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "gounidoc",
	Short: "Document conversion CLI and MCP server",
	Long: `gounidoc is a CLI tool and MCP (Model Context Protocol) server for
converting and processing PDF and Office documents using the UniDoc library.

Running without a subcommand starts the MCP server (default behavior).

Features:
  - PDF to DOCX conversion with table, image, and layout extraction
  - DOCX to PDF conversion
  - PDF page to PNG rendering
  - Office document metadata reading

Environment Variables:
  UNIDOC_KEY        - API key for both UniPDF and UniOffice
  UNIDOC_KEY_PDF    - API key specifically for UniPDF
  UNIDOC_KEY_OFFICE - API key specifically for UniOffice`,
	Example: `  # Start MCP server
  gounidoc serve

  # Convert PDF to DOCX
  gounidoc pdf2docx -i input.pdf -o output.docx

  # Convert with options
  gounidoc pdf2docx -i input.pdf --no-tables --no-images

  # Convert DOCX to PDF
  gounidoc docx2pdf -i document.docx -o output.pdf

  # Convert PDF page to PNG
  gounidoc pdf2png -i document.pdf -o page1.png --page 1 --width 1200

  # Read document metadata
  gounidoc metadata -f document.docx`,
	SilenceUsage: true,
	RunE:         runServer, // Default: run MCP server
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the MCP server",
	Long:  "Start the MCP server using stdio transport for communication with MCP clients.",
	RunE:  runServer,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s %s\n", serverName, serverVersion)
	},
}

// pdf2docx command
var (
	pdf2docxInput     string
	pdf2docxOutput    string
	pdf2docxNoTables  bool
	pdf2docxNoImages  bool
	pdf2docxNoLayout  bool
	pdf2docxBasicMode bool
)

var pdf2docxCmd = &cobra.Command{
	Use:   "pdf2docx",
	Short: "Convert PDF to DOCX",
	Long: `Convert a PDF file to DOCX format with enhanced formatting options.

Features:
  - Table extraction: Auto-detects tables and converts to DOCX tables with borders
  - Image detection: Detects images and notes their dimensions (placeholder text)
  - Layout detection: Detects headings (numbered sections, ALL CAPS, short lines) and applies bold formatting`,
	Example: `  # Full conversion with all features
  gounidoc pdf2docx -i input.pdf -o output.docx

  # Basic mode (text only, faster)
  gounidoc pdf2docx -i input.pdf -b

  # Disable specific features
  gounidoc pdf2docx -i input.pdf --no-tables
  gounidoc pdf2docx -i input.pdf --no-images
  gounidoc pdf2docx -i input.pdf --no-layout`,
	RunE: runPDF2DOCX,
}

// docx2pdf command
var (
	docx2pdfInput  string
	docx2pdfOutput string
)

var docx2pdfCmd = &cobra.Command{
	Use:     "docx2pdf",
	Short:   "Convert DOCX to PDF",
	Long:    "Convert a DOCX file to PDF format.",
	Example: "  gounidoc docx2pdf -i document.docx -o output.pdf",
	RunE:    runDOCX2PDF,
}

// pdf2png command
var (
	pdf2pngInput  string
	pdf2pngOutput string
	pdf2pngPage   int
	pdf2pngWidth  int
)

var pdf2pngCmd = &cobra.Command{
	Use:     "pdf2png",
	Short:   "Convert PDF page to PNG",
	Long:    "Convert a specific page from a PDF file to PNG image format.",
	Example: "  gounidoc pdf2png -i document.pdf -o page1.png --page 1 --width 1200",
	RunE:    runPDF2PNG,
}

// metadata command
var metadataFile string

var metadataCmd = &cobra.Command{
	Use:     "metadata",
	Short:   "Read Office document metadata",
	Long:    "Read metadata from an Office document (DOCX, PPTX, XLSX).",
	Example: "  gounidoc metadata -f document.docx",
	RunE:    runMetadata,
}

// merge command
var (
	mergeInputs []string
	mergeOutput string
)

var mergeCmd = &cobra.Command{
	Use:   "merge",
	Short: "Merge multiple PDF files",
	Long:  "Merge multiple PDF files into a single output file. Files are merged in the order provided.",
	Example: `  gounidoc merge -i file1.pdf -i file2.pdf -i file3.pdf -o merged.pdf
  gounidoc merge -i "*.pdf" -o merged.pdf`,
	RunE: runMerge,
}

// split command
var (
	splitInput     string
	splitOutputDir string
)

var splitCmd = &cobra.Command{
	Use:     "split",
	Short:   "Split PDF into individual pages",
	Long:    "Split a PDF file into individual pages. Each page is saved as a separate PDF file.",
	Example: "  gounidoc split -i document.pdf -d ./pages/",
	RunE:    runSplit,
}

// extract command
var (
	extractInput  string
	extractOutput string
	extractPages  []int
)

var extractCmd = &cobra.Command{
	Use:     "extract",
	Short:   "Extract specific pages from PDF",
	Long:    "Extract specific pages from a PDF file into a new PDF.",
	Example: "  gounidoc extract -i document.pdf -o extracted.pdf --pages 1,3,5",
	RunE:    runExtract,
}

// pagecount command
var pagecountInput string

var pagecountCmd = &cobra.Command{
	Use:     "pagecount",
	Short:   "Get PDF page count",
	Long:    "Get the number of pages in a PDF file.",
	Example: "  gounidoc pagecount -i document.pdf",
	RunE:    runPagecount,
}

// batch-pdf2docx command
var (
	batchPdf2docxInputs    []string
	batchPdf2docxOutputDir string
	batchPdf2docxNoTables  bool
	batchPdf2docxNoImages  bool
	batchPdf2docxNoLayout  bool
	batchPdf2docxBasicMode bool
)

var batchPdf2docxCmd = &cobra.Command{
	Use:   "batch-pdf2docx",
	Short: "Convert multiple PDF files to DOCX",
	Long:  "Convert multiple PDF files to DOCX format in batch. Supports glob patterns.",
	Example: `  gounidoc batch-pdf2docx -i "*.pdf" -d ./output/
  gounidoc batch-pdf2docx -i file1.pdf -i file2.pdf -d ./output/`,
	RunE: runBatchPdf2docx,
}

// batch-docx2pdf command
var (
	batchDocx2pdfInputs    []string
	batchDocx2pdfOutputDir string
)

var batchDocx2pdfCmd = &cobra.Command{
	Use:   "batch-docx2pdf",
	Short: "Convert multiple DOCX files to PDF",
	Long:  "Convert multiple DOCX files to PDF format in batch. Supports glob patterns.",
	Example: `  gounidoc batch-docx2pdf -i "*.docx" -d ./output/
  gounidoc batch-docx2pdf -i file1.docx -i file2.docx -d ./output/`,
	RunE: runBatchDocx2pdf,
}

// xlsx-info command
var xlsxInfoFile string

var xlsxInfoCmd = &cobra.Command{
	Use:     "xlsx-info",
	Short:   "Get XLSX workbook information",
	Long:    "Get information about an XLSX workbook including metadata and sheet names.",
	Example: "  gounidoc xlsx-info -f workbook.xlsx",
	RunE:    runXlsxInfo,
}

// pptx-info command
var pptxInfoFile string

var pptxInfoCmd = &cobra.Command{
	Use:     "pptx-info",
	Short:   "Get PPTX presentation information",
	Long:    "Get information about a PPTX presentation including metadata and slide count.",
	Example: "  gounidoc pptx-info -f presentation.pptx",
	RunE:    runPptxInfo,
}

// pdf-text command
var (
	pdfTextInput  string
	pdfTextOutput string
)

var pdfTextCmd = &cobra.Command{
	Use:     "pdf-text",
	Short:   "Extract text from PDF",
	Long:    "Extract all text from a PDF file.",
	Example: "  gounidoc pdf-text -i document.pdf -o text.txt",
	RunE:    runPdfText,
}

// pdf-images command
var (
	pdfImagesInput     string
	pdfImagesOutputDir string
)

var pdfImagesCmd = &cobra.Command{
	Use:     "pdf-images",
	Short:   "Extract images from PDF",
	Long:    "Extract all images from a PDF file to an output directory.",
	Example: "  gounidoc pdf-images -i document.pdf -d ./images/",
	RunE:    runPdfImages,
}

// pdf2images command
var (
	pdf2imagesInput     string
	pdf2imagesOutputDir string
	pdf2imagesFormat    string
	pdf2imagesDPI       int
	pdf2imagesQuality   int
)

var pdf2imagesCmd = &cobra.Command{
	Use:     "pdf2images",
	Short:   "Render PDF pages as images",
	Long:    "Render all pages of a PDF as images (PNG or JPEG).",
	Example: "  gounidoc pdf2images -i document.pdf -d ./pages/ --format png --dpi 150",
	RunE:    runPdf2images,
}

// images2pdf command
var (
	images2pdfInputs   []string
	images2pdfOutput   string
	images2pdfPageSize string
	images2pdfMargin   float64
)

var images2pdfCmd = &cobra.Command{
	Use:   "images2pdf",
	Short: "Combine images into PDF",
	Long:  "Combine multiple images into a single PDF file.",
	Example: `  gounidoc images2pdf -i image1.png -i image2.jpg -o output.pdf
  gounidoc images2pdf -i "*.png" -o output.pdf --page-size a4`,
	RunE: runImages2pdf,
}

// pdf-watermark command
var (
	pdfWatermarkInput    string
	pdfWatermarkOutput   string
	pdfWatermarkText     string
	pdfWatermarkImage    string
	pdfWatermarkFontSize float64
	pdfWatermarkColor    string
	pdfWatermarkOpacity  float64
	pdfWatermarkRotation float64
	pdfWatermarkPosition string
)

var pdfWatermarkCmd = &cobra.Command{
	Use:   "pdf-watermark",
	Short: "Add watermark to PDF",
	Long:  "Add a text or image watermark to a PDF file.",
	Example: `  gounidoc pdf-watermark -i document.pdf -o watermarked.pdf --text "CONFIDENTIAL"
  gounidoc pdf-watermark -i document.pdf -o watermarked.pdf --image logo.png`,
	RunE: runPdfWatermark,
}

// pdf-rotate command
var (
	pdfRotateInput   string
	pdfRotateOutput  string
	pdfRotateDegrees int
	pdfRotatePages   []int
)

var pdfRotateCmd = &cobra.Command{
	Use:     "pdf-rotate",
	Short:   "Rotate PDF pages",
	Long:    "Rotate pages in a PDF file by 90, 180, or 270 degrees.",
	Example: "  gounidoc pdf-rotate -i document.pdf -o rotated.pdf --degrees 90",
	RunE:    runPdfRotate,
}

// pdf-compress command
var (
	pdfCompressInput   string
	pdfCompressOutput  string
	pdfCompressQuality int
)

var pdfCompressCmd = &cobra.Command{
	Use:     "pdf-compress",
	Short:   "Compress PDF file",
	Long:    "Compress a PDF file to reduce its size.",
	Example: "  gounidoc pdf-compress -i large.pdf -o compressed.pdf --quality 70",
	RunE:    runPdfCompress,
}

// pdf-metadata command
var (
	pdfMetadataInput    string
	pdfMetadataOutput   string
	pdfMetadataTitle    string
	pdfMetadataAuthor   string
	pdfMetadataSubject  string
	pdfMetadataKeywords string
)

var pdfMetadataCmd = &cobra.Command{
	Use:   "pdf-metadata",
	Short: "Read or write PDF metadata",
	Long:  "Read or write PDF document metadata (title, author, subject, keywords).",
	Example: `  gounidoc pdf-metadata -i document.pdf
  gounidoc pdf-metadata -i document.pdf -o updated.pdf --title "My Document"`,
	RunE: runPdfMetadata,
}

// pdf-encrypt command
var (
	pdfEncryptInput         string
	pdfEncryptOutput        string
	pdfEncryptUserPassword  string
	pdfEncryptOwnerPassword string
	pdfEncryptNoPrint       bool
	pdfEncryptNoCopy        bool
	pdfEncryptNoModify      bool
)

var pdfEncryptCmd = &cobra.Command{
	Use:   "pdf-encrypt",
	Short: "Encrypt PDF with password",
	Long:  "Encrypt a PDF file with password protection.",
	Example: `  gounidoc pdf-encrypt -i document.pdf -o encrypted.pdf --user-password secret
  gounidoc pdf-encrypt -i document.pdf -o encrypted.pdf --owner-password admin --no-print`,
	RunE: runPdfEncrypt,
}

// pdf-decrypt command
var (
	pdfDecryptInput    string
	pdfDecryptOutput   string
	pdfDecryptPassword string
)

var pdfDecryptCmd = &cobra.Command{
	Use:     "pdf-decrypt",
	Short:   "Decrypt PDF file",
	Long:    "Remove password protection from a PDF file.",
	Example: "  gounidoc pdf-decrypt -i encrypted.pdf -o decrypted.pdf -p secret",
	RunE:    runPdfDecrypt,
}

// docx-text command
var (
	docxTextInput  string
	docxTextOutput string
)

var docxTextCmd = &cobra.Command{
	Use:     "docx-text",
	Short:   "Extract text from DOCX",
	Long:    "Extract all text from a DOCX file.",
	Example: "  gounidoc docx-text -i document.docx -o text.txt",
	RunE:    runDocxText,
}

// xlsx2csv command
var (
	xlsx2csvInput     string
	xlsx2csvOutput    string
	xlsx2csvSheet     string
	xlsx2csvDelimiter string
	xlsx2csvAllSheets bool
)

var xlsx2csvCmd = &cobra.Command{
	Use:   "xlsx2csv",
	Short: "Convert XLSX to CSV",
	Long:  "Convert an XLSX spreadsheet to CSV format.",
	Example: `  gounidoc xlsx2csv -i workbook.xlsx -o data.csv
  gounidoc xlsx2csv -i workbook.xlsx -o ./output/ --all-sheets`,
	RunE: runXlsx2csv,
}

// docx-replace command
var (
	docxReplaceInput         string
	docxReplaceOutput        string
	docxReplaceFind          string
	docxReplaceReplace       string
	docxReplaceCaseSensitive bool
	docxReplaceWholeWord     bool
)

var docxReplaceCmd = &cobra.Command{
	Use:     "docx-replace",
	Short:   "Find and replace in DOCX",
	Long:    "Find and replace text in a DOCX file.",
	Example: "  gounidoc docx-replace -i document.docx -o output.docx --find \"old\" --replace \"new\"",
	RunE:    runDocxReplace,
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "json",
		"output format: json, pretty (default: json)")

	// Add commands
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(pdf2docxCmd)
	rootCmd.AddCommand(docx2pdfCmd)
	rootCmd.AddCommand(pdf2pngCmd)
	rootCmd.AddCommand(metadataCmd)
	rootCmd.AddCommand(mergeCmd)
	rootCmd.AddCommand(splitCmd)
	rootCmd.AddCommand(extractCmd)
	rootCmd.AddCommand(pagecountCmd)
	rootCmd.AddCommand(batchPdf2docxCmd)
	rootCmd.AddCommand(batchDocx2pdfCmd)
	rootCmd.AddCommand(xlsxInfoCmd)
	rootCmd.AddCommand(pptxInfoCmd)

	// pdf2docx flags
	pdf2docxCmd.Flags().StringVarP(&pdf2docxInput, "input", "i", "", "Input PDF file path (required)")
	pdf2docxCmd.Flags().StringVarP(&pdf2docxOutput, "output", "o", "", "Output DOCX file path (optional)")
	pdf2docxCmd.Flags().BoolVar(&pdf2docxNoTables, "no-tables", false, "Disable table extraction")
	pdf2docxCmd.Flags().BoolVar(&pdf2docxNoImages, "no-images", false, "Disable image extraction")
	pdf2docxCmd.Flags().BoolVar(&pdf2docxNoLayout, "no-layout", false, "Disable layout/heading detection")
	pdf2docxCmd.Flags().BoolVarP(&pdf2docxBasicMode, "basic", "b", false, "Basic mode: text only")
	_ = pdf2docxCmd.MarkFlagRequired("input")

	// docx2pdf flags
	docx2pdfCmd.Flags().StringVarP(&docx2pdfInput, "input", "i", "", "Input DOCX file path (required)")
	docx2pdfCmd.Flags().StringVarP(&docx2pdfOutput, "output", "o", "", "Output PDF file path (optional)")
	_ = docx2pdfCmd.MarkFlagRequired("input")

	// pdf2png flags
	pdf2pngCmd.Flags().StringVarP(&pdf2pngInput, "input", "i", "", "Input PDF file path (required)")
	pdf2pngCmd.Flags().StringVarP(&pdf2pngOutput, "output", "o", "", "Output PNG file path (required)")
	pdf2pngCmd.Flags().IntVar(&pdf2pngPage, "page", 1, "Page number (1-indexed)")
	pdf2pngCmd.Flags().IntVar(&pdf2pngWidth, "width", 1200, "Output width in pixels")
	_ = pdf2pngCmd.MarkFlagRequired("input")
	_ = pdf2pngCmd.MarkFlagRequired("output")

	// metadata flags
	metadataCmd.Flags().StringVarP(&metadataFile, "file", "f", "", "File path (required)")
	_ = metadataCmd.MarkFlagRequired("file")

	// merge flags
	mergeCmd.Flags().StringArrayVarP(&mergeInputs, "input", "i", nil, "Input PDF files (can be specified multiple times)")
	mergeCmd.Flags().StringVarP(&mergeOutput, "output", "o", "", "Output merged PDF file (required)")
	_ = mergeCmd.MarkFlagRequired("input")
	_ = mergeCmd.MarkFlagRequired("output")

	// split flags
	splitCmd.Flags().StringVarP(&splitInput, "input", "i", "", "Input PDF file path (required)")
	splitCmd.Flags().StringVarP(&splitOutputDir, "dir", "d", "", "Output directory (optional, defaults to input file directory)")
	_ = splitCmd.MarkFlagRequired("input")

	// extract flags
	extractCmd.Flags().StringVarP(&extractInput, "input", "i", "", "Input PDF file path (required)")
	extractCmd.Flags().StringVarP(&extractOutput, "output", "o", "", "Output PDF file path (required)")
	extractCmd.Flags().IntSliceVar(&extractPages, "pages", nil, "Page numbers to extract (1-indexed, comma-separated)")
	_ = extractCmd.MarkFlagRequired("input")
	_ = extractCmd.MarkFlagRequired("output")
	_ = extractCmd.MarkFlagRequired("pages")

	// pagecount flags
	pagecountCmd.Flags().StringVarP(&pagecountInput, "input", "i", "", "Input PDF file path (required)")
	_ = pagecountCmd.MarkFlagRequired("input")

	// batch-pdf2docx flags
	batchPdf2docxCmd.Flags().StringArrayVarP(&batchPdf2docxInputs, "input", "i", nil, "Input PDF files (can be specified multiple times, supports glob patterns)")
	batchPdf2docxCmd.Flags().StringVarP(&batchPdf2docxOutputDir, "dir", "d", "", "Output directory (optional)")
	batchPdf2docxCmd.Flags().BoolVar(&batchPdf2docxNoTables, "no-tables", false, "Disable table extraction")
	batchPdf2docxCmd.Flags().BoolVar(&batchPdf2docxNoImages, "no-images", false, "Disable image extraction")
	batchPdf2docxCmd.Flags().BoolVar(&batchPdf2docxNoLayout, "no-layout", false, "Disable layout/heading detection")
	batchPdf2docxCmd.Flags().BoolVarP(&batchPdf2docxBasicMode, "basic", "b", false, "Basic mode: text only")
	_ = batchPdf2docxCmd.MarkFlagRequired("input")

	// batch-docx2pdf flags
	batchDocx2pdfCmd.Flags().StringArrayVarP(&batchDocx2pdfInputs, "input", "i", nil, "Input DOCX files (can be specified multiple times, supports glob patterns)")
	batchDocx2pdfCmd.Flags().StringVarP(&batchDocx2pdfOutputDir, "dir", "d", "", "Output directory (optional)")
	_ = batchDocx2pdfCmd.MarkFlagRequired("input")

	// xlsx-info flags
	xlsxInfoCmd.Flags().StringVarP(&xlsxInfoFile, "file", "f", "", "XLSX file path (required)")
	_ = xlsxInfoCmd.MarkFlagRequired("file")

	// pptx-info flags
	pptxInfoCmd.Flags().StringVarP(&pptxInfoFile, "file", "f", "", "PPTX file path (required)")
	_ = pptxInfoCmd.MarkFlagRequired("file")

	// Add Phase 1-4 commands
	rootCmd.AddCommand(pdfTextCmd)
	rootCmd.AddCommand(pdfImagesCmd)
	rootCmd.AddCommand(pdf2imagesCmd)
	rootCmd.AddCommand(images2pdfCmd)
	rootCmd.AddCommand(pdfWatermarkCmd)
	rootCmd.AddCommand(pdfRotateCmd)
	rootCmd.AddCommand(pdfCompressCmd)
	rootCmd.AddCommand(pdfMetadataCmd)
	rootCmd.AddCommand(pdfEncryptCmd)
	rootCmd.AddCommand(pdfDecryptCmd)
	rootCmd.AddCommand(docxTextCmd)
	rootCmd.AddCommand(xlsx2csvCmd)
	rootCmd.AddCommand(docxReplaceCmd)

	// pdf-text flags
	pdfTextCmd.Flags().StringVarP(&pdfTextInput, "input", "i", "", "Input PDF file path (required)")
	pdfTextCmd.Flags().StringVarP(&pdfTextOutput, "output", "o", "", "Output text file path (optional)")
	_ = pdfTextCmd.MarkFlagRequired("input")

	// pdf-images flags
	pdfImagesCmd.Flags().StringVarP(&pdfImagesInput, "input", "i", "", "Input PDF file path (required)")
	pdfImagesCmd.Flags().StringVarP(&pdfImagesOutputDir, "dir", "d", "", "Output directory (optional)")
	_ = pdfImagesCmd.MarkFlagRequired("input")

	// pdf2images flags
	pdf2imagesCmd.Flags().StringVarP(&pdf2imagesInput, "input", "i", "", "Input PDF file path (required)")
	pdf2imagesCmd.Flags().StringVarP(&pdf2imagesOutputDir, "dir", "d", "", "Output directory (optional)")
	pdf2imagesCmd.Flags().StringVar(&pdf2imagesFormat, "format", "png", "Output format: png, jpeg")
	pdf2imagesCmd.Flags().IntVar(&pdf2imagesDPI, "dpi", 150, "Resolution in DPI")
	pdf2imagesCmd.Flags().IntVar(&pdf2imagesQuality, "quality", 85, "JPEG quality 1-100")
	_ = pdf2imagesCmd.MarkFlagRequired("input")

	// images2pdf flags
	images2pdfCmd.Flags().StringArrayVarP(&images2pdfInputs, "input", "i", nil, "Input image files (can be specified multiple times)")
	images2pdfCmd.Flags().StringVarP(&images2pdfOutput, "output", "o", "", "Output PDF file path (required)")
	images2pdfCmd.Flags().StringVar(&images2pdfPageSize, "page-size", "letter", "Page size: letter, a4, fit")
	images2pdfCmd.Flags().Float64Var(&images2pdfMargin, "margin", 36.0, "Margin in points")
	_ = images2pdfCmd.MarkFlagRequired("input")
	_ = images2pdfCmd.MarkFlagRequired("output")

	// pdf-watermark flags
	pdfWatermarkCmd.Flags().StringVarP(&pdfWatermarkInput, "input", "i", "", "Input PDF file path (required)")
	pdfWatermarkCmd.Flags().StringVarP(&pdfWatermarkOutput, "output", "o", "", "Output PDF file path (required)")
	pdfWatermarkCmd.Flags().StringVar(&pdfWatermarkText, "text", "", "Watermark text")
	pdfWatermarkCmd.Flags().StringVar(&pdfWatermarkImage, "image", "", "Watermark image path")
	pdfWatermarkCmd.Flags().Float64Var(&pdfWatermarkFontSize, "font-size", 48, "Font size")
	pdfWatermarkCmd.Flags().StringVar(&pdfWatermarkColor, "color", "#888888", "Text color (hex)")
	pdfWatermarkCmd.Flags().Float64Var(&pdfWatermarkOpacity, "opacity", 0.3, "Opacity 0.0-1.0")
	pdfWatermarkCmd.Flags().Float64Var(&pdfWatermarkRotation, "rotation", -45, "Rotation angle in degrees")
	pdfWatermarkCmd.Flags().StringVar(&pdfWatermarkPosition, "position", "center", "Position: center, top-left, top-right, bottom-left, bottom-right")
	_ = pdfWatermarkCmd.MarkFlagRequired("input")
	_ = pdfWatermarkCmd.MarkFlagRequired("output")

	// pdf-rotate flags
	pdfRotateCmd.Flags().StringVarP(&pdfRotateInput, "input", "i", "", "Input PDF file path (required)")
	pdfRotateCmd.Flags().StringVarP(&pdfRotateOutput, "output", "o", "", "Output PDF file path (required)")
	pdfRotateCmd.Flags().IntVar(&pdfRotateDegrees, "degrees", 90, "Rotation: 90, 180, 270")
	pdfRotateCmd.Flags().IntSliceVar(&pdfRotatePages, "pages", nil, "Pages to rotate (1-indexed, comma-separated)")
	_ = pdfRotateCmd.MarkFlagRequired("input")
	_ = pdfRotateCmd.MarkFlagRequired("output")

	// pdf-compress flags
	pdfCompressCmd.Flags().StringVarP(&pdfCompressInput, "input", "i", "", "Input PDF file path (required)")
	pdfCompressCmd.Flags().StringVarP(&pdfCompressOutput, "output", "o", "", "Output PDF file path (required)")
	pdfCompressCmd.Flags().IntVar(&pdfCompressQuality, "quality", 80, "Image quality 1-100")
	_ = pdfCompressCmd.MarkFlagRequired("input")
	_ = pdfCompressCmd.MarkFlagRequired("output")

	// pdf-metadata flags
	pdfMetadataCmd.Flags().StringVarP(&pdfMetadataInput, "input", "i", "", "Input PDF file path (required)")
	pdfMetadataCmd.Flags().StringVarP(&pdfMetadataOutput, "output", "o", "", "Output PDF file path (for writing)")
	pdfMetadataCmd.Flags().StringVar(&pdfMetadataTitle, "title", "", "Document title")
	pdfMetadataCmd.Flags().StringVar(&pdfMetadataAuthor, "author", "", "Document author")
	pdfMetadataCmd.Flags().StringVar(&pdfMetadataSubject, "subject", "", "Document subject")
	pdfMetadataCmd.Flags().StringVar(&pdfMetadataKeywords, "keywords", "", "Document keywords")
	_ = pdfMetadataCmd.MarkFlagRequired("input")

	// pdf-encrypt flags
	pdfEncryptCmd.Flags().StringVarP(&pdfEncryptInput, "input", "i", "", "Input PDF file path (required)")
	pdfEncryptCmd.Flags().StringVarP(&pdfEncryptOutput, "output", "o", "", "Output PDF file path (required)")
	pdfEncryptCmd.Flags().StringVar(&pdfEncryptUserPassword, "user-password", "", "Password to open PDF")
	pdfEncryptCmd.Flags().StringVar(&pdfEncryptOwnerPassword, "owner-password", "", "Password for full access")
	pdfEncryptCmd.Flags().BoolVar(&pdfEncryptNoPrint, "no-print", false, "Disable printing")
	pdfEncryptCmd.Flags().BoolVar(&pdfEncryptNoCopy, "no-copy", false, "Disable copying")
	pdfEncryptCmd.Flags().BoolVar(&pdfEncryptNoModify, "no-modify", false, "Disable modifying")
	_ = pdfEncryptCmd.MarkFlagRequired("input")
	_ = pdfEncryptCmd.MarkFlagRequired("output")

	// pdf-decrypt flags
	pdfDecryptCmd.Flags().StringVarP(&pdfDecryptInput, "input", "i", "", "Input PDF file path (required)")
	pdfDecryptCmd.Flags().StringVarP(&pdfDecryptOutput, "output", "o", "", "Output PDF file path (required)")
	pdfDecryptCmd.Flags().StringVarP(&pdfDecryptPassword, "password", "p", "", "Password to decrypt")
	_ = pdfDecryptCmd.MarkFlagRequired("input")
	_ = pdfDecryptCmd.MarkFlagRequired("output")
	_ = pdfDecryptCmd.MarkFlagRequired("password")

	// docx-text flags
	docxTextCmd.Flags().StringVarP(&docxTextInput, "input", "i", "", "Input DOCX file path (required)")
	docxTextCmd.Flags().StringVarP(&docxTextOutput, "output", "o", "", "Output text file path (optional)")
	_ = docxTextCmd.MarkFlagRequired("input")

	// xlsx2csv flags
	xlsx2csvCmd.Flags().StringVarP(&xlsx2csvInput, "input", "i", "", "Input XLSX file path (required)")
	xlsx2csvCmd.Flags().StringVarP(&xlsx2csvOutput, "output", "o", "", "Output CSV file or directory (required)")
	xlsx2csvCmd.Flags().StringVar(&xlsx2csvSheet, "sheet", "", "Sheet name to export")
	xlsx2csvCmd.Flags().StringVar(&xlsx2csvDelimiter, "delimiter", ",", "CSV delimiter")
	xlsx2csvCmd.Flags().BoolVar(&xlsx2csvAllSheets, "all-sheets", false, "Export all sheets")
	_ = xlsx2csvCmd.MarkFlagRequired("input")
	_ = xlsx2csvCmd.MarkFlagRequired("output")

	// docx-replace flags
	docxReplaceCmd.Flags().StringVarP(&docxReplaceInput, "input", "i", "", "Input DOCX file path (required)")
	docxReplaceCmd.Flags().StringVarP(&docxReplaceOutput, "output", "o", "", "Output DOCX file path (required)")
	docxReplaceCmd.Flags().StringVar(&docxReplaceFind, "find", "", "Text to find (required)")
	docxReplaceCmd.Flags().StringVar(&docxReplaceReplace, "replace", "", "Replacement text (required)")
	docxReplaceCmd.Flags().BoolVar(&docxReplaceCaseSensitive, "case-sensitive", true, "Case-sensitive search")
	docxReplaceCmd.Flags().BoolVar(&docxReplaceWholeWord, "whole-word", false, "Match whole words only")
	_ = docxReplaceCmd.MarkFlagRequired("input")
	_ = docxReplaceCmd.MarkFlagRequired("output")
	_ = docxReplaceCmd.MarkFlagRequired("find")
	_ = docxReplaceCmd.MarkFlagRequired("replace")
}

func initUniDoc() error {
	return gounidoc.SetMeteredKeyEnv()
}

func outputResult(result any) error {
	var data []byte
	var err error

	switch outputFormat {
	case "pretty":
		data, err = json.MarshalIndent(result, "", "  ")
	default:
		data, err = json.Marshal(result)
	}

	if err != nil {
		return fmt.Errorf("failed to marshal result: %w", err)
	}

	fmt.Println(string(data))
	return nil
}

func runServer(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	if err := initUniDoc(); err != nil {
		return fmt.Errorf("failed to initialize UniDoc: %w", err)
	}

	// Create omniskill Runtime
	rt := runtime.New(&mcp.Implementation{
		Name:    serverName,
		Version: serverVersion,
	}, nil)

	// Create and initialize UniDoc skill
	unidocSkill := unidoc.New()
	if err := unidocSkill.Init(ctx); err != nil {
		return fmt.Errorf("failed to initialize UniDoc skill: %w", err)
	}
	defer unidocSkill.Close()

	// Register skill with the runtime
	rt.RegisterSkill(unidocSkill)

	// Run server with stdio transport
	if err := rt.ServeStdio(ctx); err != nil {
		return fmt.Errorf("server error: %w", err)
	}

	return nil
}

func runPDF2DOCX(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	if err := initUniDoc(); err != nil {
		return fmt.Errorf("failed to initialize UniDoc: %w", err)
	}

	skill := unidoc.New()
	if err := skill.Init(ctx); err != nil {
		return fmt.Errorf("failed to initialize skill: %w", err)
	}
	defer skill.Close()

	opts := unipdfutil.DefaultConversionOptions()
	if pdf2docxBasicMode {
		opts.ExtractTables = false
		opts.ExtractImages = false
		opts.DetectLayout = false
	} else {
		if pdf2docxNoTables {
			opts.ExtractTables = false
		}
		if pdf2docxNoImages {
			opts.ExtractImages = false
		}
		if pdf2docxNoLayout {
			opts.DetectLayout = false
		}
	}

	fmt.Printf("Converting: %s -> %s\n", pdf2docxInput, pdf2docxOutput)
	fmt.Printf("Options: tables=%v, images=%v, layout=%v\n",
		opts.ExtractTables, opts.ExtractImages, opts.DetectLayout)

	outputPath, err := skill.ConvertPDFToDocx(ctx, pdf2docxInput, pdf2docxOutput, opts)
	if err != nil {
		return fmt.Errorf("conversion failed: %w", err)
	}

	fmt.Printf("Successfully converted to: %s\n", outputPath)
	return nil
}

func runDOCX2PDF(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	if err := initUniDoc(); err != nil {
		return fmt.Errorf("failed to initialize UniDoc: %w", err)
	}

	skill := unidoc.New()
	if err := skill.Init(ctx); err != nil {
		return fmt.Errorf("failed to initialize skill: %w", err)
	}
	defer skill.Close()

	fmt.Printf("Converting: %s -> %s\n", docx2pdfInput, docx2pdfOutput)

	outputPath, err := skill.ConvertDocxToPDF(ctx, docx2pdfInput, docx2pdfOutput)
	if err != nil {
		return fmt.Errorf("conversion failed: %w", err)
	}

	fmt.Printf("Successfully converted to: %s\n", outputPath)
	return nil
}

func runPDF2PNG(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	if err := initUniDoc(); err != nil {
		return fmt.Errorf("failed to initialize UniDoc: %w", err)
	}

	skill := unidoc.New()
	if err := skill.Init(ctx); err != nil {
		return fmt.Errorf("failed to initialize skill: %w", err)
	}
	defer skill.Close()

	fmt.Printf("Converting: %s (page %d) -> %s\n", pdf2pngInput, pdf2pngPage, pdf2pngOutput)

	if err := skill.ConvertPDFPageToPNG(ctx, pdf2pngInput, pdf2pngOutput, uint32(pdf2pngPage), uint32(pdf2pngWidth)); err != nil {
		return fmt.Errorf("conversion failed: %w", err)
	}

	fmt.Printf("Successfully converted to: %s\n", pdf2pngOutput)
	return nil
}

func runMetadata(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	if err := initUniDoc(); err != nil {
		return fmt.Errorf("failed to initialize UniDoc: %w", err)
	}

	skill := unidoc.New()
	if err := skill.Init(ctx); err != nil {
		return fmt.Errorf("failed to initialize skill: %w", err)
	}
	defer skill.Close()

	metadata, err := skill.ReadMetadata(ctx, metadataFile)
	if err != nil {
		return fmt.Errorf("failed to read metadata: %w", err)
	}

	return outputResult(map[string]any{
		"file_path": metadataFile,
		"metadata":  metadata,
	})
}

func runMerge(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	if err := initUniDoc(); err != nil {
		return fmt.Errorf("failed to initialize UniDoc: %w", err)
	}

	skill := unidoc.New()
	if err := skill.Init(ctx); err != nil {
		return fmt.Errorf("failed to initialize skill: %w", err)
	}
	defer skill.Close()

	fmt.Printf("Merging %d files into: %s\n", len(mergeInputs), mergeOutput)
	for i, f := range mergeInputs {
		fmt.Printf("  [%d] %s\n", i+1, f)
	}

	if err := skill.MergePDFs(ctx, mergeInputs, mergeOutput); err != nil {
		return fmt.Errorf("merge failed: %w", err)
	}

	fmt.Printf("Successfully merged to: %s\n", mergeOutput)
	return nil
}

func runSplit(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	if err := initUniDoc(); err != nil {
		return fmt.Errorf("failed to initialize UniDoc: %w", err)
	}

	skill := unidoc.New()
	if err := skill.Init(ctx); err != nil {
		return fmt.Errorf("failed to initialize skill: %w", err)
	}
	defer skill.Close()

	fmt.Printf("Splitting: %s\n", splitInput)
	if splitOutputDir != "" {
		fmt.Printf("Output directory: %s\n", splitOutputDir)
	}

	outputPaths, err := skill.SplitPDF(ctx, splitInput, splitOutputDir)
	if err != nil {
		return fmt.Errorf("split failed: %w", err)
	}

	fmt.Printf("Successfully split into %d pages:\n", len(outputPaths))
	for _, p := range outputPaths {
		fmt.Printf("  %s\n", p)
	}
	return nil
}

func runExtract(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	if err := initUniDoc(); err != nil {
		return fmt.Errorf("failed to initialize UniDoc: %w", err)
	}

	skill := unidoc.New()
	if err := skill.Init(ctx); err != nil {
		return fmt.Errorf("failed to initialize skill: %w", err)
	}
	defer skill.Close()

	fmt.Printf("Extracting pages %v from: %s\n", extractPages, extractInput)

	if err := skill.ExtractPDFPages(ctx, extractInput, extractOutput, extractPages); err != nil {
		return fmt.Errorf("extraction failed: %w", err)
	}

	fmt.Printf("Successfully extracted to: %s\n", extractOutput)
	return nil
}

func runPagecount(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	if err := initUniDoc(); err != nil {
		return fmt.Errorf("failed to initialize UniDoc: %w", err)
	}

	skill := unidoc.New()
	if err := skill.Init(ctx); err != nil {
		return fmt.Errorf("failed to initialize skill: %w", err)
	}
	defer skill.Close()

	count, err := skill.GetPDFPageCount(ctx, pagecountInput)
	if err != nil {
		return fmt.Errorf("failed to get page count: %w", err)
	}

	return outputResult(map[string]any{
		"file_path":  pagecountInput,
		"page_count": count,
	})
}

func runBatchPdf2docx(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	if err := initUniDoc(); err != nil {
		return fmt.Errorf("failed to initialize UniDoc: %w", err)
	}

	skill := unidoc.New()
	if err := skill.Init(ctx); err != nil {
		return fmt.Errorf("failed to initialize skill: %w", err)
	}
	defer skill.Close()

	// Expand glob patterns
	inputPaths, err := unipdfutil.ExpandGlobPatterns(batchPdf2docxInputs)
	if err != nil {
		return fmt.Errorf("failed to expand patterns: %w", err)
	}

	// Filter to PDF files only
	inputPaths = unipdfutil.FilterByExtension(inputPaths, ".pdf")

	if len(inputPaths) == 0 {
		return fmt.Errorf("no PDF files found matching the specified patterns")
	}

	fmt.Printf("Converting %d PDF files to DOCX\n", len(inputPaths))
	if batchPdf2docxOutputDir != "" {
		fmt.Printf("Output directory: %s\n", batchPdf2docxOutputDir)
	}

	opts := unipdfutil.DefaultConversionOptions()
	if batchPdf2docxBasicMode {
		opts.ExtractTables = false
		opts.ExtractImages = false
		opts.DetectLayout = false
	} else {
		if batchPdf2docxNoTables {
			opts.ExtractTables = false
		}
		if batchPdf2docxNoImages {
			opts.ExtractImages = false
		}
		if batchPdf2docxNoLayout {
			opts.DetectLayout = false
		}
	}

	// Create progress renderer
	renderer := progress.NewSingleStageRenderer(os.Stdout).
		WithBarWidth(30).
		WithTextWidth(40)

	progressCallback := func(current, total int, inputPath string) {
		renderer.Update(current, total, filepath.Base(inputPath))
	}

	result := skill.BatchConvertPDFToDocxWithProgress(ctx, inputPaths, batchPdf2docxOutputDir, opts, progressCallback)

	// Clear progress line and show summary
	renderer.Done("")
	fmt.Printf("Batch conversion complete:\n")
	fmt.Printf("  Total:     %d\n", result.TotalFiles)
	fmt.Printf("  Succeeded: %d\n", result.Succeeded)
	fmt.Printf("  Failed:    %d\n", result.Failed)

	if result.Failed > 0 {
		fmt.Printf("\nFailed files:\n")
		for _, r := range result.Results {
			if !r.Success {
				fmt.Printf("  %s: %s\n", r.InputPath, r.Error)
			}
		}
	}

	return nil
}

func runBatchDocx2pdf(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	if err := initUniDoc(); err != nil {
		return fmt.Errorf("failed to initialize UniDoc: %w", err)
	}

	skill := unidoc.New()
	if err := skill.Init(ctx); err != nil {
		return fmt.Errorf("failed to initialize skill: %w", err)
	}
	defer skill.Close()

	// Expand glob patterns
	inputPaths, err := unipdfutil.ExpandGlobPatterns(batchDocx2pdfInputs)
	if err != nil {
		return fmt.Errorf("failed to expand patterns: %w", err)
	}

	// Filter to DOCX files only
	inputPaths = unipdfutil.FilterByExtension(inputPaths, ".docx")

	if len(inputPaths) == 0 {
		return fmt.Errorf("no DOCX files found matching the specified patterns")
	}

	fmt.Printf("Converting %d DOCX files to PDF\n", len(inputPaths))
	if batchDocx2pdfOutputDir != "" {
		fmt.Printf("Output directory: %s\n", batchDocx2pdfOutputDir)
	}

	// Create progress renderer
	renderer := progress.NewSingleStageRenderer(os.Stdout).
		WithBarWidth(30).
		WithTextWidth(40)

	progressCallback := func(current, total int, inputPath string) {
		renderer.Update(current, total, filepath.Base(inputPath))
	}

	result := skill.BatchConvertDocxToPDFWithProgress(ctx, inputPaths, batchDocx2pdfOutputDir, progressCallback)

	// Clear progress line and show summary
	renderer.Done("")
	fmt.Printf("Batch conversion complete:\n")
	fmt.Printf("  Total:     %d\n", result.TotalFiles)
	fmt.Printf("  Succeeded: %d\n", result.Succeeded)
	fmt.Printf("  Failed:    %d\n", result.Failed)

	if result.Failed > 0 {
		fmt.Printf("\nFailed files:\n")
		for _, r := range result.Results {
			if !r.Success {
				fmt.Printf("  %s: %s\n", r.InputPath, r.Error)
			}
		}
	}

	return nil
}

func runXlsxInfo(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	if err := initUniDoc(); err != nil {
		return fmt.Errorf("failed to initialize UniDoc: %w", err)
	}

	skill := unidoc.New()
	if err := skill.Init(ctx); err != nil {
		return fmt.Errorf("failed to initialize skill: %w", err)
	}
	defer skill.Close()

	info, err := skill.GetXlsxInfo(ctx, xlsxInfoFile)
	if err != nil {
		return fmt.Errorf("failed to get XLSX info: %w", err)
	}

	return outputResult(map[string]any{
		"file_path":   xlsxInfoFile,
		"metadata":    info.Metadata,
		"sheet_count": info.SheetCount,
		"sheet_names": info.SheetNames,
	})
}

func runPptxInfo(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	if err := initUniDoc(); err != nil {
		return fmt.Errorf("failed to initialize UniDoc: %w", err)
	}

	skill := unidoc.New()
	if err := skill.Init(ctx); err != nil {
		return fmt.Errorf("failed to initialize skill: %w", err)
	}
	defer skill.Close()

	info, err := skill.GetPptxInfo(ctx, pptxInfoFile)
	if err != nil {
		return fmt.Errorf("failed to get PPTX info: %w", err)
	}

	return outputResult(map[string]any{
		"file_path":   pptxInfoFile,
		"metadata":    info.Metadata,
		"slide_count": info.SlideCount,
	})
}

// ---------------------------------------------------------------------------
// Phase 1: PDF Text & Image Operations - Handlers
// ---------------------------------------------------------------------------

func runPdfText(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	if err := initUniDoc(); err != nil {
		return fmt.Errorf("failed to initialize UniDoc: %w", err)
	}

	skill := unidoc.New()
	if err := skill.Init(ctx); err != nil {
		return fmt.Errorf("failed to initialize skill: %w", err)
	}
	defer skill.Close()

	text, err := skill.ExtractPDFText(ctx, pdfTextInput)
	if err != nil {
		return fmt.Errorf("failed to extract text: %w", err)
	}

	if pdfTextOutput != "" {
		if err := unipdfutil.ExtractPDFTextToFile(pdfTextInput, pdfTextOutput); err != nil {
			return fmt.Errorf("failed to write output: %w", err)
		}
		fmt.Printf("Text extracted to: %s\n", pdfTextOutput)
		return nil
	}

	fmt.Println(text)
	return nil
}

func runPdfImages(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	if err := initUniDoc(); err != nil {
		return fmt.Errorf("failed to initialize UniDoc: %w", err)
	}

	skill := unidoc.New()
	if err := skill.Init(ctx); err != nil {
		return fmt.Errorf("failed to initialize skill: %w", err)
	}
	defer skill.Close()

	images, err := skill.ExtractPDFImages(ctx, pdfImagesInput, pdfImagesOutputDir)
	if err != nil {
		return fmt.Errorf("failed to extract images: %w", err)
	}

	fmt.Printf("Extracted %d images:\n", len(images))
	for _, img := range images {
		fmt.Printf("  Page %d, Image %d: %s (%dx%d)\n", img.PageNumber, img.Index, img.FilePath, img.Width, img.Height)
	}
	return nil
}

func runPdf2images(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	if err := initUniDoc(); err != nil {
		return fmt.Errorf("failed to initialize UniDoc: %w", err)
	}

	skill := unidoc.New()
	if err := skill.Init(ctx); err != nil {
		return fmt.Errorf("failed to initialize skill: %w", err)
	}
	defer skill.Close()

	opts := unipdfutil.RenderOptions{
		Format:  pdf2imagesFormat,
		DPI:     float64(pdf2imagesDPI),
		Quality: pdf2imagesQuality,
	}

	outputPaths, err := skill.RenderPDFToImages(ctx, pdf2imagesInput, pdf2imagesOutputDir, opts)
	if err != nil {
		return fmt.Errorf("failed to render PDF: %w", err)
	}

	fmt.Printf("Rendered %d pages:\n", len(outputPaths))
	for _, p := range outputPaths {
		fmt.Printf("  %s\n", p)
	}
	return nil
}

func runImages2pdf(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	if err := initUniDoc(); err != nil {
		return fmt.Errorf("failed to initialize UniDoc: %w", err)
	}

	skill := unidoc.New()
	if err := skill.Init(ctx); err != nil {
		return fmt.Errorf("failed to initialize skill: %w", err)
	}
	defer skill.Close()

	// Expand glob patterns
	imagePaths, err := unipdfutil.ExpandGlobPatterns(images2pdfInputs)
	if err != nil {
		return fmt.Errorf("failed to expand patterns: %w", err)
	}

	if len(imagePaths) == 0 {
		return fmt.Errorf("no images found matching the specified patterns")
	}

	opts := unipdfutil.ImageToPDFOptions{
		PageSize: images2pdfPageSize,
		Margin:   images2pdfMargin,
	}

	fmt.Printf("Combining %d images into PDF\n", len(imagePaths))

	if err := skill.ImagesToPDF(ctx, imagePaths, images2pdfOutput, opts); err != nil {
		return fmt.Errorf("failed to create PDF: %w", err)
	}

	fmt.Printf("Successfully created: %s\n", images2pdfOutput)
	return nil
}

// ---------------------------------------------------------------------------
// Phase 2: PDF Manipulation - Handlers
// ---------------------------------------------------------------------------

func runPdfWatermark(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	if err := initUniDoc(); err != nil {
		return fmt.Errorf("failed to initialize UniDoc: %w", err)
	}

	skill := unidoc.New()
	if err := skill.Init(ctx); err != nil {
		return fmt.Errorf("failed to initialize skill: %w", err)
	}
	defer skill.Close()

	if pdfWatermarkText == "" && pdfWatermarkImage == "" {
		return fmt.Errorf("either --text or --image is required")
	}

	opts := unipdfutil.WatermarkOptions{
		Text:      pdfWatermarkText,
		ImagePath: pdfWatermarkImage,
		FontSize:  pdfWatermarkFontSize,
		Color:     pdfWatermarkColor,
		Opacity:   pdfWatermarkOpacity,
		Rotation:  pdfWatermarkRotation,
		Position:  pdfWatermarkPosition,
	}

	if err := skill.AddWatermark(ctx, pdfWatermarkInput, pdfWatermarkOutput, opts); err != nil {
		return fmt.Errorf("failed to add watermark: %w", err)
	}

	fmt.Printf("Watermark added: %s\n", pdfWatermarkOutput)
	return nil
}

func runPdfRotate(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	if err := initUniDoc(); err != nil {
		return fmt.Errorf("failed to initialize UniDoc: %w", err)
	}

	skill := unidoc.New()
	if err := skill.Init(ctx); err != nil {
		return fmt.Errorf("failed to initialize skill: %w", err)
	}
	defer skill.Close()

	if err := skill.RotatePDF(ctx, pdfRotateInput, pdfRotateOutput, pdfRotateDegrees, pdfRotatePages); err != nil {
		return fmt.Errorf("failed to rotate PDF: %w", err)
	}

	fmt.Printf("Rotated %d degrees: %s\n", pdfRotateDegrees, pdfRotateOutput)
	return nil
}

func runPdfCompress(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	if err := initUniDoc(); err != nil {
		return fmt.Errorf("failed to initialize UniDoc: %w", err)
	}

	skill := unidoc.New()
	if err := skill.Init(ctx); err != nil {
		return fmt.Errorf("failed to initialize skill: %w", err)
	}
	defer skill.Close()

	opts := unipdfutil.CompressOptions{
		ImageQuality: pdfCompressQuality,
	}

	result, err := skill.CompressPDF(ctx, pdfCompressInput, pdfCompressOutput, opts)
	if err != nil {
		return fmt.Errorf("failed to compress PDF: %w", err)
	}

	fmt.Printf("Compression complete:\n")
	fmt.Printf("  Original:   %d bytes\n", result.OriginalSize)
	fmt.Printf("  Compressed: %d bytes\n", result.CompressedSize)
	fmt.Printf("  Savings:    %.1f%%\n", result.SavingsPercent)
	fmt.Printf("  Output:     %s\n", pdfCompressOutput)
	return nil
}

func runPdfMetadata(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	if err := initUniDoc(); err != nil {
		return fmt.Errorf("failed to initialize UniDoc: %w", err)
	}

	skill := unidoc.New()
	if err := skill.Init(ctx); err != nil {
		return fmt.Errorf("failed to initialize skill: %w", err)
	}
	defer skill.Close()

	// If no output and no metadata fields, just read
	if pdfMetadataOutput == "" && pdfMetadataTitle == "" && pdfMetadataAuthor == "" && pdfMetadataSubject == "" && pdfMetadataKeywords == "" {
		meta, err := skill.ReadPDFMetadata(ctx, pdfMetadataInput)
		if err != nil {
			return fmt.Errorf("failed to read metadata: %w", err)
		}
		return outputResult(meta)
	}

	// Writing metadata
	if pdfMetadataOutput == "" {
		return fmt.Errorf("--output is required when writing metadata")
	}

	meta := unipdfutil.PDFMetadata{
		Title:    pdfMetadataTitle,
		Author:   pdfMetadataAuthor,
		Subject:  pdfMetadataSubject,
		Keywords: pdfMetadataKeywords,
	}

	if err := skill.WritePDFMetadata(ctx, pdfMetadataInput, pdfMetadataOutput, meta); err != nil {
		return fmt.Errorf("failed to write metadata: %w", err)
	}

	fmt.Printf("Metadata updated: %s\n", pdfMetadataOutput)
	return nil
}

// ---------------------------------------------------------------------------
// Phase 3: PDF Security - Handlers
// ---------------------------------------------------------------------------

func runPdfEncrypt(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	if err := initUniDoc(); err != nil {
		return fmt.Errorf("failed to initialize UniDoc: %w", err)
	}

	skill := unidoc.New()
	if err := skill.Init(ctx); err != nil {
		return fmt.Errorf("failed to initialize skill: %w", err)
	}
	defer skill.Close()

	if pdfEncryptUserPassword == "" && pdfEncryptOwnerPassword == "" {
		return fmt.Errorf("at least one password (--user-password or --owner-password) is required")
	}

	opts := unipdfutil.EncryptOptions{
		UserPassword:  pdfEncryptUserPassword,
		OwnerPassword: pdfEncryptOwnerPassword,
		Permissions: unipdfutil.PDFPermissions{
			Printing:       !pdfEncryptNoPrint,
			CopyContents:   !pdfEncryptNoCopy,
			ModifyContents: !pdfEncryptNoModify,
			ModifyAnnots:   !pdfEncryptNoModify,
		},
	}

	if err := skill.EncryptPDF(ctx, pdfEncryptInput, pdfEncryptOutput, opts); err != nil {
		return fmt.Errorf("failed to encrypt PDF: %w", err)
	}

	fmt.Printf("PDF encrypted: %s\n", pdfEncryptOutput)
	return nil
}

func runPdfDecrypt(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	if err := initUniDoc(); err != nil {
		return fmt.Errorf("failed to initialize UniDoc: %w", err)
	}

	skill := unidoc.New()
	if err := skill.Init(ctx); err != nil {
		return fmt.Errorf("failed to initialize skill: %w", err)
	}
	defer skill.Close()

	if err := skill.DecryptPDF(ctx, pdfDecryptInput, pdfDecryptOutput, pdfDecryptPassword); err != nil {
		return fmt.Errorf("failed to decrypt PDF: %w", err)
	}

	fmt.Printf("PDF decrypted: %s\n", pdfDecryptOutput)
	return nil
}

// ---------------------------------------------------------------------------
// Phase 4: Office Document Operations - Handlers
// ---------------------------------------------------------------------------

func runDocxText(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	if err := initUniDoc(); err != nil {
		return fmt.Errorf("failed to initialize UniDoc: %w", err)
	}

	skill := unidoc.New()
	if err := skill.Init(ctx); err != nil {
		return fmt.Errorf("failed to initialize skill: %w", err)
	}
	defer skill.Close()

	text, err := skill.ExtractDocxText(ctx, docxTextInput)
	if err != nil {
		return fmt.Errorf("failed to extract text: %w", err)
	}

	if docxTextOutput != "" {
		if err := uniofficeutil.ExtractDocxTextToFile(docxTextInput, docxTextOutput); err != nil {
			return fmt.Errorf("failed to write output: %w", err)
		}
		fmt.Printf("Text extracted to: %s\n", docxTextOutput)
		return nil
	}

	fmt.Println(text)
	return nil
}

func runXlsx2csv(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	if err := initUniDoc(); err != nil {
		return fmt.Errorf("failed to initialize UniDoc: %w", err)
	}

	skill := unidoc.New()
	if err := skill.Init(ctx); err != nil {
		return fmt.Errorf("failed to initialize skill: %w", err)
	}
	defer skill.Close()

	opts := uniofficeutil.CSVOptions{
		SheetName: xlsx2csvSheet,
		Delimiter: xlsx2csvDelimiter,
		AllSheets: xlsx2csvAllSheets,
	}

	if xlsx2csvAllSheets {
		outputPaths, err := skill.ConvertXlsxToCSVAllSheets(ctx, xlsx2csvInput, xlsx2csvOutput, opts)
		if err != nil {
			return fmt.Errorf("failed to convert XLSX: %w", err)
		}
		fmt.Printf("Exported %d sheets:\n", len(outputPaths))
		for _, p := range outputPaths {
			fmt.Printf("  %s\n", p)
		}
		return nil
	}

	if err := skill.ConvertXlsxToCSV(ctx, xlsx2csvInput, xlsx2csvOutput, opts); err != nil {
		return fmt.Errorf("failed to convert XLSX: %w", err)
	}

	fmt.Printf("Converted to: %s\n", xlsx2csvOutput)
	return nil
}

func runDocxReplace(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	if err := initUniDoc(); err != nil {
		return fmt.Errorf("failed to initialize UniDoc: %w", err)
	}

	skill := unidoc.New()
	if err := skill.Init(ctx); err != nil {
		return fmt.Errorf("failed to initialize skill: %w", err)
	}
	defer skill.Close()

	opts := uniofficeutil.ReplaceOptions{
		CaseSensitive: docxReplaceCaseSensitive,
		WholeWord:     docxReplaceWholeWord,
		ReplaceAll:    true,
	}

	count, err := skill.ReplaceInDocx(ctx, docxReplaceInput, docxReplaceOutput, docxReplaceFind, docxReplaceReplace, opts)
	if err != nil {
		return fmt.Errorf("failed to replace text: %w", err)
	}

	fmt.Printf("Replaced %d occurrences: %s\n", count, docxReplaceOutput)
	return nil
}
