package archivewalker

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
)

type File interface {
	io.Reader
	io.ReaderAt
	io.Seeker
}

// NewFile reads the while file into memory and provides a File interface.
func NewFile(fi io.Reader, size int64) (File, error) {
	buf := bytes.NewBuffer(make([]byte, 0, size))
	written, err := io.Copy(buf, fi)
	if err != nil {
		return nil, err
	}
	if written != size {
		return nil, fmt.Errorf("size mismatch: expected %d, got %d", size, written)
	}

	return bytes.NewReader(buf.Bytes()), nil
}

type WalkFileFunc func(path string, info fs.FileInfo, r io.ReaderAt, err error) error

// WalkFiles walks over all archive files, directories and symlinks and reads them into memory.
// Allowing for reading at arbitrary positions in the extracted file.
func WalkFiles(path string, walkcFunc WalkFileFunc) error {
	return Walk(path, func(path string, info fs.FileInfo, r io.Reader, err error) error {
		if err != nil {
			return walkcFunc(path, info, nil, err)
		}

		f, err := NewFile(r, info.Size())
		if err != nil {
			return fmt.Errorf("failed to read file %s into memory: %w", path, err)
		}

		return walkcFunc(path, info, f, nil)
	})
}
