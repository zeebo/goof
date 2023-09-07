package goof

import (
	"debug/dwarf"
	"reflect"
	"sync"
)

type Troop struct {
	once      sync.Once
	err       error
	data      *dwarf.Data
	types     map[string]reflect.Type
	globals   map[string]reflect.Value
	functions map[string]functionCacheEntry
	failures  map[string]error
}

type functionCacheEntry struct {
	pc     uint64
	dtypes []dwarf.Type
}

func (t *Troop) init() {
	t.data, t.err = openProc()
	if t.err != nil {
		return
	}

	t.failures = make(map[string]error)

	t.types = make(map[string]reflect.Type)
	t.err = t.addTypes()
	if t.err != nil {
		return
	}

	t.globals = make(map[string]reflect.Value)
	t.err = t.addGlobals()
	if t.err != nil {
		return
	}

	t.functions = make(map[string]functionCacheEntry)
	t.err = t.addFunctions()
	if t.err != nil {
		return
	}
}

func (t *Troop) check() error {
	t.once.Do(t.init)
	return t.err
}
