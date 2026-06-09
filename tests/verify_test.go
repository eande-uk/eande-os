package tests

import (
	"testing"

	"eande.uk/os-conf/tests/testutil"
)

// Port of verify-phase1.sh — validates repo structure, content, and deployment.

func TestBashrc(t *testing.T) {
	tc, err := testutil.NewTestContext()
	if err != nil {
		t.Fatal(err)
	}

	testutil.RunVerify(t, ".bashrc exists at $HOME", func() error {
		return testutil.FileExists(tc.HomePath(".bashrc"))
	})

	testutil.RunVerify(t, ".bashrc sources Omarchy defaults", func() error {
		return testutil.FileContains(tc.HomePath(".bashrc"),
			"source ~/.local/share/omarchy/default/bash/rc")
	})

	testutil.RunVerify(t, ".bashrc sources .bashrc.d/ for local overrides", func() error {
		return testutil.FileContains(tc.HomePath(".bashrc"),
			`.bashrc.d`)
	})
}

func TestBashrcD(t *testing.T) {
	tc, err := testutil.NewTestContext()
	if err != nil {
		t.Fatal(err)
	}

	testutil.RunVerify(t, "10-env.sh exists", func() error {
		return testutil.FileExists(tc.DotfilesPath("home", ".bashrc.d", "10-env.sh"))
	})

	testutil.RunVerify(t, "10-env.sh sets EDITOR", func() error {
		return testutil.FileContains(
			tc.DotfilesPath("home", ".bashrc.d", "10-env.sh"), "EDITOR=nvim")
	})

	testutil.RunVerify(t, "10-env.sh sets TERMINAL", func() error {
		return testutil.FileContains(
			tc.DotfilesPath("home", ".bashrc.d", "10-env.sh"), "TERMINAL=xdg-terminal-exec")
	})

	testutil.RunVerify(t, "50-aliases.sh exists", func() error {
		return testutil.FileExists(tc.DotfilesPath("home", ".bashrc.d", "50-aliases.sh"))
	})

	testutil.RunVerify(t, "50-aliases.sh has ltt alias", func() error {
		return testutil.FileContains(
			tc.DotfilesPath("home", ".bashrc.d", "50-aliases.sh"), "alias ltt=")
	})

	testutil.RunVerify(t, "50-aliases.sh has ls alias", func() error {
		return testutil.FileContains(
			tc.DotfilesPath("home", ".bashrc.d", "50-aliases.sh"), "alias ls=")
	})

	testutil.RunVerify(t, "50-aliases.sh has grep alias", func() error {
		return testutil.FileContains(
			tc.DotfilesPath("home", ".bashrc.d", "50-aliases.sh"), "alias grep=")
	})

	testutil.RunVerify(t, "60-functions.sh exists", func() error {
		return testutil.FileExists(tc.DotfilesPath("home", ".bashrc.d", "60-functions.sh"))
	})

	testutil.RunVerify(t, "60-functions.sh has mkcd", func() error {
		return testutil.FileContains(
			tc.DotfilesPath("home", ".bashrc.d", "60-functions.sh"), "mkcd()")
	})

	testutil.RunVerify(t, "60-functions.sh has extract", func() error {
		return testutil.FileContains(
			tc.DotfilesPath("home", ".bashrc.d", "60-functions.sh"), "extract()")
	})

	testutil.RunVerify(t, "60-functions.sh has path", func() error {
		return testutil.FileContains(
			tc.DotfilesPath("home", ".bashrc.d", "60-functions.sh"), "path()")
	})
}

func TestTmux(t *testing.T) {
	tc, err := testutil.NewTestContext()
	if err != nil {
		t.Fatal(err)
	}

	tmuxPath := tc.RepoConfigPath("tmux", "tmux.conf")

	testutil.RunVerify(t, "tmux/tmux.conf exists", func() error {
		return testutil.FileExists(tmuxPath)
	})

	testutil.RunVerify(t, "tmux/tmux.conf has >= 90 lines", func() error {
		return testutil.FileMinLines(tmuxPath, 90)
	})

	testutil.RunVerify(t, "tmux/tmux.conf has real config (prefix)", func() error {
		return testutil.FileContains(tmuxPath, "set -g prefix C-Space")
	})
}

func TestKitty(t *testing.T) {
	tc, err := testutil.NewTestContext()
	if err != nil {
		t.Fatal(err)
	}

	kittyPath := tc.RepoConfigPath("kitty", "kitty.conf")

	testutil.RunVerify(t, "kitty/kitty.conf exists", func() error {
		return testutil.FileExists(kittyPath)
	})

	testutil.RunVerify(t, "kitty/kitty.conf has >= 25 lines", func() error {
		return testutil.FileMinLines(kittyPath, 25)
	})

	testutil.RunVerify(t, "kitty/kitty.conf has real config (font_family)", func() error {
		return testutil.FileContains(kittyPath, "font_family")
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

func TestNoEmptyStubs(t *testing.T) {
	tc, err := testutil.NewTestContext()
	if err != nil {
		t.Fatal(err)
	}

	testutil.RunVerify(t, "No 'Full config goes here' stubs in repo", func() error {
		return testutil.StubFree(tc.DotfilesPath("home"))
	})
}

func TestCustomBranding(t *testing.T) {
	tc, err := testutil.NewTestContext()
	if err != nil {
		t.Fatal(err)
	}

	testutil.RunVerify(t, "src-pngs/ee-logo.png exists", func() error {
		return testutil.FileExists(tc.DotfilesPath("home", ".config", "custom-branding", "src-pngs", "ee-logo.png"))
	})

	testutil.RunVerify(t, "src-pngs/ee-mark.png exists", func() error {
		return testutil.FileExists(tc.DotfilesPath("home", ".config", "custom-branding", "src-pngs", "ee-mark.png"))
	})

	testutil.RunVerify(t, "src-pngs/water-mark.png exists", func() error {
		return testutil.FileExists(tc.DotfilesPath("home", ".config", "custom-branding", "src-pngs", "water-mark.png"))
	})

	testutil.RunVerify(t, "custom-branding/about.txt >= 20 lines", func() error {
		return testutil.FileMinLines(
			tc.DotfilesPath("home", ".config", "custom-branding", "about.txt"), 20)
	})

	testutil.RunVerify(t, "custom-branding/screensaver.txt >= 8 lines", func() error {
		return testutil.FileMinLines(
			tc.DotfilesPath("home", ".config", "custom-branding", "screensaver.txt"), 8)
	})

	testutil.RunVerify(t, "omarchy/branding/about.txt exists", func() error {
		return testutil.FileExists(tc.DotfilesPath("home", ".config", "omarchy", "branding", "about.txt"))
	})

	testutil.RunVerify(t, "omarchy/branding/screensaver.txt exists", func() error {
		return testutil.FileExists(tc.DotfilesPath("home", ".config", "omarchy", "branding", "screensaver.txt"))
	})

	testutil.RunVerify(t, "omarchy/branding/about.txt is a real file (not symlink)", func() error {
		return testutil.IsNotSymlink(tc.DotfilesPath("home", ".config", "omarchy", "branding", "about.txt"))
	})

	testutil.RunVerify(t, "omarchy/branding/screensaver.txt is a real file (not symlink)", func() error {
		return testutil.IsNotSymlink(tc.DotfilesPath("home", ".config", "omarchy", "branding", "screensaver.txt"))
	})

	testutil.RunVerify(t, "custom-branding and omarchy/branding about.txt in sync", func() error {
		return testutil.FilesIdentical(
			tc.DotfilesPath("home", ".config", "custom-branding", "about.txt"),
			tc.DotfilesPath("home", ".config", "omarchy", "branding", "about.txt"))
	})

	testutil.RunVerify(t, "custom-branding and omarchy/branding screensaver.txt in sync", func() error {
		return testutil.FilesIdentical(
			tc.DotfilesPath("home", ".config", "custom-branding", "screensaver.txt"),
			tc.DotfilesPath("home", ".config", "omarchy", "branding", "screensaver.txt"))
	})
}

func TestStowEngine(t *testing.T) {
	tc, err := testutil.NewTestContext()
	if err != nil {
		t.Fatal(err)
	}

	deployPath := tc.DotfilesPath("scripts", "deploy.sh")

	testutil.RunVerify(t, "deploy.sh uses stow --no-folding", func() error {
		return testutil.FileContains(deployPath, "stow --no-folding")
	})

	testutil.RunVerify(t, "deploy.sh supports --adopt", func() error {
		return testutil.FileContains(deployPath, "ADOPT=true")
	})

	testutil.RunVerify(t, "deploy.sh supports --dry-run", func() error {
		return testutil.FileContains(deployPath, "stow -n -v")
	})

	testutil.RunVerify(t, "deploy.sh guards against master branch", func() error {
		return testutil.FileContains(deployPath, "Create a user branch first")
	})

	testutil.RunVerify(t, "stow is installed", func() error {
		return testutil.CommandExists("stow")
	})

	testutil.RunVerify(t, "gum is installed", func() error {
		return testutil.CommandExists("gum")
	})

	if tc.IsWSL() {
		t.Log("Skipping deploy --dry-run on WSL")
		return
	}

	testutil.RunVerify(t, "deploy --dry-run runs without error", func() error {
		_, err := testutil.RunMake(tc, "deploy/dry-run")
		return err
	})
}

func TestRemovedScripts(t *testing.T) {
	tc, err := testutil.NewTestContext()
	if err != nil {
		t.Fatal(err)
	}

	testutil.RunVerify(t, "pull.sh removed", func() error {
		return testutil.FileNotExists(tc.DotfilesPath("scripts", "pull.sh"))
	})

	testutil.RunVerify(t, "clean-orphans.sh removed", func() error {
		return testutil.FileNotExists(tc.DotfilesPath("scripts", "clean-orphans.sh"))
	})
}

func TestForkReconciliation(t *testing.T) {
	tc, err := testutil.NewTestContext()
	if err != nil {
		t.Fatal(err)
	}

	hypridlePath := tc.RepoConfigPath("hypr", "hypridle.conf")

	testutil.RunVerify(t, "hypridle.conf has fork's DPMS customization", func() error {
		return testutil.FileContains(hypridlePath, "hyprctl dispatch dpms off")
	})

	testutil.RunVerify(t, "hypridle.conf has fork's wake customization", func() error {
		return testutil.FileContains(hypridlePath, "omarchy-system-wake")
	})

	systemIdlePath := tc.LocalBinPath("omarchy-os-conf-idle")
	testutil.RunVerify(t, "omarchy-os-conf-idle has fork's lock-session", func() error {
		return testutil.FileContains(systemIdlePath, "loginctl lock-session")
	})

	fastfetchPath := tc.RepoConfigPath("fastfetch", "config.jsonc")
	testutil.RunVerify(t, "fastfetch uses fork's cleaner OS Age command", func() error {
		return testutil.FileContains(fastfetchPath, "$(stat -c %W /)")
	})
}

func TestOmarchyMenuExtension(t *testing.T) {
	tc, err := testutil.NewTestContext()
	if err != nil {
		t.Fatal(err)
	}

	testutil.RunVerify(t, "omarchy-menu extension exists in repo", func() error {
		return testutil.FileExists(
			tc.DotfilesPath("home", ".config", "omarchy", "extensions", "menu.sh"))
	})

	testutil.RunVerify(t, "extension adds Tiling Mode to toggle menu", func() error {
		return testutil.FileContains(
			tc.DotfilesPath("home", ".config", "omarchy", "extensions", "menu.sh"),
			"Tiling Mode")
	})

	testutil.RunVerify(t, "extension redefines show_toggle_menu", func() error {
		return testutil.FileContains(
			tc.DotfilesPath("home", ".config", "omarchy", "extensions", "menu.sh"),
			"show_toggle_menu()")
	})
}

func TestTilingMode(t *testing.T) {
	tc, err := testutil.NewTestContext()
	if err != nil {
		t.Fatal(err)
	}

	tilingConf := tc.RepoConfigPath("hypr", "tiling.conf")

	testutil.RunVerify(t, "tiling.conf exists", func() error {
		return testutil.FileExists(tilingConf)
	})

	testutil.RunVerify(t, "tiling.conf uses native omarchy toggle", func() error {
		return testutil.FileContains(tilingConf, "omarchy hyprland toggle tiling-mode")
	})

	testutil.RunVerify(t, "legacy tiling-mode-toggle removed", func() error {
		return testutil.FileNotExists(tc.HomeLocalBinPath("tiling-mode-toggle"))
	})

	toggleSrc := tc.DotfilesPath("omarchy-default", "toggles", "tiling-mode.conf")

	testutil.RunVerify(t, "Toggle source file exists", func() error {
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

	scalingCycle := tc.LocalBinPath("omarchy-os-conf-scaling-cycle")
	testutil.RunVerify(t, "scaling-cycle script exists in repo", func() error {
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

	hookPath := tc.DotfilesPath("home", ".config", "omarchy", "hooks", "post-update")

	testutil.RunVerify(t, "post-update hook exists", func() error {
		return testutil.FileExists(hookPath)
	})

	testutil.RunVerify(t, "post-update restores toggle template", func() error {
		return testutil.FileContains(hookPath, "tiling-mode.conf")
	})

	testutil.RunVerify(t, "post-update restores os-conf wrappers", func() error {
		return testutil.FileContains(hookPath, "omarchy-os-conf-")
	})

	testutil.RunVerify(t, "post-update restores scaling-cycle override", func() error {
		return testutil.FileContains(hookPath, "omarchy-hyprland-monitor-scaling-cycle")
	})

	testutil.RunVerify(t, "post-update hides stock themes", func() error {
		return testutil.FileContains(hookPath, "stock-themes")
	})
}
