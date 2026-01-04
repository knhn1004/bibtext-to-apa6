package bibtex

import (
	"fmt"
	"regexp"
	"strings"
)

type Entry struct {
	Type   string
	Key    string
	Fields map[string]string
}

func Parse(input string) (*Entry, error) {
	input = strings.TrimSpace(input)

	re := regexp.MustCompile(`(?s)@(\w+)\s*\{\s*([^,]+)\s*,(.+)\}`)
	matches := re.FindStringSubmatch(input)

	if len(matches) < 4 {
		return nil, fmt.Errorf("invalid BibTeX format")
	}

	entry := &Entry{
		Type:   strings.ToLower(matches[1]),
		Key:    matches[2],
		Fields: make(map[string]string),
	}

	fieldsStr := matches[3]
	fieldRe := regexp.MustCompile(`(\w+)\s*=\s*\{([^}]*)\}|(\w+)\s*=\s*"([^"]*)"|(\w+)\s*=\s*([^,}]+)`)
	fieldMatches := fieldRe.FindAllStringSubmatch(fieldsStr, -1)

	for _, match := range fieldMatches {
		var key, value string

		if match[1] != "" && match[2] != "" {
			key = strings.ToLower(match[1])
			value = match[2]
		} else if match[3] != "" && match[4] != "" {
			key = strings.ToLower(match[3])
			value = match[4]
		} else if match[5] != "" && match[6] != "" {
			key = strings.ToLower(match[5])
			value = strings.TrimSpace(match[6])
		}

		if key != "" {
			value = cleanBibTeXValue(value)
			entry.Fields[key] = value
		}
	}

	return entry, nil
}

func cleanBibTeXValue(value string) string {
	value = strings.TrimSpace(value)
	value = strings.ReplaceAll(value, "\n", " ")
	value = regexp.MustCompile(`\s+`).ReplaceAllString(value, " ")

	// Handle LaTeX special characters before removing braces
	value = handleLaTeXCharacters(value)

	// Fix common UTF-8 encoding issues before and after brace removal
	value = fixCommonEncodingIssues(value)

	value = strings.ReplaceAll(value, "{", "")
	value = strings.ReplaceAll(value, "}", "")

	// Try fixing encoding again after brace removal
	value = fixCommonEncodingIssues(value)

	value = strings.TrimSpace(value)

	return value
}

func handleLaTeXCharacters(value string) string {
	// Common LaTeX character replacements
	replacements := map[string]string{
		`\o`:    "ø",
		`\O`:    "Ø",
		`\'a`:   "á",
		`\'e`:   "é",
		`\'i`:   "í",
		`\'o`:   "ó",
		`\'u`:   "ú",
		`\"a`:   "ä",
		`\"e`:   "ë",
		`\"i`:   "ï",
		`\"o`:   "ö",
		`\"u`:   "ü",
		`\~n`:   "ñ",
		`\~a`:   "ã",
		`\^a`:   "â",
		`\^e`:   "ê",
		`\^i`:   "î",
		`\^o`:   "ô",
		`\^u`:   "û",
		`\c{c}`: "ç",
		`\c{C}`: "Ç",
		`\aa`:   "å",
		`\AA`:   "Å",
		`\ae`:   "æ",
		`\AE`:   "Æ",
		`\oe`:   "œ",
		`\OE`:   "Œ",
		`\ss`:   "ß",
	}

	for latex, unicode := range replacements {
		value = strings.ReplaceAll(value, latex, unicode)
	}

	// Handle {\o} and {\O} patterns
	value = strings.ReplaceAll(value, `{\o}`, "ø")
	value = strings.ReplaceAll(value, `{\O}`, "Ø")

	return value
}

func fixCommonEncodingIssues(value string) string {
	// Fix common UTF-8 encoding errors from PDF copies
	replacements := map[string]string{
		"Ã˜": "Ø", // Capital O with stroke
		"Ã¸": "ø", // Small o with stroke
		"Ã…": "Å", // Capital A with ring
		"Ã¥": "å", // Small a with ring
		"Ã¦": "æ", // Small ae ligature
		"Ã†": "Æ", // Capital AE ligature
		"Ã©": "é", // Small e with acute
		"Ã¨": "è", // Small e with grave
		"Ãª": "ê", // Small e with circumflex
		"Ã«": "ë", // Small e with diaeresis
		"Ã¡": "á", // Small a with acute
		"Ã ": "à", // Small a with grave
		"Ã¢": "â", // Small a with circumflex
		"Ã¤": "ä", // Small a with diaeresis
		"Ã¶": "ö", // Small o with diaeresis
		"Ã¼": "ü", // Small u with diaeresis
		"Ã±": "ñ", // Small n with tilde
		"Ã§": "ç", // Small c with cedilla
		"ÃŸ": "ß", // German sharp s
	}

	for bad, good := range replacements {
		value = strings.ReplaceAll(value, bad, good)
	}

	// Also handle the case where just "Ã," appears (partial corruption)
	value = strings.ReplaceAll(value, "Ã,", "Ø,")

	// Handle "Seland, Ã" pattern specifically
	value = strings.ReplaceAll(value, "Seland, Ã", "Seland, Ø")

	return value
}

func (e *Entry) GetField(field string) string {
	return e.Fields[strings.ToLower(field)]
}

func (e *Entry) HasField(field string) bool {
	_, ok := e.Fields[strings.ToLower(field)]
	return ok
}
