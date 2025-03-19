package archivewalker_test

import (
	"io"
	"io/fs"
	"testing"

	"github.com/jxsl13/archivewalker"
)

func TestWalk(t *testing.T) {
	err := archivewalker.Walk("testdata/test.tar", func(name string, fi fs.FileInfo, r io.Reader, err error) error {
		if err != nil {
			return err
		}

		if fi.IsDir() {
			return nil
		}

		if fi.Mode()|fs.ModeSymlink != 0 {
			return nil
		}
		t.Logf("file: %s, size: %d", name, fi.Size())

		data, err := io.ReadAll(r)
		if err != nil {
			return err
		}

		t.Logf("File content:\n%s", string(data))
		return nil
	})

	if err != nil {
		t.Error(err)
		return
	}

}
