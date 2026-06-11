package tests

import (
	"os"
	"testing"

	"eande.uk/eande-os/tests/testutil"
)

func TestEnvironmentVariables(t *testing.T) {
	tc, err := testutil.NewTestContext()
	if err != nil {
		t.Fatal(err)
	}

	expected := map[string]string{
		"EDITOR":   "nvim",
		"TERMINAL": "xdg-terminal-exec",
	}

	for name, want := range expected {
		name := name
		want := want
		testutil.RunVerify(t, "$"+name+" is set", func() error {
			got := os.Getenv(name)
			if got == "" {
				return testutil.FileContains(
					tc.DotfilesPath("home", ".bashrc.d", "10-env.sh"),
					name+"="+want)
			}
			return nil
		})
	}

	testutil.RunVerify(t, "HOME is set", func() error {
		if os.Getenv("HOME") == "" {
			return errFail("HOME not set")
		}
		return nil
	})

	testutil.RunVerify(t, "erch is on PATH", func() error {
		return testutil.CommandExists("erch")
	})

	testutil.RunVerify(t, "make is on PATH", func() error {
		return testutil.CommandExists("make")
	})

	testutil.RunVerify(t, "stow is on PATH", func() error {
		return testutil.CommandExists("stow")
	})
}

type errFail string

func (e errFail) Error() string { return string(e) }
