package uniofficeutil

import (
	"fmt"
	"strings"

	"github.com/unidoc/unioffice/document"
)

// ReplaceOptions configures find/replace behavior.
type ReplaceOptions struct {
	CaseSensitive bool `json:"case_sensitive"`
	WholeWord     bool `json:"whole_word"`
	ReplaceAll    bool `json:"replace_all"` // Replace all occurrences vs first only
}

// DefaultReplaceOptions returns default replace options.
func DefaultReplaceOptions() ReplaceOptions {
	return ReplaceOptions{
		CaseSensitive: true,
		ReplaceAll:    true,
	}
}

// ReplaceInDocx performs find/replace in a DOCX file.
// Returns the number of replacements made.
func ReplaceInDocx(inputPath, outputPath, find, replace string, opts ReplaceOptions) (int, error) {
	if inputPath == "" {
		return 0, fmt.Errorf("input path is required")
	}
	if outputPath == "" {
		return 0, fmt.Errorf("output path is required")
	}
	if find == "" {
		return 0, fmt.Errorf("find string is required")
	}

	doc, err := document.Open(inputPath)
	if err != nil {
		return 0, fmt.Errorf("failed to open DOCX: %w", err)
	}
	defer doc.Close()

	count := 0

	// Process paragraphs
	for _, para := range doc.Paragraphs() {
		for _, run := range para.Runs() {
			text := run.Text()
			newText, n := replaceText(text, find, replace, opts)
			if n > 0 {
				run.Clear()
				run.AddText(newText)
				count += n
				if !opts.ReplaceAll && count > 0 {
					break
				}
			}
		}
		if !opts.ReplaceAll && count > 0 {
			break
		}
	}

	// Process tables if we should continue
	if opts.ReplaceAll || count == 0 {
		for _, table := range doc.Tables() {
			for _, row := range table.Rows() {
				for _, cell := range row.Cells() {
					for _, para := range cell.Paragraphs() {
						for _, run := range para.Runs() {
							text := run.Text()
							newText, n := replaceText(text, find, replace, opts)
							if n > 0 {
								run.Clear()
								run.AddText(newText)
								count += n
								if !opts.ReplaceAll && count > 0 {
									break
								}
							}
						}
						if !opts.ReplaceAll && count > 0 {
							break
						}
					}
					if !opts.ReplaceAll && count > 0 {
						break
					}
				}
				if !opts.ReplaceAll && count > 0 {
					break
				}
			}
			if !opts.ReplaceAll && count > 0 {
				break
			}
		}
	}

	// Process headers and footers
	if opts.ReplaceAll || count == 0 {
		for _, hdr := range doc.Headers() {
			for _, para := range hdr.Paragraphs() {
				for _, run := range para.Runs() {
					text := run.Text()
					newText, n := replaceText(text, find, replace, opts)
					if n > 0 {
						run.Clear()
						run.AddText(newText)
						count += n
						if !opts.ReplaceAll && count > 0 {
							break
						}
					}
				}
				if !opts.ReplaceAll && count > 0 {
					break
				}
			}
			if !opts.ReplaceAll && count > 0 {
				break
			}
		}
	}

	if opts.ReplaceAll || count == 0 {
		for _, ftr := range doc.Footers() {
			for _, para := range ftr.Paragraphs() {
				for _, run := range para.Runs() {
					text := run.Text()
					newText, n := replaceText(text, find, replace, opts)
					if n > 0 {
						run.Clear()
						run.AddText(newText)
						count += n
						if !opts.ReplaceAll && count > 0 {
							break
						}
					}
				}
				if !opts.ReplaceAll && count > 0 {
					break
				}
			}
			if !opts.ReplaceAll && count > 0 {
				break
			}
		}
	}

	if err := doc.SaveToFile(outputPath); err != nil {
		return count, fmt.Errorf("failed to save DOCX: %w", err)
	}

	return count, nil
}

// ReplaceMultipleInDocx performs multiple find/replace operations in a DOCX file.
// Returns a map of find strings to replacement counts.
func ReplaceMultipleInDocx(inputPath, outputPath string, replacements map[string]string, opts ReplaceOptions) (map[string]int, error) {
	if inputPath == "" {
		return nil, fmt.Errorf("input path is required")
	}
	if outputPath == "" {
		return nil, fmt.Errorf("output path is required")
	}
	if len(replacements) == 0 {
		return nil, fmt.Errorf("no replacements specified")
	}

	doc, err := document.Open(inputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open DOCX: %w", err)
	}
	defer doc.Close()

	counts := make(map[string]int)
	for find := range replacements {
		counts[find] = 0
	}

	// Process paragraphs
	for _, para := range doc.Paragraphs() {
		for _, run := range para.Runs() {
			text := run.Text()
			newText := text
			for find, replace := range replacements {
				result, n := replaceText(newText, find, replace, opts)
				if n > 0 {
					newText = result
					counts[find] += n
				}
			}
			if newText != text {
				run.Clear()
				run.AddText(newText)
			}
		}
	}

	// Process tables
	for _, table := range doc.Tables() {
		for _, row := range table.Rows() {
			for _, cell := range row.Cells() {
				for _, para := range cell.Paragraphs() {
					for _, run := range para.Runs() {
						text := run.Text()
						newText := text
						for find, replace := range replacements {
							result, n := replaceText(newText, find, replace, opts)
							if n > 0 {
								newText = result
								counts[find] += n
							}
						}
						if newText != text {
							run.Clear()
							run.AddText(newText)
						}
					}
				}
			}
		}
	}

	// Process headers
	for _, hdr := range doc.Headers() {
		for _, para := range hdr.Paragraphs() {
			for _, run := range para.Runs() {
				text := run.Text()
				newText := text
				for find, replace := range replacements {
					result, n := replaceText(newText, find, replace, opts)
					if n > 0 {
						newText = result
						counts[find] += n
					}
				}
				if newText != text {
					run.Clear()
					run.AddText(newText)
				}
			}
		}
	}

	// Process footers
	for _, ftr := range doc.Footers() {
		for _, para := range ftr.Paragraphs() {
			for _, run := range para.Runs() {
				text := run.Text()
				newText := text
				for find, replace := range replacements {
					result, n := replaceText(newText, find, replace, opts)
					if n > 0 {
						newText = result
						counts[find] += n
					}
				}
				if newText != text {
					run.Clear()
					run.AddText(newText)
				}
			}
		}
	}

	if err := doc.SaveToFile(outputPath); err != nil {
		return counts, fmt.Errorf("failed to save DOCX: %w", err)
	}

	return counts, nil
}

// replaceText performs find/replace on a string with the given options.
// Returns the modified string and the number of replacements made.
func replaceText(text, find, replace string, opts ReplaceOptions) (string, int) {
	if text == "" || find == "" {
		return text, 0
	}

	searchText := text
	searchFind := find

	if !opts.CaseSensitive {
		searchText = strings.ToLower(text)
		searchFind = strings.ToLower(find)
	}

	if opts.WholeWord {
		// For whole word matching, we need to check word boundaries
		count := 0
		result := text
		idx := 0
		for {
			pos := strings.Index(searchText[idx:], searchFind)
			if pos == -1 {
				break
			}
			absolutePos := idx + pos

			// Check word boundaries
			isWordStart := absolutePos == 0 || !isWordChar(rune(text[absolutePos-1]))
			isWordEnd := absolutePos+len(find) >= len(text) || !isWordChar(rune(text[absolutePos+len(find)]))

			if isWordStart && isWordEnd {
				// Replace at this position
				result = result[:absolutePos] + replace + result[absolutePos+len(find):]
				searchText = searchText[:absolutePos] + replace + searchText[absolutePos+len(find):]
				count++
				idx = absolutePos + len(replace)
				if !opts.ReplaceAll {
					return result, count
				}
			} else {
				idx = absolutePos + 1
			}
		}
		return result, count
	}

	// Simple replacement
	var count int
	if opts.ReplaceAll {
		if opts.CaseSensitive {
			count = strings.Count(text, find)
			return strings.ReplaceAll(text, find, replace), count
		}
		// Case insensitive replace all
		result := text
		idx := 0
		for {
			pos := strings.Index(strings.ToLower(result[idx:]), strings.ToLower(find))
			if pos == -1 {
				break
			}
			absolutePos := idx + pos
			result = result[:absolutePos] + replace + result[absolutePos+len(find):]
			count++
			idx = absolutePos + len(replace)
		}
		return result, count
	}

	// Replace first only
	if opts.CaseSensitive {
		pos := strings.Index(text, find)
		if pos == -1 {
			return text, 0
		}
		return text[:pos] + replace + text[pos+len(find):], 1
	}

	// Case insensitive replace first
	pos := strings.Index(strings.ToLower(text), strings.ToLower(find))
	if pos == -1 {
		return text, 0
	}
	return text[:pos] + replace + text[pos+len(find):], 1
}

// isWordChar returns true if the rune is a word character (letter, digit, or underscore).
func isWordChar(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_'
}
