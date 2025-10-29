package utils

import (
	"archive/zip"
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var (
	ErrInvalidPath          = errors.New("zip entry has invalid path")
	ErrTooLargeFile         = errors.New("zip entry is too large")
	ErrTooLargeZipTotalSize = errors.New("zip total size after unzipping is too large")
)

// SecureUnzip безопасно разархивирует zip-архив в targetDir.
// maxTotalSize — лимит общего распакованного размера (например, 500 МБ).
// maxFileSize — лимит для одного файла (например, 50 МБ).
func SecureUnzip(zipBytes []byte, targetDir string, maxTotalSize, maxFileSize int64) error {
	r, err := zip.NewReader(bytes.NewReader(zipBytes), int64(len(zipBytes)))
	if err != nil {
		return err
	}

	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return err
	}

	var totalSize int64

	for _, f := range r.File {
		// zip slip protect
		cleanPath := filepath.Clean(f.Name)
		destPath := filepath.Join(targetDir, cleanPath)
		if !strings.HasPrefix(destPath, filepath.Clean(targetDir)+string(os.PathSeparator)) {
			return ErrInvalidPath
		}

		// check file size
		if int64(f.UncompressedSize64) > maxFileSize {
			return ErrTooLargeFile
		}
		totalSize += int64(f.UncompressedSize64)
		// check total size
		if totalSize > maxTotalSize {
			return ErrTooLargeZipTotalSize
		}

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(destPath, f.Mode()); err != nil {
				return err
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			return err
		}

		outFile, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			rc.Close()
			return err
		}

		written, err := io.CopyN(outFile, rc, maxFileSize+1)
		if err != nil && err != io.EOF {
			outFile.Close()
			rc.Close()
			return err
		}
		if written > maxFileSize {
			outFile.Close()
			rc.Close()
			return ErrTooLargeFile
		}

		outFile.Close()
		rc.Close()
	}

	return nil
}
