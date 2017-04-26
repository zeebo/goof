// +build !go1.8,!darwin

package goof

import (
	"debug/dwarf"
	"debug/elf"
)

func openProc() (*dwarf.Data, error) {
	fh, err := elf.Open("/proc/self/exe")
	if err != nil {
		return nil, err
	}
	defer fh.Close()

	return fh.DWARF()
}
