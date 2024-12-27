package cmd

import (
	"io"
	"io/fs"
	"os"
)

type ReadWriterWithStat interface {
	io.ReadWriter

	Stat() (fs.FileInfo, error)
}

func isPipe(f ReadWriterWithStat) (bool, error) {
	stat, err := f.Stat()
	if err != nil {
		return false, err
	}

	if stat.Mode()&os.ModeCharDevice == 0 {
		return true, nil
	}
	return false, nil
}
