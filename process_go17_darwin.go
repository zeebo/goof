// +build !go1.8,darwin

package goof

import (
	"debug/dwarf"
	"debug/macho"
	"unsafe"

	"github.com/zeebo/errs"
)

// #include <mach-o/dyld.h>
// #include <stdlib.h>
import "C"

func openProc() (*dwarf.Data, error) {
	const bufsize = 4096

	buf := (*C.char)(C.malloc(bufsize))
	defer C.free(unsafe.Pointer(buf))

	size := C.uint32_t(bufsize)
	if rc := C._NSGetExecutablePath(buf, &size); rc != 0 {
		return nil, errs.New("error in cgo call to get path: %d", rc)
	}

	fh, err := macho.Open(C.GoString(buf))
	if err != nil {
		return nil, errs.Wrap(err)
	}
	defer fh.Close()

	dwarf, err := fh.DWARF()
	return dwarf, errs.Wrap(err)
}
