package unidoc

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/grokify/gounidoc"
	"github.com/grokify/gounidoc/unipdfutil"
)

func init() {
	// Initialize UniDoc license from environment
	_ = gounidoc.SetMeteredKeyEnv()
}

func TestSkillInterface(t *testing.T) {
	s := New()

	if s.Name() != "unidoc" {
		t.Errorf("expected name 'unidoc', got '%s'", s.Name())
	}

	if s.Description() == "" {
		t.Error("expected non-empty description")
	}

	tools := s.Tools()
	if len(tools) != 25 {
		t.Errorf("expected 25 tools, got %d", len(tools))
	}

	expectedTools := map[string]bool{
		"pdf2docx":          false,
		"docx2pdf":          false,
		"pdf2png":           false,
		"metadata":          false,
		"pdf_merge":         false,
		"pdf_split":         false,
		"pdf_extract_pages": false,
		"pdf_page_count":    false,
		"batch_pdf2docx":    false,
		"batch_docx2pdf":    false,
		"xlsx_info":         false,
		"pptx_info":         false,
		// Phase 1: PDF Text & Image Operations
		"pdf_text":   false,
		"pdf_images": false,
		"pdf2images": false,
		"images2pdf": false,
		// Phase 2: PDF Manipulation
		"pdf_watermark": false,
		"pdf_rotate":    false,
		"pdf_compress":  false,
		"pdf_metadata":  false,
		// Phase 3: PDF Security
		"pdf_encrypt": false,
		"pdf_decrypt": false,
		// Phase 4: Office Document Operations
		"docx_text":    false,
		"xlsx2csv":     false,
		"docx_replace": false,
	}

	for _, tool := range tools {
		if _, ok := expectedTools[tool.Name()]; ok {
			expectedTools[tool.Name()] = true
		} else {
			t.Errorf("unexpected tool: %s", tool.Name())
		}
	}

	for name, found := range expectedTools {
		if !found {
			t.Errorf("missing tool: %s", name)
		}
	}
}

func TestConvertPDFToDocx(t *testing.T) {
	if os.Getenv("UNIDOC_KEY") == "" {
		t.Skip("UNIDOC_KEY not set, skipping integration test")
	}

	ctx := context.Background()
	s := New()
	if err := s.Init(ctx); err != nil {
		t.Fatalf("failed to init skill: %v", err)
	}
	defer s.Close()

	inputPDF := filepath.Join("..", "..", "testdata", "lorem_ipsum.pdf")
	if _, err := os.Stat(inputPDF); os.IsNotExist(err) {
		t.Skip("testdata/lorem_ipsum.pdf not found, skipping test")
	}

	outputDocx := filepath.Join(t.TempDir(), "output.docx")

	opts := unipdfutil.DefaultConversionOptions()
	output, err := s.ConvertPDFToDocx(ctx, inputPDF, outputDocx, opts)
	if err != nil {
		t.Fatalf("ConvertPDFToDocx failed: %v", err)
	}

	if output != outputDocx {
		t.Errorf("expected output '%s', got '%s'", outputDocx, output)
	}

	if _, err := os.Stat(output); os.IsNotExist(err) {
		t.Error("output file was not created")
	}
}

func TestConvertDocxToPDF(t *testing.T) {
	if os.Getenv("UNIDOC_KEY") == "" {
		t.Skip("UNIDOC_KEY not set, skipping integration test")
	}

	ctx := context.Background()
	s := New()
	if err := s.Init(ctx); err != nil {
		t.Fatalf("failed to init skill: %v", err)
	}
	defer s.Close()

	inputDocx := filepath.Join("..", "..", "testdata", "lorem_ipsum.docx")
	if _, err := os.Stat(inputDocx); os.IsNotExist(err) {
		t.Skip("testdata/lorem_ipsum.docx not found, skipping test")
	}

	outputPDF := filepath.Join(t.TempDir(), "output.pdf")

	output, err := s.ConvertDocxToPDF(ctx, inputDocx, outputPDF)
	if err != nil {
		t.Fatalf("ConvertDocxToPDF failed: %v", err)
	}

	if output != outputPDF {
		t.Errorf("expected output '%s', got '%s'", outputPDF, output)
	}

	if _, err := os.Stat(output); os.IsNotExist(err) {
		t.Error("output file was not created")
	}
}

func TestConvertPDFToDocx_DefaultOutput(t *testing.T) {
	if os.Getenv("UNIDOC_KEY") == "" {
		t.Skip("UNIDOC_KEY not set, skipping integration test")
	}

	ctx := context.Background()
	s := New()
	if err := s.Init(ctx); err != nil {
		t.Fatalf("failed to init skill: %v", err)
	}
	defer s.Close()

	// Copy test PDF to temp dir
	inputPDF := filepath.Join("..", "..", "testdata", "lorem_ipsum.pdf")
	if _, err := os.Stat(inputPDF); os.IsNotExist(err) {
		t.Skip("testdata/lorem_ipsum.pdf not found, skipping test")
	}

	tempDir := t.TempDir()
	tempPDF := filepath.Join(tempDir, "test.pdf")

	data, err := os.ReadFile(inputPDF)
	if err != nil {
		t.Fatalf("failed to read input: %v", err)
	}
	if err := os.WriteFile(tempPDF, data, 0600); err != nil { //nolint:gosec // G703: Test fixture path
		t.Fatalf("failed to write temp file: %v", err)
	}

	// Convert with empty output path (should default to .docx)
	opts := unipdfutil.DefaultConversionOptions()
	output, err := s.ConvertPDFToDocx(ctx, tempPDF, "", opts)
	if err != nil {
		t.Fatalf("ConvertPDFToDocx failed: %v", err)
	}

	expectedOutput := filepath.Join(tempDir, "test.docx")
	if output != expectedOutput {
		t.Errorf("expected output '%s', got '%s'", expectedOutput, output)
	}
}

func TestConvertPDFToDocx_EmptyInput(t *testing.T) {
	ctx := context.Background()
	s := New()

	opts := unipdfutil.DefaultConversionOptions()
	_, err := s.ConvertPDFToDocx(ctx, "", "output.docx", opts)
	if err == nil {
		t.Error("expected error for empty input path")
	}
}

func TestConvertDocxToPDF_EmptyInput(t *testing.T) {
	ctx := context.Background()
	s := New()

	_, err := s.ConvertDocxToPDF(ctx, "", "output.pdf")
	if err == nil {
		t.Error("expected error for empty input path")
	}
}

func TestReadMetadata_EmptyPath(t *testing.T) {
	ctx := context.Background()
	s := New()

	_, err := s.ReadMetadata(ctx, "")
	if err == nil {
		t.Error("expected error for empty file path")
	}
}

func TestMergePDFs_EmptyInputs(t *testing.T) {
	ctx := context.Background()
	s := New()

	err := s.MergePDFs(ctx, nil, "output.pdf")
	if err == nil {
		t.Error("expected error for empty input paths")
	}

	err = s.MergePDFs(ctx, []string{}, "output.pdf")
	if err == nil {
		t.Error("expected error for empty input paths slice")
	}
}

func TestMergePDFs_EmptyOutput(t *testing.T) {
	ctx := context.Background()
	s := New()

	err := s.MergePDFs(ctx, []string{"file1.pdf", "file2.pdf"}, "")
	if err == nil {
		t.Error("expected error for empty output path")
	}
}

func TestSplitPDF_EmptyInput(t *testing.T) {
	ctx := context.Background()
	s := New()

	_, err := s.SplitPDF(ctx, "", "")
	if err == nil {
		t.Error("expected error for empty input path")
	}
}

func TestExtractPDFPages_EmptyInput(t *testing.T) {
	ctx := context.Background()
	s := New()

	err := s.ExtractPDFPages(ctx, "", "output.pdf", []int{1, 2})
	if err == nil {
		t.Error("expected error for empty input path")
	}
}

func TestExtractPDFPages_EmptyOutput(t *testing.T) {
	ctx := context.Background()
	s := New()

	err := s.ExtractPDFPages(ctx, "input.pdf", "", []int{1, 2})
	if err == nil {
		t.Error("expected error for empty output path")
	}
}

func TestExtractPDFPages_EmptyPages(t *testing.T) {
	ctx := context.Background()
	s := New()

	err := s.ExtractPDFPages(ctx, "input.pdf", "output.pdf", nil)
	if err == nil {
		t.Error("expected error for empty pages")
	}

	err = s.ExtractPDFPages(ctx, "input.pdf", "output.pdf", []int{})
	if err == nil {
		t.Error("expected error for empty pages slice")
	}
}

func TestGetPDFPageCount_EmptyInput(t *testing.T) {
	ctx := context.Background()
	s := New()

	_, err := s.GetPDFPageCount(ctx, "")
	if err == nil {
		t.Error("expected error for empty input path")
	}
}

func TestMergePDFs_Integration(t *testing.T) {
	if os.Getenv("UNIDOC_KEY") == "" {
		t.Skip("UNIDOC_KEY not set, skipping integration test")
	}

	ctx := context.Background()
	s := New()
	if err := s.Init(ctx); err != nil {
		t.Fatalf("failed to init skill: %v", err)
	}
	defer s.Close()

	inputPDF := filepath.Join("..", "..", "testdata", "lorem_ipsum.pdf")
	if _, err := os.Stat(inputPDF); os.IsNotExist(err) {
		t.Skip("testdata/lorem_ipsum.pdf not found, skipping test")
	}

	// Merge the same file twice
	outputPath := filepath.Join(t.TempDir(), "merged.pdf")
	err := s.MergePDFs(ctx, []string{inputPDF, inputPDF}, outputPath)
	if err != nil {
		t.Fatalf("MergePDFs failed: %v", err)
	}

	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Error("merged output file was not created")
	}

	// Verify merged file has twice the pages
	originalCount, _ := s.GetPDFPageCount(ctx, inputPDF)
	mergedCount, err := s.GetPDFPageCount(ctx, outputPath)
	if err != nil {
		t.Fatalf("failed to get merged page count: %v", err)
	}

	if mergedCount != originalCount*2 {
		t.Errorf("expected %d pages in merged PDF, got %d", originalCount*2, mergedCount)
	}
}

func TestSplitPDF_Integration(t *testing.T) {
	if os.Getenv("UNIDOC_KEY") == "" {
		t.Skip("UNIDOC_KEY not set, skipping integration test")
	}

	ctx := context.Background()
	s := New()
	if err := s.Init(ctx); err != nil {
		t.Fatalf("failed to init skill: %v", err)
	}
	defer s.Close()

	inputPDF := filepath.Join("..", "..", "testdata", "lorem_ipsum.pdf")
	if _, err := os.Stat(inputPDF); os.IsNotExist(err) {
		t.Skip("testdata/lorem_ipsum.pdf not found, skipping test")
	}

	outputDir := t.TempDir()
	outputPaths, err := s.SplitPDF(ctx, inputPDF, outputDir)
	if err != nil {
		t.Fatalf("SplitPDF failed: %v", err)
	}

	pageCount, _ := s.GetPDFPageCount(ctx, inputPDF)
	if len(outputPaths) != pageCount {
		t.Errorf("expected %d split files, got %d", pageCount, len(outputPaths))
	}

	// Verify each split file exists
	for _, p := range outputPaths {
		if _, err := os.Stat(p); os.IsNotExist(err) {
			t.Errorf("split file not created: %s", p)
		}
	}
}

func TestGetPDFPageCount_Integration(t *testing.T) {
	if os.Getenv("UNIDOC_KEY") == "" {
		t.Skip("UNIDOC_KEY not set, skipping integration test")
	}

	ctx := context.Background()
	s := New()
	if err := s.Init(ctx); err != nil {
		t.Fatalf("failed to init skill: %v", err)
	}
	defer s.Close()

	inputPDF := filepath.Join("..", "..", "testdata", "lorem_ipsum.pdf")
	if _, err := os.Stat(inputPDF); os.IsNotExist(err) {
		t.Skip("testdata/lorem_ipsum.pdf not found, skipping test")
	}

	count, err := s.GetPDFPageCount(ctx, inputPDF)
	if err != nil {
		t.Fatalf("GetPDFPageCount failed: %v", err)
	}

	if count < 1 {
		t.Errorf("expected at least 1 page, got %d", count)
	}
}

func TestNewWithLogger(t *testing.T) {
	// Test with nil logger (should use default)
	s1 := NewWithLogger(nil)
	if s1.Name() != "unidoc" {
		t.Error("skill should work with nil logger")
	}

	// Test with custom logger
	logger := slog.Default()
	s2 := NewWithLogger(logger)
	if s2.Name() != "unidoc" {
		t.Error("skill should work with custom logger")
	}
}

func TestSetLogger(t *testing.T) {
	s := New()

	// SetLogger with nil should not change logger
	s.SetLogger(nil)

	// SetLogger with valid logger should work
	logger := slog.Default()
	s.SetLogger(logger)
}

func TestGetXlsxInfo_EmptyPath(t *testing.T) {
	ctx := context.Background()
	s := New()

	_, err := s.GetXlsxInfo(ctx, "")
	if err == nil {
		t.Error("expected error for empty file path")
	}
}

func TestGetPptxInfo_EmptyPath(t *testing.T) {
	ctx := context.Background()
	s := New()

	_, err := s.GetPptxInfo(ctx, "")
	if err == nil {
		t.Error("expected error for empty file path")
	}
}
