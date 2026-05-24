package unipdfutil

import (
	"fmt"
	"path/filepath"
	"strings"
)

// ProgressCallback is called during batch operations to report progress.
// current is the 1-based index of the current file, total is the total count.
// inputPath is the file being processed.
type ProgressCallback func(current, total int, inputPath string)

// BatchResult represents the result of a single file conversion in a batch operation.
type BatchResult struct {
	InputPath  string `json:"input_path"`
	OutputPath string `json:"output_path"`
	Success    bool   `json:"success"`
	Error      string `json:"error,omitempty"`
}

// BatchConversionResult represents the result of a batch conversion operation.
type BatchConversionResult struct {
	Results    []BatchResult `json:"results"`
	TotalFiles int           `json:"total_files"`
	Succeeded  int           `json:"succeeded"`
	Failed     int           `json:"failed"`
}

// BatchConvertPDFToDocx converts multiple PDF files to DOCX format.
// If outputDir is empty, output files are created in the same directory as the inputs.
func BatchConvertPDFToDocx(inputPaths []string, outputDir string, opts ConversionOptions) BatchConversionResult {
	return BatchConvertPDFToDocxWithProgress(inputPaths, outputDir, opts, nil)
}

// BatchConvertPDFToDocxWithProgress converts multiple PDF files to DOCX format with progress reporting.
// If outputDir is empty, output files are created in the same directory as the inputs.
func BatchConvertPDFToDocxWithProgress(inputPaths []string, outputDir string, opts ConversionOptions, progress ProgressCallback) BatchConversionResult {
	result := BatchConversionResult{
		Results:    make([]BatchResult, 0, len(inputPaths)),
		TotalFiles: len(inputPaths),
	}

	for i, inputPath := range inputPaths {
		if progress != nil {
			progress(i+1, len(inputPaths), inputPath)
		}

		outputPath := generateOutputPath(inputPath, outputDir, ".docx")
		br := BatchResult{
			InputPath:  inputPath,
			OutputPath: outputPath,
		}

		if err := ConvertPDFFileToDocxFileWithOptions(inputPath, outputPath, opts); err != nil {
			br.Success = false
			br.Error = err.Error()
			result.Failed++
		} else {
			br.Success = true
			result.Succeeded++
		}

		result.Results = append(result.Results, br)
	}

	return result
}

// BatchConvertDocxToPDF converts multiple DOCX files to PDF format.
// If outputDir is empty, output files are created in the same directory as the inputs.
func BatchConvertDocxToPDF(inputPaths []string, outputDir string, converter func(input, output string) error) BatchConversionResult {
	return BatchConvertDocxToPDFWithProgress(inputPaths, outputDir, converter, nil)
}

// BatchConvertDocxToPDFWithProgress converts multiple DOCX files to PDF format with progress reporting.
// If outputDir is empty, output files are created in the same directory as the inputs.
func BatchConvertDocxToPDFWithProgress(inputPaths []string, outputDir string, converter func(input, output string) error, progress ProgressCallback) BatchConversionResult {
	result := BatchConversionResult{
		Results:    make([]BatchResult, 0, len(inputPaths)),
		TotalFiles: len(inputPaths),
	}

	for i, inputPath := range inputPaths {
		if progress != nil {
			progress(i+1, len(inputPaths), inputPath)
		}

		outputPath := generateOutputPath(inputPath, outputDir, ".pdf")
		br := BatchResult{
			InputPath:  inputPath,
			OutputPath: outputPath,
		}

		if err := converter(inputPath, outputPath); err != nil {
			br.Success = false
			br.Error = err.Error()
			result.Failed++
		} else {
			br.Success = true
			result.Succeeded++
		}

		result.Results = append(result.Results, br)
	}

	return result
}

// ExpandGlobPatterns expands glob patterns to a list of file paths.
// Non-glob patterns are returned as-is if they exist.
func ExpandGlobPatterns(patterns []string) ([]string, error) {
	var files []string
	seen := make(map[string]bool)

	for _, pattern := range patterns {
		// Check if it's a glob pattern
		if strings.ContainsAny(pattern, "*?[") {
			matches, err := filepath.Glob(pattern)
			if err != nil {
				return nil, fmt.Errorf("invalid glob pattern %q: %w", pattern, err)
			}
			for _, match := range matches {
				if !seen[match] {
					seen[match] = true
					files = append(files, match)
				}
			}
		} else {
			// Not a glob, add as-is
			if !seen[pattern] {
				seen[pattern] = true
				files = append(files, pattern)
			}
		}
	}

	return files, nil
}

// generateOutputPath generates an output path for a converted file.
func generateOutputPath(inputPath, outputDir, newExt string) string {
	baseName := filepath.Base(inputPath)
	ext := filepath.Ext(baseName)
	nameWithoutExt := baseName[:len(baseName)-len(ext)]
	outputName := nameWithoutExt + newExt

	if outputDir != "" {
		return filepath.Join(outputDir, outputName)
	}
	return filepath.Join(filepath.Dir(inputPath), outputName)
}

// FilterByExtension filters a list of paths to include only those with the specified extension.
func FilterByExtension(paths []string, ext string) []string {
	ext = strings.ToLower(ext)
	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}

	var filtered []string
	for _, p := range paths {
		if strings.ToLower(filepath.Ext(p)) == ext {
			filtered = append(filtered, p)
		}
	}
	return filtered
}
