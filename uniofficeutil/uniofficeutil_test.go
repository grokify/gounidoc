//nolint:dupl // Test functions have similar table-driven structure by design
package uniofficeutil

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// ---------------------------------------------------------------------------
// Unit Tests - replaceText function
// ---------------------------------------------------------------------------

func TestReplaceText_CaseSensitive(t *testing.T) {
	tests := []struct {
		name      string
		text      string
		find      string
		replace   string
		opts      ReplaceOptions
		wantText  string
		wantCount int
	}{
		{
			name:      "simple replace",
			text:      "Hello World",
			find:      "World",
			replace:   "Go",
			opts:      ReplaceOptions{CaseSensitive: true, ReplaceAll: true},
			wantText:  "Hello Go",
			wantCount: 1,
		},
		{
			name:      "multiple occurrences",
			text:      "foo bar foo baz foo",
			find:      "foo",
			replace:   "qux",
			opts:      ReplaceOptions{CaseSensitive: true, ReplaceAll: true},
			wantText:  "qux bar qux baz qux",
			wantCount: 3,
		},
		{
			name:      "no match",
			text:      "Hello World",
			find:      "Goodbye",
			replace:   "Hi",
			opts:      ReplaceOptions{CaseSensitive: true, ReplaceAll: true},
			wantText:  "Hello World",
			wantCount: 0,
		},
		{
			name:      "case mismatch",
			text:      "Hello World",
			find:      "world",
			replace:   "Go",
			opts:      ReplaceOptions{CaseSensitive: true, ReplaceAll: true},
			wantText:  "Hello World",
			wantCount: 0,
		},
		{
			name:      "empty text",
			text:      "",
			find:      "foo",
			replace:   "bar",
			opts:      ReplaceOptions{CaseSensitive: true, ReplaceAll: true},
			wantText:  "",
			wantCount: 0,
		},
		{
			name:      "empty find",
			text:      "Hello",
			find:      "",
			replace:   "bar",
			opts:      ReplaceOptions{CaseSensitive: true, ReplaceAll: true},
			wantText:  "Hello",
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotText, gotCount := replaceText(tt.text, tt.find, tt.replace, tt.opts)
			if gotText != tt.wantText {
				t.Errorf("replaceText() text = %q, want %q", gotText, tt.wantText)
			}
			if gotCount != tt.wantCount {
				t.Errorf("replaceText() count = %d, want %d", gotCount, tt.wantCount)
			}
		})
	}
}

func TestReplaceText_CaseInsensitive(t *testing.T) {
	tests := []struct {
		name      string
		text      string
		find      string
		replace   string
		wantText  string
		wantCount int
	}{
		{
			name:      "lowercase find uppercase text",
			text:      "Hello WORLD",
			find:      "world",
			replace:   "Go",
			wantText:  "Hello Go",
			wantCount: 1,
		},
		{
			name:      "mixed case multiple",
			text:      "Foo FOO foo fOO",
			find:      "foo",
			replace:   "bar",
			wantText:  "bar bar bar bar",
			wantCount: 4,
		},
	}

	opts := ReplaceOptions{CaseSensitive: false, ReplaceAll: true}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotText, gotCount := replaceText(tt.text, tt.find, tt.replace, opts)
			if gotText != tt.wantText {
				t.Errorf("replaceText() text = %q, want %q", gotText, tt.wantText)
			}
			if gotCount != tt.wantCount {
				t.Errorf("replaceText() count = %d, want %d", gotCount, tt.wantCount)
			}
		})
	}
}

func TestReplaceText_ReplaceFirst(t *testing.T) {
	tests := []struct {
		name      string
		text      string
		find      string
		replace   string
		wantText  string
		wantCount int
	}{
		{
			name:      "replace first only",
			text:      "foo bar foo baz foo",
			find:      "foo",
			replace:   "qux",
			wantText:  "qux bar foo baz foo",
			wantCount: 1,
		},
		{
			name:      "no match",
			text:      "Hello World",
			find:      "Goodbye",
			replace:   "Hi",
			wantText:  "Hello World",
			wantCount: 0,
		},
	}

	opts := ReplaceOptions{CaseSensitive: true, ReplaceAll: false}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotText, gotCount := replaceText(tt.text, tt.find, tt.replace, opts)
			if gotText != tt.wantText {
				t.Errorf("replaceText() text = %q, want %q", gotText, tt.wantText)
			}
			if gotCount != tt.wantCount {
				t.Errorf("replaceText() count = %d, want %d", gotCount, tt.wantCount)
			}
		})
	}
}

func TestReplaceText_WholeWord(t *testing.T) {
	tests := []struct {
		name      string
		text      string
		find      string
		replace   string
		wantText  string
		wantCount int
	}{
		{
			name:      "whole word match",
			text:      "the cat sat",
			find:      "cat",
			replace:   "dog",
			wantText:  "the dog sat",
			wantCount: 1,
		},
		{
			name:      "no match - part of word",
			text:      "category catalog",
			find:      "cat",
			replace:   "dog",
			wantText:  "category catalog",
			wantCount: 0,
		},
		{
			name:      "word at start",
			text:      "cat is here",
			find:      "cat",
			replace:   "dog",
			wantText:  "dog is here",
			wantCount: 1,
		},
		{
			name:      "word at end",
			text:      "see the cat",
			find:      "cat",
			replace:   "dog",
			wantText:  "see the dog",
			wantCount: 1,
		},
		{
			name:      "multiple whole words",
			text:      "cat and cat and cat",
			find:      "cat",
			replace:   "dog",
			wantText:  "dog and dog and dog",
			wantCount: 3,
		},
	}

	opts := ReplaceOptions{CaseSensitive: true, WholeWord: true, ReplaceAll: true}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotText, gotCount := replaceText(tt.text, tt.find, tt.replace, opts)
			if gotText != tt.wantText {
				t.Errorf("replaceText() text = %q, want %q", gotText, tt.wantText)
			}
			if gotCount != tt.wantCount {
				t.Errorf("replaceText() count = %d, want %d", gotCount, tt.wantCount)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Unit Tests - isWordChar function
// ---------------------------------------------------------------------------

func TestIsWordChar(t *testing.T) {
	tests := []struct {
		r    rune
		want bool
	}{
		{'a', true},
		{'z', true},
		{'A', true},
		{'Z', true},
		{'0', true},
		{'9', true},
		{'_', true},
		{' ', false},
		{'-', false},
		{'.', false},
		{'!', false},
		{'@', false},
	}

	for _, tt := range tests {
		t.Run(string(tt.r), func(t *testing.T) {
			if got := isWordChar(tt.r); got != tt.want {
				t.Errorf("isWordChar(%q) = %v, want %v", tt.r, got, tt.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Unit Tests - sanitizeFilename function
// ---------------------------------------------------------------------------

func TestSanitizeFilename(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"simple", "document", "document"},
		{"with spaces", "my document", "my_document"},
		{"with slash", "path/name", "path_name"},
		{"with backslash", "path\\name", "path_name"},
		{"with colon", "file:name", "file_name"},
		{"with asterisk", "file*name", "file_name"},
		{"with question", "file?name", "file_name"},
		{"with quotes", "file\"name", "file_name"},
		{"with angle brackets", "file<name>", "file_name_"},
		{"with pipe", "file|name", "file_name"},
		{"multiple special", "a/b\\c:d*e?f\"g<h>i|j k", "a_b_c_d_e_f_g_h_i_j_k"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sanitizeFilename(tt.input); got != tt.want {
				t.Errorf("sanitizeFilename(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Unit Tests - Default Options
// ---------------------------------------------------------------------------

func TestDefaultReplaceOptions(t *testing.T) {
	opts := DefaultReplaceOptions()

	if !opts.CaseSensitive {
		t.Error("DefaultReplaceOptions().CaseSensitive should be true")
	}
	if !opts.ReplaceAll {
		t.Error("DefaultReplaceOptions().ReplaceAll should be true")
	}
	if opts.WholeWord {
		t.Error("DefaultReplaceOptions().WholeWord should be false")
	}
}

func TestDefaultCSVOptions(t *testing.T) {
	opts := DefaultCSVOptions()

	if opts.Delimiter != "," {
		t.Errorf("DefaultCSVOptions().Delimiter = %q, want %q", opts.Delimiter, ",")
	}
	if opts.SheetName != "" {
		t.Errorf("DefaultCSVOptions().SheetName should be empty, got %q", opts.SheetName)
	}
	if opts.AllSheets {
		t.Error("DefaultCSVOptions().AllSheets should be false")
	}
}

// ---------------------------------------------------------------------------
// Input Validation Tests
// ---------------------------------------------------------------------------

func TestExtractDocxText_InputValidation(t *testing.T) {
	_, err := ExtractDocxText("")
	if err == nil {
		t.Error("ExtractDocxText(\"\") should return error")
	}
	if !strings.Contains(err.Error(), "input path is required") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestExtractDocxTextToFile_InputValidation(t *testing.T) {
	// Empty output path
	err := ExtractDocxTextToFile("input.docx", "")
	if err == nil {
		t.Error("ExtractDocxTextToFile with empty output should return error")
	}
	if !strings.Contains(err.Error(), "output path is required") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestReplaceInDocx_InputValidation(t *testing.T) {
	tests := []struct {
		name       string
		inputPath  string
		outputPath string
		find       string
		wantErr    string
	}{
		{"empty input", "", "output.docx", "find", "input path is required"},
		{"empty output", "input.docx", "", "find", "output path is required"},
		{"empty find", "input.docx", "output.docx", "", "find string is required"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ReplaceInDocx(tt.inputPath, tt.outputPath, tt.find, "replace", DefaultReplaceOptions())
			if err == nil {
				t.Errorf("ReplaceInDocx should return error for %s", tt.name)
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("unexpected error: %v, want containing %q", err, tt.wantErr)
			}
		})
	}
}

func TestReplaceMultipleInDocx_InputValidation(t *testing.T) {
	tests := []struct {
		name         string
		inputPath    string
		outputPath   string
		replacements map[string]string
		wantErr      string
	}{
		{"empty input", "", "output.docx", map[string]string{"a": "b"}, "input path is required"},
		{"empty output", "input.docx", "", map[string]string{"a": "b"}, "output path is required"},
		{"empty replacements", "input.docx", "output.docx", map[string]string{}, "no replacements specified"},
		{"nil replacements", "input.docx", "output.docx", nil, "no replacements specified"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ReplaceMultipleInDocx(tt.inputPath, tt.outputPath, tt.replacements, DefaultReplaceOptions())
			if err == nil {
				t.Errorf("ReplaceMultipleInDocx should return error for %s", tt.name)
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("unexpected error: %v, want containing %q", err, tt.wantErr)
			}
		})
	}
}

func TestConvertXlsxToCSV_InputValidation(t *testing.T) {
	tests := []struct {
		name       string
		inputPath  string
		outputPath string
		wantErr    string
	}{
		{"empty input", "", "output.csv", "input path is required"},
		{"empty output", "input.xlsx", "", "output path is required"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ConvertXlsxToCSV(tt.inputPath, tt.outputPath, DefaultCSVOptions())
			if err == nil {
				t.Errorf("ConvertXlsxToCSV should return error for %s", tt.name)
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("unexpected error: %v, want containing %q", err, tt.wantErr)
			}
		})
	}
}

func TestConvertXlsxToCSVAllSheets_InputValidation(t *testing.T) {
	_, err := ConvertXlsxToCSVAllSheets("", "", DefaultCSVOptions())
	if err == nil {
		t.Error("ConvertXlsxToCSVAllSheets with empty input should return error")
	}
	if !strings.Contains(err.Error(), "input path is required") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestReadFileMetadata_UnsupportedType(t *testing.T) {
	_, err := ReadFileMetadata("document.txt")
	if err == nil {
		t.Error("ReadFileMetadata should return error for unsupported type")
	}
	if !strings.Contains(err.Error(), "unsupported file type") {
		t.Errorf("unexpected error: %v", err)
	}
}

// ---------------------------------------------------------------------------
// Integration Tests - require test fixtures
// ---------------------------------------------------------------------------

func TestExtractDocxText_Integration(t *testing.T) {
	testFile := filepath.Join("..", "testdata", "lorem_ipsum.docx")
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Skip("test fixture not found:", testFile)
	}

	text, err := ExtractDocxText(testFile)
	if err != nil {
		if strings.Contains(err.Error(), "license required") {
			t.Skip("UniDoc license required for integration tests")
		}
		t.Fatalf("ExtractDocxText failed: %v", err)
	}

	if text == "" {
		t.Error("ExtractDocxText returned empty text")
	}

	// Check for expected content (lorem ipsum should contain these)
	if !strings.Contains(strings.ToLower(text), "lorem") {
		t.Error("extracted text should contain 'lorem'")
	}
}

func TestReadFileMetadata_Integration(t *testing.T) {
	testFile := filepath.Join("..", "testdata", "lorem_ipsum.docx")
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Skip("test fixture not found:", testFile)
	}

	meta, err := ReadFileMetadata(testFile)
	if err != nil {
		if strings.Contains(err.Error(), "license required") {
			t.Skip("UniDoc license required for integration tests")
		}
		t.Fatalf("ReadFileMetadata failed: %v", err)
	}

	if meta.FileType != "docx" {
		t.Errorf("FileType = %q, want %q", meta.FileType, "docx")
	}
}

func TestReadFileMetadataDocx_Integration(t *testing.T) {
	testFile := filepath.Join("..", "testdata", "lorem_ipsum.docx")
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Skip("test fixture not found:", testFile)
	}

	meta, err := ReadFileMetadataDocx(testFile)
	if err != nil {
		if strings.Contains(err.Error(), "license required") {
			t.Skip("UniDoc license required for integration tests")
		}
		t.Fatalf("ReadFileMetadataDocx failed: %v", err)
	}

	if meta.FileType != "docx" {
		t.Errorf("FileType = %q, want %q", meta.FileType, "docx")
	}
}

func TestExtractDocxTextToFile_Integration(t *testing.T) {
	testFile := filepath.Join("..", "testdata", "lorem_ipsum.docx")
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Skip("test fixture not found:", testFile)
	}

	// Create temp output file
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.txt")

	err := ExtractDocxTextToFile(testFile, outputFile)
	if err != nil {
		if strings.Contains(err.Error(), "license required") {
			t.Skip("UniDoc license required for integration tests")
		}
		t.Fatalf("ExtractDocxTextToFile failed: %v", err)
	}

	// Verify output file exists and has content
	data, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}

	if len(data) == 0 {
		t.Error("output file is empty")
	}

	if !strings.Contains(strings.ToLower(string(data)), "lorem") {
		t.Error("output file should contain 'lorem'")
	}
}

// ---------------------------------------------------------------------------
// Benchmark Tests
// ---------------------------------------------------------------------------

func BenchmarkReplaceText_Simple(b *testing.B) {
	text := "The quick brown fox jumps over the lazy dog"
	opts := ReplaceOptions{CaseSensitive: true, ReplaceAll: true}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		replaceText(text, "fox", "cat", opts)
	}
}

func BenchmarkReplaceText_CaseInsensitive(b *testing.B) {
	text := "The Quick Brown FOX jumps over the lazy dog"
	opts := ReplaceOptions{CaseSensitive: false, ReplaceAll: true}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		replaceText(text, "fox", "cat", opts)
	}
}

func BenchmarkReplaceText_WholeWord(b *testing.B) {
	text := "The quick brown fox jumps over the lazy dog"
	opts := ReplaceOptions{CaseSensitive: true, WholeWord: true, ReplaceAll: true}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		replaceText(text, "fox", "cat", opts)
	}
}

func BenchmarkSanitizeFilename(b *testing.B) {
	input := "My Document: Version 1.0 <draft>"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sanitizeFilename(input)
	}
}
