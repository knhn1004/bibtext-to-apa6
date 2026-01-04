package richtext

import (
	"fmt"
	"regexp"
	"strings"
)

// ConvertToRTF converts APA format with *italic* markers to RTF format
func ConvertToRTF(text string) string {
	// RTF header
	rtf := `{\rtf1\ansi\ansicpg1252\cocoartf2639
\cocoatextscaling0\cocoaplatform0{\fonttbl\f0\fnil\fcharset0 HelveticaNeue;}
{\colortbl;\red255\green255\blue255;\red0\green0\blue0;}
{\*\expandedcolortbl;;\cssrgb\c0\c0\c0;}
\paperw11900\paperh16840\margl1440\margr1440\vieww11520\viewh8400\viewkind0
\deftab720
\pard\pardeftab720\partightenfactor0
\f0\fs24 \cf2 `

	// Process the text to handle italics
	processedText := processItalics(text)
	
	// RTF footer
	rtf += processedText + "}"
	
	return rtf
}

// processItalics converts *text* to RTF italic format
func processItalics(text string) string {
	// Replace *text* with RTF italic markers
	re := regexp.MustCompile(`\*([^*]+)\*`)
	result := re.ReplaceAllStringFunc(text, func(match string) string {
		// Remove the asterisks and wrap in RTF italic tags
		content := match[1 : len(match)-1]
		return fmt.Sprintf(`\i %s\i0 `, escapeRTF(content))
	})
	
	return escapeRTF(result)
}

// escapeRTF escapes special characters for RTF
func escapeRTF(s string) string {
	// In the replacement function above, we already escape the content
	// This is for the non-italic parts
	if strings.Contains(s, `\i `) {
		// Already processed, don't escape
		return s
	}
	
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `{`, `\{`)
	s = strings.ReplaceAll(s, `}`, `\}`)
	
	// Handle unicode characters
	var result strings.Builder
	for _, r := range s {
		if r > 127 {
			// Unicode character - use RTF unicode escape
			result.WriteString(fmt.Sprintf(`\u%d?`, r))
		} else {
			result.WriteRune(r)
		}
	}
	
	return result.String()
}

// ConvertToHTML converts APA format with *italic* markers to HTML format
func ConvertToHTML(text string) string {
	// Replace *text* with HTML italic tags
	re := regexp.MustCompile(`\*([^*]+)\*`)
	result := re.ReplaceAllString(text, `<i>$1</i>`)
	
	// Escape HTML entities (but preserve our italic tags)
	result = strings.ReplaceAll(result, "&", "&amp;")
	// Don't escape < and > that are part of our tags
	temp := strings.ReplaceAll(result, "<i>", "§ITALIC_START§")
	temp = strings.ReplaceAll(temp, "</i>", "§ITALIC_END§")
	temp = strings.ReplaceAll(temp, "<", "&lt;")
	temp = strings.ReplaceAll(temp, ">", "&gt;")
	temp = strings.ReplaceAll(temp, "§ITALIC_START§", "<i>")
	result = strings.ReplaceAll(temp, "§ITALIC_END§", "</i>")
	
	// Convert double newlines to paragraph breaks for proper formatting in word processors
	// Each reference becomes its own paragraph with APA hanging indent style
	paragraphs := strings.Split(result, "\n\n")
	if len(paragraphs) > 1 {
		var htmlParagraphs []string
		for _, p := range paragraphs {
			p = strings.TrimSpace(p)
			if p != "" {
				// Apply hanging indent: 0.5 inch indent except first line
				// line-height: 1.15 for tight spacing, margin-bottom for between-reference spacing
				htmlParagraphs = append(htmlParagraphs, fmt.Sprintf(
					`<p style="margin-left: 0.5in; text-indent: -0.5in; margin-top: 0; margin-bottom: 12pt; line-height: 1.15;">%s</p>`,
					p,
				))
			}
		}
		result = strings.Join(htmlParagraphs, "\n")
	}
	
	return result
}

// StripFormatting removes all formatting markers for plain text display
func StripFormatting(text string) string {
	// Remove italic markers
	re := regexp.MustCompile(`\*([^*]+)\*`)
	return re.ReplaceAllString(text, `$1`)
}