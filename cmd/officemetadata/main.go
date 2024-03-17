package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/grokify/gounidoc"
	"github.com/grokify/mogo/config"
	"github.com/grokify/mogo/fmt/fmtutil"
)

type Options struct {
	Env  string `short:"e" long:"env" description:"Env file"`
	Key  string `short:"k" long:"key" description:"UniDoc metered API key"`
	File string `short:"f" long:"file" description:"File to read metadata for"`
}

func (opts *Options) DotEnvFilename() string {
	return strings.TrimSpace(opts.Env)
}

func (opts *Options) APIKey() string {
	if key := strings.TrimSpace(opts.Key); key != "" {
		return key
	} else {
		return os.Getenv(gounidoc.EnvUnidocKey)
	}
}

func (opts *Options) Filename() string {
	if filename := strings.TrimSpace(opts.File); filename != "" {
		return filename
	} else {
		return os.Getenv(gounidoc.EnvUnidocFile)
	}
}

func main() {
	opts := Options{}
	err := config.ParseFlagsAndLoadDotEnv(&opts)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf(opts.APIKey())

	err = gounidoc.SetMeteredKey(opts.APIKey())
	if err != nil {
		log.Fatal(err)
	}
	m, err := gounidoc.ReadFileMetadataPresentation(opts.Filename())
	if err != nil {
		log.Fatal(err)
	}
	fmtutil.PrintJSON(m)
	fmt.Println("DONE")
}
