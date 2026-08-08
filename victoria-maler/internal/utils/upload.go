package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var allowedExtensions = map[string]bool{
	".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true,
	".pdf": true, ".csv": true, ".txt": true,
}

const maxUploadSize = 10 << 20 // 10 MB

func SaveUpload(uploadDir, fileName string, src io.Reader) (string, error) {
	ext := strings.ToLower(filepath.Ext(fileName))
	if !allowedExtensions[ext] {
		return "", fmt.Errorf("file type %q not allowed", ext)
	}

	if err := os.MkdirAll(uploadDir, 0o755); err != nil {
		return "", err
	}

	dst, err := os.CreateTemp(uploadDir, "upload-*"+ext)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	written, err := io.Copy(dst, io.LimitReader(src, maxUploadSize+1))
	if err != nil {
		os.Remove(dst.Name())
		return "", err
	}
	if written > maxUploadSize {
		os.Remove(dst.Name())
		return "", fmt.Errorf("file exceeds 10 MB limit")
	}

	return filepath.Base(dst.Name()), nil
}
