package compressor

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)


type ToolAvailability struct {
	FFmpeg      bool
	ImageMagick bool
	Ghostscript bool
	SevenZip    bool
}


func (t ToolAvailability) SupportedFormats() []string {
	formats := []string{}

	if t.FFmpeg {
		formats = append(formats, ".mp4", ".mov", ".avi", ".mkv", ".mp3", ".wav")
	}
	if t.ImageMagick {
		formats = append(formats, ".jpg", ".jpeg", ".png", ".gif", ".bmp")
	}
	if t.Ghostscript {
		formats = append(formats, ".pdf")
	}
	if t.SevenZip {
		formats = append(formats, ".docx", ".xlsx", ".pptx", ".zip")
	}

	return formats
}


func (t ToolAvailability) HasAnyTool() bool {
	return t.FFmpeg || t.ImageMagick || t.Ghostscript || t.SevenZip
}


func (t ToolAvailability) IsFileSupported(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))

	switch ext {
	case ".mp4", ".mov", ".avi", ".mkv", ".mp3", ".wav":
		return t.FFmpeg
	case ".jpg", ".jpeg", ".png", ".gif", ".bmp":
		return t.ImageMagick
	case ".pdf":
		return t.Ghostscript
	case ".docx", ".xlsx", ".pptx", ".zip":
		return t.SevenZip
	default:
		return false
	}
}


func checkExternalTools() ToolAvailability {
	tools := ToolAvailability{}

	
	if _, err := exec.LookPath("ffmpeg"); err == nil {
		tools.FFmpeg = true
	}

	
	if _, err := exec.LookPath("convert"); err == nil {
		tools.ImageMagick = true
	} else if _, err := exec.LookPath("magick"); err == nil {
		
		tools.ImageMagick = true
	}

	
	if _, err := exec.LookPath("gs"); err == nil {
		tools.Ghostscript = true
	} else if _, err := exec.LookPath("gswin64c"); err == nil {
		
		tools.Ghostscript = true
	}

	
	if _, err := exec.LookPath("7z"); err == nil {
		tools.SevenZip = true
	} else if _, err := exec.LookPath("7za"); err == nil {
		tools.SevenZip = true
	}

	return tools
}

func compressFile(inputPath, outputPath string) error {
	ext := strings.ToLower(filepath.Ext(inputPath))

	var cmd *exec.Cmd

	switch ext {
	
	case ".mp4", ".mov", ".avi", ".mkv":
		cmd = exec.Command("ffmpeg", "-i", inputPath, "-vcodec", "libx265", "-crf", "28", outputPath)

	
	case ".mp3", ".wav":
		cmd = exec.Command("ffmpeg", "-i", inputPath, "-codec:a", "libmp3lame", "-qscale:a", "2", outputPath)

	
	case ".jpg", ".jpeg", ".png", ".gif", ".bmp":
		
		convertCmd := "convert"
		if _, err := exec.LookPath("magick"); err == nil {
			convertCmd = "magick"
		}
		cmd = exec.Command(convertCmd, inputPath, "-quality", "75", outputPath)

	
	case ".pdf":
		
		gsCmd := "gs"
		if _, err := exec.LookPath("gswin64c"); err == nil {
			gsCmd = "gswin64c"
		}
		cmd = exec.Command(gsCmd,
			"-sDEVICE=pdfwrite",
			"-dCompatibilityLevel=1.4",
			"-dPDFSETTINGS=/ebook",
			"-dNOPAUSE",
			"-dQUIET",
			"-dBATCH",
			"-sOutputFile="+outputPath,
			inputPath)

	
	case ".docx", ".xlsx", ".pptx", ".zip":
		
		sevenZipCmd := "7z"
		if _, err := exec.LookPath("7za"); err == nil {
			sevenZipCmd = "7za"
		}
		cmd = exec.Command(sevenZipCmd, "a", "-tzip", "-mx=9", outputPath, inputPath)

	default:
		return fmt.Errorf("Nicht unterstützter Dateityp: %s", ext)
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Fehler bei der Komprimierung von %s: %w", filepath.Base(inputPath), err)
	}
	return nil
}
