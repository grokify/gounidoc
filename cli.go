package gounidoc

import (
	uniLicenseOffice "github.com/unidoc/unioffice/common/license"
	uniLicensePDF "github.com/unidoc/unipdf/v3/common/license"
)

const (
	EnvUnidocKey       = "UNIDOC_KEY"
	EnvUnidocKeyOffice = "UNIDOC_KEY_OFFICE"
	EnvUnidocKeyPDF    = "UNIDOC_KEY_PDF"
	EnvUnidocFile      = "UNIDOC_FILE"
)

type Options struct {
	Key  string `short:"k" long:"key" description:"UniDoc metered API key"`
	File string `short:"f" long:"file" description:"File to read metadata for"`
}

func SetMeteredKey(officeKey, pdfKey string) error {
	if officeKey != "" {
		if err := uniLicenseOffice.SetMeteredKey(officeKey); err != nil {
			return err
		}
	}
	if pdfKey != "" {
		if err := uniLicensePDF.SetMeteredKey(officeKey); err != nil {
			return err
		}
	}
	return nil
}
