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

func TestISOStructure(t *testing.T) {
	tc, err := testutil.NewTestContext()
	if err != nil {
		t.Fatal(err)
	}

	isoProfiles := []string{"erch", "e-os-console", "e-os-school", "e-os-uni", "e-os-org"}

	for _, profile := range isoProfiles {
		profile := profile
		isoDir := filepath.Join(tc.RepoRoot, "iso", profile)

		testutil.RunVerify(t, "iso/"+profile+"/profiledef.sh exists", func() error {
			return testutil.FileExists(filepath.Join(isoDir, "profiledef.sh"))
		})

		testutil.RunVerify(t, "iso/"+profile+"/packages.x86_64 exists", func() error {
			return testutil.FileExists(filepath.Join(isoDir, "packages.x86_64"))
		})

		testutil.RunVerify(t, "iso/"+profile+"/pacman.conf exists", func() error {
			return testutil.FileExists(filepath.Join(isoDir, "pacman.conf"))
		})

		testutil.RunVerify(t, "iso/"+profile+"/airootfs/root/installer.sh exists", func() error {
			return testutil.FileExists(filepath.Join(isoDir, "airootfs", "root", "installer.sh"))
		})

		testutil.RunVerify(t, "iso/"+profile+"/airootfs/etc/systemd/system/getty@tty1.service.d/autologin.conf exists", func() error {
			return testutil.FileExists(filepath.Join(isoDir, "airootfs", "etc", "systemd", "system", "getty@tty1.service.d", "autologin.conf"))
		})
	}
}

func TestISOInstallerScripts(t *testing.T) {
	tc, err := testutil.NewTestContext()
	if err != nil {
		t.Fatal(err)
	}

	isoProfiles := []string{"erch", "e-os-console", "e-os-school", "e-os-uni", "e-os-org"}

	for _, profile := range isoProfiles {
		profile := profile
		installerPath := filepath.Join(tc.RepoRoot, "iso", profile, "airootfs", "root", "installer.sh")

		testutil.RunVerify(t, "iso/"+profile+" installer has bash shebang", func() error {
			return testutil.FileContains(installerPath, "#!/bin/bash")
		})

		testutil.RunVerify(t, "iso/"+profile+" installer has set -eEo pipefail", func() error {
			return testutil.FileContains(installerPath, "set -eEo pipefail")
		})

		testutil.RunVerify(t, "iso/"+profile+" installer has main function", func() error {
			return testutil.FileContains(installerPath, "main()")
		})

		testutil.RunVerify(t, "iso/"+profile+" installer partitions disk", func() error {
			return testutil.FileContains(installerPath, "sgdisk")
		})

		testutil.RunVerify(t, "iso/"+profile+" installer uses pacstrap", func() error {
			return testutil.FileContains(installerPath, "pacstrap")
		})

		testutil.RunVerify(t, "iso/"+profile+" installer sets up Limine", func() error {
			return testutil.FileContains(installerPath, "limine")
		})

		testutil.RunVerify(t, "iso/"+profile+" installer creates user", func() error {
			return testutil.FileContains(installerPath, "useradd")
		})
	}
}

func TestISOProfiledef(t *testing.T) {
	tc, err := testutil.NewTestContext()
	if err != nil {
		t.Fatal(err)
	}

	isoProfiles := map[string]string{
		"erch":          "erch",
		"e-os-console":  "e-os-console",
		"e-os-school":   "e-os-school",
		"e-os-uni":      "e-os-uni",
		"e-os-org":      "e-os-org",
	}

	for profile, isoName := range isoProfiles {
		profile, isoName := profile, isoName
		profiledefPath := filepath.Join(tc.RepoRoot, "iso", profile, "profiledef.sh")

		testutil.RunVerify(t, "iso/"+profile+" profiledef.sh has iso_name", func() error {
			return testutil.FileContains(profiledefPath, "iso_name=\""+isoName+"\"")
		})

		testutil.RunVerify(t, "iso/"+profile+" profiledef.sh has iso_publisher", func() error {
			return testutil.FileContains(profiledefPath, "iso_publisher=")
		})

		testutil.RunVerify(t, "iso/"+profile+" profiledef.sh has bootmodes", func() error {
			return testutil.FileContains(profiledefPath, "bootmodes=")
		})

		testutil.RunVerify(t, "iso/"+profile+" profiledef.sh has file_permissions for installer", func() error {
			return testutil.FileContains(profiledefPath, "/root/installer.sh")
		})
	}
}

func TestISOPackages(t *testing.T) {
	tc, err := testutil.NewTestContext()
	if err != nil {
		t.Fatal(err)
	}

	isoProfiles := []string{"erch", "e-os-console", "e-os-school", "e-os-uni", "e-os-org"}

	for _, profile := range isoProfiles {
		profile := profile
		packagesPath := filepath.Join(tc.RepoRoot, "iso", profile, "packages.x86_64")

		testutil.RunVerify(t, "iso/"+profile+" packages has base", func() error {
			return testutil.FileContains(packagesPath, "base")
		})

		testutil.RunVerify(t, "iso/"+profile+" packages has linux", func() error {
			return testutil.FileContains(packagesPath, "linux")
		})

		testutil.RunVerify(t, "iso/"+profile+" packages has mkinitcpio-archiso", func() error {
			return testutil.FileContains(packagesPath, "mkinitcpio-archiso")
		})

		testutil.RunVerify(t, "iso/"+profile+" packages has gum", func() error {
			return testutil.FileContains(packagesPath, "gum")
		})

		testutil.RunVerify(t, "iso/"+profile+" packages has btrfs-progs", func() error {
			return testutil.FileContains(packagesPath, "btrfs-progs")
		})

		testutil.RunVerify(t, "iso/"+profile+" packages has limine", func() error {
			return testutil.FileContains(packagesPath, "limine")
		})
	}
}

func TestISODocs(t *testing.T) {
	tc, err := testutil.NewTestContext()
	if err != nil {
		t.Fatal(err)
	}

	testutil.RunVerify(t, "docs/ISO.md exists", func() error {
		return testutil.FileExists(filepath.Join(tc.RepoRoot, "docs", "ISO.md"))
	})

	testutil.RunVerify(t, "docs/ISO.md has build instructions", func() error {
		return testutil.FileContains(filepath.Join(tc.RepoRoot, "docs", "ISO.md"), "make iso/build")
	})
}

func TestMakefileISOTargets(t *testing.T) {
	tc, err := testutil.NewTestContext()
	if err != nil {
		t.Fatal(err)
	}

	makefilePath := filepath.Join(tc.RepoRoot, "Makefile")

	targets := []string{
		"iso/build:",
		"iso/build/erch:",
		"iso/build/e-os:",
		"iso/build/e-os-console:",
		"iso/build/e-os-school:",
		"iso/build/e-os-uni:",
		"iso/build/e-os-org:",
		"iso/clean:",
		"iso/test:",
	}

	for _, target := range targets {
		target := target
		testutil.RunVerify(t, "Makefile has target "+target, func() error {
			return testutil.FileContains(makefilePath, target)
		})
	}
}

func TestEOSBooksStructure(t *testing.T) {
	tc, err := testutil.NewTestContext()
	if err != nil {
		t.Fatal(err)
	}

	booksDir := filepath.Join(tc.RepoRoot, "E-OS", "books")

	testutil.RunVerify(t, "books/book.toml exists", func() error {
		return testutil.FileExists(filepath.Join(booksDir, "book.toml"))
	})

	testutil.RunVerify(t, "books/.gitignore exists", func() error {
		return testutil.FileExists(filepath.Join(booksDir, ".gitignore"))
	})

	testutil.RunVerify(t, "books/src/SUMMARY.md exists", func() error {
		return testutil.FileExists(filepath.Join(booksDir, "src", "SUMMARY.md"))
	})

	testutil.RunVerify(t, "books/src/README.md exists", func() error {
		return testutil.FileExists(filepath.Join(booksDir, "src", "README.md"))
	})

	testutil.RunVerify(t, "book.toml has title", func() error {
		return testutil.FileContains(filepath.Join(booksDir, "book.toml"), "title = \"E&E OS Knowledge Base\"")
	})

	testutil.RunVerify(t, "book.toml has language", func() error {
		return testutil.FileContains(filepath.Join(booksDir, "book.toml"), "language = \"en\"")
	})
}

func TestEOSBookChapters(t *testing.T) {
	tc, err := testutil.NewTestContext()
	if err != nil {
		t.Fatal(err)
	}

	booksSrc := filepath.Join(tc.RepoRoot, "E-OS", "books", "src")

	chapters := []string{
		"getting-started/installation.md",
		"getting-started/profiles.md",
		"getting-started/post-install.md",
		"system/package-management.md",
		"system/systemd.md",
		"system/networking.md",
		"system/storage.md",
		"hyprland/getting-started.md",
		"hyprland/configuration.md",
		"hyprland/keybinds.md",
		"hyprland/plugins.md",
		"shell/bash-fundamentals.md",
		"shell/scripting.md",
		"security/hardening-checklist.md",
		"security/kernel-hardening.md",
		"security/systemd-sandboxing.md",
		"theming/color-schemes.md",
		"theming/waybar.md",
		"troubleshooting/common-issues.md",
		"troubleshooting/recovery.md",
	}

	for _, ch := range chapters {
		ch := ch
		testutil.RunVerify(t, "chapter "+ch+" exists", func() error {
			return testutil.FileExists(filepath.Join(booksSrc, ch))
		})
	}

	testutil.RunVerify(t, "all chapters have content (>100 bytes)", func() error {
		for _, ch := range chapters {
			path := filepath.Join(booksSrc, ch)
			info, err := os.Stat(path)
			if err != nil {
				return errFail("cannot stat " + ch + ": " + err.Error())
			}
			if info.Size() < 100 {
				return errFail(ch + " is only " + fmt.Sprintf("%d", info.Size()) + " bytes")
			}
		}
		return nil
	})
}

func TestEOSBooksCLI(t *testing.T) {
	tc, err := testutil.NewTestContext()
	if err != nil {
		t.Fatal(err)
	}

	booksBin := filepath.Join(tc.RepoRoot, "E-OS", "bin", "e-os-books")

	testutil.RunVerify(t, "e-os-books has bash shebang", func() error {
		return testutil.FileContains(booksBin, "#!/bin/bash")
	})

	testutil.RunVerify(t, "e-os-books has set -eEo pipefail", func() error {
		return testutil.FileContains(booksBin, "set -eEo pipefail")
	})

	testutil.RunVerify(t, "e-os-books has build command", func() error {
		return testutil.FileContains(booksBin, "cmd_build")
	})

	testutil.RunVerify(t, "e-os-books has serve command", func() error {
		return testutil.FileContains(booksBin, "cmd_serve")
	})

	testutil.RunVerify(t, "e-os-books has install command", func() error {
		return testutil.FileContains(booksBin, "cmd_install")
	})

	testutil.RunVerify(t, "e-os-books has clean command", func() error {
		return testutil.FileContains(booksBin, "cmd_clean")
	})

	testutil.RunVerify(t, "e-os-books has e-os:summary metadata", func() error {
		return testutil.FileContains(booksBin, "e-os:summary=")
	})
}

func TestMakefileDocsTargets(t *testing.T) {
	tc, err := testutil.NewTestContext()
	if err != nil {
		t.Fatal(err)
	}

	makefilePath := filepath.Join(tc.RepoRoot, "Makefile")

	targets := []string{
		"docs/build:",
		"docs/serve:",
		"docs/clean:",
		"test/qemu:",
		"test/e2e:",
	}

	for _, target := range targets {
		target := target
		testutil.RunVerify(t, "Makefile has target "+target, func() error {
			return testutil.FileContains(makefilePath, target)
		})
	}
}
