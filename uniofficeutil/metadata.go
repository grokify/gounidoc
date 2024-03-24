package uniofficeutil

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

func ExportMetadata(cp common.CoreProperties) Metadata {
	return Metadata{
		Author:         cp.Author(),
		Category:       cp.Category(),
		ContentStatus:  cp.ContentStatus(),
		Created:        cp.Created(),
		Description:    cp.Description(),
		LastModifiedBy: cp.LastModifiedBy(),
		Modified:       cp.Modified(),
		Title:          cp.Title(),
	}
}

func ReadFileMetadataPresentation(filename string) (Metadata, error) {
	if pres, err := presentation.Open(filename); err != nil {
		return Metadata{}, err
	} else {
		return ExportMetadata(pres.DocBase.CoreProperties), nil
	}
}
