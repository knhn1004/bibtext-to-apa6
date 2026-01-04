package richtext

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// CopyToClipboard copies text with rich formatting to the clipboard
// It handles plain text fallback if rich text copying fails
func CopyToClipboard(text string) error {
	switch runtime.GOOS {
	case "darwin":
		return copyToClipboardDarwin(text)
	case "windows":
		return copyToClipboardWindows(text)
	case "linux":
		return copyToClipboardLinux(text)
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
}

// copyToClipboardDarwin copies rich text to clipboard on macOS
func copyToClipboardDarwin(text string) error {
	// For macOS, let's use a simpler approach with RTF
	// First, let's try with a temporary HTML file and textutil
	html := ConvertToHTML(text)

	// Create HTML content
	// If HTML already contains <p> tags (multiple references), use as-is
	// Otherwise wrap in a single paragraph
	var bodyContent string
	if strings.Contains(html, "<p ") || strings.Contains(html, "<p>") {
		bodyContent = html
	} else {
		bodyContent = fmt.Sprintf(`<p style="margin-left: 0.5in; text-indent: -0.5in; margin-top: 0; margin-bottom: 0; line-height: 1.15;">%s</p>`, html)
	}

	htmlContent := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
</head>
<body>
%s
</body>
</html>`, bodyContent)

	// Write HTML to temp file
	tmpFile := "/tmp/bibapa-clipboard.html"
	if err := os.WriteFile(tmpFile, []byte(htmlContent), 0644); err != nil {
		return copyPlainText(text)
	}
	defer os.Remove(tmpFile)

	// Convert HTML to RTF using textutil
	rtfFile := "/tmp/bibapa-clipboard.rtf"
	cmd := exec.Command("textutil", "-convert", "rtf", "-output", rtfFile, tmpFile)
	if err := cmd.Run(); err != nil {
		return copyPlainText(text)
	}
	defer os.Remove(rtfFile)

	// Use osascript to set the clipboard with RTF data
	script := fmt.Sprintf(`
set rtfFile to POSIX file "%s"
set rtfData to read rtfFile as «class RTF »
set the clipboard to {Unicode text:"%s", «class RTF »:rtfData}
`, rtfFile, escapeAppleScript(StripFormatting(text)))

	cmd = exec.Command("osascript", "-e", script)
	output, err := cmd.CombinedOutput()
	if err != nil {
		// Log error for debugging
		fmt.Fprintf(os.Stderr, "osascript error: %v, output: %s\n", err, output)
		return copyPlainText(text)
	}

	return nil
}

// copyToClipboardWindows copies rich text to clipboard on Windows
func copyToClipboardWindows(text string) error {
	// For Windows, we'll use PowerShell to set clipboard with HTML
	html := ConvertToHTML(text)

	// Create a PowerShell script that sets both HTML and text formats
	script := fmt.Sprintf(`
$html = @"
<html>
<body>
%s
</body>
</html>
"@

$text = @"
%s
"@

Add-Type -Assembly System.Windows.Forms
$dataObject = New-Object System.Windows.Forms.DataObject
$dataObject.SetData([System.Windows.Forms.DataFormats]::Html, $html)
$dataObject.SetData([System.Windows.Forms.DataFormats]::Text, $text)
[System.Windows.Forms.Clipboard]::SetDataObject($dataObject, $true)
`, html, StripFormatting(text))

	cmd := exec.Command("powershell", "-Command", script)
	if err := cmd.Run(); err != nil {
		// Fallback to plain text
		return copyPlainText(text)
	}

	return nil
}

// copyToClipboardLinux copies rich text to clipboard on Linux
func copyToClipboardLinux(text string) error {
	// For Linux, we'll try xclip with HTML format
	html := ConvertToHTML(text)

	// Try to use xclip with HTML target
	cmd := exec.Command("xclip", "-selection", "clipboard", "-t", "text/html")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return copyPlainText(text)
	}

	if err := cmd.Start(); err != nil {
		return copyPlainText(text)
	}

	_, err = stdin.Write([]byte(html))
	if err != nil {
		return copyPlainText(text)
	}
	stdin.Close()

	if err := cmd.Wait(); err != nil {
		return copyPlainText(text)
	}

	// Also set plain text
	return copyPlainText(text)
}

// copyPlainText is a fallback that copies plain text without formatting
func copyPlainText(text string) error {
	plainText := StripFormatting(text)

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("pbcopy")
	case "windows":
		cmd = exec.Command("powershell", "-Command", "Set-Clipboard", "-Value", plainText)
	case "linux":
		cmd = exec.Command("xclip", "-selection", "clipboard")
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	if runtime.GOOS != "windows" {
		stdin, err := cmd.StdinPipe()
		if err != nil {
			return err
		}

		if err := cmd.Start(); err != nil {
			return err
		}

		_, err = stdin.Write([]byte(plainText))
		if err != nil {
			return err
		}
		stdin.Close()

		return cmd.Wait()
	}

	return cmd.Run()
}

// escapeAppleScript escapes text for use in AppleScript
func escapeAppleScript(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `"`, `\"`)
	s = strings.ReplaceAll(s, "\n", `\n`)
	s = strings.ReplaceAll(s, "\r", `\r`)
	s = strings.ReplaceAll(s, "\t", `\t`)
	return s
}
