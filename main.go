package main // Define the main package for the executable program

import (
	"io"            // Provides functions to copy data between readers and writers
	"log"           // Provides simple logging functions
	"net/http"      // Supports HTTP client and server functionality
	"net/url"       // Handles parsing and formatting URLs
	"os"            // Enables file and directory manipulation
	"path"          // Contains functions to manipulate URL-style paths
	"path/filepath" // Contains functions to manipulate native file system paths
	"strings"       // Provides string manipulation utilities
)

// downloadPDF downloads a PDF from the provided raw URL and saves it to the specified output directory.
// It handles redirects, constructs a valid filename, and logs any errors without returning them.
func downloadPDF(rawURL, outputDir string) {
	// Parse the input URL to validate and work with it
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		log.Printf("Invalid URL %q: %v", rawURL, err) // Log the error if URL is invalid
		return                                        // Stop processing this URL
	}

	// Perform an HTTP GET request; follows redirects automatically
	resp, err := http.Get(rawURL)
	if err != nil {
		log.Printf("Failed to fetch %s: %v", rawURL, err) // Log the fetch error
		return
	}
	defer resp.Body.Close() // Ensure the response body is closed after reading

	// Check if the HTTP response status is 200 OK
	if resp.StatusCode != http.StatusOK {
		log.Printf("Bad status for %s: %s", rawURL, resp.Status) // Log non-OK status codes
		return
	}

	// Get the final URL after any redirects (some sites redirect to CDN or file hosting services)
	finalURL := resp.Request.URL.String()

	// Log both original and final URLs for traceability
	log.Printf("Original URL: %s", parsedURL)
	log.Printf("Final URL:    %s", finalURL)

	// Parse the final URL to extract the file path and name
	finalParsed, err := url.Parse(finalURL)
	if err != nil {
		log.Printf("Cannot parse final URL %q: %v", finalURL, err) // Log error in parsing redirected URL
		return
	}

	// Extract the file name from the final URL path
	fileName := path.Base(finalParsed.Path)
	if fileName == "" || fileName == "/" {
		log.Printf("Could not determine file name from %q", finalURL) // Log if file name is invalid
		return
	}

	// Append .pdf extension if the file doesn't have one
	if !strings.HasSuffix(strings.ToLower(fileName), ".pdf") {
		fileName += ".pdf"
	}

	// Create the output directory if it doesn't exist
	err = os.MkdirAll(outputDir, 0755)
	if err != nil {
		log.Printf("Failed to create directory %s: %v", outputDir, err) // Log directory creation failure
		return
	}

	// Build the complete path to where the file will be saved
	filePath := filepath.Join(outputDir, fileName)

	// Create a file at the target location
	outFile, err := os.Create(filePath)
	if err != nil {
		log.Printf("Failed to create file %s: %v", filePath, err) // Log file creation failure
		return
	}
	defer outFile.Close() // Ensure the file is closed after writing

	// Copy the contents of the response body to the output file
	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		log.Printf("Failed to write to %s: %v", filePath, err) // Log file writing errors
		return
	}

	// Log the successful download
	log.Printf("Downloaded to %s\n", filePath)
}

// main is the entry point of the program
func main() {
	// Define a list of PDF URLs to be downloaded
	pdfURLs := []string{
		"https://nypdonline.org/files/948580_01142022_2022007.pdf",
		"https://nypdonline.org/files/965915_10102023_2023071.pdf",
	}

	// Set the output directory where the PDFs will be saved
	outputDir := "nypd_pdfs"

	// Loop through each URL and attempt to download the PDF
	for _, url := range pdfURLs {
		downloadPDF(url, outputDir)
	}
}
