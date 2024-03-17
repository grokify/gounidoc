package gounidoc

import (
	"time"

	"github.com/unidoc/unioffice/common"
	"github.com/unidoc/unioffice/presentation"
)

type Metadata struct {
	Author         string
	Category       string
	ContentStatus  string
	Created        time.Time
	Description    string
	LastModifiedBy string
	Modified       time.Time
	Title          string
}

func ExportMetadata(cpo common.CoreProperties) Metadata {
	return Metadata{
		Author:         cpo.Author(),
		Category:       cpo.Category(),
		ContentStatus:  cpo.ContentStatus(),
		Created:        cpo.Created(),
		Description:    cpo.Description(),
		LastModifiedBy: cpo.LastModifiedBy(),
		Modified:       cpo.Modified(),
		Title:          cpo.Title(),
	}
}

func ReadFilePresentationMetadata(filename string) (Metadata, error) {
	if pres, err := presentation.Open(filename); err != nil {
		return Metadata{}, err
	} else {
		return ExportMetadata(pres.DocBase.CoreProperties), nil
	}
}
