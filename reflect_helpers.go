package goof

import (
	"fmt"
	"path"
	"reflect"
	"strings"
	"unsafe"
)

func reflectCanBeNil(rtyp reflect.Type) bool {
	switch rtyp.Kind() {
	case reflect.Interface,
		reflect.Ptr,
		reflect.Map,
		reflect.Chan,
		reflect.Slice:
		return true
	}
	return false
}

func typesByString(types []reflect.Type) sortTypesByString {
	cache := make([]string, 0, len(types))
	for _, typ := range types {
		cache = append(cache, typ.String())
	}
	return sortTypesByString{
		types: types,
		cache: cache,
	}
}

type sortTypesByString struct {
	types []reflect.Type
	cache []string
}

func (s sortTypesByString) Len() int { return len(s.types) }

func (s sortTypesByString) Less(i, j int) bool {
	return s.cache[i] < s.cache[j]
}

func (s sortTypesByString) Swap(i, j int) {
	s.cache[i], s.cache[j] = s.cache[j], s.cache[i]
	s.types[i], s.types[j] = s.types[j], s.types[i]
}

var (
	unsafePointerType = reflect.TypeOf((*unsafe.Pointer)(nil)).Elem()
)

// dwarfName does a best effort to return the dwarf entry name for the reflect
// type so that we can map between them. here's hoping it doesn't do it wrong
func dwarfName(rtyp reflect.Type) (out string) {
	pkg_path := rtyp.PkgPath()
	name := rtyp.Name()

	switch {
	// this type's PkgPath returns "" instead of "unsafe". hah.
	case rtyp == unsafePointerType:
		return "unsafe.Pointer"

	case pkg_path != "":
		// this is crazy, but sometimes a dot is encoded as %2e, but only when
		// it's in the last path component. i wonder if this is sufficient.
		if strings.Contains(pkg_path, "/") {
			dir := path.Dir(pkg_path)
			base := strings.Replace(path.Base(pkg_path), ".", "%2e", -1)
			pkg_path = dir + "/" + base
		}

		return pkg_path + "." + name

	case name != "":
		return name

	default:
		switch rtyp.Kind() {
		case reflect.Ptr:
			return "*" + dwarfName(rtyp.Elem())

		case reflect.Slice:
			return "[]" + dwarfName(rtyp.Elem())

		case reflect.Array:
			return fmt.Sprintf("[%d]%s",
				rtyp.Len(),
				dwarfName(rtyp.Elem()))

		case reflect.Map:
			return fmt.Sprintf("map[%s]%s",
				dwarfName(rtyp.Key()),
				dwarfName(rtyp.Elem()))

		case reflect.Chan:
			prefix := "chan"
			switch rtyp.ChanDir() {
			case reflect.SendDir:
				prefix = "chan<-"
			case reflect.RecvDir:
				prefix = "<-chan"
			}
			return fmt.Sprintf("%s %s",
				prefix,
				dwarfName(rtyp.Elem()))

		// TODO: func, struct

		default:
			// oh well. this sometimes works.
			return rtyp.String()
		}
	}
}
