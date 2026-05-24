package unipdfutil

import (
	"fmt"
	"os"

	"github.com/unidoc/unipdf/v3/model"
)

// RotatePDF rotates all pages in a PDF by the specified degrees.
// Valid rotation values are 0, 90, 180, or 270.
func RotatePDF(inputPath, outputPath string, degrees int) error {
	return RotatePDFPages(inputPath, outputPath, degrees, nil)
}

// RotatePDFPages rotates specific pages in a PDF by the specified degrees.
// If pages is empty or nil, all pages are rotated.
// Valid rotation values are 0, 90, 180, or 270.
func RotatePDFPages(inputPath, outputPath string, degrees int, pages []int) error {
	if inputPath == "" {
		return fmt.Errorf("input path is required")
	}
	if outputPath == "" {
		return fmt.Errorf("output path is required")
	}

	// Normalize rotation to valid values
	degrees = degrees % 360
	if degrees < 0 {
		degrees += 360
	}
	if degrees != 0 && degrees != 90 && degrees != 180 && degrees != 270 {
		return fmt.Errorf("rotation must be 0, 90, 180, or 270 degrees, got %d", degrees)
	}

	f, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
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

	// Build a set of pages to rotate
	pageSet := make(map[int]bool)
	if len(pages) == 0 {
		// All pages
		for i := 1; i <= numPages; i++ {
			pageSet[i] = true
		}
	} else {
		for _, p := range pages {
			if p >= 1 && p <= numPages {
				pageSet[p] = true
			}
		}
	}

	writer := model.NewPdfWriter()

	for pageNum := 1; pageNum <= numPages; pageNum++ {
		page, err := reader.GetPage(pageNum)
		if err != nil {
			return fmt.Errorf("failed to get page %d: %w", pageNum, err)
		}

		if pageSet[pageNum] && degrees != 0 {
			// Get current rotation
			currentRotation := 0
			if page.Rotate != nil {
				currentRotation = int(*page.Rotate)
			}

			// Apply new rotation
			newRotation := int64((currentRotation + degrees) % 360)
			page.Rotate = &newRotation
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
		return fmt.Errorf("failed to write output: %w", err)
	}

	return nil
}
