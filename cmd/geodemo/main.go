package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// License Key from MaxMind account
var LicenseKey = ""

func init() {
	LicenseKey = os.Getenv("MAXMIND_KEY")
}

// URL for GeoIP2-City database (commercial version)
const downloadURL = "https://download.maxmind.com/app/geoip_download?edition_id=GeoIP2-City&license_key=%s&suffix=tar.gz"

// Path to save the downloaded file
const downloadPath = "GeoIP2-City.tar.gz"

// Function to download the file
func downloadFile(url string, filepath string) error {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

// Function to extract the .tar.gz file to a specific directory
func extractTarGz(gzipStream io.Reader, dest string) error {
	uncompressedStream, err := gzip.NewReader(gzipStream)
	if err != nil {
		return fmt.Errorf("gzip.NewReader failed: %w", err)
	}
	defer uncompressedStream.Close()

	tarReader := tar.NewReader(uncompressedStream)

	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			return fmt.Errorf("tarReader.Next() failed: %w", err)
		}

		// Create the target path inside the destination directory
		target := filepath.Join(dest, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			// Create directory if it doesn't exist
			if err := os.MkdirAll(target, 0755); err != nil {
				return fmt.Errorf("os.MkdirAll() failed: %w", err)
			}
		case tar.TypeReg:
			// Create file
			outFile, err := os.Create(target)
			if err != nil {
				return fmt.Errorf("os.Create() failed: %w", err)
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				return fmt.Errorf("io.Copy() failed: %w", err)
			}
			outFile.Close()
		default:
			return fmt.Errorf("unknown type: %v in %s", header.Typeflag, header.Name)
		}
	}
	return nil
}

func main() {
	// Construct the URL for download
	url := fmt.Sprintf(downloadURL, LicenseKey)

	// Download the file
	fmt.Println("Downloading GeoIP2 database...")
	if err := downloadFile(url, downloadPath); err != nil {
		fmt.Println("Error downloading file:", err)
		return
	}
	fmt.Println("Download complete.")

	// Open the downloaded tar.gz file
	file, err := os.Open(downloadPath)
	if err != nil {
		fmt.Println("Error opening downloaded file:", err)
		return
	}
	defer file.Close()

	// Define the target directory for extraction
	targetDir := "./GeoIP2-Database"

	// Extract the tar.gz file to the specific directory
	fmt.Println("Extracting GeoIP2 database...")
	if err := extractTarGz(file, targetDir); err != nil {
		fmt.Println("Error extracting file:", err)
		return
	}
	fmt.Println("Extraction complete. Files extracted to:", targetDir)
}
