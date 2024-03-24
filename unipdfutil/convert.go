package unipdfutil

import (
	"errors"
	"fmt"
	"os"

	"github.com/grokify/mogo/image/imageutil"
	"github.com/unidoc/unipdf/v3/model"
	"github.com/unidoc/unipdf/v3/render"
)

// ConvertPDFFilePageToPNGFile reads a PDF file and converts the specified page to a PNG file
// at the specified output path.
func ConvertPDFFilePageToPNGFile(inputPath, outputPath string, pageNum, outputWidth uint) error {
	if pageNum == 0 {
		return errors.New("page num (1-indexed) cannot be zero")
	} else if outputWidth == 0 {
		return errors.New("output width cannot be zero")
	}

	device := render.NewImageDevice()
	device.OutputWidth = int(outputWidth)

	f, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer f.Close()
	if pdfReader, err := model.NewPdfReader(f); err != nil {
		return err
	} else if numPages, err := pdfReader.GetNumPages(); err != nil {
		return err
	} else if pageNum > uint(numPages) {
		return fmt.Errorf("requested page (%d) greater than page count (%d)", pageNum, numPages)
	} else if pg, err := pdfReader.GetPage(int(pageNum)); err != nil {
		return err
	} else if img, err := device.Render(pg); err != nil {
		return err
	} else {
		img2 := imageutil.Image{Image: img}
		return img2.WritePNGFile(outputPath)
	}
}
