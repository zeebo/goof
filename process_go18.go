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
)

func openProc() (*dwarf.Data, error) {
	path, err := os.Executable()
	if err != nil {
		return nil, err
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
		return nil, err
	}
	defer fh.Close()

	return fh.DWARF()
}
