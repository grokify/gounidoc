package unipdfutil

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/unidoc/unipdf/v3/model"
)

// MergePDFFiles merges multiple PDF files into a single output file.
// The files are merged in the order provided.
func MergePDFFiles(inputPaths []string, outputPath string) error {
	if len(inputPaths) == 0 {
		return fmt.Errorf("no input files provided")
	}
	if outputPath == "" {
		return fmt.Errorf("output path is required")
	}

	pdfWriter := model.NewPdfWriter()

	for _, inputPath := range inputPaths {
		if err := appendPDFToWriter(inputPath, &pdfWriter); err != nil {
			return fmt.Errorf("failed to append %s: %w", inputPath, err)
		}
	}

	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	if err := pdfWriter.Write(outFile); err != nil {
		return fmt.Errorf("failed to write merged PDF: %w", err)
	}

	return nil
}

// appendPDFToWriter adds all pages from a PDF file to the writer.
func appendPDFToWriter(inputPath string, writer *model.PdfWriter) error {
	f, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer f.Close()

	reader, err := model.NewPdfReader(f)
	if err != nil {
		return err
	}

	numPages, err := reader.GetNumPages()
	if err != nil {
		return err
	}

	for i := 1; i <= numPages; i++ {
		page, err := reader.GetPage(i)
		if err != nil {
			return fmt.Errorf("failed to get page %d: %w", i, err)
		}
		if err := writer.AddPage(page); err != nil {
			return fmt.Errorf("failed to add page %d: %w", i, err)
		}
	}

	return nil
}

// SplitPDFFile splits a PDF file into individual pages or page ranges.
// If outputDir is empty, files are created in the same directory as the input.
// Returns the list of created file paths.
func SplitPDFFile(inputPath, outputDir string) ([]string, error) {
	if inputPath == "" {
		return nil, fmt.Errorf("input path is required")
	}

	if outputDir == "" {
		outputDir = filepath.Dir(inputPath)
	}

	f, err := os.Open(inputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open input file: %w", err)
	}
	defer f.Close()

	reader, err := model.NewPdfReader(f)
	if err != nil {
		return nil, fmt.Errorf("failed to create PDF reader: %w", err)
	}

	numPages, err := reader.GetNumPages()
	if err != nil {
		return nil, fmt.Errorf("failed to get page count: %w", err)
	}

	baseName := filepath.Base(inputPath)
	ext := filepath.Ext(baseName)
	nameWithoutExt := baseName[:len(baseName)-len(ext)]

	var outputPaths []string

	for i := 1; i <= numPages; i++ {
		page, err := reader.GetPage(i)
		if err != nil {
			return outputPaths, fmt.Errorf("failed to get page %d: %w", i, err)
		}

		writer := model.NewPdfWriter()
		if err := writer.AddPage(page); err != nil {
			return outputPaths, fmt.Errorf("failed to add page %d to writer: %w", i, err)
		}

		outputPath := filepath.Join(outputDir, fmt.Sprintf("%s_page_%03d.pdf", nameWithoutExt, i))
		outFile, err := os.Create(outputPath)
		if err != nil {
			return outputPaths, fmt.Errorf("failed to create output file for page %d: %w", i, err)
		}

		if err := writer.Write(outFile); err != nil {
			outFile.Close()
			return outputPaths, fmt.Errorf("failed to write page %d: %w", i, err)
		}
		outFile.Close()

		outputPaths = append(outputPaths, outputPath)
	}

	return outputPaths, nil
}

// ExtractPDFPages extracts specific pages from a PDF file.
// Pages are 1-indexed. Returns the output file path.
func ExtractPDFPages(inputPath, outputPath string, pages []int) error {
	if inputPath == "" {
		return fmt.Errorf("input path is required")
	}
	if outputPath == "" {
		return fmt.Errorf("output path is required")
	}
	if len(pages) == 0 {
		return fmt.Errorf("no pages specified")
	}

	f, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer f.Close()

	reader, err := model.NewPdfReader(f)
	if err != nil {
		return fmt.Errorf("failed to create PDF reader: %w", err)
	}

	numPages, err := reader.GetNumPages()
	if err != nil {
		return fmt.Errorf("failed to get page count: %w", err)
	}

	writer := model.NewPdfWriter()

	for _, pageNum := range pages {
		if pageNum < 1 || pageNum > numPages {
			return fmt.Errorf("page %d out of range (1-%d)", pageNum, numPages)
		}

		page, err := reader.GetPage(pageNum)
		if err != nil {
			return fmt.Errorf("failed to get page %d: %w", pageNum, err)
		}

		if err := writer.AddPage(page); err != nil {
			return fmt.Errorf("failed to add page %d: %w", pageNum, err)
		}
	}

	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	if err := writer.Write(outFile); err != nil {
		return fmt.Errorf("failed to write output PDF: %w", err)
	}

	return nil
}

// GetPDFPageCount returns the number of pages in a PDF file.
func GetPDFPageCount(inputPath string) (int, error) {
	f, err := os.Open(inputPath)
	if err != nil {
		return 0, fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	reader, err := model.NewPdfReader(f)
	if err != nil {
		return 0, fmt.Errorf("failed to create PDF reader: %w", err)
	}

	return reader.GetNumPages()
}
