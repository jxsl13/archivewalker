package archivewalker

import (
	"fmt"
	"io"
	"io/fs"
	"os"

	"github.com/gabriel-vasile/mimetype"
)

var (
	ErrUnsupportedArchive = fmt.Errorf("unsupported archive")
)

// WalkFunc defines the function in order to efficiently walk over the archive
type WalkFunc func(path string, info fs.FileInfo, r io.Reader, err error) error

func Walk(path string, walkcFunc WalkFunc) error {

	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		return err
	}

	mime, err := mimetype.DetectReader(f)
	if err != nil {
		return fmt.Errorf("could not detect mime type: %w", err)
	}
	_, err = f.Seek(0, io.SeekStart)
	if err != nil {
		return fmt.Errorf("could not seek to start of file: %w", err)
	}

	switch mime.Extension() {
	case ".7z":
		return Walk7Zip(f, stat.Size(), walkcFunc)
	case ".gz":
		return WalkTarGzip(f, walkcFunc)
	case ".tar":
		return WalkTar(f, walkcFunc)
	case ".zip":
		return WalkZip(f, stat.Size(), walkcFunc)
	case ".xz":
		return WalkTarXz(f, walkcFunc)
	case ".zst":
		return WalkTarZstd(f, walkcFunc)
	case ".bz2":
		return WalkTarBzip2(f, walkcFunc)
	case ".lz":
		return WalkTarLz(f, walkcFunc)
	}
	return fmt.Errorf("%w: %s", ErrUnsupportedArchive, mime.Extension())
}
