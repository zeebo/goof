package goof

import (
	"errors"
	"fmt"
	"os"
	"testing"
)

//go:noinline
func Barf(err error) int {
	if err != nil {
		return len(err.Error())
	}
	return 5
}

// no dead code elim
func init() { Barf(nil) }

func (suite) TestCall(t *testing.T) {
	symbol := fmt.Sprintf("%s.Barf", importPath)
	if _, err := troop.Call(symbol, errors.New("hi")); err != nil {
		t.Fatal(err)
	}
	if _, err := troop.Call(symbol, nil); err != nil {
		t.Fatal(err)
	}
}

func (suite) TestCallFprintf(t *testing.T) {
	fmt.Fprintf(os.Stdout, "hello world %d\n", 2)
	if _, err := troop.Call("fmt.Fprintf", os.Stdout, "hello world %d\n", []interface{}{2}); err != nil {
		t.Fatal(err)
	}
}

func (suite) TestCallFailures(t *testing.T) {
	symbol := fmt.Sprintf("%s.Barf", importPath)
	type is []interface{}

	cases := []struct {
		name string
		args []interface{}
	}{
		{symbol, is{nil, nil}},   // too many args
		{symbol, is{false}},      // wrong arg kind
		{symbol, is{"hello", 2}}, // wrong arg kind
	}

	for i, c := range cases {
		out, err := troop.Call(c.name, c.args...)
		if err == nil {
			t.Logf("%d: %+v", i, c)
			t.Errorf("expected an error. out: %#v", out)
		}
	}
}

func (suite) TestGlobals(t *testing.T) {
	troop_rv, err := troop.Global(fmt.Sprintf("%s.troop", importPath))
	if err != nil {
		t.Fatal(err)
	}
	troop2 := troop_rv.Addr().Interface().(*Troop)
	if troop2 != &troop {
		t.Fatal("got a different troop")
	}
}
