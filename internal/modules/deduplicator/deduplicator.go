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


func findDuplicates(files []string) tea.Cmd {
	return func() tea.Msg {
		
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

		
		var filesToHash []string
		for _, group := range sizeGroups {
			if len(group) > 1 {
				filesToHash = append(filesToHash, group...)
			}
		}

		
		numWorkers := runtime.NumCPU()
		jobs := make(chan string, len(filesToHash))
		results := make(chan hashResult, len(filesToHash))

		var wg sync.WaitGroup

		
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

		
		go func() {
			for _, file := range filesToHash {
				jobs <- file
			}
			close(jobs)
		}()

		
		go func() {
			wg.Wait()
			close(results)
		}()

		
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

		
		var duplicates []DuplicateGroup
		duplicateSize := int64(0)

		for hash, group := range hashGroups {
			if len(group) > 1 {
				
				wastedSpace := group[0].Size * int64(len(group)-1)
				duplicateSize += wastedSpace

				duplicates = append(duplicates, DuplicateGroup{
					Hash:  hash,
					Files: group,
					Size:  group[0].Size,
				})
			}
		}

		
		
		
		similarImages, err := findSimilarImages(files, 10)
		if err != nil {
			
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

