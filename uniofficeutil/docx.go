package uniofficeutil

import (
	"github.com/grokify/mogo/errors/errorsutil"
	"github.com/unidoc/unioffice/document"
	"github.com/unidoc/unioffice/document/convert"
)

func ConvertDOCXFileToPDFFile(inputFilenameDOCX, outputFilenamePDF string) error {
	doc, err := document.Open(inputFilenameDOCX)
	if err != nil {
		err = errorsutil.NewErrorWithLocation(err.Error())
		return errorsutil.Wrapf(err, "error opening document (%s)", inputFilenameDOCX)
	}
	defer doc.Close()

	c := convert.ConvertToPdf(doc)

	err = c.WriteToFile(outputFilenamePDF)
	if err != nil {
		err = errorsutil.NewErrorWithLocation(err.Error())
		return errorsutil.Wrapf(err, "error converting/writing document (%s) from (%s)", outputFilenamePDF, inputFilenameDOCX)
	} else {
		return nil
	}
}
