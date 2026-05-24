package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/grokify/gounidoc"
	"github.com/grokify/gounidoc/unipdfutil"
	flags "github.com/jessevdk/go-flags"
)

type Options struct {
	Input     string `short:"i" long:"input" description:"Input PDF file path" required:"true"`
	Output    string `short:"o" long:"output" description:"Output DOCX file path (optional, defaults to input name with .docx)"`
	NoTables  bool   `long:"no-tables" description:"Disable table extraction"`
	NoImages  bool   `long:"no-images" description:"Disable image extraction"`
	NoLayout  bool   `long:"no-layout" description:"Disable layout/heading detection"`
	BasicMode bool   `short:"b" long:"basic" description:"Basic mode: text only, no tables/images/layout"`
}

func main() {
	opts := Options{}
	_, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(1)
	}

	// Set UniDoc metered key from environment
	if err := gounidoc.SetMeteredKeyEnv(); err != nil {
		fmt.Fprintf(os.Stderr, "Error setting UniDoc key: %v\n", err)
		os.Exit(1)
	}

	// Determine output path
	outputPath := opts.Output
	if outputPath == "" {
		ext := filepath.Ext(opts.Input)
		outputPath = strings.TrimSuffix(opts.Input, ext) + ".docx"
	}

	// Configure conversion options
	convOpts := unipdfutil.DefaultConversionOptions()
	if opts.BasicMode {
		convOpts.ExtractTables = false
		convOpts.ExtractImages = false
		convOpts.DetectLayout = false
	} else {
		if opts.NoTables {
			convOpts.ExtractTables = false
		}
		if opts.NoImages {
			convOpts.ExtractImages = false
		}
		if opts.NoLayout {
			convOpts.DetectLayout = false
		}
	}

	fmt.Printf("Converting: %s -> %s\n", opts.Input, outputPath)
	fmt.Printf("Options: tables=%v, images=%v, layout=%v\n",
		convOpts.ExtractTables, convOpts.ExtractImages, convOpts.DetectLayout)

	if err := unipdfutil.ConvertPDFFileToDocxFileWithOptions(opts.Input, outputPath, convOpts); err != nil {
		fmt.Fprintf(os.Stderr, "Error converting PDF to DOCX: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully converted to: %s\n", outputPath)
}
