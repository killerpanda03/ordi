package deduplicator

import (
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"

	"github.com/nfnt/resize"
)

func computeDHash(imagePath string) (uint64, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return 0, fmt.Errorf("failed to decode image: %w", err)
	}

	resized := resize.Resize(9, 8, img, resize.Lanczos3)

	// Convert to grayscale and compute hash
	var hash uint64
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			leftPixel := resized.At(x, y)
			rightPixel := resized.At(x+1, y)

			leftGray := rgbaToGray(leftPixel)
			rightGray := rgbaToGray(rightPixel)

			bitIndex := y*8 + x
			if leftGray < rightGray {
				hash |= 1 << bitIndex
			}
		}
	}

	return hash, nil
}

// rgbaToGray converts RGBA color to grayscale
func rgbaToGray(c color.Color) uint32 {
	r, g, b, _ := c.RGBA()
	// Standard grayscale conversion formula
	return (r*299 + g*587 + b*114) / 1000
}


func hammingDistance(hash1, hash2 uint64) int {
	xor := hash1 ^ hash2
	distance := 0
	for xor != 0 {
		distance++
		xor &= xor - 1 // Clear the lowest set bit
	}
	return distance
}

func findSimilarImages(files []string, threshold int) ([]SimilarGroup, error) {
	var imageFiles []string
	for _, file := range files {
		if isImageFile(file) {
			imageFiles = append(imageFiles, file)
		}
	}

	if len(imageFiles) == 0 {
		return nil, nil
	}

	type imageHash struct {
		path string
		hash uint64
		size int64
	}

	var hashes []imageHash
	for _, imagePath := range imageFiles {
		hash, err := computeDHash(imagePath)
		if err != nil {
			continue
		}

		info, err := os.Stat(imagePath)
		if err != nil {
			continue
		}

		hashes = append(hashes, imageHash{
			path: imagePath,
			hash: hash,
			size: info.Size(),
		})
	}

	visited := make(map[int]bool)
	var similarGroups []SimilarGroup

	for i, img1 := range hashes {
		if visited[i] {
			continue
		}

		var group []FileInfo
		group = append(group, FileInfo{
			Path: img1.path,
			Size: img1.size,
		})

		visited[i] = true
		totalDistance := 0
		comparisons := 0

		for j := i + 1; j < len(hashes); j++ {
			if visited[j] {
				continue
			}

			distance := hammingDistance(img1.hash, hashes[j].hash)
			if distance <= threshold {
				group = append(group, FileInfo{
					Path: hashes[j].path,
					Size: hashes[j].size,
				})
				visited[j] = true
				totalDistance += distance
				comparisons++
			}
		}

		if len(group) > 1 {
			avgDistance := 0
			if comparisons > 0 {
				avgDistance = totalDistance / comparisons
			}

			similarity := 100.0 - (float64(avgDistance)/64.0)*100.0

			similarGroups = append(similarGroups, SimilarGroup{
				Files:      group,
				Similarity: similarity,
			})
		}
	}

	return similarGroups, nil
}


func isImageFile(path string) bool {
	ext := filepath.Ext(path)
	imageExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".bmp":  true,
		".tiff": true,
		".webp": true,
	}
	return imageExts[ext]
}
