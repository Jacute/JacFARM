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
	ErrInvalidPath  = errors.New("zip entry has invalid path")
	ErrTooLargeFile = errors.New("zip entry is too large")
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
		// Защита от zip slip
		if strings.Contains(f.Name, "..") {
			return ErrInvalidPath
		}

		cleanPath := filepath.Clean(f.Name)
		destPath := filepath.Join(targetDir, cleanPath)
		if !strings.HasPrefix(destPath, filepath.Clean(targetDir)+string(os.PathSeparator)) {
			return ErrInvalidPath
		}

		// Проверка размера
		if int64(f.UncompressedSize64) > maxFileSize {
			return ErrTooLargeFile
		}
		totalSize += int64(f.UncompressedSize64)
		if totalSize > maxTotalSize {
			return errors.New("unzip exceeds max allowed total size")
		}

		// Создание директорий
		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(destPath, f.Mode()); err != nil {
				return err
			}
			continue
		}

		// Создаём директорию, если надо
		if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			return err
		}

		// Открываем файл внутри архива
		rc, err := f.Open()
		if err != nil {
			return err
		}

		outFile, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			rc.Close()
			return err
		}

		// Ограничим запись maxFileSize
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
