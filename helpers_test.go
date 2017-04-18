package goof

import (
	"os"
	"os/exec"
	"reflect"
	"strings"
	"testing"
)

const importPath = "github.com/zeebo/goof"

var troop Troop

type suite struct{}

func Test(t *testing.T) {
	// normal test invocations don't come with dwarf info, so we can't actually
	// run the tests :(
	//
	// it's ok. skip it here and hope TestCompile checks it.
	if err := troop.check(); err != nil {
		t.Log("skipping test due to error:", troop.check())
		t.SkipNow()
	}

	// this way we don't have to remember to add the tests to this function
	// to run.
	s := reflect.TypeOf(suite{})
	for i := 0; i < s.NumMethod(); i++ {
		method := s.Method(i)
		if !strings.HasPrefix(method.Name, "Test") {
			continue
		}
		t.Run(method.Name[4:], func(t *testing.T) {
			method.Func.Call([]reflect.Value{
				reflect.ValueOf(suite{}),
				reflect.ValueOf(t),
			})
		})
	}
}

// TestCompile tries to actually run the tests on an environment that can
// compile the test. HEH.
func TestCompile(t *testing.T) {
	if os.Getenv("GOOF_SKIP_COMPILE") != "" {
		t.SkipNow()
	}
	if err := exec.Command("go", "test", "-c", importPath).Run(); err != nil {
		t.Logf("skipping compile test due to error: %v", err)
		t.SkipNow()
	}
	cmd := exec.Command("./goof.test", "-test.v")
	cmd.Env = append(cmd.Env, "GOOF_SKIP_COMPILE=true")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("%s", output)
		t.Fatal(err)
	}
	t.Log("\n" + string(output))
}
