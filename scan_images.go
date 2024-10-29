package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// List of common image file extensions, including mobile-specific and RAW formats
var imageExtensions = map[string]bool{
	".jpg":   true,
	".jpeg":  true,
	".png":   true,
	".gif":   true,
	".bmp":   true,
	".tiff":  true,
	".svg":   true,
	".webp":  true,
	".heic":  true, // High Efficiency Image Format used on iOS devices
	".heif":  true, // Another High Efficiency Image Format extension
	".raw":   true, // RAW image format
	".cr2":   true, // Canon RAW format
	".nef":   true, // Nikon RAW format
	".orf":   true, // Olympus RAW format
	".sr2":   true, // Sony RAW format
	".arw":   true, // Sony RAW format
	".dng":   true, // Adobe Digital Negative format
	".rw2":   true, // Panasonic RAW format
}

// List of directories to ignore
var ignoreDirs = []string{
	"Windows",
	"Program Files",
	"System Volume Information",
	"$Recycle.Bin",
	"Users",
	"SmartPSS",
	"Python312",
	"ProgramData",
}

// scanDirectory scans the specified directory, prints each image file's path and size, 
// adds up the total size of all found images, and tracks image counts per directory.
func scanDirectory(root string) (int64, map[string]int, error) {
	var totalSize int64
	dirFileCount := make(map[string]int)

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip ignored directories
		for _, ignoreDir := range ignoreDirs {
			if info.IsDir() && strings.Contains(path, ignoreDir) {
				return filepath.SkipDir
			}
		}

		// Check if the file has an image extension and process it
		if !info.IsDir() && imageExtensions[strings.ToLower(filepath.Ext(path))] {
			fileSize := info.Size()
			fmt.Printf("File: %s | Size: %d bytes\n", path, fileSize)
			totalSize += fileSize // Add the file size to the total

			// Track the count of images in each directory
			dir := filepath.Dir(path)
			dirFileCount[dir]++
		}
		return nil
	})

	return totalSize, dirFileCount, err
}

// printDirectoryFileCounts sorts and prints directory paths by file count in descending order
func printDirectoryFileCounts(dirFileCount map[string]int) []string {
	var acceptedDirs []string 
	type dirCount struct {
		path  string
		count int
	}

	var sortedDirs []dirCount
	for dir, count := range dirFileCount {
		sortedDirs = append(sortedDirs, dirCount{path: dir, count: count})
	}

	// Sort by file count in descending order
	sort.Slice(sortedDirs, func(i, j int) bool {
		return sortedDirs[i].count > sortedDirs[j].count
	})

	// Print the sorted directory counts
	fmt.Println("\nDirectories sorted by number of image files:")
	for _, dc := range sortedDirs {
		if dc.count > 5 {
			fmt.Printf("Directory: %s | Image Files: %d\n", dc.path, dc.count)
			acceptedDirs = append(acceptedDirs, dc.path)
		}
	}
	return acceptedDirs
}

func main() {
	// Get the directory to scan from the command line
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run scan_images.go <directory>")
		return
	}

	root := os.Args[1] // Use the provided directory

	fmt.Println("Scanning for image files in:", root)
	totalSize, dirFileCount, err := scanDirectory(root)
	if err != nil {
		fmt.Println("Error scanning directory:", err)
		return
	}

	// Print the total size summary
	fmt.Printf("\nTotal Size: %d bytes\n", totalSize)
	fmt.Printf("Total Size: %.2f KB\n", float64(totalSize)/1024)
	fmt.Printf("Total Size: %.2f MB\n", float64(totalSize)/(1024*1024))
	fmt.Printf("Total Size: %.2f GB\n", float64(totalSize)/(1024*1024*1024))

	// Print directories sorted by the number of image files
	acceptedDirs := printDirectoryFileCounts(dirFileCount)
	for _, dir := range(acceptedDirs) {
		fmt.Printf("%s\n", dir)
	} 
}
