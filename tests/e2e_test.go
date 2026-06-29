package tests

import (
	"path/filepath"
	"testing"

	"eande.uk/eande-os/tests/testutil"
)

func TestMakeTargets(t *testing.T) {
	tc, err := testutil.NewTestContext()
	if err != nil {
		t.Fatal(err)
	}

	targets := []struct {
		name   string
		target string
	}{
		{"help runs without error", "help"},
		{"status shows branch info", "status"},
	}

	for _, tt := range targets {
		tt := tt
		testutil.RunVerify(t, tt.name, func() error {
			out, err := testutil.RunMake(tc, tt.target)
			if err != nil {
				return err
			}
			if len(out) == 0 && tt.name != "help runs without error" {
				return errFail(tt.target + " produced no output")
			}
			return nil
		})
	}
}

func TestErchSubmoduleInit(t *testing.T) {
	tc, err := testutil.NewTestContext()
	if err != nil {
		t.Fatal(err)
	}

	testutil.RunVerify(t, "erch submodule directory exists", func() error {
		return testutil.DirExists(filepath.Join(tc.RepoRoot, "erch"))
	})

	testutil.RunVerify(t, "erch/.git exists (submodule initialized)", func() error {
		return testutil.FileExists(filepath.Join(tc.RepoRoot, "erch", ".git"))
	})

	testutil.RunVerify(t, "erch/install.sh exists", func() error {
		return testutil.FileExists(filepath.Join(tc.RepoRoot, "erch", "install.sh"))
	})

	testutil.RunVerify(t, "erch/AGENTS.md exists", func() error {
		return testutil.FileExists(filepath.Join(tc.RepoRoot, "erch", "AGENTS.md"))
	})
}
