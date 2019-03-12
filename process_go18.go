// +build go1.8

package goof

import (
	"debug/dwarf"
	"debug/elf"
	"debug/macho"
	"debug/pe"
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/zeebo/errs"
)

func openProc() (*dwarf.Data, error) {
	path, err := os.Executable()
	if err != nil {
		return nil, errs.Wrap(err)
	}

	var fh interface {
		io.Closer
		DWARF() (*dwarf.Data, error)
	}

	switch runtime.GOOS {
	case "linux":
		fh, err = elf.Open(path)
	case "darwin":
		fh, err = macho.Open(path)
	case "windows":
		fh, err = pe.Open(path)
	default:
		return nil, fmt.Errorf("unknown goos: %q", runtime.GOOS)
	}
	if err != nil {
		return nil, errs.Wrap(err)
	}
	defer fh.Close()

	data, err := fh.DWARF()
	return data, errs.Wrap(err)
}
