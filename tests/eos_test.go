package tests

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"eande.uk/eande-os/tests/testutil"
)

func TestEOSHubStructure(t *testing.T) {
	tc, err := testutil.NewTestContext()
	if err != nil {
		t.Fatal(err)
	}

	testutil.RunVerify(t, "E-OS submodule exists", func() error {
		return testutil.DirExists(filepath.Join(tc.RepoRoot, "E-OS"))
	})

	requiredFiles := []struct {
		path string
		name string
	}{
		{filepath.Join(tc.RepoRoot, "E-OS", "install.sh"), "E-OS/install.sh"},
		{filepath.Join(tc.RepoRoot, "E-OS", "boot.sh"), "E-OS/boot.sh"},
		{filepath.Join(tc.RepoRoot, "E-OS", "version"), "E-OS/version"},
		{filepath.Join(tc.RepoRoot, "E-OS", "README.md"), "E-OS/README.md"},
		{filepath.Join(tc.RepoRoot, "E-OS", "AGENTS.md"), "E-OS/AGENTS.md"},
		{filepath.Join(tc.RepoRoot, "E-OS", "bin", "e-os"), "E-OS/bin/e-os"},
	}

	for _, f := range requiredFiles {
		f := f
		testutil.RunVerify(t, f.name+" exists", func() error {
			return testutil.FileExists(f.path)
		})
	}

	requiredDirs := []struct {
		path string
		name string
	}{
		{filepath.Join(tc.RepoRoot, "E-OS", "install"), "E-OS/install/"},
		{filepath.Join(tc.RepoRoot, "E-OS", "default"), "E-OS/default/"},
		{filepath.Join(tc.RepoRoot, "E-OS", "config"), "E-OS/config/"},
		{filepath.Join(tc.RepoRoot, "E-OS", "themes"), "E-OS/themes/"},
		{filepath.Join(tc.RepoRoot, "E-OS", "bin"), "E-OS/bin/"},
		{filepath.Join(tc.RepoRoot, "E-OS", "docs"), "E-OS/docs/"},
	}

	for _, d := range requiredDirs {
		d := d
		testutil.RunVerify(t, d.name+" exists", func() error {
			return testutil.DirExists(d.path)
		})
	}
}

func TestEOSInstallPipeline(t *testing.T) {
	tc, err := testutil.NewTestContext()
	if err != nil {
		t.Fatal(err)
	}

	eosInstall := filepath.Join(tc.RepoRoot, "E-OS", "install")

	stages := []string{
		"helpers/all.sh",
		"preflight/all.sh",
		"packaging/all.sh",
		"config/all.sh",
		"login/all.sh",
		"post-install/all.sh",
	}

	for _, stage := range stages {
		stage := stage
		testutil.RunVerify(t, "install stage "+stage+" exists", func() error {
			return testutil.FileExists(filepath.Join(eosInstall, stage))
		})
	}
}

func TestEOSPackages(t *testing.T) {
	tc, err := testutil.NewTestContext()
	if err != nil {
		t.Fatal(err)
	}

	pkgDir := filepath.Join(tc.RepoRoot, "E-OS", "install", "packages")

	packages := []string{
		"common.pkgs",
		"console.pkgs",
		"school.pkgs",
		"uni.pkgs",
		"org.pkgs",
	}

	for _, pkg := range packages {
		pkg := pkg
		testutil.RunVerify(t, "package list "+pkg+" exists", func() error {
			return testutil.FileExists(filepath.Join(pkgDir, pkg))
		})
	}

	testutil.RunVerify(t, "Each profile has >= 3 package entries", func() error {
		for _, pkg := range packages {
			path := filepath.Join(pkgDir, pkg)
			data, err := os.ReadFile(path)
			if err != nil {
				return errFail("cannot read " + pkg + ": " + err.Error())
			}
			lines := 0
			for _, line := range strings.Split(string(data), "\n") {
				line = strings.TrimSpace(line)
				if line != "" && !strings.HasPrefix(line, "#") {
					lines++
				}
			}
			if lines < 3 {
				return errFail(pkg + " has only " + fmt.Sprintf("%d", lines) + " packages, expected >= 3")
			}
		}
		return nil
	})
}

func TestEOSThemes(t *testing.T) {
	tc, err := testutil.NewTestContext()
	if err != nil {
		t.Fatal(err)
	}

	themesDir := filepath.Join(tc.RepoRoot, "E-OS", "themes")

	themes := []string{
		"hackerman",
		"tokyo-night",
		"catppuccin-latte",
		"everforest",
		"nord",
		"tokyo-night-storm",
		"catppuccin-mocha",
		"dracula",
	}

	for _, theme := range themes {
		theme := theme
		testutil.RunVerify(t, "theme "+theme+" has colors.toml", func() error {
			return testutil.FileExists(filepath.Join(themesDir, theme, "colors.toml"))
		})
	}
}

func TestEOSDefaultConfigs(t *testing.T) {
	tc, err := testutil.NewTestContext()
	if err != nil {
		t.Fatal(err)
	}

	defaultDir := filepath.Join(tc.RepoRoot, "E-OS", "default")

	configs := []string{
		"alacritty/alacritty.toml",
		"hypr/hyprland.conf",
		"waybar/config.jsonc",
		"waybar/style.css",
		"mako/config",
		"sddm/sddm.conf",
		"bash/bashrc",
	}

	for _, cfg := range configs {
		cfg := cfg
		testutil.RunVerify(t, "default config "+cfg+" exists", func() error {
			return testutil.FileExists(filepath.Join(defaultDir, cfg))
		})
	}
}

func TestEOSThemedTemplates(t *testing.T) {
	tc, err := testutil.NewTestContext()
	if err != nil {
		t.Fatal(err)
	}

	themedDir := filepath.Join(tc.RepoRoot, "E-OS", "default", "themed")

	templates := []string{
		"alacritty.toml.tpl",
		"hyprland.lua.tpl",
		"waybar.css.tpl",
		"mako.ini.tpl",
	}

	for _, tpl := range templates {
		tpl := tpl
		testutil.RunVerify(t, "themed template "+tpl+" exists", func() error {
			return testutil.FileExists(filepath.Join(themedDir, tpl))
		})
	}
}

func TestEOSBinCommands(t *testing.T) {
	tc, err := testutil.NewTestContext()
	if err != nil {
		t.Fatal(err)
	}

	binDir := filepath.Join(tc.RepoRoot, "E-OS", "bin")

	commands := []string{
		"e-os",
		"e-os-version",
		"e-os-theme-list",
		"e-os-theme-set",
		"e-os-gpu-detect",
		"e-os-books",
		"e-os-refresh-config",
		"e-os-restart",
		"e-os-toggle",
		"e-os-update",
		"e-os-debug",
		"e-os-show-logo",
		"e-os-pkg-add",
		"e-os-pkg-present",
		"e-os-pkg-missing",
		"e-os-pkg-remove",
		"e-os-install",
	}

	for _, cmd := range commands {
		cmd := cmd
		testutil.RunVerify(t, "command "+cmd+" exists", func() error {
			return testutil.FileExists(filepath.Join(binDir, cmd))
		})
	}
}

func TestEOSDocs(t *testing.T) {
	tc, err := testutil.NewTestContext()
	if err != nil {
		t.Fatal(err)
	}

	docsDir := filepath.Join(tc.RepoRoot, "E-OS", "docs")

	docs := []string{
		"ARCHITECTURE.md",
		"PROFILES.md",
		"QUICKSTART.md",
	}

	for _, doc := range docs {
		doc := doc
		testutil.RunVerify(t, "doc "+doc+" exists", func() error {
			return testutil.FileExists(filepath.Join(docsDir, doc))
		})
	}
}

func TestEOSInstallShHasShebang(t *testing.T) {
	tc, err := testutil.NewTestContext()
	if err != nil {
		t.Fatal(err)
	}

	testutil.RunVerify(t, "install.sh has bash shebang", func() error {
		return testutil.FileContains(filepath.Join(tc.RepoRoot, "E-OS", "install.sh"), "#!/bin/bash")
	})

	testutil.RunVerify(t, "boot.sh has bash shebang", func() error {
		return testutil.FileContains(filepath.Join(tc.RepoRoot, "E-OS", "boot.sh"), "#!/bin/bash")
	})
}

func TestEOSBinMetadata(t *testing.T) {
	tc, err := testutil.NewTestContext()
	if err != nil {
		t.Fatal(err)
	}

	binDir := filepath.Join(tc.RepoRoot, "E-OS", "bin")

	entries, err := os.ReadDir(binDir)
	if err != nil {
		t.Fatal(err)
	}

	var missing []string
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasPrefix(entry.Name(), "e-os-") {
			continue
		}
		path := filepath.Join(binDir, entry.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			missing = append(missing, entry.Name()+": read error")
			continue
		}
		if !strings.Contains(string(data), "e-os:summary=") {
			missing = append(missing, entry.Name()+": missing e-os:summary")
		}
	}

	if len(missing) > 0 {
		t.Errorf("Missing e-os:summary metadata: %s", strings.Join(missing, "; "))
	}
}
