// Package unidoc provides an omniskill Skill for PDF and Office document operations.
//
// This package provides document conversion tools using the UniDoc library:
//   - pdf2docx: Convert PDF to DOCX with optional table, image, and layout extraction
//   - docx2pdf: Convert DOCX to PDF
//   - pdf2png: Convert PDF page to PNG image
//   - metadata: Read Office document metadata
//
// The skill exposes both typed Go methods for direct use and MCP tools for
// AI agent integration. CLI tools should use the typed methods for compile-time
// type safety.
package unidoc

import (
	"context"
	"fmt"
	"log/slog"
	"path/filepath"
	"strings"
	"time"

	"github.com/grokify/gounidoc/unipdfutil"
	"github.com/grokify/gounidoc/uniofficeutil"
	"github.com/plexusone/omniskill/skill"
)

// Skill provides UniDoc document conversion tools.
type Skill struct {
	skill.BaseSkill
	logger *slog.Logger
}

// New creates a new UniDoc skill with default logger.
func New() *Skill {
	return &Skill{
		logger: slog.Default(),
	}
}

// NewWithLogger creates a new UniDoc skill with a custom logger.
func NewWithLogger(logger *slog.Logger) *Skill {
	if logger == nil {
		logger = slog.Default()
	}
	return &Skill{
		logger: logger,
	}
}

// SetLogger sets the logger for this skill.
func (s *Skill) SetLogger(logger *slog.Logger) {
	if logger != nil {
		s.logger = logger
	}
}

// Name returns the skill identifier.
func (s *Skill) Name() string {
	return "unidoc"
}

// Description returns what this skill does.
func (s *Skill) Description() string {
	return "Convert and process PDF and Office documents using UniDoc"
}

// Init initializes the skill.
func (s *Skill) Init(ctx context.Context) error {
	s.logger.Debug("initializing unidoc skill")
	return nil
}

// Close releases resources.
func (s *Skill) Close() error {
	s.logger.Debug("closing unidoc skill")
	return nil
}

// Tools returns all tools provided by this skill.
func (s *Skill) Tools() []skill.Tool {
	return []skill.Tool{
		s.pdf2docxTool(),
		s.docx2pdfTool(),
		s.pdf2pngTool(),
		s.metadataTool(),
		s.pdfMergeTool(),
		s.pdfSplitTool(),
		s.pdfExtractPagesTool(),
		s.pdfPageCountTool(),
		s.batchPdf2docxTool(),
		s.batchDocx2pdfTool(),
		s.xlsxInfoTool(),
		s.pptxInfoTool(),
		// Phase 1: PDF Text & Image Operations
		s.pdfTextTool(),
		s.pdfImagesTool(),
		s.pdf2imagesTool(),
		s.images2pdfTool(),
		// Phase 2: PDF Manipulation
		s.pdfWatermarkTool(),
		s.pdfRotateTool(),
		s.pdfCompressTool(),
		s.pdfMetadataTool(),
		// Phase 3: PDF Security
		s.pdfEncryptTool(),
		s.pdfDecryptTool(),
		// Phase 4: Office Document Operations
		s.docxTextTool(),
		s.xlsx2csvTool(),
		s.docxReplaceTool(),
	}
}

// Ensure Skill implements skill.Skill.
var _ skill.Skill = (*Skill)(nil)

// ---------------------------------------------------------------------------
// Typed methods for Go callers (CLI, library usage)
// ---------------------------------------------------------------------------

// ConvertPDFToDocx converts a PDF file to DOCX format with the specified options.
// If outputPath is empty, it defaults to the input filename with .docx extension.
func (s *Skill) ConvertPDFToDocx(ctx context.Context, inputPath, outputPath string, opts unipdfutil.ConversionOptions) (string, error) {
	if inputPath == "" {
		return "", fmt.Errorf("input_path is required")
	}
	if outputPath == "" {
		ext := filepath.Ext(inputPath)
		outputPath = strings.TrimSuffix(inputPath, ext) + ".docx"
	}

	s.logger.Info("converting PDF to DOCX",
		slog.String("input", inputPath),
		slog.String("output", outputPath),
		slog.Bool("extract_tables", opts.ExtractTables),
		slog.Bool("extract_images", opts.ExtractImages),
		slog.Bool("detect_layout", opts.DetectLayout),
	)

	start := time.Now()
	if err := unipdfutil.ConvertPDFFileToDocxFileWithOptions(inputPath, outputPath, opts); err != nil {
		s.logger.Error("PDF to DOCX conversion failed",
			slog.String("input", inputPath),
			slog.String("error", err.Error()),
			slog.Duration("duration", time.Since(start)),
		)
		return "", fmt.Errorf("failed to convert PDF to DOCX: %w", err)
	}

	s.logger.Info("PDF to DOCX conversion complete",
		slog.String("output", outputPath),
		slog.Duration("duration", time.Since(start)),
	)
	return outputPath, nil
}

// ConvertDocxToPDF converts a DOCX file to PDF format.
// If outputPath is empty, it defaults to the input filename with .pdf extension.
func (s *Skill) ConvertDocxToPDF(ctx context.Context, inputPath, outputPath string) (string, error) {
	if inputPath == "" {
		return "", fmt.Errorf("input_path is required")
	}
	if outputPath == "" {
		outputPath = inputPath + ".pdf"
	}

	s.logger.Info("converting DOCX to PDF",
		slog.String("input", inputPath),
		slog.String("output", outputPath),
	)

	start := time.Now()
	if err := uniofficeutil.ConvertDOCXFileToPDFFile(inputPath, outputPath); err != nil {
		s.logger.Error("DOCX to PDF conversion failed",
			slog.String("input", inputPath),
			slog.String("error", err.Error()),
			slog.Duration("duration", time.Since(start)),
		)
		return "", fmt.Errorf("failed to convert DOCX to PDF: %w", err)
	}

	s.logger.Info("DOCX to PDF conversion complete",
		slog.String("output", outputPath),
		slog.Duration("duration", time.Since(start)),
	)
	return outputPath, nil
}

// ConvertPDFPageToPNG converts a specific page from a PDF file to PNG format.
func (s *Skill) ConvertPDFPageToPNG(ctx context.Context, inputPath, outputPath string, pageNum, width uint32) error {
	if inputPath == "" {
		return fmt.Errorf("input_path is required")
	}
	if outputPath == "" {
		return fmt.Errorf("output_path is required")
	}
	if pageNum == 0 {
		pageNum = 1
	}
	if width == 0 {
		width = 1200
	}

	s.logger.Info("converting PDF page to PNG",
		slog.String("input", inputPath),
		slog.String("output", outputPath),
		slog.Uint64("page", uint64(pageNum)),
		slog.Uint64("width", uint64(width)),
	)

	start := time.Now()
	if err := unipdfutil.ConvertPDFFilePageToPNGFile(inputPath, outputPath, pageNum, width); err != nil {
		s.logger.Error("PDF to PNG conversion failed",
			slog.String("input", inputPath),
			slog.String("error", err.Error()),
			slog.Duration("duration", time.Since(start)),
		)
		return fmt.Errorf("failed to convert PDF to PNG: %w", err)
	}

	s.logger.Info("PDF to PNG conversion complete",
		slog.String("output", outputPath),
		slog.Duration("duration", time.Since(start)),
	)
	return nil
}

// ReadMetadata reads metadata from an Office document (DOCX, XLSX, PPTX).
func (s *Skill) ReadMetadata(ctx context.Context, filePath string) (any, error) {
	if filePath == "" {
		return nil, fmt.Errorf("file_path is required")
	}

	s.logger.Info("reading document metadata",
		slog.String("file", filePath),
	)

	start := time.Now()
	metadata, err := uniofficeutil.ReadFileMetadata(filePath)
	if err != nil {
		s.logger.Error("metadata read failed",
			slog.String("file", filePath),
			slog.String("error", err.Error()),
			slog.Duration("duration", time.Since(start)),
		)
		return nil, fmt.Errorf("failed to read metadata: %w", err)
	}

	s.logger.Info("metadata read complete",
		slog.String("file", filePath),
		slog.String("file_type", metadata.FileType),
		slog.Duration("duration", time.Since(start)),
	)
	return metadata, nil
}

// GetXlsxInfo returns information about an XLSX workbook including sheet names.
func (s *Skill) GetXlsxInfo(ctx context.Context, filePath string) (uniofficeutil.XlsxInfo, error) {
	if filePath == "" {
		return uniofficeutil.XlsxInfo{}, fmt.Errorf("file_path is required")
	}

	s.logger.Info("getting XLSX info",
		slog.String("file", filePath),
	)

	start := time.Now()
	info, err := uniofficeutil.GetXlsxInfo(filePath)
	if err != nil {
		s.logger.Error("failed to get XLSX info",
			slog.String("file", filePath),
			slog.String("error", err.Error()),
			slog.Duration("duration", time.Since(start)),
		)
		return uniofficeutil.XlsxInfo{}, fmt.Errorf("failed to get XLSX info: %w", err)
	}

	s.logger.Info("XLSX info retrieved",
		slog.String("file", filePath),
		slog.Int("sheet_count", info.SheetCount),
		slog.Duration("duration", time.Since(start)),
	)
	return info, nil
}

// GetPptxInfo returns information about a PPTX presentation.
func (s *Skill) GetPptxInfo(ctx context.Context, filePath string) (uniofficeutil.PptxInfo, error) {
	if filePath == "" {
		return uniofficeutil.PptxInfo{}, fmt.Errorf("file_path is required")
	}

	s.logger.Info("getting PPTX info",
		slog.String("file", filePath),
	)

	start := time.Now()
	info, err := uniofficeutil.GetPptxInfo(filePath)
	if err != nil {
		s.logger.Error("failed to get PPTX info",
			slog.String("file", filePath),
			slog.String("error", err.Error()),
			slog.Duration("duration", time.Since(start)),
		)
		return uniofficeutil.PptxInfo{}, fmt.Errorf("failed to get PPTX info: %w", err)
	}

	s.logger.Info("PPTX info retrieved",
		slog.String("file", filePath),
		slog.Int("slide_count", info.SlideCount),
		slog.Duration("duration", time.Since(start)),
	)
	return info, nil
}

// MergePDFs merges multiple PDF files into a single output file.
func (s *Skill) MergePDFs(ctx context.Context, inputPaths []string, outputPath string) error {
	if len(inputPaths) == 0 {
		return fmt.Errorf("no input files provided")
	}
	if outputPath == "" {
		return fmt.Errorf("output_path is required")
	}

	s.logger.Info("merging PDF files",
		slog.Int("input_count", len(inputPaths)),
		slog.String("output", outputPath),
	)

	start := time.Now()
	if err := unipdfutil.MergePDFFiles(inputPaths, outputPath); err != nil {
		s.logger.Error("PDF merge failed",
			slog.String("error", err.Error()),
			slog.Duration("duration", time.Since(start)),
		)
		return fmt.Errorf("failed to merge PDFs: %w", err)
	}

	s.logger.Info("PDF merge complete",
		slog.String("output", outputPath),
		slog.Int("merged_files", len(inputPaths)),
		slog.Duration("duration", time.Since(start)),
	)
	return nil
}

// SplitPDF splits a PDF file into individual pages.
// Returns the list of created file paths.
func (s *Skill) SplitPDF(ctx context.Context, inputPath, outputDir string) ([]string, error) {
	if inputPath == "" {
		return nil, fmt.Errorf("input_path is required")
	}

	s.logger.Info("splitting PDF file",
		slog.String("input", inputPath),
		slog.String("output_dir", outputDir),
	)

	start := time.Now()
	outputPaths, err := unipdfutil.SplitPDFFile(inputPath, outputDir)
	if err != nil {
		s.logger.Error("PDF split failed",
			slog.String("input", inputPath),
			slog.String("error", err.Error()),
			slog.Duration("duration", time.Since(start)),
		)
		return outputPaths, fmt.Errorf("failed to split PDF: %w", err)
	}

	s.logger.Info("PDF split complete",
		slog.Int("pages_created", len(outputPaths)),
		slog.Duration("duration", time.Since(start)),
	)
	return outputPaths, nil
}

// ExtractPDFPages extracts specific pages from a PDF file.
func (s *Skill) ExtractPDFPages(ctx context.Context, inputPath, outputPath string, pages []int) error {
	if inputPath == "" {
		return fmt.Errorf("input_path is required")
	}
	if outputPath == "" {
		return fmt.Errorf("output_path is required")
	}
	if len(pages) == 0 {
		return fmt.Errorf("no pages specified")
	}

	s.logger.Info("extracting PDF pages",
		slog.String("input", inputPath),
		slog.String("output", outputPath),
		slog.Any("pages", pages),
	)

	start := time.Now()
	if err := unipdfutil.ExtractPDFPages(inputPath, outputPath, pages); err != nil {
		s.logger.Error("PDF page extraction failed",
			slog.String("input", inputPath),
			slog.String("error", err.Error()),
			slog.Duration("duration", time.Since(start)),
		)
		return fmt.Errorf("failed to extract pages: %w", err)
	}

	s.logger.Info("PDF page extraction complete",
		slog.String("output", outputPath),
		slog.Int("pages_extracted", len(pages)),
		slog.Duration("duration", time.Since(start)),
	)
	return nil
}

// GetPDFPageCount returns the number of pages in a PDF file.
func (s *Skill) GetPDFPageCount(ctx context.Context, inputPath string) (int, error) {
	if inputPath == "" {
		return 0, fmt.Errorf("input_path is required")
	}

	s.logger.Debug("getting PDF page count",
		slog.String("input", inputPath),
	)

	count, err := unipdfutil.GetPDFPageCount(inputPath)
	if err != nil {
		s.logger.Error("failed to get page count",
			slog.String("input", inputPath),
			slog.String("error", err.Error()),
		)
		return 0, fmt.Errorf("failed to get page count: %w", err)
	}

	s.logger.Debug("PDF page count retrieved",
		slog.String("input", inputPath),
		slog.Int("pages", count),
	)
	return count, nil
}

// BatchConvertPDFToDocx converts multiple PDF files to DOCX format.
// If outputDir is empty, output files are created in the same directory as the inputs.
func (s *Skill) BatchConvertPDFToDocx(ctx context.Context, inputPaths []string, outputDir string, opts unipdfutil.ConversionOptions) unipdfutil.BatchConversionResult {
	return s.BatchConvertPDFToDocxWithProgress(ctx, inputPaths, outputDir, opts, nil)
}

// BatchConvertPDFToDocxWithProgress converts multiple PDF files to DOCX format with progress reporting.
// If outputDir is empty, output files are created in the same directory as the inputs.
func (s *Skill) BatchConvertPDFToDocxWithProgress(ctx context.Context, inputPaths []string, outputDir string, opts unipdfutil.ConversionOptions, progress unipdfutil.ProgressCallback) unipdfutil.BatchConversionResult {
	s.logger.Info("starting batch PDF to DOCX conversion",
		slog.Int("file_count", len(inputPaths)),
		slog.String("output_dir", outputDir),
	)

	start := time.Now()
	result := unipdfutil.BatchConvertPDFToDocxWithProgress(inputPaths, outputDir, opts, progress)

	s.logger.Info("batch PDF to DOCX conversion complete",
		slog.Int("total", result.TotalFiles),
		slog.Int("succeeded", result.Succeeded),
		slog.Int("failed", result.Failed),
		slog.Duration("duration", time.Since(start)),
	)
	return result
}

// BatchConvertDocxToPDF converts multiple DOCX files to PDF format.
// If outputDir is empty, output files are created in the same directory as the inputs.
func (s *Skill) BatchConvertDocxToPDF(ctx context.Context, inputPaths []string, outputDir string) unipdfutil.BatchConversionResult {
	return s.BatchConvertDocxToPDFWithProgress(ctx, inputPaths, outputDir, nil)
}

// BatchConvertDocxToPDFWithProgress converts multiple DOCX files to PDF format with progress reporting.
// If outputDir is empty, output files are created in the same directory as the inputs.
func (s *Skill) BatchConvertDocxToPDFWithProgress(ctx context.Context, inputPaths []string, outputDir string, progress unipdfutil.ProgressCallback) unipdfutil.BatchConversionResult {
	s.logger.Info("starting batch DOCX to PDF conversion",
		slog.Int("file_count", len(inputPaths)),
		slog.String("output_dir", outputDir),
	)

	start := time.Now()
	result := unipdfutil.BatchConvertDocxToPDFWithProgress(inputPaths, outputDir, uniofficeutil.ConvertDOCXFileToPDFFile, progress)

	s.logger.Info("batch DOCX to PDF conversion complete",
		slog.Int("total", result.TotalFiles),
		slog.Int("succeeded", result.Succeeded),
		slog.Int("failed", result.Failed),
		slog.Duration("duration", time.Since(start)),
	)
	return result
}

func (s *Skill) pdf2docxTool() skill.Tool {
	return skill.NewTool(
		"pdf2docx",
		"Convert a PDF file to DOCX format with optional table extraction, image detection, and layout-aware formatting",
		map[string]skill.Parameter{
			"input_path": {
				Type:        "string",
				Description: "Path to the input PDF file",
				Required:    true,
			},
			"output_path": {
				Type:        "string",
				Description: "Path for the output DOCX file (optional, defaults to input filename with .docx extension)",
				Required:    false,
			},
			"extract_tables": {
				Type:        "boolean",
				Description: "Extract and convert tables to DOCX tables (default: true)",
				Required:    false,
				Default:     true,
			},
			"extract_images": {
				Type:        "boolean",
				Description: "Detect images and add placeholder text (default: true)",
				Required:    false,
				Default:     true,
			},
			"detect_layout": {
				Type:        "boolean",
				Description: "Detect headings and text structure (default: true)",
				Required:    false,
				Default:     true,
			},
		},
		func(ctx context.Context, params map[string]any) (any, error) {
			inputPath, _ := params["input_path"].(string)
			outputPath, _ := params["output_path"].(string)

			opts := unipdfutil.DefaultConversionOptions()
			if extractTables, ok := params["extract_tables"].(bool); ok {
				opts.ExtractTables = extractTables
			}
			if extractImages, ok := params["extract_images"].(bool); ok {
				opts.ExtractImages = extractImages
			}
			if detectLayout, ok := params["detect_layout"].(bool); ok {
				opts.DetectLayout = detectLayout
			}

			finalOutput, err := s.ConvertPDFToDocx(ctx, inputPath, outputPath, opts)
			if err != nil {
				return nil, err
			}

			return map[string]any{
				"success":     true,
				"input_path":  inputPath,
				"output_path": finalOutput,
				"options": map[string]any{
					"extract_tables": opts.ExtractTables,
					"extract_images": opts.ExtractImages,
					"detect_layout":  opts.DetectLayout,
				},
			}, nil
		},
	)
}

func (s *Skill) docx2pdfTool() skill.Tool {
	return skill.NewTool(
		"docx2pdf",
		"Convert a DOCX file to PDF format",
		map[string]skill.Parameter{
			"input_path": {
				Type:        "string",
				Description: "Path to the input DOCX file",
				Required:    true,
			},
			"output_path": {
				Type:        "string",
				Description: "Path for the output PDF file (optional, defaults to input filename with .pdf extension)",
				Required:    false,
			},
		},
		func(ctx context.Context, params map[string]any) (any, error) {
			inputPath, _ := params["input_path"].(string)
			outputPath, _ := params["output_path"].(string)

			finalOutput, err := s.ConvertDocxToPDF(ctx, inputPath, outputPath)
			if err != nil {
				return nil, err
			}

			return map[string]any{
				"success":     true,
				"input_path":  inputPath,
				"output_path": finalOutput,
			}, nil
		},
	)
}

func (s *Skill) pdf2pngTool() skill.Tool {
	return skill.NewTool(
		"pdf2png",
		"Convert a PDF page to PNG image format",
		map[string]skill.Parameter{
			"input_path": {
				Type:        "string",
				Description: "Path to the input PDF file",
				Required:    true,
			},
			"output_path": {
				Type:        "string",
				Description: "Path for the output PNG file",
				Required:    true,
			},
			"page_number": {
				Type:        "integer",
				Description: "Page number to convert (1-indexed, default: 1)",
				Required:    false,
				Default:     1,
			},
			"width": {
				Type:        "integer",
				Description: "Output image width in pixels (default: 1200)",
				Required:    false,
				Default:     1200,
			},
		},
		func(ctx context.Context, params map[string]any) (any, error) {
			inputPath, _ := params["input_path"].(string)
			outputPath, _ := params["output_path"].(string)

			pageNum := uint32(1)
			if pn, ok := params["page_number"].(float64); ok && pn > 0 {
				pageNum = uint32(pn)
			}

			width := uint32(1200)
			if w, ok := params["width"].(float64); ok && w > 0 {
				width = uint32(w)
			}

			if err := s.ConvertPDFPageToPNG(ctx, inputPath, outputPath, pageNum, width); err != nil {
				return nil, err
			}

			return map[string]any{
				"success":     true,
				"input_path":  inputPath,
				"output_path": outputPath,
				"page_number": pageNum,
				"width":       width,
			}, nil
		},
	)
}

func (s *Skill) metadataTool() skill.Tool {
	return skill.NewTool(
		"metadata",
		"Read metadata from an Office document (DOCX, PPTX, XLSX)",
		map[string]skill.Parameter{
			"file_path": {
				Type:        "string",
				Description: "Path to the Office document file",
				Required:    true,
			},
		},
		func(ctx context.Context, params map[string]any) (any, error) {
			filePath, _ := params["file_path"].(string)

			metadata, err := s.ReadMetadata(ctx, filePath)
			if err != nil {
				return nil, err
			}

			return map[string]any{
				"success":   true,
				"file_path": filePath,
				"metadata":  metadata,
			}, nil
		},
	)
}

func (s *Skill) pdfMergeTool() skill.Tool {
	return skill.NewTool(
		"pdf_merge",
		"Merge multiple PDF files into a single output file",
		map[string]skill.Parameter{
			"input_paths": {
				Type:        "array",
				Description: "List of input PDF file paths to merge (in order)",
				Required:    true,
			},
			"output_path": {
				Type:        "string",
				Description: "Path for the merged output PDF file",
				Required:    true,
			},
		},
		func(ctx context.Context, params map[string]any) (any, error) {
			outputPath, _ := params["output_path"].(string)

			// Parse input_paths from array
			var inputPaths []string
			if paths, ok := params["input_paths"].([]any); ok {
				for _, p := range paths {
					if pathStr, ok := p.(string); ok {
						inputPaths = append(inputPaths, pathStr)
					}
				}
			}

			if err := s.MergePDFs(ctx, inputPaths, outputPath); err != nil {
				return nil, err
			}

			return map[string]any{
				"success":      true,
				"input_paths":  inputPaths,
				"output_path":  outputPath,
				"merged_count": len(inputPaths),
			}, nil
		},
	)
}

func (s *Skill) pdfSplitTool() skill.Tool {
	return skill.NewTool(
		"pdf_split",
		"Split a PDF file into individual pages",
		map[string]skill.Parameter{
			"input_path": {
				Type:        "string",
				Description: "Path to the input PDF file to split",
				Required:    true,
			},
			"output_dir": {
				Type:        "string",
				Description: "Directory for output files (optional, defaults to input file directory)",
				Required:    false,
			},
		},
		func(ctx context.Context, params map[string]any) (any, error) {
			inputPath, _ := params["input_path"].(string)
			outputDir, _ := params["output_dir"].(string)

			outputPaths, err := s.SplitPDF(ctx, inputPath, outputDir)
			if err != nil {
				return nil, err
			}

			return map[string]any{
				"success":      true,
				"input_path":   inputPath,
				"output_paths": outputPaths,
				"pages_count":  len(outputPaths),
			}, nil
		},
	)
}

func (s *Skill) pdfExtractPagesTool() skill.Tool {
	return skill.NewTool(
		"pdf_extract_pages",
		"Extract specific pages from a PDF file into a new PDF",
		map[string]skill.Parameter{
			"input_path": {
				Type:        "string",
				Description: "Path to the input PDF file",
				Required:    true,
			},
			"output_path": {
				Type:        "string",
				Description: "Path for the output PDF file",
				Required:    true,
			},
			"pages": {
				Type:        "array",
				Description: "List of page numbers to extract (1-indexed)",
				Required:    true,
			},
		},
		func(ctx context.Context, params map[string]any) (any, error) {
			inputPath, _ := params["input_path"].(string)
			outputPath, _ := params["output_path"].(string)

			// Parse pages from array
			var pages []int
			if pageList, ok := params["pages"].([]any); ok {
				for _, p := range pageList {
					if pageNum, ok := p.(float64); ok {
						pages = append(pages, int(pageNum))
					}
				}
			}

			if err := s.ExtractPDFPages(ctx, inputPath, outputPath, pages); err != nil {
				return nil, err
			}

			return map[string]any{
				"success":         true,
				"input_path":      inputPath,
				"output_path":     outputPath,
				"pages_extracted": pages,
			}, nil
		},
	)
}

func (s *Skill) pdfPageCountTool() skill.Tool {
	return skill.NewTool(
		"pdf_page_count",
		"Get the number of pages in a PDF file",
		map[string]skill.Parameter{
			"input_path": {
				Type:        "string",
				Description: "Path to the PDF file",
				Required:    true,
			},
		},
		func(ctx context.Context, params map[string]any) (any, error) {
			inputPath, _ := params["input_path"].(string)

			count, err := s.GetPDFPageCount(ctx, inputPath)
			if err != nil {
				return nil, err
			}

			return map[string]any{
				"success":    true,
				"input_path": inputPath,
				"page_count": count,
			}, nil
		},
	)
}

func (s *Skill) batchPdf2docxTool() skill.Tool {
	return skill.NewTool(
		"batch_pdf2docx",
		"Convert multiple PDF files to DOCX format in batch",
		map[string]skill.Parameter{
			"input_paths": {
				Type:        "array",
				Description: "List of input PDF file paths (supports glob patterns)",
				Required:    true,
			},
			"output_dir": {
				Type:        "string",
				Description: "Output directory for converted files (optional, defaults to same directory as input)",
				Required:    false,
			},
			"extract_tables": {
				Type:        "boolean",
				Description: "Extract and convert tables (default: true)",
				Required:    false,
				Default:     true,
			},
			"extract_images": {
				Type:        "boolean",
				Description: "Detect images (default: true)",
				Required:    false,
				Default:     true,
			},
			"detect_layout": {
				Type:        "boolean",
				Description: "Detect headings and text structure (default: true)",
				Required:    false,
				Default:     true,
			},
		},
		func(ctx context.Context, params map[string]any) (any, error) {
			outputDir, _ := params["output_dir"].(string)

			// Parse input_paths from array
			var inputPatterns []string
			if paths, ok := params["input_paths"].([]any); ok {
				for _, p := range paths {
					if pathStr, ok := p.(string); ok {
						inputPatterns = append(inputPatterns, pathStr)
					}
				}
			}

			// Expand glob patterns
			inputPaths, err := unipdfutil.ExpandGlobPatterns(inputPatterns)
			if err != nil {
				return nil, err
			}

			// Filter to PDF files only
			inputPaths = unipdfutil.FilterByExtension(inputPaths, ".pdf")

			opts := unipdfutil.DefaultConversionOptions()
			if extractTables, ok := params["extract_tables"].(bool); ok {
				opts.ExtractTables = extractTables
			}
			if extractImages, ok := params["extract_images"].(bool); ok {
				opts.ExtractImages = extractImages
			}
			if detectLayout, ok := params["detect_layout"].(bool); ok {
				opts.DetectLayout = detectLayout
			}

			result := s.BatchConvertPDFToDocx(ctx, inputPaths, outputDir, opts)

			return map[string]any{
				"success":    result.Failed == 0,
				"total":      result.TotalFiles,
				"succeeded":  result.Succeeded,
				"failed":     result.Failed,
				"results":    result.Results,
				"output_dir": outputDir,
			}, nil
		},
	)
}

func (s *Skill) batchDocx2pdfTool() skill.Tool {
	return skill.NewTool(
		"batch_docx2pdf",
		"Convert multiple DOCX files to PDF format in batch",
		map[string]skill.Parameter{
			"input_paths": {
				Type:        "array",
				Description: "List of input DOCX file paths (supports glob patterns)",
				Required:    true,
			},
			"output_dir": {
				Type:        "string",
				Description: "Output directory for converted files (optional, defaults to same directory as input)",
				Required:    false,
			},
		},
		func(ctx context.Context, params map[string]any) (any, error) {
			outputDir, _ := params["output_dir"].(string)

			// Parse input_paths from array
			var inputPatterns []string
			if paths, ok := params["input_paths"].([]any); ok {
				for _, p := range paths {
					if pathStr, ok := p.(string); ok {
						inputPatterns = append(inputPatterns, pathStr)
					}
				}
			}

			// Expand glob patterns
			inputPaths, err := unipdfutil.ExpandGlobPatterns(inputPatterns)
			if err != nil {
				return nil, err
			}

			// Filter to DOCX files only
			inputPaths = unipdfutil.FilterByExtension(inputPaths, ".docx")

			result := s.BatchConvertDocxToPDF(ctx, inputPaths, outputDir)

			return map[string]any{
				"success":    result.Failed == 0,
				"total":      result.TotalFiles,
				"succeeded":  result.Succeeded,
				"failed":     result.Failed,
				"results":    result.Results,
				"output_dir": outputDir,
			}, nil
		},
	)
}

func (s *Skill) xlsxInfoTool() skill.Tool {
	return skill.NewTool(
		"xlsx_info",
		"Get information about an XLSX workbook including metadata and sheet names",
		map[string]skill.Parameter{
			"file_path": {
				Type:        "string",
				Description: "Path to the XLSX file",
				Required:    true,
			},
		},
		func(ctx context.Context, params map[string]any) (any, error) {
			filePath, _ := params["file_path"].(string)

			info, err := s.GetXlsxInfo(ctx, filePath)
			if err != nil {
				return nil, err
			}

			return map[string]any{
				"success":     true,
				"file_path":   filePath,
				"metadata":    info.Metadata,
				"sheet_count": info.SheetCount,
				"sheet_names": info.SheetNames,
			}, nil
		},
	)
}

func (s *Skill) pptxInfoTool() skill.Tool {
	return skill.NewTool(
		"pptx_info",
		"Get information about a PPTX presentation including metadata and slide count",
		map[string]skill.Parameter{
			"file_path": {
				Type:        "string",
				Description: "Path to the PPTX file",
				Required:    true,
			},
		},
		func(ctx context.Context, params map[string]any) (any, error) {
			filePath, _ := params["file_path"].(string)

			info, err := s.GetPptxInfo(ctx, filePath)
			if err != nil {
				return nil, err
			}

			return map[string]any{
				"success":     true,
				"file_path":   filePath,
				"metadata":    info.Metadata,
				"slide_count": info.SlideCount,
			}, nil
		},
	)
}

// ---------------------------------------------------------------------------
// Phase 1: PDF Text & Image Operations - Typed Methods
// ---------------------------------------------------------------------------

// ExtractPDFText extracts all text from a PDF file.
func (s *Skill) ExtractPDFText(ctx context.Context, inputPath string) (string, error) {
	if inputPath == "" {
		return "", fmt.Errorf("input_path is required")
	}

	s.logger.Info("extracting text from PDF",
		slog.String("input", inputPath),
	)

	start := time.Now()
	text, err := unipdfutil.ExtractPDFText(inputPath)
	if err != nil {
		s.logger.Error("PDF text extraction failed",
			slog.String("input", inputPath),
			slog.String("error", err.Error()),
			slog.Duration("duration", time.Since(start)),
		)
		return "", fmt.Errorf("failed to extract text: %w", err)
	}

	s.logger.Info("PDF text extraction complete",
		slog.String("input", inputPath),
		slog.Int("text_length", len(text)),
		slog.Duration("duration", time.Since(start)),
	)
	return text, nil
}

// ExtractPDFImages extracts all images from a PDF file to an output directory.
func (s *Skill) ExtractPDFImages(ctx context.Context, inputPath, outputDir string) ([]unipdfutil.ExtractedImage, error) {
	if inputPath == "" {
		return nil, fmt.Errorf("input_path is required")
	}

	s.logger.Info("extracting images from PDF",
		slog.String("input", inputPath),
		slog.String("output_dir", outputDir),
	)

	start := time.Now()
	images, err := unipdfutil.ExtractPDFImages(inputPath, outputDir)
	if err != nil {
		s.logger.Error("PDF image extraction failed",
			slog.String("input", inputPath),
			slog.String("error", err.Error()),
			slog.Duration("duration", time.Since(start)),
		)
		return images, fmt.Errorf("failed to extract images: %w", err)
	}

	s.logger.Info("PDF image extraction complete",
		slog.String("input", inputPath),
		slog.Int("images_extracted", len(images)),
		slog.Duration("duration", time.Since(start)),
	)
	return images, nil
}

// RenderPDFToImages renders all pages of a PDF as images.
func (s *Skill) RenderPDFToImages(ctx context.Context, inputPath, outputDir string, opts unipdfutil.RenderOptions) ([]string, error) {
	if inputPath == "" {
		return nil, fmt.Errorf("input_path is required")
	}

	s.logger.Info("rendering PDF to images",
		slog.String("input", inputPath),
		slog.String("output_dir", outputDir),
		slog.String("format", opts.Format),
		slog.Float64("dpi", opts.DPI),
	)

	start := time.Now()
	outputPaths, err := unipdfutil.RenderPDFToImages(inputPath, outputDir, opts)
	if err != nil {
		s.logger.Error("PDF rendering failed",
			slog.String("input", inputPath),
			slog.String("error", err.Error()),
			slog.Duration("duration", time.Since(start)),
		)
		return outputPaths, fmt.Errorf("failed to render PDF: %w", err)
	}

	s.logger.Info("PDF rendering complete",
		slog.String("input", inputPath),
		slog.Int("pages_rendered", len(outputPaths)),
		slog.Duration("duration", time.Since(start)),
	)
	return outputPaths, nil
}

// ImagesToPDF combines multiple images into a single PDF file.
func (s *Skill) ImagesToPDF(ctx context.Context, imagePaths []string, outputPath string, opts unipdfutil.ImageToPDFOptions) error {
	if len(imagePaths) == 0 {
		return fmt.Errorf("no images provided")
	}
	if outputPath == "" {
		return fmt.Errorf("output_path is required")
	}

	s.logger.Info("converting images to PDF",
		slog.Int("image_count", len(imagePaths)),
		slog.String("output", outputPath),
		slog.String("page_size", opts.PageSize),
	)

	start := time.Now()
	if err := unipdfutil.ImagesToPDF(imagePaths, outputPath, opts); err != nil {
		s.logger.Error("images to PDF conversion failed",
			slog.String("error", err.Error()),
			slog.Duration("duration", time.Since(start)),
		)
		return fmt.Errorf("failed to convert images to PDF: %w", err)
	}

	s.logger.Info("images to PDF conversion complete",
		slog.String("output", outputPath),
		slog.Int("images_combined", len(imagePaths)),
		slog.Duration("duration", time.Since(start)),
	)
	return nil
}

// ---------------------------------------------------------------------------
// Phase 2: PDF Manipulation - Typed Methods
// ---------------------------------------------------------------------------

// AddWatermark adds a text or image watermark to a PDF.
func (s *Skill) AddWatermark(ctx context.Context, inputPath, outputPath string, opts unipdfutil.WatermarkOptions) error {
	if inputPath == "" {
		return fmt.Errorf("input_path is required")
	}
	if outputPath == "" {
		return fmt.Errorf("output_path is required")
	}

	s.logger.Info("adding watermark to PDF",
		slog.String("input", inputPath),
		slog.String("output", outputPath),
		slog.String("text", opts.Text),
		slog.String("position", opts.Position),
	)

	start := time.Now()
	if err := unipdfutil.AddWatermark(inputPath, outputPath, opts); err != nil {
		s.logger.Error("watermark addition failed",
			slog.String("input", inputPath),
			slog.String("error", err.Error()),
			slog.Duration("duration", time.Since(start)),
		)
		return fmt.Errorf("failed to add watermark: %w", err)
	}

	s.logger.Info("watermark addition complete",
		slog.String("output", outputPath),
		slog.Duration("duration", time.Since(start)),
	)
	return nil
}

// RotatePDF rotates pages in a PDF by the specified degrees.
func (s *Skill) RotatePDF(ctx context.Context, inputPath, outputPath string, degrees int, pages []int) error {
	if inputPath == "" {
		return fmt.Errorf("input_path is required")
	}
	if outputPath == "" {
		return fmt.Errorf("output_path is required")
	}

	s.logger.Info("rotating PDF pages",
		slog.String("input", inputPath),
		slog.String("output", outputPath),
		slog.Int("degrees", degrees),
		slog.Any("pages", pages),
	)

	start := time.Now()
	if err := unipdfutil.RotatePDFPages(inputPath, outputPath, degrees, pages); err != nil {
		s.logger.Error("PDF rotation failed",
			slog.String("input", inputPath),
			slog.String("error", err.Error()),
			slog.Duration("duration", time.Since(start)),
		)
		return fmt.Errorf("failed to rotate PDF: %w", err)
	}

	s.logger.Info("PDF rotation complete",
		slog.String("output", outputPath),
		slog.Duration("duration", time.Since(start)),
	)
	return nil
}

// CompressPDF compresses a PDF file to reduce its size.
func (s *Skill) CompressPDF(ctx context.Context, inputPath, outputPath string, opts unipdfutil.CompressOptions) (unipdfutil.CompressionResult, error) {
	if inputPath == "" {
		return unipdfutil.CompressionResult{}, fmt.Errorf("input_path is required")
	}
	if outputPath == "" {
		return unipdfutil.CompressionResult{}, fmt.Errorf("output_path is required")
	}

	s.logger.Info("compressing PDF",
		slog.String("input", inputPath),
		slog.String("output", outputPath),
		slog.Int("image_quality", opts.ImageQuality),
	)

	start := time.Now()
	result, err := unipdfutil.CompressPDFWithStats(inputPath, outputPath, opts)
	if err != nil {
		s.logger.Error("PDF compression failed",
			slog.String("input", inputPath),
			slog.String("error", err.Error()),
			slog.Duration("duration", time.Since(start)),
		)
		return result, fmt.Errorf("failed to compress PDF: %w", err)
	}

	s.logger.Info("PDF compression complete",
		slog.String("output", outputPath),
		slog.Int64("original_size", result.OriginalSize),
		slog.Int64("compressed_size", result.CompressedSize),
		slog.Float64("savings_percent", result.SavingsPercent),
		slog.Duration("duration", time.Since(start)),
	)
	return result, nil
}

// ReadPDFMetadata reads metadata from a PDF file.
func (s *Skill) ReadPDFMetadata(ctx context.Context, inputPath string) (unipdfutil.PDFMetadata, error) {
	if inputPath == "" {
		return unipdfutil.PDFMetadata{}, fmt.Errorf("input_path is required")
	}

	s.logger.Info("reading PDF metadata",
		slog.String("input", inputPath),
	)

	start := time.Now()
	meta, err := unipdfutil.ReadPDFMetadata(inputPath)
	if err != nil {
		s.logger.Error("PDF metadata read failed",
			slog.String("input", inputPath),
			slog.String("error", err.Error()),
			slog.Duration("duration", time.Since(start)),
		)
		return meta, fmt.Errorf("failed to read PDF metadata: %w", err)
	}

	s.logger.Info("PDF metadata read complete",
		slog.String("input", inputPath),
		slog.String("title", meta.Title),
		slog.Duration("duration", time.Since(start)),
	)
	return meta, nil
}

// WritePDFMetadata updates metadata in a PDF file.
func (s *Skill) WritePDFMetadata(ctx context.Context, inputPath, outputPath string, meta unipdfutil.PDFMetadata) error {
	if inputPath == "" {
		return fmt.Errorf("input_path is required")
	}
	if outputPath == "" {
		return fmt.Errorf("output_path is required")
	}

	s.logger.Info("writing PDF metadata",
		slog.String("input", inputPath),
		slog.String("output", outputPath),
	)

	start := time.Now()
	if err := unipdfutil.WritePDFMetadata(inputPath, outputPath, meta); err != nil {
		s.logger.Error("PDF metadata write failed",
			slog.String("input", inputPath),
			slog.String("error", err.Error()),
			slog.Duration("duration", time.Since(start)),
		)
		return fmt.Errorf("failed to write PDF metadata: %w", err)
	}

	s.logger.Info("PDF metadata write complete",
		slog.String("output", outputPath),
		slog.Duration("duration", time.Since(start)),
	)
	return nil
}

// ---------------------------------------------------------------------------
// Phase 3: PDF Security - Typed Methods
// ---------------------------------------------------------------------------

// EncryptPDF encrypts a PDF file with password protection.
func (s *Skill) EncryptPDF(ctx context.Context, inputPath, outputPath string, opts unipdfutil.EncryptOptions) error {
	if inputPath == "" {
		return fmt.Errorf("input_path is required")
	}
	if outputPath == "" {
		return fmt.Errorf("output_path is required")
	}

	s.logger.Info("encrypting PDF",
		slog.String("input", inputPath),
		slog.String("output", outputPath),
	)

	start := time.Now()
	if err := unipdfutil.EncryptPDF(inputPath, outputPath, opts); err != nil {
		s.logger.Error("PDF encryption failed",
			slog.String("input", inputPath),
			slog.String("error", err.Error()),
			slog.Duration("duration", time.Since(start)),
		)
		return fmt.Errorf("failed to encrypt PDF: %w", err)
	}

	s.logger.Info("PDF encryption complete",
		slog.String("output", outputPath),
		slog.Duration("duration", time.Since(start)),
	)
	return nil
}

// DecryptPDF removes password protection from a PDF file.
func (s *Skill) DecryptPDF(ctx context.Context, inputPath, outputPath, password string) error {
	if inputPath == "" {
		return fmt.Errorf("input_path is required")
	}
	if outputPath == "" {
		return fmt.Errorf("output_path is required")
	}

	s.logger.Info("decrypting PDF",
		slog.String("input", inputPath),
		slog.String("output", outputPath),
	)

	start := time.Now()
	if err := unipdfutil.DecryptPDF(inputPath, outputPath, password); err != nil {
		s.logger.Error("PDF decryption failed",
			slog.String("input", inputPath),
			slog.String("error", err.Error()),
			slog.Duration("duration", time.Since(start)),
		)
		return fmt.Errorf("failed to decrypt PDF: %w", err)
	}

	s.logger.Info("PDF decryption complete",
		slog.String("output", outputPath),
		slog.Duration("duration", time.Since(start)),
	)
	return nil
}

// ---------------------------------------------------------------------------
// Phase 4: Office Document Operations - Typed Methods
// ---------------------------------------------------------------------------

// ExtractDocxText extracts all text from a DOCX file.
func (s *Skill) ExtractDocxText(ctx context.Context, inputPath string) (string, error) {
	if inputPath == "" {
		return "", fmt.Errorf("input_path is required")
	}

	s.logger.Info("extracting text from DOCX",
		slog.String("input", inputPath),
	)

	start := time.Now()
	text, err := uniofficeutil.ExtractDocxText(inputPath)
	if err != nil {
		s.logger.Error("DOCX text extraction failed",
			slog.String("input", inputPath),
			slog.String("error", err.Error()),
			slog.Duration("duration", time.Since(start)),
		)
		return "", fmt.Errorf("failed to extract text: %w", err)
	}

	s.logger.Info("DOCX text extraction complete",
		slog.String("input", inputPath),
		slog.Int("text_length", len(text)),
		slog.Duration("duration", time.Since(start)),
	)
	return text, nil
}

// ConvertXlsxToCSV exports a spreadsheet to CSV format.
func (s *Skill) ConvertXlsxToCSV(ctx context.Context, inputPath, outputPath string, opts uniofficeutil.CSVOptions) error {
	if inputPath == "" {
		return fmt.Errorf("input_path is required")
	}
	if outputPath == "" {
		return fmt.Errorf("output_path is required")
	}

	s.logger.Info("converting XLSX to CSV",
		slog.String("input", inputPath),
		slog.String("output", outputPath),
		slog.String("sheet", opts.SheetName),
	)

	start := time.Now()
	if err := uniofficeutil.ConvertXlsxToCSV(inputPath, outputPath, opts); err != nil {
		s.logger.Error("XLSX to CSV conversion failed",
			slog.String("input", inputPath),
			slog.String("error", err.Error()),
			slog.Duration("duration", time.Since(start)),
		)
		return fmt.Errorf("failed to convert XLSX to CSV: %w", err)
	}

	s.logger.Info("XLSX to CSV conversion complete",
		slog.String("output", outputPath),
		slog.Duration("duration", time.Since(start)),
	)
	return nil
}

// ConvertXlsxToCSVAllSheets exports all sheets to separate CSV files.
func (s *Skill) ConvertXlsxToCSVAllSheets(ctx context.Context, inputPath, outputDir string, opts uniofficeutil.CSVOptions) ([]string, error) {
	if inputPath == "" {
		return nil, fmt.Errorf("input_path is required")
	}

	s.logger.Info("converting XLSX to CSV (all sheets)",
		slog.String("input", inputPath),
		slog.String("output_dir", outputDir),
	)

	start := time.Now()
	outputPaths, err := uniofficeutil.ConvertXlsxToCSVAllSheets(inputPath, outputDir, opts)
	if err != nil {
		s.logger.Error("XLSX to CSV conversion failed",
			slog.String("input", inputPath),
			slog.String("error", err.Error()),
			slog.Duration("duration", time.Since(start)),
		)
		return outputPaths, fmt.Errorf("failed to convert XLSX to CSV: %w", err)
	}

	s.logger.Info("XLSX to CSV conversion complete",
		slog.Int("sheets_exported", len(outputPaths)),
		slog.Duration("duration", time.Since(start)),
	)
	return outputPaths, nil
}

// ReplaceInDocx performs find/replace in a DOCX file.
func (s *Skill) ReplaceInDocx(ctx context.Context, inputPath, outputPath, find, replace string, opts uniofficeutil.ReplaceOptions) (int, error) {
	if inputPath == "" {
		return 0, fmt.Errorf("input_path is required")
	}
	if outputPath == "" {
		return 0, fmt.Errorf("output_path is required")
	}
	if find == "" {
		return 0, fmt.Errorf("find string is required")
	}

	s.logger.Info("replacing text in DOCX",
		slog.String("input", inputPath),
		slog.String("output", outputPath),
		slog.String("find", find),
	)

	start := time.Now()
	count, err := uniofficeutil.ReplaceInDocx(inputPath, outputPath, find, replace, opts)
	if err != nil {
		s.logger.Error("DOCX replace failed",
			slog.String("input", inputPath),
			slog.String("error", err.Error()),
			slog.Duration("duration", time.Since(start)),
		)
		return count, fmt.Errorf("failed to replace text: %w", err)
	}

	s.logger.Info("DOCX replace complete",
		slog.String("output", outputPath),
		slog.Int("replacements", count),
		slog.Duration("duration", time.Since(start)),
	)
	return count, nil
}

// ---------------------------------------------------------------------------
// Phase 1: PDF Text & Image Operations - MCP Tools
// ---------------------------------------------------------------------------

func (s *Skill) pdfTextTool() skill.Tool {
	return skill.NewTool(
		"pdf_text",
		"Extract all text from a PDF file",
		map[string]skill.Parameter{
			"input_path": {
				Type:        "string",
				Description: "Path to the input PDF file",
				Required:    true,
			},
			"output_path": {
				Type:        "string",
				Description: "Path to save extracted text (optional, if not provided text is returned directly)",
				Required:    false,
			},
			"page": {
				Type:        "integer",
				Description: "Specific page to extract (optional, 1-indexed). If not specified, extracts all pages.",
				Required:    false,
			},
		},
		func(ctx context.Context, params map[string]any) (any, error) {
			inputPath, _ := params["input_path"].(string)
			outputPath, _ := params["output_path"].(string)

			text, err := s.ExtractPDFText(ctx, inputPath)
			if err != nil {
				return nil, err
			}

			if outputPath != "" {
				if err := unipdfutil.ExtractPDFTextToFile(inputPath, outputPath); err != nil {
					return nil, err
				}
				return map[string]any{
					"success":     true,
					"input_path":  inputPath,
					"output_path": outputPath,
					"text_length": len(text),
				}, nil
			}

			return map[string]any{
				"success":     true,
				"input_path":  inputPath,
				"text":        text,
				"text_length": len(text),
			}, nil
		},
	)
}

func (s *Skill) pdfImagesTool() skill.Tool {
	return skill.NewTool(
		"pdf_images",
		"Extract all images from a PDF file",
		map[string]skill.Parameter{
			"input_path": {
				Type:        "string",
				Description: "Path to the input PDF file",
				Required:    true,
			},
			"output_dir": {
				Type:        "string",
				Description: "Directory to save extracted images (optional, defaults to input file directory)",
				Required:    false,
			},
		},
		func(ctx context.Context, params map[string]any) (any, error) {
			inputPath, _ := params["input_path"].(string)
			outputDir, _ := params["output_dir"].(string)

			images, err := s.ExtractPDFImages(ctx, inputPath, outputDir)
			if err != nil {
				return nil, err
			}

			return map[string]any{
				"success":     true,
				"input_path":  inputPath,
				"output_dir":  outputDir,
				"images":      images,
				"image_count": len(images),
			}, nil
		},
	)
}

func (s *Skill) pdf2imagesTool() skill.Tool {
	return skill.NewTool(
		"pdf2images",
		"Render all pages of a PDF as images",
		map[string]skill.Parameter{
			"input_path": {
				Type:        "string",
				Description: "Path to the input PDF file",
				Required:    true,
			},
			"output_dir": {
				Type:        "string",
				Description: "Directory to save rendered images (optional, defaults to input file directory)",
				Required:    false,
			},
			"format": {
				Type:        "string",
				Description: "Output image format: png or jpeg (default: png)",
				Required:    false,
				Default:     "png",
			},
			"dpi": {
				Type:        "integer",
				Description: "Resolution in DPI (default: 150)",
				Required:    false,
				Default:     150,
			},
			"quality": {
				Type:        "integer",
				Description: "JPEG quality 1-100 (default: 85)",
				Required:    false,
				Default:     85,
			},
		},
		func(ctx context.Context, params map[string]any) (any, error) {
			inputPath, _ := params["input_path"].(string)
			outputDir, _ := params["output_dir"].(string)

			opts := unipdfutil.DefaultRenderOptions()
			if format, ok := params["format"].(string); ok && format != "" {
				opts.Format = format
			}
			if dpi, ok := params["dpi"].(float64); ok && dpi > 0 {
				opts.DPI = dpi
			}
			if quality, ok := params["quality"].(float64); ok && quality > 0 {
				opts.Quality = int(quality)
			}

			outputPaths, err := s.RenderPDFToImages(ctx, inputPath, outputDir, opts)
			if err != nil {
				return nil, err
			}

			return map[string]any{
				"success":      true,
				"input_path":   inputPath,
				"output_dir":   outputDir,
				"output_paths": outputPaths,
				"pages_count":  len(outputPaths),
				"format":       opts.Format,
				"dpi":          opts.DPI,
			}, nil
		},
	)
}

func (s *Skill) images2pdfTool() skill.Tool {
	return skill.NewTool(
		"images2pdf",
		"Combine multiple images into a single PDF file",
		map[string]skill.Parameter{
			"input_paths": {
				Type:        "array",
				Description: "List of image file paths to combine",
				Required:    true,
			},
			"output_path": {
				Type:        "string",
				Description: "Path for the output PDF file",
				Required:    true,
			},
			"page_size": {
				Type:        "string",
				Description: "Page size: letter, a4, or fit (default: letter)",
				Required:    false,
				Default:     "letter",
			},
			"margin": {
				Type:        "number",
				Description: "Margin in points (default: 36)",
				Required:    false,
				Default:     36.0,
			},
		},
		func(ctx context.Context, params map[string]any) (any, error) {
			outputPath, _ := params["output_path"].(string)

			var imagePaths []string
			if paths, ok := params["input_paths"].([]any); ok {
				for _, p := range paths {
					if pathStr, ok := p.(string); ok {
						imagePaths = append(imagePaths, pathStr)
					}
				}
			}

			opts := unipdfutil.DefaultImageToPDFOptions()
			if pageSize, ok := params["page_size"].(string); ok && pageSize != "" {
				opts.PageSize = pageSize
			}
			if margin, ok := params["margin"].(float64); ok {
				opts.Margin = margin
			}

			if err := s.ImagesToPDF(ctx, imagePaths, outputPath, opts); err != nil {
				return nil, err
			}

			return map[string]any{
				"success":      true,
				"input_paths":  imagePaths,
				"output_path":  outputPath,
				"image_count":  len(imagePaths),
				"page_size":    opts.PageSize,
			}, nil
		},
	)
}

// ---------------------------------------------------------------------------
// Phase 2: PDF Manipulation - MCP Tools
// ---------------------------------------------------------------------------

func (s *Skill) pdfWatermarkTool() skill.Tool {
	return skill.NewTool(
		"pdf_watermark",
		"Add a text or image watermark to a PDF",
		map[string]skill.Parameter{
			"input_path": {
				Type:        "string",
				Description: "Path to the input PDF file",
				Required:    true,
			},
			"output_path": {
				Type:        "string",
				Description: "Path for the output PDF file",
				Required:    true,
			},
			"text": {
				Type:        "string",
				Description: "Watermark text (required if image_path not provided)",
				Required:    false,
			},
			"image_path": {
				Type:        "string",
				Description: "Path to watermark image (alternative to text)",
				Required:    false,
			},
			"font_size": {
				Type:        "number",
				Description: "Font size for text watermark (default: 48)",
				Required:    false,
				Default:     48.0,
			},
			"color": {
				Type:        "string",
				Description: "Hex color for text watermark (default: #888888)",
				Required:    false,
				Default:     "#888888",
			},
			"opacity": {
				Type:        "number",
				Description: "Watermark opacity 0.0-1.0 (default: 0.3)",
				Required:    false,
				Default:     0.3,
			},
			"rotation": {
				Type:        "number",
				Description: "Rotation angle in degrees (default: -45)",
				Required:    false,
				Default:     -45.0,
			},
			"position": {
				Type:        "string",
				Description: "Position: center, top-left, top-right, bottom-left, bottom-right (default: center)",
				Required:    false,
				Default:     "center",
			},
		},
		func(ctx context.Context, params map[string]any) (any, error) {
			inputPath, _ := params["input_path"].(string)
			outputPath, _ := params["output_path"].(string)

			opts := unipdfutil.DefaultWatermarkOptions()
			if text, ok := params["text"].(string); ok {
				opts.Text = text
			}
			if imagePath, ok := params["image_path"].(string); ok {
				opts.ImagePath = imagePath
			}
			if fontSize, ok := params["font_size"].(float64); ok {
				opts.FontSize = fontSize
			}
			if color, ok := params["color"].(string); ok {
				opts.Color = color
			}
			if opacity, ok := params["opacity"].(float64); ok {
				opts.Opacity = opacity
			}
			if rotation, ok := params["rotation"].(float64); ok {
				opts.Rotation = rotation
			}
			if position, ok := params["position"].(string); ok {
				opts.Position = position
			}

			if err := s.AddWatermark(ctx, inputPath, outputPath, opts); err != nil {
				return nil, err
			}

			return map[string]any{
				"success":     true,
				"input_path":  inputPath,
				"output_path": outputPath,
				"watermark":   opts.Text,
			}, nil
		},
	)
}

func (s *Skill) pdfRotateTool() skill.Tool {
	return skill.NewTool(
		"pdf_rotate",
		"Rotate pages in a PDF file",
		map[string]skill.Parameter{
			"input_path": {
				Type:        "string",
				Description: "Path to the input PDF file",
				Required:    true,
			},
			"output_path": {
				Type:        "string",
				Description: "Path for the output PDF file",
				Required:    true,
			},
			"degrees": {
				Type:        "integer",
				Description: "Rotation angle: 90, 180, or 270 degrees",
				Required:    true,
			},
			"pages": {
				Type:        "array",
				Description: "List of page numbers to rotate (1-indexed). If not specified, all pages are rotated.",
				Required:    false,
			},
		},
		func(ctx context.Context, params map[string]any) (any, error) {
			inputPath, _ := params["input_path"].(string)
			outputPath, _ := params["output_path"].(string)
			degrees := 0
			if d, ok := params["degrees"].(float64); ok {
				degrees = int(d)
			}

			var pages []int
			if pageList, ok := params["pages"].([]any); ok {
				for _, p := range pageList {
					if pageNum, ok := p.(float64); ok {
						pages = append(pages, int(pageNum))
					}
				}
			}

			if err := s.RotatePDF(ctx, inputPath, outputPath, degrees, pages); err != nil {
				return nil, err
			}

			return map[string]any{
				"success":     true,
				"input_path":  inputPath,
				"output_path": outputPath,
				"degrees":     degrees,
				"pages":       pages,
			}, nil
		},
	)
}

func (s *Skill) pdfCompressTool() skill.Tool {
	return skill.NewTool(
		"pdf_compress",
		"Compress a PDF file to reduce its size",
		map[string]skill.Parameter{
			"input_path": {
				Type:        "string",
				Description: "Path to the input PDF file",
				Required:    true,
			},
			"output_path": {
				Type:        "string",
				Description: "Path for the output PDF file",
				Required:    true,
			},
			"quality": {
				Type:        "integer",
				Description: "Image quality 1-100 (lower = more compression, default: 80)",
				Required:    false,
				Default:     80,
			},
		},
		func(ctx context.Context, params map[string]any) (any, error) {
			inputPath, _ := params["input_path"].(string)
			outputPath, _ := params["output_path"].(string)

			opts := unipdfutil.DefaultCompressOptions()
			if quality, ok := params["quality"].(float64); ok {
				opts.ImageQuality = int(quality)
			}

			result, err := s.CompressPDF(ctx, inputPath, outputPath, opts)
			if err != nil {
				return nil, err
			}

			return map[string]any{
				"success":         true,
				"input_path":      inputPath,
				"output_path":     outputPath,
				"original_size":   result.OriginalSize,
				"compressed_size": result.CompressedSize,
				"savings_percent": result.SavingsPercent,
			}, nil
		},
	)
}

func (s *Skill) pdfMetadataTool() skill.Tool {
	return skill.NewTool(
		"pdf_metadata",
		"Read or write PDF document metadata",
		map[string]skill.Parameter{
			"input_path": {
				Type:        "string",
				Description: "Path to the input PDF file",
				Required:    true,
			},
			"output_path": {
				Type:        "string",
				Description: "Path for output PDF (required when writing metadata)",
				Required:    false,
			},
			"title": {
				Type:        "string",
				Description: "Document title (for writing)",
				Required:    false,
			},
			"author": {
				Type:        "string",
				Description: "Document author (for writing)",
				Required:    false,
			},
			"subject": {
				Type:        "string",
				Description: "Document subject (for writing)",
				Required:    false,
			},
			"keywords": {
				Type:        "string",
				Description: "Document keywords (for writing)",
				Required:    false,
			},
		},
		func(ctx context.Context, params map[string]any) (any, error) {
			inputPath, _ := params["input_path"].(string)
			outputPath, _ := params["output_path"].(string)
			title, _ := params["title"].(string)
			author, _ := params["author"].(string)
			subject, _ := params["subject"].(string)
			keywords, _ := params["keywords"].(string)

			// If no output_path and no metadata fields, just read
			if outputPath == "" && title == "" && author == "" && subject == "" && keywords == "" {
				meta, err := s.ReadPDFMetadata(ctx, inputPath)
				if err != nil {
					return nil, err
				}

				return map[string]any{
					"success":    true,
					"input_path": inputPath,
					"metadata":   meta,
				}, nil
			}

			// Writing metadata
			if outputPath == "" {
				return nil, fmt.Errorf("output_path is required when writing metadata")
			}

			meta := unipdfutil.PDFMetadata{
				Title:    title,
				Author:   author,
				Subject:  subject,
				Keywords: keywords,
			}

			if err := s.WritePDFMetadata(ctx, inputPath, outputPath, meta); err != nil {
				return nil, err
			}

			return map[string]any{
				"success":     true,
				"input_path":  inputPath,
				"output_path": outputPath,
				"metadata":    meta,
			}, nil
		},
	)
}

// ---------------------------------------------------------------------------
// Phase 3: PDF Security - MCP Tools
// ---------------------------------------------------------------------------

func (s *Skill) pdfEncryptTool() skill.Tool {
	return skill.NewTool(
		"pdf_encrypt",
		"Encrypt a PDF file with password protection",
		map[string]skill.Parameter{
			"input_path": {
				Type:        "string",
				Description: "Path to the input PDF file",
				Required:    true,
			},
			"output_path": {
				Type:        "string",
				Description: "Path for the output encrypted PDF file",
				Required:    true,
			},
			"user_password": {
				Type:        "string",
				Description: "Password required to open the PDF",
				Required:    false,
			},
			"owner_password": {
				Type:        "string",
				Description: "Password for full access (modifying, printing, etc.)",
				Required:    false,
			},
			"allow_printing": {
				Type:        "boolean",
				Description: "Allow printing (default: true)",
				Required:    false,
				Default:     true,
			},
			"allow_copying": {
				Type:        "boolean",
				Description: "Allow copying content (default: true)",
				Required:    false,
				Default:     true,
			},
			"allow_modifying": {
				Type:        "boolean",
				Description: "Allow modifying content (default: true)",
				Required:    false,
				Default:     true,
			},
		},
		func(ctx context.Context, params map[string]any) (any, error) {
			inputPath, _ := params["input_path"].(string)
			outputPath, _ := params["output_path"].(string)

			opts := unipdfutil.EncryptOptions{
				Permissions: unipdfutil.DefaultPDFPermissions(),
			}
			if userPwd, ok := params["user_password"].(string); ok {
				opts.UserPassword = userPwd
			}
			if ownerPwd, ok := params["owner_password"].(string); ok {
				opts.OwnerPassword = ownerPwd
			}
			if allowPrint, ok := params["allow_printing"].(bool); ok {
				opts.Permissions.Printing = allowPrint
			}
			if allowCopy, ok := params["allow_copying"].(bool); ok {
				opts.Permissions.CopyContents = allowCopy
			}
			if allowModify, ok := params["allow_modifying"].(bool); ok {
				opts.Permissions.ModifyContents = allowModify
			}

			if err := s.EncryptPDF(ctx, inputPath, outputPath, opts); err != nil {
				return nil, err
			}

			return map[string]any{
				"success":     true,
				"input_path":  inputPath,
				"output_path": outputPath,
				"encrypted":   true,
			}, nil
		},
	)
}

func (s *Skill) pdfDecryptTool() skill.Tool {
	return skill.NewTool(
		"pdf_decrypt",
		"Remove password protection from a PDF file",
		map[string]skill.Parameter{
			"input_path": {
				Type:        "string",
				Description: "Path to the encrypted PDF file",
				Required:    true,
			},
			"output_path": {
				Type:        "string",
				Description: "Path for the output decrypted PDF file",
				Required:    true,
			},
			"password": {
				Type:        "string",
				Description: "Password to decrypt the PDF",
				Required:    true,
			},
		},
		func(ctx context.Context, params map[string]any) (any, error) {
			inputPath, _ := params["input_path"].(string)
			outputPath, _ := params["output_path"].(string)
			password, _ := params["password"].(string)

			if err := s.DecryptPDF(ctx, inputPath, outputPath, password); err != nil {
				return nil, err
			}

			return map[string]any{
				"success":     true,
				"input_path":  inputPath,
				"output_path": outputPath,
				"decrypted":   true,
			}, nil
		},
	)
}

// ---------------------------------------------------------------------------
// Phase 4: Office Document Operations - MCP Tools
// ---------------------------------------------------------------------------

func (s *Skill) docxTextTool() skill.Tool {
	return skill.NewTool(
		"docx_text",
		"Extract all text from a DOCX file",
		map[string]skill.Parameter{
			"input_path": {
				Type:        "string",
				Description: "Path to the input DOCX file",
				Required:    true,
			},
			"output_path": {
				Type:        "string",
				Description: "Path to save extracted text (optional)",
				Required:    false,
			},
		},
		func(ctx context.Context, params map[string]any) (any, error) {
			inputPath, _ := params["input_path"].(string)
			outputPath, _ := params["output_path"].(string)

			text, err := s.ExtractDocxText(ctx, inputPath)
			if err != nil {
				return nil, err
			}

			if outputPath != "" {
				if err := uniofficeutil.ExtractDocxTextToFile(inputPath, outputPath); err != nil {
					return nil, err
				}
				return map[string]any{
					"success":     true,
					"input_path":  inputPath,
					"output_path": outputPath,
					"text_length": len(text),
				}, nil
			}

			return map[string]any{
				"success":     true,
				"input_path":  inputPath,
				"text":        text,
				"text_length": len(text),
			}, nil
		},
	)
}

func (s *Skill) xlsx2csvTool() skill.Tool {
	return skill.NewTool(
		"xlsx2csv",
		"Convert an XLSX spreadsheet to CSV format",
		map[string]skill.Parameter{
			"input_path": {
				Type:        "string",
				Description: "Path to the input XLSX file",
				Required:    true,
			},
			"output_path": {
				Type:        "string",
				Description: "Path for the output CSV file (or directory for all sheets)",
				Required:    true,
			},
			"sheet": {
				Type:        "string",
				Description: "Sheet name to export (optional, defaults to first sheet)",
				Required:    false,
			},
			"delimiter": {
				Type:        "string",
				Description: "CSV delimiter character (default: comma)",
				Required:    false,
				Default:     ",",
			},
			"all_sheets": {
				Type:        "boolean",
				Description: "Export all sheets to separate CSV files (default: false)",
				Required:    false,
				Default:     false,
			},
		},
		func(ctx context.Context, params map[string]any) (any, error) {
			inputPath, _ := params["input_path"].(string)
			outputPath, _ := params["output_path"].(string)

			opts := uniofficeutil.DefaultCSVOptions()
			if sheet, ok := params["sheet"].(string); ok {
				opts.SheetName = sheet
			}
			if delimiter, ok := params["delimiter"].(string); ok && delimiter != "" {
				opts.Delimiter = delimiter
			}
			allSheets := false
			if all, ok := params["all_sheets"].(bool); ok {
				allSheets = all
			}

			if allSheets {
				outputPaths, err := s.ConvertXlsxToCSVAllSheets(ctx, inputPath, outputPath, opts)
				if err != nil {
					return nil, err
				}
				return map[string]any{
					"success":      true,
					"input_path":   inputPath,
					"output_paths": outputPaths,
					"sheets_count": len(outputPaths),
				}, nil
			}

			if err := s.ConvertXlsxToCSV(ctx, inputPath, outputPath, opts); err != nil {
				return nil, err
			}

			return map[string]any{
				"success":     true,
				"input_path":  inputPath,
				"output_path": outputPath,
				"sheet":       opts.SheetName,
			}, nil
		},
	)
}

func (s *Skill) docxReplaceTool() skill.Tool {
	return skill.NewTool(
		"docx_replace",
		"Find and replace text in a DOCX file",
		map[string]skill.Parameter{
			"input_path": {
				Type:        "string",
				Description: "Path to the input DOCX file",
				Required:    true,
			},
			"output_path": {
				Type:        "string",
				Description: "Path for the output DOCX file",
				Required:    true,
			},
			"find": {
				Type:        "string",
				Description: "Text to find",
				Required:    true,
			},
			"replace": {
				Type:        "string",
				Description: "Replacement text",
				Required:    true,
			},
			"case_sensitive": {
				Type:        "boolean",
				Description: "Case-sensitive search (default: true)",
				Required:    false,
				Default:     true,
			},
			"whole_word": {
				Type:        "boolean",
				Description: "Match whole words only (default: false)",
				Required:    false,
				Default:     false,
			},
		},
		func(ctx context.Context, params map[string]any) (any, error) {
			inputPath, _ := params["input_path"].(string)
			outputPath, _ := params["output_path"].(string)
			find, _ := params["find"].(string)
			replace, _ := params["replace"].(string)

			opts := uniofficeutil.DefaultReplaceOptions()
			if caseSensitive, ok := params["case_sensitive"].(bool); ok {
				opts.CaseSensitive = caseSensitive
			}
			if wholeWord, ok := params["whole_word"].(bool); ok {
				opts.WholeWord = wholeWord
			}

			count, err := s.ReplaceInDocx(ctx, inputPath, outputPath, find, replace, opts)
			if err != nil {
				return nil, err
			}

			return map[string]any{
				"success":      true,
				"input_path":   inputPath,
				"output_path":  outputPath,
				"find":         find,
				"replace":      replace,
				"replacements": count,
			}, nil
		},
	)
}
