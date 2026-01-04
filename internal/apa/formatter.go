package apa

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/knhn1004/bibtext-to-apa6/internal/bibtex"
)

func Format(entry *bibtex.Entry) (string, error) {
	switch entry.Type {
	case "article":
		return formatArticle(entry), nil
	case "book":
		return formatBook(entry), nil
	case "inproceedings", "conference":
		return formatInProceedings(entry), nil
	case "inbook", "incollection":
		return formatInBook(entry), nil
	case "misc", "online":
		return formatMisc(entry), nil
	case "phdthesis", "mastersthesis":
		return formatThesis(entry), nil
	default:
		return formatGeneric(entry), nil
	}
}

func formatArticle(entry *bibtex.Entry) string {
	authors := formatAuthors(entry.GetField("author"))
	year := formatYear(entry.GetField("year"))
	title := sentenceCase(entry.GetField("title"))
	journal := entry.GetField("journal")
	volume := entry.GetField("volume")
	number := entry.GetField("number")
	pages := formatPages(entry.GetField("pages"))
	doi := entry.GetField("doi")
	
	result := fmt.Sprintf("%s (%s). %s.", authors, year, title)
	
	if journal != "" {
		result += fmt.Sprintf(" %s", italicize(journal))
		if volume != "" {
			result += fmt.Sprintf(", %s", italicize(volume))
			if number != "" {
				result += fmt.Sprintf("(%s)", number)
			}
		}
		if pages != "" {
			result += fmt.Sprintf(", %s", pages)
		}
		result += "."
	}
	
	if doi != "" {
		result += fmt.Sprintf(" https://doi.org/%s", doi)
	}
	
	return result
}

func formatBook(entry *bibtex.Entry) string {
	authors := formatAuthors(entry.GetField("author"))
	year := formatYear(entry.GetField("year"))
	title := italicize(sentenceCase(entry.GetField("title")))
	publisher := entry.GetField("publisher")
	address := entry.GetField("address")
	
	result := fmt.Sprintf("%s (%s). %s", authors, year, title)
	
	if address != "" && publisher != "" {
		result += fmt.Sprintf(". %s: %s", address, publisher)
	} else if publisher != "" {
		result += fmt.Sprintf(". %s", publisher)
	}
	
	result += "."
	return result
}

func formatInProceedings(entry *bibtex.Entry) string {
	authors := formatAuthors(entry.GetField("author"))
	year := formatYear(entry.GetField("year"))
	title := sentenceCase(entry.GetField("title"))
	booktitle := italicize(entry.GetField("booktitle"))
	pages := formatPages(entry.GetField("pages"))
	publisher := entry.GetField("publisher")
	
	result := fmt.Sprintf("%s (%s). %s. In %s", authors, year, title, booktitle)
	
	if pages != "" {
		result += fmt.Sprintf(" (pp. %s)", pages)
	}
	
	if publisher != "" {
		result += fmt.Sprintf(". %s", publisher)
	}
	
	result += "."
	return result
}

func formatInBook(entry *bibtex.Entry) string {
	authors := formatAuthors(entry.GetField("author"))
	year := formatYear(entry.GetField("year"))
	title := sentenceCase(entry.GetField("title"))
	booktitle := italicize(entry.GetField("booktitle"))
	editor := entry.GetField("editor")
	pages := formatPages(entry.GetField("pages"))
	publisher := entry.GetField("publisher")
	
	result := fmt.Sprintf("%s (%s). %s", authors, year, title)
	
	if editor != "" {
		result += fmt.Sprintf(". In %s (Ed.), %s", formatAuthors(editor), booktitle)
	} else {
		result += fmt.Sprintf(". In %s", booktitle)
	}
	
	if pages != "" {
		result += fmt.Sprintf(" (pp. %s)", pages)
	}
	
	if publisher != "" {
		result += fmt.Sprintf(". %s", publisher)
	}
	
	result += "."
	return result
}

func formatThesis(entry *bibtex.Entry) string {
	authors := formatAuthors(entry.GetField("author"))
	year := formatYear(entry.GetField("year"))
	title := italicize(sentenceCase(entry.GetField("title")))
	school := entry.GetField("school")
	thesisType := "Doctoral dissertation"
	
	if entry.Type == "mastersthesis" {
		thesisType = "Master's thesis"
	}
	
	result := fmt.Sprintf("%s (%s). %s [%s]", authors, year, title, thesisType)
	
	if school != "" {
		result += fmt.Sprintf(". %s", school)
	}
	
	result += "."
	return result
}

func formatMisc(entry *bibtex.Entry) string {
	authors := formatAuthors(entry.GetField("author"))
	year := formatYear(entry.GetField("year"))
	title := italicize(sentenceCase(entry.GetField("title")))
	url := entry.GetField("url")
	
	result := fmt.Sprintf("%s (%s). %s", authors, year, title)
	
	if url != "" {
		result += fmt.Sprintf(". Retrieved from %s", url)
	} else {
		result += "."
	}
	
	return result
}

func formatGeneric(entry *bibtex.Entry) string {
	authors := formatAuthors(entry.GetField("author"))
	year := formatYear(entry.GetField("year"))
	title := sentenceCase(entry.GetField("title"))
	
	return fmt.Sprintf("%s (%s). %s.", authors, year, title)
}

func formatAuthors(authors string) string {
	if authors == "" {
		return "Unknown"
	}
	
	// Fix encoding issues that might have survived the parser
	authors = fixAuthorEncoding(authors)
	
	authorList := strings.Split(authors, " and ")
	formatted := []string{}
	
	for i, author := range authorList {
		// Handle et al. for more than 7 authors
		if len(authorList) > 7 && i == 6 {
			formatted = append(formatted, "…")
			// Now process the last author
			lastAuthor := authorList[len(authorList)-1]
			lastAuthor = strings.TrimSpace(lastAuthor)
			parts := strings.Split(lastAuthor, ",")
			
			if len(parts) == 2 {
				lastName := strings.TrimSpace(parts[0])
				firstName := strings.TrimSpace(parts[1])
				initials := getInitials(firstName)
				formatted = append(formatted, fmt.Sprintf("%s, %s", lastName, initials))
			} else {
				words := strings.Fields(lastAuthor)
				if len(words) > 1 {
					lastName := words[len(words)-1]
					firstName := strings.Join(words[:len(words)-1], " ")
					initials := getInitials(firstName)
					formatted = append(formatted, fmt.Sprintf("%s, %s", lastName, initials))
				} else {
					formatted = append(formatted, lastAuthor)
				}
			}
			break // Exit loop after processing last author
		} else if len(authorList) > 7 && i > 6 {
			continue // Skip all authors after the 6th except the last (already handled)
		}
		
		// Process regular authors (first 6 when >7 authors, or all when <=7)
		author = strings.TrimSpace(author)
		parts := strings.Split(author, ",")
		
		if len(parts) == 2 {
			lastName := strings.TrimSpace(parts[0])
			firstName := strings.TrimSpace(parts[1])
			initials := getInitials(firstName)
			
			formatted = append(formatted, fmt.Sprintf("%s, %s", lastName, initials))
		} else {
			words := strings.Fields(author)
			if len(words) > 1 {
				lastName := words[len(words)-1]
				firstName := strings.Join(words[:len(words)-1], " ")
				initials := getInitials(firstName)
				
				formatted = append(formatted, fmt.Sprintf("%s, %s", lastName, initials))
			} else {
				formatted = append(formatted, author)
			}
		}
	}
	
	if len(formatted) == 1 {
		return formatted[0]
	} else if len(formatted) == 2 {
		return strings.Join(formatted, ", & ")
	} else {
		// Check if we have ellipsis
		hasEllipsis := false
		for _, f := range formatted {
			if f == "…" {
				hasEllipsis = true
				break
			}
		}
		
		if hasEllipsis {
			// Join all parts with spaces, the ellipsis is already in the array
			result := ""
			for i, part := range formatted {
				if i > 0 {
					if part == "…" {
						result += " "
					} else if i == len(formatted)-1 {
						result += " "
					} else {
						result += ", "
					}
				}
				result += part
			}
			return result
		} else {
			// Normal case without ellipsis
			lastAuthor := formatted[len(formatted)-1]
			otherAuthors := strings.Join(formatted[:len(formatted)-1], ", ")
			return fmt.Sprintf("%s, & %s", otherAuthors, lastAuthor)
		}
	}
}

func getInitials(firstName string) string {
	parts := strings.Fields(firstName)
	initials := []string{}
	
	for _, part := range parts {
		if len(part) > 0 {
			// Just the initial without period after it
			initials = append(initials, string(part[0]))
		}
	}
	
	// Join initials without spaces, then add period at end
	return strings.Join(initials, "") + "."
}

func formatYear(year string) string {
	if year == "" {
		return "n.d."
	}
	return year
}

func formatPages(pages string) string {
	if pages == "" {
		return ""
	}
	// Replace hyphens and double hyphens with en dash
	pages = strings.ReplaceAll(pages, "--", "–")
	pages = strings.ReplaceAll(pages, "-", "–")
	// Remove spaces around en dash
	pages = strings.ReplaceAll(pages, " – ", "–")
	pages = strings.ReplaceAll(pages, " —", "–") // Also handle em dash
	pages = strings.ReplaceAll(pages, "— ", "–")
	return pages
}

func italicize(text string) string {
	return fmt.Sprintf("*%s*", text)
}

func fixAuthorEncoding(authors string) string {
	// Fix common UTF-8 encoding issues specific to author names
	authors = strings.ReplaceAll(authors, "Ã˜", "Ø")
	authors = strings.ReplaceAll(authors, "Ã¸", "ø")
	
	// Use regex to catch Ã as a single initial (followed by comma, period, or whitespace)
	re := regexp.MustCompile(`\bÃ\b`)
	authors = re.ReplaceAllString(authors, "Ø")
	
	return authors
}

func sentenceCase(text string) string {
	if text == "" {
		return ""
	}
	
	// Split on common delimiters
	words := strings.Fields(text)
	if len(words) == 0 {
		return ""
	}
	
	// First word is always capitalized
	result := []string{words[0]}
	
	// Rest are lowercase unless they're proper nouns (which we can't reliably detect)
	// or after certain punctuation
	for i := 1; i < len(words); i++ {
		word := strings.ToLower(words[i])
		
		// Check if previous word ended with sentence-ending punctuation
		if i > 0 && (strings.HasSuffix(words[i-1], ".") || 
		             strings.HasSuffix(words[i-1], "?") || 
		             strings.HasSuffix(words[i-1], "!") ||
		             strings.HasSuffix(words[i-1], ":")) {
			word = strings.Title(word)
		}
		
		result = append(result, word)
	}
	
	return strings.Join(result, " ")
}