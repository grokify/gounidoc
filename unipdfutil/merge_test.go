package unipdfutil

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/grokify/gounidoc"
)

func init() {
	_ = gounidoc.SetMeteredKeyEnv()
}

func TestMergePDFFiles_EmptyInputs(t *testing.T) {
	err := MergePDFFiles(nil, "output.pdf")
	if err == nil {
		t.Error("expected error for nil input paths")
	}

	err = MergePDFFiles([]string{}, "output.pdf")
	if err == nil {
		t.Error("expected error for empty input paths")
	}
}

func TestMergePDFFiles_EmptyOutput(t *testing.T) {
	err := MergePDFFiles([]string{"file1.pdf"}, "")
	if err == nil {
		t.Error("expected error for empty output path")
	}
}

func TestSplitPDFFile_EmptyInput(t *testing.T) {
	_, err := SplitPDFFile("", "")
	if err == nil {
		t.Error("expected error for empty input path")
	}
}

func TestExtractPDFPages_EmptyInput(t *testing.T) {
	err := ExtractPDFPages("", "output.pdf", []int{1})
	if err == nil {
		t.Error("expected error for empty input path")
	}
}

func TestExtractPDFPages_EmptyOutput(t *testing.T) {
	err := ExtractPDFPages("input.pdf", "", []int{1})
	if err == nil {
		t.Error("expected error for empty output path")
	}
}

func TestExtractPDFPages_EmptyPages(t *testing.T) {
	err := ExtractPDFPages("input.pdf", "output.pdf", nil)
	if err == nil {
		t.Error("expected error for nil pages")
	}

	err = ExtractPDFPages("input.pdf", "output.pdf", []int{})
	if err == nil {
		t.Error("expected error for empty pages")
	}
}

func TestMergePDFFiles_Integration(t *testing.T) {
	if os.Getenv("UNIDOC_KEY") == "" {
		t.Skip("UNIDOC_KEY not set, skipping integration test")
	}

	inputPDF := filepath.Join("..", "testdata", "lorem_ipsum.pdf")
	if _, err := os.Stat(inputPDF); os.IsNotExist(err) {
		t.Skip("testdata/lorem_ipsum.pdf not found, skipping test")
	}

	outputPath := filepath.Join(t.TempDir(), "merged.pdf")
	err := MergePDFFiles([]string{inputPDF, inputPDF}, outputPath)
	if err != nil {
		t.Fatalf("MergePDFFiles failed: %v", err)
	}

	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Error("merged output file was not created")
	}

	// Verify merged file has pages
	count, err := GetPDFPageCount(outputPath)
	if err != nil {
		t.Fatalf("failed to get page count: %v", err)
	}
	if count < 2 {
		t.Errorf("expected at least 2 pages in merged PDF, got %d", count)
	}
}

func TestSplitPDFFile_Integration(t *testing.T) {
	if os.Getenv("UNIDOC_KEY") == "" {
		t.Skip("UNIDOC_KEY not set, skipping integration test")
	}

	inputPDF := filepath.Join("..", "testdata", "lorem_ipsum.pdf")
	if _, err := os.Stat(inputPDF); os.IsNotExist(err) {
		t.Skip("testdata/lorem_ipsum.pdf not found, skipping test")
	}

	outputDir := t.TempDir()
	outputPaths, err := SplitPDFFile(inputPDF, outputDir)
	if err != nil {
		t.Fatalf("SplitPDFFile failed: %v", err)
	}

	if len(outputPaths) == 0 {
		t.Error("expected at least one output file")
	}

	// Verify each split file exists and has 1 page
	for _, p := range outputPaths {
		if _, err := os.Stat(p); os.IsNotExist(err) {
			t.Errorf("split file not created: %s", p)
		}
		count, _ := GetPDFPageCount(p)
		if count != 1 {
			t.Errorf("expected 1 page per split file, got %d for %s", count, p)
		}
	}
}

func TestExtractPDFPages_Integration(t *testing.T) {
	if os.Getenv("UNIDOC_KEY") == "" {
		t.Skip("UNIDOC_KEY not set, skipping integration test")
	}

	inputPDF := filepath.Join("..", "testdata", "lorem_ipsum.pdf")
	if _, err := os.Stat(inputPDF); os.IsNotExist(err) {
		t.Skip("testdata/lorem_ipsum.pdf not found, skipping test")
	}

	outputPath := filepath.Join(t.TempDir(), "extracted.pdf")
	err := ExtractPDFPages(inputPDF, outputPath, []int{1})
	if err != nil {
		t.Fatalf("ExtractPDFPages failed: %v", err)
	}

	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Error("extracted output file was not created")
	}

	count, err := GetPDFPageCount(outputPath)
	if err != nil {
		t.Fatalf("failed to get page count: %v", err)
	}
	if count != 1 {
		t.Errorf("expected 1 page in extracted PDF, got %d", count)
	}
}

func TestGetPDFPageCount_Integration(t *testing.T) {
	if os.Getenv("UNIDOC_KEY") == "" {
		t.Skip("UNIDOC_KEY not set, skipping integration test")
	}

	inputPDF := filepath.Join("..", "testdata", "lorem_ipsum.pdf")
	if _, err := os.Stat(inputPDF); os.IsNotExist(err) {
		t.Skip("testdata/lorem_ipsum.pdf not found, skipping test")
	}

	count, err := GetPDFPageCount(inputPDF)
	if err != nil {
		t.Fatalf("GetPDFPageCount failed: %v", err)
	}

	if count < 1 {
		t.Errorf("expected at least 1 page, got %d", count)
	}
}

func TestExtractPDFPages_OutOfRange(t *testing.T) {
	if os.Getenv("UNIDOC_KEY") == "" {
		t.Skip("UNIDOC_KEY not set, skipping integration test")
	}

	inputPDF := filepath.Join("..", "testdata", "lorem_ipsum.pdf")
	if _, err := os.Stat(inputPDF); os.IsNotExist(err) {
		t.Skip("testdata/lorem_ipsum.pdf not found, skipping test")
	}

	outputPath := filepath.Join(t.TempDir(), "extracted.pdf")

	// Try to extract page 999 which shouldn't exist
	err := ExtractPDFPages(inputPDF, outputPath, []int{999})
	if err == nil {
		t.Error("expected error for out of range page number")
	}
}
