package goof

import (
	"debug/dwarf"
	"debug/macho"
	"fmt"
	"unsafe"
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
		return nil, fmt.Errorf("error in cgo call to get path: %d", rc)
	}

	fh, err := macho.Open(C.GoString(buf))
	if err != nil {
		return nil, err
	}
	defer fh.Close()

	return fh.DWARF()
}
