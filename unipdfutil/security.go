package unipdfutil

import (
	"fmt"
	"os"

	"github.com/unidoc/unipdf/v3/model"
)

// PDFPermissions configures PDF access permissions.
type PDFPermissions struct {
	Printing       bool `json:"printing"`
	ModifyContents bool `json:"modify_contents"`
	CopyContents   bool `json:"copy_contents"`
	ModifyAnnots   bool `json:"modify_annots"`
}

// DefaultPDFPermissions returns default permissions (all allowed).
func DefaultPDFPermissions() PDFPermissions {
	return PDFPermissions{
		Printing:       true,
		ModifyContents: true,
		CopyContents:   true,
		ModifyAnnots:   true,
	}
}

// EncryptOptions configures PDF encryption.
type EncryptOptions struct {
	UserPassword  string         `json:"user_password"`  // Password to open PDF
	OwnerPassword string         `json:"owner_password"` // Password for full access
	Permissions   PDFPermissions `json:"permissions"`
}

// EncryptPDF encrypts a PDF file with password protection.
func EncryptPDF(inputPath, outputPath string, opts EncryptOptions) error {
	if inputPath == "" {
		return fmt.Errorf("input path is required")
	}
	if outputPath == "" {
		return fmt.Errorf("output path is required")
	}
	if opts.UserPassword == "" && opts.OwnerPassword == "" {
		return fmt.Errorf("at least one password (user or owner) is required")
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

	// Check if already encrypted
	isEncrypted, err := reader.IsEncrypted()
	if err != nil {
		return fmt.Errorf("failed to check encryption status: %w", err)
	}
	if isEncrypted {
		return fmt.Errorf("PDF is already encrypted")
	}

	numPages, err := reader.GetNumPages()
	if err != nil {
		return fmt.Errorf("failed to get page count: %w", err)
	}

	writer := model.NewPdfWriter()

	for pageNum := 1; pageNum <= numPages; pageNum++ {
		page, err := reader.GetPage(pageNum)
		if err != nil {
			return fmt.Errorf("failed to get page %d: %w", pageNum, err)
		}
		if err := writer.AddPage(page); err != nil {
			return fmt.Errorf("failed to add page %d: %w", pageNum, err)
		}
	}

	// Set encryption
	// Use simple encryption without custom permissions for compatibility
	err = writer.Encrypt([]byte(opts.UserPassword), []byte(opts.OwnerPassword), nil)
	if err != nil {
		return fmt.Errorf("failed to set encryption: %w", err)
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

// DecryptPDF removes password protection from a PDF file.
func DecryptPDF(inputPath, outputPath, password string) error {
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

	// Check if encrypted
	isEncrypted, err := reader.IsEncrypted()
	if err != nil {
		return fmt.Errorf("failed to check encryption status: %w", err)
	}
	if !isEncrypted {
		return fmt.Errorf("PDF is not encrypted")
	}

	// Try to decrypt
	auth, err := reader.Decrypt([]byte(password))
	if err != nil {
		return fmt.Errorf("failed to decrypt: %w", err)
	}
	if !auth {
		return fmt.Errorf("incorrect password")
	}

	numPages, err := reader.GetNumPages()
	if err != nil {
		return fmt.Errorf("failed to get page count: %w", err)
	}

	writer := model.NewPdfWriter()

	for pageNum := 1; pageNum <= numPages; pageNum++ {
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
		return fmt.Errorf("failed to write output: %w", err)
	}

	return nil
}

// IsPDFEncrypted checks if a PDF file is password protected.
func IsPDFEncrypted(inputPath string) (bool, error) {
	if inputPath == "" {
		return false, fmt.Errorf("input path is required")
	}

	f, err := os.Open(inputPath)
	if err != nil {
		return false, fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	reader, err := model.NewPdfReader(f)
	if err != nil {
		return false, fmt.Errorf("failed to create PDF reader: %w", err)
	}

	return reader.IsEncrypted()
}
