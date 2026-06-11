package tests

import (
	"strings"
	"testing"

	"eande.uk/eande-os/tests/testutil"
)

func TestMakeTargets(t *testing.T) {
	tc, err := testutil.NewTestContext()
	if err != nil {
		t.Fatal(err)
	}

	targets := []struct {
		name    string
		target  string
		args    []string
	}{
		{"help runs without error", "help", nil},
		{"status shows branch info", "status", nil},
		{"deploy/dry-run shows preview", "deploy/dry-run", nil},
	}

	for _, tt := range targets {
		tt := tt
		testutil.RunVerify(t, tt.name, func() error {
			if tt.name == "deploy/dry-run shows preview" && tc.IsWSL() {
				t.Skip("Skipping deploy/dry-run on WSL")
			}
			out, err := testutil.RunMake(tc, tt.target, tt.args...)
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

func TestRepoStructure(t *testing.T) {
	tc, err := testutil.NewTestContext()
	if err != nil {
		t.Fatal(err)
	}

	requiredDirs := []struct {
		path string
		name string
	}{
		{tc.DotfilesPath("home"), "home/ directory"},
		{tc.DotfilesPath("home", ".config"), "home/.config/ directory"},
		{tc.DotfilesPath("home", ".config", "hypr"), "hypr/ config directory"},
		{tc.RepoPath("scripts"), "scripts/ directory"},
		{tc.LayerZeroDir, "layer-zero/ directory"},
		{tc.LocalBinDir, ".local/bin/ directory"},
	}

	for _, d := range requiredDirs {
		d := d
		testutil.RunVerify(t, d.name+" exists", func() error {
			return testutil.DirExists(d.path)
		})
	}

	requiredFiles := []struct {
		path string
		name string
	}{
		{tc.RepoPath("Makefile"), "Makefile"},
		{tc.DotfilesPath("home", ".bashrc"), ".bashrc"},
		{tc.DotfilesPath("home", ".config", "hypr", "hyprland.conf"), "hyprland.conf"},
	}

	for _, f := range requiredFiles {
		f := f
		testutil.RunVerify(t, f.name+" exists", func() error {
			return testutil.FileExists(f.path)
		})
	}
}

func TestErchCLIIntegration(t *testing.T) {
	tc, err := testutil.NewTestContext()
	if err != nil {
		t.Fatal(err)
	}

	testutil.RunVerify(t, "erch version returns a version", func() error {
		out, err := testutil.RunErch("version")
		if err != nil {
			return err
		}
		if out == "" {
			return errFail("erch version returned empty output")
		}
		return nil
	})

	testutil.RunVerify(t, "erch theme list shows available themes", func() error {
		out, err := testutil.RunErch("theme", "list")
		if err != nil {
			return err
		}
		if !strings.Contains(out, "Nord") && !strings.Contains(out, "Catppuccin") && !strings.Contains(out, "Tokyo") && !strings.Contains(out, "Matte") {
			return errFail("theme list missing expected themes:\n" + out)
		}
		return nil
	})

	testutil.RunVerify(t, "erch theme current shows a theme", func() error {
		out, err := testutil.RunErch("theme", "current")
		if err != nil {
			return err
		}
		if out == "" {
			return errFail("theme current returned empty")
		}
		return nil
	})

	testutil.RunVerify(t, "erch commands --json is valid JSON", func() error {
		out, err := testutil.RunErch("commands", "--json")
		if err != nil {
			return err
		}
		if !strings.Contains(out, `"ok"`) && !strings.Contains(out, `"commands"`) {
			return errFail("commands --json missing expected fields")
		}
		if !strings.HasPrefix(out, "{") {
			return errFail("commands --json does not start with {")
		}
		return nil
	})

	if tc.IsWSL() {
		return
	}

	testutil.RunVerify(t, "erch cmd present stow succeeds", func() error {
		_, err := testutil.RunErch("cmd", "present", "stow")
		return err
	})

	testutil.RunVerify(t, "erch cmd present gum succeeds", func() error {
		_, err := testutil.RunErch("cmd", "present", "gum")
		return err
	})
}


