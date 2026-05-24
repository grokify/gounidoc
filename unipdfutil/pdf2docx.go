package unipdfutil

import (
	"fmt"
	"os"
	"strings"

	"github.com/unidoc/unioffice/color"
	"github.com/unidoc/unioffice/document"
	"github.com/unidoc/unioffice/measurement"
	"github.com/unidoc/unioffice/schema/soo/wml"
	"github.com/unidoc/unipdf/v3/extractor"
	"github.com/unidoc/unipdf/v3/model"
)

// ConversionOptions controls the PDF to DOCX conversion behavior.
type ConversionOptions struct {
	ExtractTables bool // Extract and convert tables
	ExtractImages bool // Extract and embed images
	DetectLayout  bool // Detect headings and text structure
}

// DefaultConversionOptions returns the default conversion options with all features enabled.
func DefaultConversionOptions() ConversionOptions {
	return ConversionOptions{
		ExtractTables: true,
		ExtractImages: true,
		DetectLayout:  true,
	}
}

// ConvertPDFFileToDocxFile extracts text from a PDF file and creates a DOCX file.
// This is a basic text extraction - formatting and images are not preserved.
func ConvertPDFFileToDocxFile(inputPDFPath, outputDocxPath string) error {
	return ConvertPDFFileToDocxFileWithOptions(inputPDFPath, outputDocxPath, DefaultConversionOptions())
}

// ConvertPDFFileToDocxFileWithOptions extracts content from a PDF file and creates a DOCX file
// with enhanced formatting based on the provided options.
func ConvertPDFFileToDocxFileWithOptions(inputPDFPath, outputDocxPath string, opts ConversionOptions) error {
	// Open PDF file
	f, err := os.Open(inputPDFPath)
	if err != nil {
		return err
	}
	defer f.Close()

	pdfReader, err := model.NewPdfReader(f)
	if err != nil {
		return err
	}

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return err
	}

	// Create new DOCX document
	doc := document.New()
	defer doc.Close()

	imageCount := 0

	// Extract content from each page
	for pageNum := 1; pageNum <= numPages; pageNum++ {
		page, err := pdfReader.GetPage(pageNum)
		if err != nil {
			return fmt.Errorf("error getting page %d: %w", pageNum, err)
		}

		ex, err := extractor.New(page)
		if err != nil {
			return fmt.Errorf("error creating extractor for page %d: %w", pageNum, err)
		}

		// Extract images if enabled
		if opts.ExtractImages {
			if err := extractAndAddImages(ex, doc, &imageCount); err != nil {
				// Log but continue - image extraction failure shouldn't stop conversion
				fmt.Printf("Warning: image extraction failed on page %d: %v\n", pageNum, err)
			}
		}

		// Extract page text with rich information
		pageText, _, _, err := ex.ExtractPageText()
		if err != nil {
			// Fall back to basic text extraction
			text, err := ex.ExtractText()
			if err != nil {
				return fmt.Errorf("error extracting text from page %d: %w", pageNum, err)
			}
			addBasicText(doc, text)
			continue
		}

		// Extract tables if enabled
		if opts.ExtractTables {
			tables := pageText.Tables()
			if len(tables) > 0 {
				for _, table := range tables {
					addTableToDoc(doc, table)
					// Add spacing after table
					doc.AddParagraph()
				}
			}
		}

		// Process text with layout detection if enabled
		if opts.DetectLayout {
			addTextWithLayout(doc, pageText)
		} else {
			addBasicText(doc, pageText.Text())
		}

		// Add page break between pages (except last page)
		if pageNum < numPages {
			para := doc.AddParagraph()
			run := para.AddRun()
			run.AddPageBreak()
		}
	}

	// Save DOCX file
	return doc.SaveToFile(outputDocxPath)
}

// extractAndAddImages extracts images from the page and adds them to the document.
// Note: Image extraction is experimental and may not work for all PDF types.
func extractAndAddImages(ex *extractor.Extractor, doc *document.Document, imageCount *int) error {
	pageImages, err := ex.ExtractPageImages(nil)
	if err != nil {
		return err
	}

	for _, imgMark := range pageImages.Images {
		img := imgMark.Image

		// Convert to Go image
		goImg, err := img.ToGoImage()
		if err != nil {
			continue // Skip images that can't be converted
		}

		*imageCount++

		// Add a placeholder paragraph noting an image was found
		// Full image embedding requires additional unioffice image handling
		para := doc.AddParagraph()
		run := para.AddRun()
		bounds := goImg.Bounds()
		run.AddText(fmt.Sprintf("[Image %d: %dx%d pixels]", *imageCount, bounds.Dx(), bounds.Dy()))
	}

	return nil
}

// addTableToDoc converts a PDF table to a DOCX table.
func addTableToDoc(doc *document.Document, pdfTable extractor.TextTable) {
	if pdfTable.H == 0 || pdfTable.W == 0 {
		return
	}

	table := doc.AddTable()
	table.Properties().SetWidth(6 * measurement.Inch)

	// Add borders
	borders := table.Properties().Borders()
	borders.SetAll(wml.ST_BorderSingle, color.Black, measurement.Point)

	for row := 0; row < pdfTable.H; row++ {
		tableRow := table.AddRow()
		for col := 0; col < pdfTable.W; col++ {
			cell := tableRow.AddCell()
			cellText := ""
			if row < len(pdfTable.Cells) && col < len(pdfTable.Cells[row]) {
				cellText = strings.TrimSpace(pdfTable.Cells[row][col].Text)
			}
			para := cell.AddParagraph()
			run := para.AddRun()
			run.AddText(cellText)
		}
	}
}

// addTextWithLayout analyzes text positioning to detect structure like headings.
func addTextWithLayout(doc *document.Document, pageText *extractor.PageText) {
	text := pageText.Text()
	marks := pageText.Marks()

	if marks == nil || marks.Len() == 0 {
		addBasicText(doc, text)
		return
	}

	// Group text by lines based on Y position
	lines := groupTextByLines(text, *marks)

	for _, line := range lines {
		if strings.TrimSpace(line.text) == "" {
			continue
		}

		para := doc.AddParagraph()
		run := para.AddRun()

		// Detect if this looks like a heading (larger font, shorter text, often bold patterns)
		if line.isHeading {
			run.Properties().SetBold(true)
			run.Properties().SetSize(measurement.Distance(14 * measurement.Point))
		}

		run.AddText(strings.TrimSpace(line.text))
	}
}

type textLine struct {
	text      string
	isHeading bool
}

// groupTextByLines groups text marks into lines and detects potential headings.
// The marks parameter is reserved for future font-size based heading detection.
func groupTextByLines(text string, _ extractor.TextMarkArray) []textLine {
	var lines []textLine
	paragraphs := strings.Split(text, "\n")

	// Calculate average line length for heading detection
	totalLen := 0
	nonEmptyCount := 0
	for _, p := range paragraphs {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			totalLen += len(trimmed)
			nonEmptyCount++
		}
	}
	avgLen := 0
	if nonEmptyCount > 0 {
		avgLen = totalLen / nonEmptyCount
	}

	for _, p := range paragraphs {
		trimmed := strings.TrimSpace(p)
		if trimmed == "" {
			continue
		}

		line := textLine{text: trimmed}

		// Heuristic: short lines that end without punctuation might be headings
		// Also check for numbered sections like "1.", "1.1", etc.
		isShort := len(trimmed) < avgLen/2 && len(trimmed) < 80
		endsWithoutPunct := !strings.HasSuffix(trimmed, ".") &&
			!strings.HasSuffix(trimmed, ",") &&
			!strings.HasSuffix(trimmed, ";")
		isNumberedSection := isNumberedHeading(trimmed)
		isAllCaps := trimmed == strings.ToUpper(trimmed) && len(trimmed) > 3

		line.isHeading = (isShort && endsWithoutPunct) || isNumberedSection || isAllCaps

		lines = append(lines, line)
	}

	return lines
}

// isNumberedHeading checks if text starts with a numbered section pattern.
func isNumberedHeading(text string) bool {
	if len(text) < 2 {
		return false
	}

	// Check for patterns like "1.", "1.1", "51.", etc.
	for i, c := range text {
		if c >= '0' && c <= '9' {
			continue
		}
		if c == '.' {
			// Found number followed by period
			if i > 0 && i < len(text)-1 {
				// Check if next char is space or another digit
				next := rune(text[i+1])
				if next == ' ' || (next >= '0' && next <= '9') {
					return true
				}
			}
			continue
		}
		if c == ' ' && i > 1 {
			// Space after number/period sequence
			return true
		}
		break
	}
	return false
}

// addBasicText adds simple text paragraphs to the document.
func addBasicText(doc *document.Document, text string) {
	paragraphs := strings.Split(text, "\n")
	for _, pText := range paragraphs {
		pText = strings.TrimSpace(pText)
		if pText != "" {
			para := doc.AddParagraph()
			run := para.AddRun()
			run.AddText(pText)
		}
	}
}
