package uniofficeutil

import (
	"fmt"
	"os"
	"strings"

	"github.com/unidoc/unioffice/document"
)

// ExtractDocxText extracts all text from a DOCX file.
func ExtractDocxText(inputPath string) (string, error) {
	if inputPath == "" {
		return "", fmt.Errorf("input path is required")
	}

	doc, err := document.Open(inputPath)
	if err != nil {
		return "", fmt.Errorf("failed to open DOCX: %w", err)
	}
	defer doc.Close()

	var textParts []string

	for _, para := range doc.Paragraphs() {
		var paraText []string
		for _, run := range para.Runs() {
			paraText = append(paraText, run.Text())
		}
		textParts = append(textParts, strings.Join(paraText, ""))
	}

	// Also extract text from tables
	for _, table := range doc.Tables() {
		for _, row := range table.Rows() {
			var rowText []string
			for _, cell := range row.Cells() {
				for _, para := range cell.Paragraphs() {
					var cellText []string
					for _, run := range para.Runs() {
						cellText = append(cellText, run.Text())
					}
					if len(cellText) > 0 {
						rowText = append(rowText, strings.Join(cellText, ""))
					}
				}
			}
			if len(rowText) > 0 {
				textParts = append(textParts, strings.Join(rowText, "\t"))
			}
		}
	}

	return strings.Join(textParts, "\n"), nil
}

// ExtractDocxTextToFile extracts text from a DOCX and writes it to a file.
func ExtractDocxTextToFile(inputPath, outputPath string) error {
	if outputPath == "" {
		return fmt.Errorf("output path is required")
	}

	text, err := ExtractDocxText(inputPath)
	if err != nil {
		return err
	}

	if err := os.WriteFile(outputPath, []byte(text), 0600); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	return nil
}
