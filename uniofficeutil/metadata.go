package uniofficeutil

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/unidoc/unioffice/common"
	"github.com/unidoc/unioffice/document"
	"github.com/unidoc/unioffice/presentation"
	"github.com/unidoc/unioffice/spreadsheet"
)

// Metadata contains common Office document metadata fields.
type Metadata struct {
	Author         string    `json:"author,omitempty"`
	Category       string    `json:"category,omitempty"`
	ContentStatus  string    `json:"content_status,omitempty"`
	Created        time.Time `json:"created,omitempty"`
	Description    string    `json:"description,omitempty"`
	LastModifiedBy string    `json:"last_modified_by,omitempty"`
	Modified       time.Time `json:"modified,omitempty"`
	Title          string    `json:"title,omitempty"`
	FileType       string    `json:"file_type,omitempty"`
}

// ExportMetadata converts CoreProperties to Metadata.
func ExportMetadata(cp common.CoreProperties) Metadata {
	return Metadata{
		Author:         cp.Author(),
		Category:       cp.Category(),
		ContentStatus:  cp.ContentStatus(),
		Created:        cp.Created(),
		Description:    cp.Description(),
		LastModifiedBy: cp.LastModifiedBy(),
		Modified:       cp.Modified(),
		Title:          cp.Title(),
	}
}

// ReadFileMetadata reads metadata from any supported Office document.
// Supports: .docx, .xlsx, .pptx
func ReadFileMetadata(filename string) (Metadata, error) {
	ext := strings.ToLower(filepath.Ext(filename))

	switch ext {
	case ".docx":
		return ReadFileMetadataDocx(filename)
	case ".xlsx":
		return ReadFileMetadataXlsx(filename)
	case ".pptx":
		return ReadFileMetadataPptx(filename)
	default:
		return Metadata{}, fmt.Errorf("unsupported file type: %s", ext)
	}
}

// ReadFileMetadataDocx reads metadata from a DOCX file.
func ReadFileMetadataDocx(filename string) (Metadata, error) {
	doc, err := document.Open(filename)
	if err != nil {
		return Metadata{}, fmt.Errorf("failed to open DOCX: %w", err)
	}
	defer doc.Close()

	meta := ExportMetadata(doc.CoreProperties)
	meta.FileType = "docx"
	return meta, nil
}

// ReadFileMetadataXlsx reads metadata from an XLSX file.
func ReadFileMetadataXlsx(filename string) (Metadata, error) {
	wb, err := spreadsheet.Open(filename)
	if err != nil {
		return Metadata{}, fmt.Errorf("failed to open XLSX: %w", err)
	}
	defer wb.Close()

	meta := ExportMetadata(wb.CoreProperties)
	meta.FileType = "xlsx"
	return meta, nil
}

// ReadFileMetadataPptx reads metadata from a PPTX file.
func ReadFileMetadataPptx(filename string) (Metadata, error) {
	pres, err := presentation.Open(filename)
	if err != nil {
		return Metadata{}, fmt.Errorf("failed to open PPTX: %w", err)
	}
	defer pres.Close()

	meta := ExportMetadata(pres.CoreProperties)
	meta.FileType = "pptx"
	return meta, nil
}

// ReadFileMetadataPresentation is an alias for ReadFileMetadataPptx for backward compatibility.
func ReadFileMetadataPresentation(filename string) (Metadata, error) {
	return ReadFileMetadataPptx(filename)
}

// XlsxInfo contains information about an XLSX workbook.
type XlsxInfo struct {
	Metadata   Metadata `json:"metadata"`
	SheetCount int      `json:"sheet_count"`
	SheetNames []string `json:"sheet_names"`
}

// GetXlsxInfo returns information about an XLSX workbook including sheet names.
func GetXlsxInfo(filename string) (XlsxInfo, error) {
	wb, err := spreadsheet.Open(filename)
	if err != nil {
		return XlsxInfo{}, fmt.Errorf("failed to open XLSX: %w", err)
	}
	defer wb.Close()

	info := XlsxInfo{
		Metadata: ExportMetadata(wb.CoreProperties),
	}
	info.Metadata.FileType = "xlsx"

	sheets := wb.Sheets()
	info.SheetCount = len(sheets)
	info.SheetNames = make([]string, 0, len(sheets))
	for _, sheet := range sheets {
		info.SheetNames = append(info.SheetNames, sheet.Name())
	}

	return info, nil
}

// PptxInfo contains information about a PPTX presentation.
type PptxInfo struct {
	Metadata   Metadata `json:"metadata"`
	SlideCount int      `json:"slide_count"`
}

// GetPptxInfo returns information about a PPTX presentation.
func GetPptxInfo(filename string) (PptxInfo, error) {
	pres, err := presentation.Open(filename)
	if err != nil {
		return PptxInfo{}, fmt.Errorf("failed to open PPTX: %w", err)
	}
	defer pres.Close()

	info := PptxInfo{
		Metadata: ExportMetadata(pres.CoreProperties),
	}
	info.Metadata.FileType = "pptx"
	info.SlideCount = len(pres.Slides())

	return info, nil
}
