package unipdfutil

import (
	"fmt"
	"os"

	"github.com/unidoc/unipdf/v3/creator"
	"github.com/unidoc/unipdf/v3/model"
)

// WatermarkOptions configures watermark appearance.
type WatermarkOptions struct {
	Text      string  `json:"text"`       // Text watermark content
	ImagePath string  `json:"image_path"` // Path to image watermark (alternative to text)
	FontSize  float64 `json:"font_size"`  // Font size for text watermark
	FontName  string  `json:"font_name"`  // Font name (e.g., Helvetica)
	Color     string  `json:"color"`      // Hex color (e.g., "#FF0000")
	Opacity   float64 `json:"opacity"`    // 0.0-1.0
	Rotation  float64 `json:"rotation"`   // Degrees
	Position  string  `json:"position"`   // center, top-left, top-right, bottom-left, bottom-right
	Pages     []int   `json:"pages"`      // Empty = all pages
}

// DefaultWatermarkOptions returns default watermark options.
func DefaultWatermarkOptions() WatermarkOptions {
	return WatermarkOptions{
		FontSize: 48,
		FontName: "Helvetica",
		Color:    "#888888",
		Opacity:  0.3,
		Rotation: -45,
		Position: "center",
	}
}

// AddWatermark adds a text or image watermark to a PDF.
func AddWatermark(inputPath, outputPath string, opts WatermarkOptions) error {
	if inputPath == "" {
		return fmt.Errorf("input path is required")
	}
	if outputPath == "" {
		return fmt.Errorf("output path is required")
	}
	if opts.Text == "" && opts.ImagePath == "" {
		return fmt.Errorf("either text or image_path is required")
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

	c := creator.New()

	// Build a set of pages to watermark
	pageSet := make(map[int]bool)
	if len(opts.Pages) == 0 {
		// All pages
		for i := 1; i <= numPages; i++ {
			pageSet[i] = true
		}
	} else {
		for _, p := range opts.Pages {
			if p >= 1 && p <= numPages {
				pageSet[p] = true
			}
		}
	}

	// Parse color
	r, g, b := parseHexColor(opts.Color)
	color := creator.ColorRGBFromArithmetic(r, g, b)

	for pageNum := 1; pageNum <= numPages; pageNum++ {
		page, err := reader.GetPage(pageNum)
		if err != nil {
			return fmt.Errorf("failed to get page %d: %w", pageNum, err)
		}

		if err := c.AddPage(page); err != nil {
			return fmt.Errorf("failed to add page %d: %w", pageNum, err)
		}

		if !pageSet[pageNum] {
			continue
		}

		// Get page dimensions
		mediaBox, err := page.GetMediaBox()
		if err != nil {
			continue
		}
		pageWidth := mediaBox.Urx - mediaBox.Llx
		pageHeight := mediaBox.Ury - mediaBox.Lly

		if opts.ImagePath != "" {
			// Image watermark
			img, err := c.NewImageFromFile(opts.ImagePath)
			if err != nil {
				return fmt.Errorf("failed to load watermark image: %w", err)
			}

			// Scale image to fit page
			img.ScaleToWidth(pageWidth * 0.5)
			img.SetOpacity(opts.Opacity)
			img.SetAngle(opts.Rotation)

			x, y := calculatePosition(opts.Position, pageWidth, pageHeight, img.Width(), img.Height())
			img.SetPos(x, y)

			if err := c.Draw(img); err != nil {
				return fmt.Errorf("failed to draw image watermark: %w", err)
			}
		} else {
			// Text watermark
			p := c.NewParagraph(opts.Text)
			p.SetFontSize(opts.FontSize)
			p.SetColor(color)
			p.SetAngle(opts.Rotation)

			// Get text dimensions (approximate)
			textWidth := float64(len(opts.Text)) * opts.FontSize * 0.5
			textHeight := opts.FontSize * 1.2

			x, y := calculatePosition(opts.Position, pageWidth, pageHeight, textWidth, textHeight)
			p.SetPos(x, y)

			if err := c.Draw(p); err != nil {
				return fmt.Errorf("failed to draw text watermark: %w", err)
			}
		}
	}

	if err := c.WriteToFile(outputPath); err != nil {
		return fmt.Errorf("failed to write output: %w", err)
	}

	return nil
}

// parseHexColor parses a hex color string to RGB values (0-1).
func parseHexColor(hex string) (float64, float64, float64) {
	if len(hex) == 0 {
		return 0.5, 0.5, 0.5 // Default gray
	}

	// Remove # prefix if present
	if hex[0] == '#' {
		hex = hex[1:]
	}

	if len(hex) != 6 {
		return 0.5, 0.5, 0.5 // Default gray
	}

	var r, g, b uint8
	_, _ = fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b)

	return float64(r) / 255.0, float64(g) / 255.0, float64(b) / 255.0
}

// calculatePosition calculates x, y position for watermark placement.
func calculatePosition(position string, pageWidth, pageHeight, objWidth, objHeight float64) (float64, float64) {
	switch position {
	case "top-left":
		return 20, 20
	case "top-right":
		return pageWidth - objWidth - 20, 20
	case "bottom-left":
		return 20, pageHeight - objHeight - 20
	case "bottom-right":
		return pageWidth - objWidth - 20, pageHeight - objHeight - 20
	default: // center
		return (pageWidth - objWidth) / 2, (pageHeight - objHeight) / 2
	}
}
