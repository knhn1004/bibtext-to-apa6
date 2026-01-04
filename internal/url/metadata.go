package url

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/playwright-community/playwright-go"
)

type Metadata struct {
	Title       string
	Author      string
	Year        string
	Publisher   string
	URL         string
	AccessDate  time.Time
}

// ExtractMetadata tries HTTP first, then falls back to Playwright if that fails
func ExtractMetadata(urlStr string) (*Metadata, error) {
	_, err := url.Parse(urlStr)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	// Try HTTP method first
	metadata, err := extractMetadataHTTP(urlStr)
	if err == nil {
		return metadata, nil
	}

	// Fall back to Playwright for sites that block simple HTTP requests
	return extractMetadataPlaywright(urlStr)
}

// extractMetadataHTTP uses a simple HTTP client with browser-like headers
func extractMetadataHTTP(urlStr string) (*Metadata, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers to mimic a browser request
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP error: %s", resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	return parseDocument(doc, urlStr), nil
}

// extractMetadataPlaywright uses a headless browser to fetch the page
func extractMetadataPlaywright(urlStr string) (*Metadata, error) {
	pw, err := playwright.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to start playwright: %w", err)
	}
	defer pw.Stop()

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to launch browser: %w", err)
	}
	defer browser.Close()

	context, err := browser.NewContext(playwright.BrowserNewContextOptions{
		UserAgent: playwright.String("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create browser context: %w", err)
	}
	defer context.Close()

	page, err := context.NewPage()
	if err != nil {
		return nil, fmt.Errorf("failed to create page: %w", err)
	}

	_, err = page.Goto(urlStr, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
		Timeout:   playwright.Float(30000), // 30 seconds
	})
	if err != nil {
		return nil, fmt.Errorf("failed to navigate to URL: %w", err)
	}

	// Get the page content
	content, err := page.Content()
	if err != nil {
		return nil, fmt.Errorf("failed to get page content: %w", err)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	return parseDocument(doc, urlStr), nil
}

// parseDocument extracts metadata from a goquery document
func parseDocument(doc *goquery.Document, urlStr string) *Metadata {
	metadata := &Metadata{
		URL:        urlStr,
		AccessDate: time.Now(),
	}

	metadata.Title = extractTitle(doc)
	metadata.Author = extractAuthor(doc)
	metadata.Publisher = extractPublisher(doc, urlStr)
	metadata.Year = extractYear(doc)

	return metadata
}

func extractTitle(doc *goquery.Document) string {
	selectors := []string{
		`meta[property="og:title"]`,
		`meta[name="twitter:title"]`,
		`meta[name="citation_title"]`,
		`meta[name="DC.title"]`,
		`title`,
	}

	for _, selector := range selectors {
		if selector == "title" {
			title := doc.Find(selector).First().Text()
			if title != "" {
				return strings.TrimSpace(title)
			}
		} else {
			title := doc.Find(selector).First().AttrOr("content", "")
			if title != "" {
				return strings.TrimSpace(title)
			}
		}
	}

	return "Untitled"
}

func extractAuthor(doc *goquery.Document) string {
	selectors := []string{
		`meta[name="author"]`,
		`meta[property="article:author"]`,
		`meta[name="citation_author"]`,
		`meta[name="DC.creator"]`,
		`meta[name="byl"]`,
	}

	authors := []string{}

	for _, selector := range selectors {
		doc.Find(selector).Each(func(i int, s *goquery.Selection) {
			author := s.AttrOr("content", "")
			if author != "" {
				authors = append(authors, strings.TrimSpace(author))
			}
		})
		if len(authors) > 0 {
			break
		}
	}

	if len(authors) > 0 {
		return strings.Join(authors, " & ")
	}

	return ""
}

func extractPublisher(doc *goquery.Document, urlStr string) string {
	selectors := []string{
		`meta[property="og:site_name"]`,
		`meta[name="publisher"]`,
		`meta[name="DC.publisher"]`,
		`meta[name="citation_publisher"]`,
	}

	for _, selector := range selectors {
		publisher := doc.Find(selector).First().AttrOr("content", "")
		if publisher != "" {
			return strings.TrimSpace(publisher)
		}
	}

	u, err := url.Parse(urlStr)
	if err == nil {
		parts := strings.Split(u.Hostname(), ".")
		if len(parts) >= 2 {
			return strings.Title(parts[len(parts)-2])
		}
	}

	return ""
}

func extractYear(doc *goquery.Document) string {
	selectors := []string{
		`meta[name="publication_date"]`,
		`meta[property="article:published_time"]`,
		`meta[name="citation_publication_date"]`,
		`meta[name="DC.date"]`,
		`time[datetime]`,
	}

	for _, selector := range selectors {
		var dateStr string
		if selector == `time[datetime]` {
			dateStr = doc.Find(selector).First().AttrOr("datetime", "")
		} else {
			dateStr = doc.Find(selector).First().AttrOr("content", "")
		}

		if dateStr != "" {
			if year := extractYearFromDate(dateStr); year != "" {
				return year
			}
		}
	}

	return fmt.Sprintf("%d", time.Now().Year())
}

func extractYearFromDate(dateStr string) string {
	layouts := []string{
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02T15:04:05",
		"2006-01-02",
		"2006/01/02",
		"02/01/2006",
		"01/02/2006",
		"2006",
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, dateStr); err == nil {
			return fmt.Sprintf("%d", t.Year())
		}
	}

	if len(dateStr) >= 4 {
		yearStr := dateStr[:4]
		if _, err := time.Parse("2006", yearStr); err == nil {
			return yearStr
		}
	}

	return ""
}

func (m *Metadata) ToAPAFormat() string {
	author := m.Author
	if author == "" {
		author = m.Publisher
	}
	if author == "" {
		author = "Unknown"
	}

	title := m.Title
	if title == "" {
		title = "Untitled"
	}

	result := fmt.Sprintf("%s. (%s). %s", author, m.Year, title)

	if m.Publisher != "" && m.Publisher != author {
		result += fmt.Sprintf(". %s", m.Publisher)
	}

	result += fmt.Sprintf(". Retrieved %s, from %s", 
		m.AccessDate.Format("January 2, 2006"), 
		m.URL)

	return result
}