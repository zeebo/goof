package goof

import (
	"debug/dwarf"
	"fmt"
	"reflect"
	"sort"
	"unsafe"

	"github.com/zeebo/errs"
)

func (t *Troop) addFunctions() error {
	reader := t.data.Reader()

	for {
		entry, err := reader.Next()
		if err != nil {
			return errs.Wrap(err)
		}
		if entry == nil {
			break
		}

		if entry.Tag != dwarf.TagSubprogram {
			continue
		}

		name, ok := entry.Val(dwarf.AttrName).(string)
		if !ok {
			continue
		}

		pc, ok := entry.Val(dwarf.AttrLowpc).(uint64)
		if !ok {
			continue
		}

		dtypes, err := getFunctionArgTypes(t.data, entry)
		if err != nil {
			continue
		}

		_, err = t.findDwarfTypes(dtypes)
		if err != nil {
			continue
		}

		t.functions[name] = functionCacheEntry{
			pc:     pc,
			dtypes: dtypes,
		}
	}

	return nil
}

func (t *Troop) Functions() ([]string, error) {
	if err := t.check(); err != nil {
		return nil, err
	}

	out := make([]string, 0, len(t.functions))
	for name := range t.functions {
		out = append(out, name)
	}
	sort.Strings(out)
	return out, nil
}

func (t *Troop) Call(name string, args ...interface{}) ([]interface{}, error) {
	if err := t.check(); err != nil {
		return nil, err
	}

	// and so it begins. find the function, the pc, and the types of the args
	// and return values. we don't know which is which, but we assume the
	// caller passed the appropriate things.
	entry, ok := t.functions[name]
	if !ok {
		return nil, fmt.Errorf("call %s: unknown or uncallable function", name)
	}
	pc, dtypes := entry.pc, entry.dtypes

	// make sure they didn't pass more arguments than total types
	num_results := len(dtypes) - len(args)
	if num_results < 0 {
		return nil, fmt.Errorf("call %s: too many args", name)
	}

	// build the arguments, checking consistency and doing hacks to make it
	// nice to use.
	in, in_types, err := t.buildArguments(args, dtypes[:len(args)])
	if err != nil {
		return nil, fmt.Errorf("call %s: %v", name, err)
	}

	// the rest must be the return values, RIGHT? LOL
	out_types, err := t.findDwarfTypes(dtypes[len(args):])
	if err != nil {
		return nil, fmt.Errorf("call %s: %v", name, err)
	}

	// make it happen
	fn_typ := reflect.FuncOf(in_types, out_types, false)
	fn := reflect.ValueOf(makeInterface(dataPtr(fn_typ), unsafe.Pointer(&pc)))
	return ifaces(fn.Call(in)), nil
}

func (t *Troop) buildArguments(args []interface{}, dtypes []dwarf.Type) (
	[]reflect.Value, []reflect.Type, error) {

	if len(args) != len(dtypes) {
		return nil, nil, fmt.Errorf("number of arguments does not match")
	}

	// so I want the Call signatrue to have args passed as ...interface{}
	// because that makes the api nice: taking a reflect.Value is hard for the
	// consumer.
	//
	// Unfortunately, that means that if you pass an interface type in, they
	// get down-casted to just interface{}. Doubly unfortunately, the itab
	// pointer was lost when the value was converted to an interface{} instead
	// of whatever interface it was in the first place.
	//
	// So let's just look and see if we can find the interface type in the
	// types map based on the dwarf type. If not, shoot. Hopefully that's
	// rare! Maybe we can expose a CallReflect api that can get around this.
	//
	// Heh.

	in_types := make([]reflect.Type, 0, len(args))
	in := make([]reflect.Value, 0, len(args))

	for i, arg := range args {
		dtyp := dtypes[i]

		val, err := t.constructValue(arg, dtyp)
		if err != nil {
			return nil, nil, fmt.Errorf("arg %d: %v", i, err)
		}

		in_types = append(in_types, val.Type())
		in = append(in, val)
	}

	return in, in_types, nil
}

func (t *Troop) constructValue(arg interface{}, dtyp dwarf.Type) (
	val reflect.Value, err error) {

	rtyp, err := t.consistentValue(arg, dtyp)
	if err != nil {
		return val, err
	}
	if arg == nil {
		return reflect.New(rtyp).Elem(), nil
	}
	return reflect.ValueOf(arg).Convert(rtyp), nil
}

func (t *Troop) consistentValue(arg interface{}, dtyp dwarf.Type) (
	reflect.Type, error) {

	rtyp, err := t.findDwarfType(dtyp)
	if err != nil {
		return nil, err
	}
	if arg == nil {
		if !reflectCanBeNil(rtyp) {
			return nil, fmt.Errorf("passing nil to non-nillable type: %v",
				rtyp)
		}
		return rtyp, nil
	}
	if !reflect.TypeOf(arg).ConvertibleTo(rtyp) {
		return nil, fmt.Errorf("cannot convert %v to %T", rtyp, arg)
	}
	return rtyp, nil
}
