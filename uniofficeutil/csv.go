package uniofficeutil

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/unidoc/unioffice/spreadsheet"
)

// CSVOptions configures XLSX to CSV conversion.
type CSVOptions struct {
	SheetName  string `json:"sheet_name"`  // Specific sheet name or empty for first
	SheetIndex int    `json:"sheet_index"` // Alternative to name (0-indexed)
	Delimiter  string `json:"delimiter"`   // Default comma
	AllSheets  bool   `json:"all_sheets"`  // Export all sheets
}

// DefaultCSVOptions returns default CSV conversion options.
func DefaultCSVOptions() CSVOptions {
	return CSVOptions{
		Delimiter: ",",
	}
}

// ConvertXlsxToCSV exports a spreadsheet to CSV format.
func ConvertXlsxToCSV(inputPath, outputPath string, opts CSVOptions) error {
	if inputPath == "" {
		return fmt.Errorf("input path is required")
	}
	if outputPath == "" {
		return fmt.Errorf("output path is required")
	}

	wb, err := spreadsheet.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open XLSX: %w", err)
	}
	defer wb.Close()

	sheets := wb.Sheets()
	if len(sheets) == 0 {
		return fmt.Errorf("workbook has no sheets")
	}

	// Determine which sheet to export
	var sheet spreadsheet.Sheet
	if opts.SheetName != "" {
		// Find sheet by name
		found := false
		for _, s := range sheets {
			if s.Name() == opts.SheetName {
				sheet = s
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("sheet not found: %s", opts.SheetName)
		}
	} else if opts.SheetIndex >= 0 && opts.SheetIndex < len(sheets) {
		sheet = sheets[opts.SheetIndex]
	} else {
		// Default to first sheet
		sheet = sheets[0]
	}

	return writeSheetToCSV(sheet, outputPath, opts.Delimiter)
}

// ConvertXlsxToCSVAllSheets exports all sheets to separate CSV files.
// Returns the list of created file paths.
func ConvertXlsxToCSVAllSheets(inputPath, outputDir string, opts CSVOptions) ([]string, error) {
	if inputPath == "" {
		return nil, fmt.Errorf("input path is required")
	}
	if outputDir == "" {
		outputDir = filepath.Dir(inputPath)
	}

	// Ensure output directory exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create output directory: %w", err)
	}

	wb, err := spreadsheet.Open(inputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open XLSX: %w", err)
	}
	defer wb.Close()

	sheets := wb.Sheets()
	if len(sheets) == 0 {
		return nil, fmt.Errorf("workbook has no sheets")
	}

	baseName := strings.TrimSuffix(filepath.Base(inputPath), filepath.Ext(inputPath))
	var outputPaths []string

	for i, sheet := range sheets {
		// Generate output filename
		sheetName := sheet.Name()
		if sheetName == "" {
			sheetName = fmt.Sprintf("sheet%d", i+1)
		}
		// Sanitize sheet name for filename
		sheetName = sanitizeFilename(sheetName)

		outputPath := filepath.Join(outputDir, fmt.Sprintf("%s_%s.csv", baseName, sheetName))

		if err := writeSheetToCSV(sheet, outputPath, opts.Delimiter); err != nil {
			return outputPaths, fmt.Errorf("failed to export sheet %s: %w", sheet.Name(), err)
		}

		outputPaths = append(outputPaths, outputPath)
	}

	return outputPaths, nil
}

// writeSheetToCSV writes a single sheet to a CSV file.
func writeSheetToCSV(sheet spreadsheet.Sheet, outputPath, delimiter string) error {
	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	writer := csv.NewWriter(outFile)
	if delimiter != "" && len(delimiter) == 1 {
		writer.Comma = rune(delimiter[0])
	}
	defer writer.Flush()

	for _, row := range sheet.Rows() {
		var rowData []string
		for _, cell := range row.Cells() {
			rowData = append(rowData, cell.GetFormattedValue())
		}
		if err := writer.Write(rowData); err != nil {
			return fmt.Errorf("failed to write row: %w", err)
		}
	}

	return nil
}

// sanitizeFilename removes or replaces characters that are not safe in filenames.
func sanitizeFilename(name string) string {
	// Replace common problematic characters
	replacer := strings.NewReplacer(
		"/", "_",
		"\\", "_",
		":", "_",
		"*", "_",
		"?", "_",
		"\"", "_",
		"<", "_",
		">", "_",
		"|", "_",
		" ", "_",
	)
	return replacer.Replace(name)
}
