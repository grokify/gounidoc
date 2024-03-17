package gounidoc

import "github.com/unidoc/unioffice/common/license"

const (
	EnvUnidocKey  = "UNIDOC_KEY"
	EnvUnidocFile = "UNIDOC_FILE"
)

type Options struct {
	Key  string `short:"k" long:"key" description:"UniDoc metered API key"`
	File string `short:"f" long:"file" description:"File to read metadata for"`
}

func SetMeteredKey(key string) error {
	// Make sure to load your metered License API key prior to using the library.
	// If you need a key, you can sign up and create a free one at https://cloud.unidoc.io
	return license.SetMeteredKey(key)
}
