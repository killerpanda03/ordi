package deduplicator

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
)

// scanDirectory recursively scans a directory for files
func scanDirectory(dirPath string) tea.Cmd {
	return func() tea.Msg {
		var files []string

		err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && info.Size() > 0 {
				files = append(files, path)
			}
			return nil
		})

		if err != nil {
			return ScanCompleteMsg{Err: err}
		}

		return ScanCompleteMsg{Files: files}
	}
}

// findDuplicates finds duplicate files using parallel SHA256 hashing
func findDuplicates(files []string) tea.Cmd {
	return func() tea.Msg {
		// Group files by size first (optimization)
		sizeGroups := make(map[int64][]string)
		totalSize := int64(0)

		for _, file := range files {
			info, err := os.Stat(file)
			if err != nil {
				continue
			}
			size := info.Size()
			totalSize += size
			sizeGroups[size] = append(sizeGroups[size], file)
		}

		// Only hash files that have potential duplicates (same size)
		var filesToHash []string
		for _, group := range sizeGroups {
			if len(group) > 1 {
				filesToHash = append(filesToHash, group...)
			}
		}

		// Parallel hashing with worker pool
		numWorkers := runtime.NumCPU()
		jobs := make(chan string, len(filesToHash))
		results := make(chan hashResult, len(filesToHash))

		var wg sync.WaitGroup

		// Start workers
		for w := 0; w < numWorkers; w++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for path := range jobs {
					hash, size, err := hashFile(path)
					results <- hashResult{
						path: path,
						hash: hash,
						size: size,
						err:  err,
					}
				}
			}()
		}

		// Send jobs
		go func() {
			for _, file := range filesToHash {
				jobs <- file
			}
			close(jobs)
		}()

		// Wait for all workers to finish
		go func() {
			wg.Wait()
			close(results)
		}()

		// Collect results
		hashGroups := make(map[string][]FileInfo)
		processed := 0

		for result := range results {
			processed++
			if result.err != nil {
				continue
			}
			hashGroups[result.hash] = append(hashGroups[result.hash], FileInfo{
				Path: result.path,
				Size: result.size,
			})
		}

		// Find duplicate groups (hash with more than one file)
		var duplicates []DuplicateGroup
		duplicateSize := int64(0)

		for hash, group := range hashGroups {
			if len(group) > 1 {
				// Calculate wasted space (all copies except one)
				wastedSpace := group[0].Size * int64(len(group)-1)
				duplicateSize += wastedSpace

				duplicates = append(duplicates, DuplicateGroup{
					Hash:  hash,
					Files: group,
					Size:  group[0].Size,
				})
			}
		}

		// Find similar images using perceptual hashing
		// Threshold of 10 means up to 10 bits different (out of 64)
		// This allows for minor variations while avoiding false positives
		similarImages, err := findSimilarImages(files, 10)
		if err != nil {
			// Log error but don't fail - similar images is optional
			similarImages = []SimilarGroup{}
		}

		return HashCompleteMsg{
			Duplicates:    duplicates,
			SimilarImages: similarImages,
			TotalSize:     totalSize,
			DuplicateSize: duplicateSize,
		}
	}
}

type hashResult struct {
	path string
	hash string
	size int64
	err  error
}

// hashFile computes SHA256 hash of a file
func hashFile(path string) (string, int64, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", 0, err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return "", 0, err
	}

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", 0, err
	}

	return hex.EncodeToString(hasher.Sum(nil)), info.Size(), nil
}

// deleteDuplicates deletes selected duplicate files
func deleteDuplicates(groups []DuplicateGroup) tea.Cmd {
	return func() tea.Msg {
		deletedCount := 0
		freedSpace := int64(0)
		var lastErr error

		for _, group := range groups {
			for _, file := range group.Files {
				if file.Selected {
					err := os.Remove(file.Path)
					if err != nil {
						lastErr = fmt.Errorf("failed to delete %s: %w", file.Path, err)
						continue
					}
					deletedCount++
					freedSpace += file.Size
				}
			}
		}

		if lastErr != nil && deletedCount == 0 {
			return DeleteCompleteMsg{Err: lastErr}
		}

		return DeleteCompleteMsg{
			DeletedCount: deletedCount,
			FreedSpace:   freedSpace,
			Err:          lastErr,
		}
	}
}

