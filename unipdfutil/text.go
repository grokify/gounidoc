package unipdfutil

import (
	"fmt"
	"os"
	"strings"

	"github.com/unidoc/unipdf/v3/extractor"
	"github.com/unidoc/unipdf/v3/model"
)

// ExtractPDFText extracts all text from a PDF file.
func ExtractPDFText(inputPath string) (string, error) {
	pages, err := ExtractPDFTextByPage(inputPath)
	if err != nil {
		return "", err
	}
	return strings.Join(pages, "\n\n"), nil
}

// ExtractPDFTextByPage extracts text from each page of a PDF file.
// Returns a slice of strings, one per page.
func ExtractPDFTextByPage(inputPath string) ([]string, error) {
	if inputPath == "" {
		return nil, fmt.Errorf("input path is required")
	}

	f, err := os.Open(inputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
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

	var texts []string
	for i := 1; i <= numPages; i++ {
		page, err := reader.GetPage(i)
		if err != nil {
			return texts, fmt.Errorf("failed to get page %d: %w", i, err)
		}

		ex, err := extractor.New(page)
		if err != nil {
			return texts, fmt.Errorf("failed to create extractor for page %d: %w", i, err)
		}

		text, err := ex.ExtractText()
		if err != nil {
			return texts, fmt.Errorf("failed to extract text from page %d: %w", i, err)
		}

		texts = append(texts, text)
	}

	return texts, nil
}

// ExtractPDFTextToFile extracts text from a PDF and writes it to a file.
func ExtractPDFTextToFile(inputPath, outputPath string) error {
	if outputPath == "" {
		return fmt.Errorf("output path is required")
	}

	text, err := ExtractPDFText(inputPath)
	if err != nil {
		return err
	}

	if err := os.WriteFile(outputPath, []byte(text), 0644); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	return nil
}

// ExtractPDFPageText extracts text from a specific page of a PDF file.
// Page numbers are 1-indexed.
func ExtractPDFPageText(inputPath string, pageNum int) (string, error) {
	if inputPath == "" {
		return "", fmt.Errorf("input path is required")
	}
	if pageNum < 1 {
		return "", fmt.Errorf("page number must be >= 1")
	}

	f, err := os.Open(inputPath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	reader, err := model.NewPdfReader(f)
	if err != nil {
		return "", fmt.Errorf("failed to create PDF reader: %w", err)
	}

	numPages, err := reader.GetNumPages()
	if err != nil {
		return "", fmt.Errorf("failed to get page count: %w", err)
	}

	if pageNum > numPages {
		return "", fmt.Errorf("page %d out of range (1-%d)", pageNum, numPages)
	}

	page, err := reader.GetPage(pageNum)
	if err != nil {
		return "", fmt.Errorf("failed to get page %d: %w", pageNum, err)
	}

	ex, err := extractor.New(page)
	if err != nil {
		return "", fmt.Errorf("failed to create extractor: %w", err)
	}

	text, err := ex.ExtractText()
	if err != nil {
		return "", fmt.Errorf("failed to extract text: %w", err)
	}

	return text, nil
}
