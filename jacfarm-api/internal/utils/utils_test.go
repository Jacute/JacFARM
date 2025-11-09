package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSecureUnzip(t *testing.T) {
	testcases := []struct {
		name string
		file string
		maxTotalSize,
		maxFileSize int64
		expectedError  error
		archiveEntries []string
	}{
		{
			name:           "zipping file with size >195KB should return error",
			file:           "./testcases/001_file195KB.zip",
			maxTotalSize:   500 * 1024, // 500 KB
			maxFileSize:    100 * 1024, // 100 KB
			expectedError:  ErrTooLargeFile,
			archiveEntries: []string{"file.txt"}, // testing without nesting
		},
		{
			name:           "zipping file with size >195KB should return error",
			file:           "./testcases/002_total_size_390KB.zip",
			maxTotalSize:   300 * 1024, // 300 KB
			maxFileSize:    200 * 1024, // 200 KB
			expectedError:  ErrTooLargeZipTotalSize,
			archiveEntries: []string{"file.txt"},
		},
		{
			name:           "ok",
			file:           "./testcases/002_total_size_390KB.zip",
			maxTotalSize:   500 * 1024, // 500 KB
			maxFileSize:    200 * 1024, // 200 KB
			expectedError:  nil,
			archiveEntries: []string{"file.txt", "aboba"},
		},
		{
			name:           "path traversal/zip slip protect should return error",
			file:           "./testcases/003_path_traversal.zip",
			maxTotalSize:   500 * 1024, // 500 KB
			maxFileSize:    200 * 1024, // 200 KB
			expectedError:  ErrInvalidPath,
			archiveEntries: []string{"../file.txt"},
		},
		{
			name:           "ok double dot",
			file:           "./testcases/004_double_dot.zip",
			maxTotalSize:   500 * 1024, // 500 KB
			maxFileSize:    200 * 1024, // 200 KB
			expectedError:  nil,
			archiveEntries: []string{"file..txt"},
		},
	}

	for _, tc := range testcases {
		if _, err := os.Stat(tc.file); os.IsNotExist(err) {
			t.Fatalf("test file not found at %s", tc.file)
		}
		data, err := os.ReadFile(tc.file)
		if err != nil {
			t.Fatalf("error reading testfile %s", tc.file)
		}

		tmpDir, err := os.MkdirTemp("/tmp", "jacfarm_test_secury_unzip_")
		if err != nil {
			t.Fatal("error creating temp dir")
		}
		t.Run(tc.name, func(t *testing.T) {
			err := SecureUnzip(data, tmpDir, tc.maxTotalSize, tc.maxFileSize)
			require.ErrorIs(t, err, tc.expectedError)

			if tc.expectedError != nil {
				entries, err := os.ReadDir(tmpDir)
				if err != nil {
					t.Errorf("test error reading tmp dir %s", tmpDir)
				}
				for _, entry := range entries {
					require.Contains(t, tc.archiveEntries, entry.Name())
				}
			}
		})
	}
}
