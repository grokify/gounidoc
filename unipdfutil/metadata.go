package unipdfutil

import (
	"fmt"
	"os"
	"time"

	"github.com/unidoc/unipdf/v3/core"
	"github.com/unidoc/unipdf/v3/model"
)

// makePdfString creates a PdfObjectString from a Go string.
func makePdfString(s string) *core.PdfObjectString {
	obj := core.MakeString(s)
	return obj
}

// PDFMetadata contains PDF document properties.
type PDFMetadata struct {
	Title        string    `json:"title,omitempty"`
	Author       string    `json:"author,omitempty"`
	Subject      string    `json:"subject,omitempty"`
	Keywords     string    `json:"keywords,omitempty"`
	Creator      string    `json:"creator,omitempty"`
	Producer     string    `json:"producer,omitempty"`
	CreationDate time.Time `json:"creation_date,omitempty"`
	ModDate      time.Time `json:"mod_date,omitempty"`
	PageCount    int       `json:"page_count"`
	Version      string    `json:"version,omitempty"`
	Encrypted    bool      `json:"encrypted"`
	FileSize     int64     `json:"file_size"`
}

// ReadPDFMetadata reads metadata from a PDF file.
func ReadPDFMetadata(inputPath string) (PDFMetadata, error) {
	var meta PDFMetadata

	if inputPath == "" {
		return meta, fmt.Errorf("input path is required")
	}

	// Get file size
	fi, err := os.Stat(inputPath)
	if err != nil {
		return meta, fmt.Errorf("failed to stat file: %w", err)
	}
	meta.FileSize = fi.Size()

	f, err := os.Open(inputPath)
	if err != nil {
		return meta, fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	reader, err := model.NewPdfReader(f)
	if err != nil {
		return meta, fmt.Errorf("failed to create PDF reader: %w", err)
	}

	// Get page count
	numPages, err := reader.GetNumPages()
	if err != nil {
		return meta, fmt.Errorf("failed to get page count: %w", err)
	}
	meta.PageCount = numPages

	// Check encryption
	meta.Encrypted, _ = reader.IsEncrypted()

	// Get PDF version
	pdfVer := reader.PdfVersion()
	meta.Version = fmt.Sprintf("%d.%d", pdfVer.Major, pdfVer.Minor)

	// Get document info
	info, err := reader.GetPdfInfo()
	if err == nil && info != nil {
		if info.Title != nil {
			meta.Title = info.Title.Decoded()
		}
		if info.Author != nil {
			meta.Author = info.Author.Decoded()
		}
		if info.Subject != nil {
			meta.Subject = info.Subject.Decoded()
		}
		if info.Keywords != nil {
			meta.Keywords = info.Keywords.Decoded()
		}
		if info.Creator != nil {
			meta.Creator = info.Creator.Decoded()
		}
		if info.Producer != nil {
			meta.Producer = info.Producer.Decoded()
		}
		if info.CreationDate != nil {
			meta.CreationDate = info.CreationDate.ToGoTime()
		}
		if info.ModifiedDate != nil {
			meta.ModDate = info.ModifiedDate.ToGoTime()
		}
	}

	return meta, nil
}

// WritePDFMetadata updates metadata in a PDF file.
// Only non-empty fields in meta will be updated.
func WritePDFMetadata(inputPath, outputPath string, meta PDFMetadata) error {
	if inputPath == "" {
		return fmt.Errorf("input path is required")
	}
	if outputPath == "" {
		return fmt.Errorf("output path is required")
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

	writer := model.NewPdfWriter()

	// Copy all pages
	for pageNum := 1; pageNum <= numPages; pageNum++ {
		page, err := reader.GetPage(pageNum)
		if err != nil {
			return fmt.Errorf("failed to get page %d: %w", pageNum, err)
		}
		if err := writer.AddPage(page); err != nil {
			return fmt.Errorf("failed to add page %d: %w", pageNum, err)
		}
	}

	// Set metadata
	info := model.PdfInfo{}

	// Get existing info if available
	existingInfo, err := reader.GetPdfInfo()
	if err == nil && existingInfo != nil {
		info = *existingInfo
	}

	// Update with new values using core.MakeString
	if meta.Title != "" {
		info.Title = makePdfString(meta.Title)
	}
	if meta.Author != "" {
		info.Author = makePdfString(meta.Author)
	}
	if meta.Subject != "" {
		info.Subject = makePdfString(meta.Subject)
	}
	if meta.Keywords != "" {
		info.Keywords = makePdfString(meta.Keywords)
	}
	if meta.Creator != "" {
		info.Creator = makePdfString(meta.Creator)
	}
	if meta.Producer != "" {
		info.Producer = makePdfString(meta.Producer)
	}

	// Update modification date
	modDate, _ := model.NewPdfDateFromTime(time.Now())
	info.ModifiedDate = &modDate

	writer.SetDocInfo(&info)

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
