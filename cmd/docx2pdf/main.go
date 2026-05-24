package main

import (
	"fmt"
	"log"
	"os"

	"github.com/grokify/gounidoc"
	"github.com/grokify/gounidoc/uniofficeutil"
	"github.com/jessevdk/go-flags"
)

type Options struct {
	InputFile  string `short:"i" long:"input" description:"Input DOCX file"`
	OutputFile string `short:"o" long:"output" description:"Output PDF file"`
}

func main() {
	opts := Options{}
	_, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(1)
	}

	if err := gounidoc.SetMeteredKeyEnv(); err != nil {
		log.Fatalf("Error setting UniDoc key: %v", err)
	}

	if opts.OutputFile == "" && opts.InputFile != "" {
		opts.OutputFile = opts.InputFile + ".pdf"
	}

	if err := uniofficeutil.ConvertDOCXFileToPDFFile(opts.InputFile, opts.OutputFile); err != nil {
		log.Fatalf("Error converting DOCX to PDF: %v", err)
	}

	fmt.Println("DONE")
}
