package unipdfutil

import (
	"fmt"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/unidoc/unipdf/v3/model"
	"github.com/unidoc/unipdf/v3/render"
)

// RenderOptions configures PDF page rendering.
type RenderOptions struct {
	Format  string  `json:"format"`  // png, jpeg
	DPI     float64 `json:"dpi"`     // default 150
	Quality int     `json:"quality"` // jpeg quality 1-100
}

// DefaultRenderOptions returns default rendering options.
func DefaultRenderOptions() RenderOptions {
	return RenderOptions{
		Format:  "png",
		DPI:     150,
		Quality: 85,
	}
}

// RenderPDFToImages renders all pages of a PDF as images.
// Returns the list of created image file paths.
func RenderPDFToImages(inputPath, outputDir string, opts RenderOptions) ([]string, error) {
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
	var outputPaths []string

	device := render.NewImageDevice()
	// Set DPI (default is 72, so we scale width accordingly)
	// OutputWidth = (pageWidth / 72) * DPI
	// We'll calculate per page based on the page dimensions

	for pageNum := 1; pageNum <= numPages; pageNum++ {
		page, err := reader.GetPage(pageNum)
		if err != nil {
			return outputPaths, fmt.Errorf("failed to get page %d: %w", pageNum, err)
		}

		// Get page dimensions and calculate output width based on DPI
		mediaBox, err := page.GetMediaBox()
		if err != nil {
			return outputPaths, fmt.Errorf("failed to get media box for page %d: %w", pageNum, err)
		}
		pageWidth := mediaBox.Urx - mediaBox.Llx
		device.OutputWidth = int(pageWidth * opts.DPI / 72.0)

		img, err := device.Render(page)
		if err != nil {
			return outputPaths, fmt.Errorf("failed to render page %d: %w", pageNum, err)
		}

		ext := strings.ToLower(opts.Format)
		if ext != "png" && ext != "jpeg" && ext != "jpg" {
			ext = "png"
		}
		if ext == "jpg" {
			ext = "jpeg"
		}

		outputPath := filepath.Join(outputDir, fmt.Sprintf("%s_page_%03d.%s", baseName, pageNum, ext))
		outFile, err := os.Create(outputPath)
		if err != nil {
			return outputPaths, fmt.Errorf("failed to create output file for page %d: %w", pageNum, err)
		}

		switch ext {
		case "jpeg":
			err = jpeg.Encode(outFile, img, &jpeg.Options{Quality: opts.Quality})
		default:
			err = png.Encode(outFile, img)
		}
		outFile.Close()

		if err != nil {
			return outputPaths, fmt.Errorf("failed to encode page %d: %w", pageNum, err)
		}

		outputPaths = append(outputPaths, outputPath)
	}

	return outputPaths, nil
}

// RenderPDFPageToImage renders a single page from a PDF as an image.
func RenderPDFPageToImage(inputPath, outputPath string, pageNum int, opts RenderOptions) error {
	if inputPath == "" {
		return fmt.Errorf("input path is required")
	}
	if outputPath == "" {
		return fmt.Errorf("output path is required")
	}
	if pageNum < 1 {
		return fmt.Errorf("page number must be >= 1")
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

	if pageNum > numPages {
		return fmt.Errorf("page %d out of range (1-%d)", pageNum, numPages)
	}

	page, err := reader.GetPage(pageNum)
	if err != nil {
		return fmt.Errorf("failed to get page %d: %w", pageNum, err)
	}

	device := render.NewImageDevice()

	// Get page dimensions and calculate output width based on DPI
	mediaBox, err := page.GetMediaBox()
	if err != nil {
		return fmt.Errorf("failed to get media box: %w", err)
	}
	pageWidth := mediaBox.Urx - mediaBox.Llx
	device.OutputWidth = int(pageWidth * opts.DPI / 72.0)

	img, err := device.Render(page)
	if err != nil {
		return fmt.Errorf("failed to render page: %w", err)
	}

	ext := strings.ToLower(filepath.Ext(outputPath))
	if ext == "" {
		ext = "." + strings.ToLower(opts.Format)
		outputPath = outputPath + ext
	}

	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	switch ext {
	case ".jpg", ".jpeg":
		return jpeg.Encode(outFile, img, &jpeg.Options{Quality: opts.Quality})
	default:
		return png.Encode(outFile, img)
	}
}
