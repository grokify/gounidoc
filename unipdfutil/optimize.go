package unipdfutil

import (
	"fmt"
	"os"

	"github.com/unidoc/unipdf/v3/model"
	"github.com/unidoc/unipdf/v3/model/optimize"
)

// CompressOptions configures PDF compression.
type CompressOptions struct {
	ImageQuality int  `json:"image_quality"` // 1-100, lower = more compression
	Linearize    bool `json:"linearize"`     // Optimize for web viewing
}

// DefaultCompressOptions returns default compression options.
func DefaultCompressOptions() CompressOptions {
	return CompressOptions{
		ImageQuality: 80,
		Linearize:    false,
	}
}

// CompressPDF optimizes a PDF file to reduce its size.
func CompressPDF(inputPath, outputPath string, opts CompressOptions) error {
	if inputPath == "" {
		return fmt.Errorf("input path is required")
	}
	if outputPath == "" {
		return fmt.Errorf("output path is required")
	}

	// Validate image quality
	if opts.ImageQuality < 1 {
		opts.ImageQuality = 1
	}
	if opts.ImageQuality > 100 {
		opts.ImageQuality = 100
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

	// Set optimizer
	optimizer := optimize.New(optimize.Options{
		CombineDuplicateDirectObjects:   true,
		CombineDuplicateStreams:         true,
		CombineIdenticalIndirectObjects: true,
		CompressStreams:                 true,
		UseObjectStreams:                true,
		ImageQuality:                    opts.ImageQuality,
		ImageUpperPPI:                   150,
	})
	writer.SetOptimizer(optimizer)

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

// CompressionResult contains information about PDF compression.
type CompressionResult struct {
	OriginalSize   int64   `json:"original_size"`
	CompressedSize int64   `json:"compressed_size"`
	Ratio          float64 `json:"ratio"` // compressed/original
	SavingsPercent float64 `json:"savings_percent"`
}

// GetPDFFileSize returns the file size in bytes.
func GetPDFFileSize(path string) (int64, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return 0, fmt.Errorf("failed to stat file: %w", err)
	}
	return fi.Size(), nil
}

// CompressPDFWithStats compresses a PDF and returns compression statistics.
func CompressPDFWithStats(inputPath, outputPath string, opts CompressOptions) (CompressionResult, error) {
	var result CompressionResult

	// Get original size
	originalSize, err := GetPDFFileSize(inputPath)
	if err != nil {
		return result, err
	}
	result.OriginalSize = originalSize

	// Compress
	if err := CompressPDF(inputPath, outputPath, opts); err != nil {
		return result, err
	}

	// Get compressed size
	compressedSize, err := GetPDFFileSize(outputPath)
	if err != nil {
		return result, err
	}
	result.CompressedSize = compressedSize

	// Calculate ratio and savings
	if originalSize > 0 {
		result.Ratio = float64(compressedSize) / float64(originalSize)
		result.SavingsPercent = (1.0 - result.Ratio) * 100.0
	}

	return result, nil
}
