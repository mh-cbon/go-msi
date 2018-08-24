package stringexec

import (
	"testing"
)

func Test1(t *testing.T) {
	expected := "ok\nok\n"
	linux, err := Command("echo ok && echo ok")
	if err != nil {
		t.Errorf("should err=nil, got err=%q\n", err)
	}

	out, err := linux.Output()
	if err != nil {
		t.Errorf("should err=nil, got err=%q\n", err)
	}

	if expected != string(out) {
		t.Errorf("should output=\n%q\n, got output=\n%q\n", expected, out)
	}
}
