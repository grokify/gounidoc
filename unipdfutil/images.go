package unipdfutil

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/unidoc/unipdf/v3/creator"
	"github.com/unidoc/unipdf/v3/extractor"
	"github.com/unidoc/unipdf/v3/model"
)

// ExtractedImage represents an extracted image from a PDF.
type ExtractedImage struct {
	PageNumber int    `json:"page_number"`
	Index      int    `json:"index"`
	Format     string `json:"format"`
	Width      int    `json:"width"`
	Height     int    `json:"height"`
	FilePath   string `json:"file_path,omitempty"`
}

// ExtractPDFImages extracts all images from a PDF file to an output directory.
// Returns information about each extracted image.
func ExtractPDFImages(inputPath, outputDir string) ([]ExtractedImage, error) {
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

	baseName := strings.TrimSuffix(filepath.Base(inputPath), filepath.Ext(inputPath))
	var extracted []ExtractedImage

	for pageNum := 1; pageNum <= numPages; pageNum++ {
		page, err := reader.GetPage(pageNum)
		if err != nil {
			return extracted, fmt.Errorf("failed to get page %d: %w", pageNum, err)
		}

		ex, err := extractor.New(page)
		if err != nil {
			return extracted, fmt.Errorf("failed to create extractor for page %d: %w", pageNum, err)
		}

		images, err := ex.ExtractPageImages(nil)
		if err != nil {
			// Some pages may not have extractable images
			continue
		}

		for idx, imgMark := range images.Images {
			img, err := imgMark.Image.ToGoImage()
			if err != nil {
				continue
			}

			bounds := img.Bounds()
			width := bounds.Max.X - bounds.Min.X
			height := bounds.Max.Y - bounds.Min.Y

			// Determine format based on image content type or default to PNG
			format := "png"
			outputPath := filepath.Join(outputDir, fmt.Sprintf("%s_page%03d_img%03d.%s", baseName, pageNum, idx+1, format))

			outFile, err := os.Create(outputPath)
			if err != nil {
				return extracted, fmt.Errorf("failed to create image file: %w", err)
			}

			if err := png.Encode(outFile, img); err != nil {
				outFile.Close()
				return extracted, fmt.Errorf("failed to encode image: %w", err)
			}
			outFile.Close()

			extracted = append(extracted, ExtractedImage{
				PageNumber: pageNum,
				Index:      idx + 1,
				Format:     format,
				Width:      width,
				Height:     height,
				FilePath:   outputPath,
			})
		}
	}

	return extracted, nil
}

// ImageToPDFOptions configures image to PDF conversion.
type ImageToPDFOptions struct {
	PageSize string  `json:"page_size"` // letter, a4, or fit
	Margin   float64 `json:"margin"`    // margin in points
}

// DefaultImageToPDFOptions returns default options for image to PDF conversion.
func DefaultImageToPDFOptions() ImageToPDFOptions {
	return ImageToPDFOptions{
		PageSize: "letter",
		Margin:   36.0, // 0.5 inch margins
	}
}

// ImagesToPDF combines multiple images into a single PDF file.
func ImagesToPDF(imagePaths []string, outputPath string, opts ImageToPDFOptions) error {
	if len(imagePaths) == 0 {
		return fmt.Errorf("no images provided")
	}
	if outputPath == "" {
		return fmt.Errorf("output path is required")
	}

	c := creator.New()

	// Set page size
	switch strings.ToLower(opts.PageSize) {
	case "a4":
		c.SetPageSize(creator.PageSizeA4)
	case "fit":
		// Page size will be set per image
	default: // letter
		c.SetPageSize(creator.PageSizeLetter)
	}

	for _, imgPath := range imagePaths {
		img, err := loadImage(imgPath)
		if err != nil {
			return fmt.Errorf("failed to load image %s: %w", imgPath, err)
		}

		cImg, err := c.NewImageFromGoImage(img)
		if err != nil {
			return fmt.Errorf("failed to create PDF image from %s: %w", imgPath, err)
		}

		// Get page dimensions based on page size
		var pageWidth, pageHeight float64
		switch strings.ToLower(opts.PageSize) {
		case "a4":
			pageWidth, pageHeight = 595.0, 842.0 // A4 in points
		case "fit":
			pageWidth = cImg.Width() + 2*opts.Margin
			pageHeight = cImg.Height() + 2*opts.Margin
			c.SetPageSize(creator.PageSize{pageWidth, pageHeight})
		default: // letter
			pageWidth, pageHeight = 612.0, 792.0 // Letter in points
		}

		// Set image position with margins
		cImg.SetPos(opts.Margin, opts.Margin)

		// Scale image to fit page if not using "fit" mode
		if strings.ToLower(opts.PageSize) != "fit" {
			maxWidth := pageWidth - 2*opts.Margin
			maxHeight := pageHeight - 2*opts.Margin
			cImg.ScaleToWidth(maxWidth)
			if cImg.Height() > maxHeight {
				cImg.ScaleToHeight(maxHeight)
			}
		}

		c.NewPage()
		if err := c.Draw(cImg); err != nil {
			return fmt.Errorf("failed to draw image %s: %w", imgPath, err)
		}
	}

	if err := c.WriteToFile(outputPath); err != nil {
		return fmt.Errorf("failed to write PDF: %w", err)
	}

	return nil
}

// loadImage loads an image from a file path.
func loadImage(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".png":
		return png.Decode(f)
	case ".jpg", ".jpeg":
		return jpeg.Decode(f)
	default:
		// Try generic decode
		img, _, err := image.Decode(f)
		return img, err
	}
}
