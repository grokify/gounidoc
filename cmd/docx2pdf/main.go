package main

import (
	"fmt"
	"log/slog"

	"github.com/grokify/gounidoc"
	"github.com/grokify/gounidoc/uniofficeutil"
	"github.com/grokify/mogo/errors/errorsutil"
	"github.com/grokify/mogo/log/logutil"
	"github.com/jessevdk/go-flags"
)

type Options struct {
	InputFile  string `short:"d" long:"delete" description:"Delete subscription"`
	OutputFile string `short:"r" long:"recreate" description:"Recreate subscription"`
}

func main() {
	opts := Options{}
	_, err := flags.Parse(&opts)
	logutil.FatalErr(err)

	err = gounidoc.SetMeteredKeyEnv()
	logutil.FatalErr(err)

	if opts.OutputFile == "" && opts.InputFile != "" {
		opts.OutputFile = opts.InputFile + ".pdf"
	}

	err = uniofficeutil.ConvertDOCXFileToPDFFile(opts.InputFile, opts.OutputFile)
	if err != nil {
		slog.Error(errorsutil.NewErrorWithLocation(err.Error()).Error())
		logutil.FatalErr(err)
	}

	fmt.Println("DONE")
}
