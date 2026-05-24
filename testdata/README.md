# Test Data

This directory contains sample files for testing gounidoc functionality.

## Files

| File | Description |
|------|-------------|
| `lorem_ipsum.docx` | Simple DOCX with Lorem Ipsum text and headings |
| `lorem_ipsum.pdf` | PDF generated from the DOCX file |

## Regenerating Test Files

If you need to regenerate the PDF from the DOCX:

```bash
export UNIDOC_KEY="your-api-key"
gounidoc docx2pdf -i testdata/lorem_ipsum.docx -o testdata/lorem_ipsum.pdf
```

## Usage in Tests

```go
import "github.com/grokify/gounidoc/testdata"

func TestConversion(t *testing.T) {
    inputPDF := filepath.Join("testdata", "lorem_ipsum.pdf")
    // ... run conversion tests
}
```
