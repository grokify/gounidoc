package uniofficeutil

import (
	"fmt"

	"github.com/unidoc/unioffice/document"
	"github.com/unidoc/unioffice/document/convert"
)

func ConvertDOCXFileToPDFFile(inputFilenameDOCX, outputFilenamePDF string) error {
	doc, err := document.Open(inputFilenameDOCX)
	if err != nil {
		return fmt.Errorf("error opening document (%s): %w", inputFilenameDOCX, err)
	}
	defer doc.Close()

	c := convert.ConvertToPdf(doc)

	if err := c.WriteToFile(outputFilenamePDF); err != nil {
		return fmt.Errorf("error converting/writing document (%s) from (%s): %w", outputFilenamePDF, inputFilenameDOCX, err)
	}
	return nil
}
