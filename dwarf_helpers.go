package goof

import (
	"debug/dwarf"
	"encoding/binary"
	"unsafe"

	"github.com/zeebo/errs"
)

func dwarfTypeName(dtyp dwarf.Type) string {
	// for some reason the debug/dwarf package doesn't set the Name field
	// on the common type for struct fields. what is this misery?
	switch dtyp := dtyp.(type) {
	case *dwarf.StructType:
		return dtyp.StructName
	default:
		return dtyp.Common().Name
	}
}

func getFunctionArgTypes(data *dwarf.Data, entry *dwarf.Entry) (
	[]dwarf.Type, error) {

	reader := data.Reader()

	reader.Seek(entry.Offset)
	_, err := reader.Next()
	if err != nil {
		return nil, err
	}

	var args []dwarf.Type
	for {
		child, err := reader.Next()
		if err != nil {
			return nil, err
		}
		if child == nil || child.Tag == 0 {
			break
		}

		if child.Tag != dwarf.TagFormalParameter {
			continue
		}

		dtyp, err := entryType(data, child)
		if err != nil {
			return nil, err
		}

		args = append(args, dtyp)
	}

	return args, nil
}

func entryType(data *dwarf.Data, entry *dwarf.Entry) (dwarf.Type, error) {
	off, ok := entry.Val(dwarf.AttrType).(dwarf.Offset)
	if !ok {
		return nil, errs.New("unable to find type offset for entry")
	}
	return data.Type(off)
}

func entryLocation(data *dwarf.Data, entry *dwarf.Entry) (uint64, error) {
	loc, ok := entry.Val(dwarf.AttrLocation).([]byte)
	if !ok {
		return 0, errs.New("unable to find location for entry")
	}
	if len(loc) == 0 {
		return 0, errs.New("location had no data")
	}

	// only support this opcode. did you know dwarf location information is
	// a stack based programming language with opcodes and stuff? i wonder
	// how many interpreters for that have code execution bugs in them.
	if loc[0] != 0x03 {
		return 0, errs.New("can't interpret location information")
	}

	// oh man let's also just assume that the dwarf info is written in the
	// same order as the host! WHAT COULD GO WRONG?!
	switch len(loc) - 1 {
	case 4:
		return uint64(hostOrder.Uint32(loc[1:])), nil
	case 8:
		return uint64(hostOrder.Uint64(loc[1:])), nil
	default:
		return 0, errs.New("what kind of computer is this?")
	}
}

var hostOrder binary.ByteOrder

func init() {
	i := 1
	data := (*[4]byte)(unsafe.Pointer(&i))

	if data[0] == 0 {
		hostOrder = binary.BigEndian
	} else {
		hostOrder = binary.LittleEndian
	}
}
