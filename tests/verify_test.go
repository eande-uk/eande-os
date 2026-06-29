package tests

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"eande.uk/eande-os/tests/testutil"
)

func TestBashrc(t *testing.T) {
	tc, err := testutil.NewTestContext()
	if err != nil {
		t.Fatal(err)
	}

	testutil.RunVerify(t, ".bashrc exists at $HOME", func() error {
		return testutil.FileExists(tc.HomePath(".bashrc"))
	})

	testutil.RunVerify(t, ".bashrc sources erch defaults or .bashrc.d/", func() error {
		err1 := testutil.FileContains(tc.HomePath(".bashrc"), "source ~/.local/share/erch/default/bash/rc")
		err2 := testutil.FileContains(tc.HomePath(".bashrc"), `.bashrc.d`)
		if err1 != nil && err2 != nil {
			return errFail(".bashrc missing erch source or .bashrc.d/ override")
		}
		return nil
	})
}

func TestDeployedFiles(t *testing.T) {
	tc, err := testutil.NewTestContext()
	if err != nil {
		t.Fatal(err)
	}

	testutil.RunVerify(t, ".bashrc deployed", func() error {
		return testutil.FileExists(tc.HomePath(".bashrc"))
	})

	testutil.RunVerify(t, "tmux/tmux.conf deployed", func() error {
		return testutil.FileExists(tc.HomeConfigPath("tmux", "tmux.conf"))
	})

	testutil.RunVerify(t, "kitty/kitty.conf deployed", func() error {
		return testutil.FileExists(tc.HomeConfigPath("kitty", "kitty.conf"))
	})

	testutil.RunVerify(t, "hypr/hyprland.conf deployed", func() error {
		return testutil.FileExists(tc.HomeConfigPath("hypr", "hyprland.conf"))
	})

	testutil.RunVerify(t, "hypr/bindings.conf deployed", func() error {
		return testutil.FileExists(tc.HomeConfigPath("hypr", "bindings.conf"))
	})
}

func TestErchHooksAndBranding(t *testing.T) {
	tc, err := testutil.NewTestContext()
	if err != nil {
		t.Fatal(err)
	}

	erchDefault := filepath.Join(tc.RepoRoot, "erch", "default")

	hooks := []string{"theme-set", "font-set", "post-update"}
	for _, hook := range hooks {
		hook := hook
		testutil.RunVerify(t, "erch hook "+hook+" exists", func() error {
			return testutil.FileExists(filepath.Join(erchDefault, "hooks", hook))
		})
	}

	brandingFiles := []string{"about.txt", "screensaver.txt", "ee-logo.png", "ee-mark.png", "water-mark.png"}
	for _, f := range brandingFiles {
		f := f
		testutil.RunVerify(t, "erch branding "+f+" exists", func() error {
			return testutil.FileExists(filepath.Join(erchDefault, "branding", f))
		})
	}
}

func TestForkReconciliation(t *testing.T) {
	tc, err := testutil.NewTestContext()
	if err != nil {
		t.Fatal(err)
	}

	scalingCycle := tc.RepoPath("erch", "bin", "erch-hyprland-monitor-scaling-cycle")
	testutil.RunVerify(t, "erch-hyprland-monitor-scaling-cycle exists in erch", func() error {
		return testutil.FileExists(scalingCycle)
	})
}

func TestTilingMode(t *testing.T) {
	tc, err := testutil.NewTestContext()
	if err != nil {
		t.Fatal(err)
	}

	toggleSrc := tc.RepoPath("erch", "default", "hypr", "toggles", "tiling-mode.conf")

	testutil.RunVerify(t, "Toggle source file exists in erch", func() error {
		return testutil.FileExists(toggleSrc)
	})

	testutil.RunVerify(t, "Source migrates SUPER+J → SUPER+CTRL+J", func() error {
		return testutil.FileContains(toggleSrc, "SUPER CTRL, J")
	})

	testutil.RunVerify(t, "Source migrates SUPER+K → SUPER+SHIFT+K", func() error {
		return testutil.FileContains(toggleSrc, "SUPER SHIFT, K")
	})

	testutil.RunVerify(t, "Source migrates SUPER+L → SUPER+ALT+L", func() error {
		return testutil.FileContains(toggleSrc, "SUPER ALT, L")
	})

	testutil.RunVerify(t, "Source has focus binding", func() error {
		return testutil.FileContains(toggleSrc, "Focus left")
	})

	testutil.RunVerify(t, "Source has swap binding", func() error {
		return testutil.FileContains(toggleSrc, "Swap left")
	})

	testutil.RunVerify(t, "Source has group nav binding", func() error {
		return testutil.FileContains(toggleSrc, "Group left")
	})

	testutil.RunVerify(t, "Source has move workspace binding", func() error {
		return testutil.FileContains(toggleSrc, "Move ws left")
	})

	testutil.RunVerify(t, "Source has group window focus binding", func() error {
		return testutil.FileContains(toggleSrc, "Group prev")
	})

	scalingCycle := tc.RepoPath("erch", "bin", "erch-hyprland-monitor-scaling-cycle")
	testutil.RunVerify(t, "scaling-cycle script exists in erch", func() error {
		return testutil.FileExists(scalingCycle)
	})

	testutil.RunVerify(t, "scaling-cycle cycles all monitors (not just focused)", func() error {
		return testutil.FileContains(scalingCycle, "hyprctl monitors -j")
	})

	testutil.RunVerify(t, "scaling-cycle sends notification", func() error {
		return testutil.FileContains(scalingCycle, "notify-send")
	})
}

func TestPostUpdateHook(t *testing.T) {
	tc, err := testutil.NewTestContext()
	if err != nil {
		t.Fatal(err)
	}

	hookPath := filepath.Join(tc.RepoRoot, "erch", "default", "hooks", "post-update")

	testutil.RunVerify(t, "post-update hook exists in erch", func() error {
		return testutil.FileExists(hookPath)
	})

	testutil.RunVerify(t, "post-update is a shell script with shebang", func() error {
		return testutil.FileContains(hookPath, "#!/bin/bash")
	})

	testutil.RunVerify(t, "post-update hosts stock themes", func() error {
		return testutil.FileContains(hookPath, "stock-themes")
	})
}

func TestErchCLIIntegration(t *testing.T) {
	tc, err := testutil.NewTestContext()
	if err != nil {
		t.Fatal(err)
	}

	if err := testutil.CommandExists("erch"); err != nil {
		t.Log("erch not on PATH — skipping CLI integration tests (run 'make deploy')")
		return
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

func TestProfilesExist(t *testing.T) {
	tc, err := testutil.NewTestContext()
	if err != nil {
		t.Fatal(err)
	}

	profiles := map[string]string{
		"WORK:Office":      "install/packages/work.pkgs",
		"WORK:Dev":         "install/packages/work.pkgs",
		"EDUCATION:School": "install/packages/school.pkgs",
		"GAME":             "install/packages/game.pkgs",
	}

	for name, rel := range profiles {
		name, rel := name, rel
		path := filepath.Join(tc.RepoRoot, "erch", rel)
		testutil.RunVerify(t, name+" package list exists in erch", func() error {
			return testutil.FileExists(path)
		})
	}

	testutil.RunVerify(t, "Each profile has >= 3 package entries", func() error {
		for name, rel := range profiles {
			path := filepath.Join(tc.RepoRoot, "erch", rel)
			data, err := os.ReadFile(path)
			if err != nil {
				return errFail("cannot read " + name + ": " + err.Error())
			}
			lines := 0
			for _, line := range strings.Split(string(data), "\n") {
				line = strings.TrimSpace(line)
				if line != "" && !strings.HasPrefix(line, "#") {
					lines++
				}
			}
			if lines < 3 {
				return errFail(name + " has only " + fmt.Sprintf("%d", lines) + " packages, expected >= 3")
			}
		}
		return nil
	})
}

func TestNoOmarchyReferencesInEandeOS(t *testing.T) {
	tc, err := testutil.NewTestContext()
	if err != nil {
		t.Fatal(err)
	}

	testutil.RunVerify(t, "No omarchy references in eande-os docs/", func() error {
		docsDir := filepath.Join(tc.RepoRoot, "docs")
		return testutil.DirFreeOf(docsDir, "omarchy")
	})

	testutil.RunVerify(t, "AGENTS.md has no omarchy references", func() error {
		return testutil.FileFreeOf(tc.RepoPath("AGENTS.md"), "omarchy")
	})

	testutil.RunVerify(t, "README.md has no omarchy references", func() error {
		return testutil.FileFreeOf(tc.RepoPath("README.md"), "omarchy")
	})

	testutil.RunVerify(t, "Makefile has no omarchy references", func() error {
		return testutil.FileFreeOf(tc.RepoPath("Makefile"), "omarchy")
	})
}

func TestHubStructure(t *testing.T) {
	tc, err := testutil.NewTestContext()
	if err != nil {
		t.Fatal(err)
	}

	requiredDirs := []struct {
		path string
		name string
	}{
		{filepath.Join(tc.RepoRoot, "erch"), "erch/ submodule"},
		{filepath.Join(tc.RepoRoot, "docs"), "docs/ directory"},
		{filepath.Join(tc.RepoRoot, "tests"), "tests/ directory"},
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
		{tc.RepoPath("AGENTS.md"), "AGENTS.md"},
		{tc.RepoPath("README.md"), "README.md"},
		{tc.RepoPath("docs", "ARCHITECTURE.md"), "docs/ARCHITECTURE.md"},
		{tc.RepoPath(".gitmodules"), ".gitmodules"},
	}

	for _, f := range requiredFiles {
		f := f
		testutil.RunVerify(t, f.name+" exists", func() error {
			return testutil.FileExists(f.path)
		})
	}

	removedDirs := []struct {
		path string
		name string
	}{
		{filepath.Join(tc.RepoRoot, "dotfiles"), "dotfiles/ removed (absorbed into erch)"},
		{filepath.Join(tc.RepoRoot, "layer-zero"), "layer-zero/ removed (absorbed into erch)"},
		{filepath.Join(tc.RepoRoot, "scripts"), "scripts/ removed (absorbed into erch)"},
	}

	for _, d := range removedDirs {
		d := d
		testutil.RunVerify(t, d.name, func() error {
			return testutil.DirNotExists(d.path)
		})
	}
}
