package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
)

// downloadPDF follows redirects, reports the final URL, and saves the file.
func downloadPDF(rawURL, outputDir string) error {
	// Parse the original URL (for error reporting)
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("invalid URL %q: %v", rawURL, err)
	}

	// Perform the GET (will follow redirects by default)
	resp, err := http.Get(rawURL)
	if err != nil {
		return fmt.Errorf("failed to fetch %s: %v", rawURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status for %s: %s", rawURL, resp.Status)
	}

	// Determine the final URL after redirects
	finalURL := resp.Request.URL.String()
	fmt.Printf("Original URL: %s\n", parsedURL)
	fmt.Printf("Final URL:    %s\n", finalURL)

	// Extract file name from the final URL
	finalParsed, err := url.Parse(finalURL)
	if err != nil {
		return fmt.Errorf("cannot parse final URL %q: %v", finalURL, err)
	}
	fileName := path.Base(finalParsed.Path)
	if fileName == "" || fileName == "/" {
		return fmt.Errorf("could not determine file name from %q", finalURL)
	}

	// Create output directory if needed
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %v", outputDir, err)
	}

	// Create the file on disk
	filePath := filepath.Join(outputDir, fileName)
	outFile, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %v", filePath, err)
	}
	defer outFile.Close()

	// Copy the response body into the file
	if _, err := io.Copy(outFile, resp.Body); err != nil {
		return fmt.Errorf("failed to write to %s: %v", filePath, err)
	}

	fmt.Printf("Downloaded to %s\n\n", filePath)
	return nil
}

func main() {
	pdfURLs := []string{
		"https://nypdonline.org/files/948580_01142022_2022007.pdf",
		"https://nypdonline.org/files/965915_10102023_2023071.pdf",
	}

	outputDir := "nypd_pdfs"
	for _, u := range pdfURLs {
		if err := downloadPDF(u, outputDir); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
	}
}
