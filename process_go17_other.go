// +build !go1.8,!darwin

package goof

import (
	"debug/dwarf"
	"debug/elf"

	"github.com/zeebo/errs"
)

func openProc() (*dwarf.Data, error) {
	fh, err := elf.Open("/proc/self/exe")
	if err != nil {
		return nil, errs.Wrap(err)
	}
	defer fh.Close()

	data, err := fh.DWARF()
	return data, errs.Wrap(err)
}
