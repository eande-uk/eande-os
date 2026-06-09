package tests

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"eande.uk/os-conf/tests/testutil"
)

func TestDocsBehaviourCoherency(t *testing.T) {
	tc, err := testutil.NewTestContext()
	if err != nil {
		t.Fatal(err)
	}

	testutil.RunVerify(t, "AGENTS.md references make targets that exist in Makefile", func() error {
		makefile, err := os.ReadFile(tc.DotfilesPath("Makefile"))
		if err != nil {
			return err
		}
		agents, err := os.ReadFile(tc.RepoPath("AGENTS.md"))
		if err != nil {
			return err
		}

		agentsContent := string(agents)
		makefileContent := string(makefile)

		makeTargets := extractMakeTargets(makefileContent)
		var missing []string

		agentMakeRefs := extractMarkdownCodeRefs(agentsContent)
		for _, ref := range agentMakeRefs {
			if strings.HasPrefix(ref, "make ") {
				target := strings.TrimPrefix(ref, "make ")
				target = strings.Fields(target)[0]
				if !contains(makeTargets, target) {
					missing = append(missing, target)
				}
			}
		}

		if len(missing) > 0 {
			return errFail("AGENTS.md references make targets not in Makefile: " +
				strings.Join(missing, ", "))
		}
		return nil
	})

	testutil.RunVerify(t, "plan.md layer table matches actual directory structure", func() error {
		checks := []struct {
			path     string
			purpose  string
		}{
			{tc.LayerZeroDir, "Layer 0 directory"},
			{tc.DotfilesPath("omarchy-default"), "Layer 2 themes directory"},
			{tc.DotfilesPath("home", ".config", "hypr"), "Layer 3 config directory"},
			{tc.LocalBinDir, "Layer 4 scripts directory"},
		}
		for _, c := range checks {
			info, err := os.Stat(c.path)
			if err != nil {
				return errFail(c.purpose + " missing: " + c.path)
			}
			if !info.IsDir() {
				return errFail(c.purpose + " not a directory: " + c.path)
			}
		}
		return nil
	})

	testutil.RunVerify(t, "dotfiles/README.md file descriptions match actual files", func() error {
		dotfilesReadme, err := os.ReadFile(tc.DotfilesPath("README.md"))
		if err != nil {
			return err
		}

		content := string(dotfilesReadme)
		if !strings.Contains(content, "home/") {
			return errFail("dotfiles/README.md missing home/ directory reference")
		}

		checks := []struct {
			file     string
			pattern  string
		}{
			{"home/.config/hypr/hyprland.conf", "hyprland.conf"},
			{"home/.config/hypr/bindings.conf", "bindings.conf"},
			{"home/.config/hypr/tiling.conf", "tiling.conf"},
			{"omarchy-default/toggles/tiling-mode.conf", "tiling-mode.conf"},
			{"layer-zero/layer-zero.sh", "layer-zero.sh"},
			{"scripts/deploy.sh", "deploy.sh"},
		}
		for _, c := range checks {
			fullPath := filepath.Join(tc.DotfilesRoot, filepath.FromSlash(c.file))
			info, err := os.Stat(fullPath)
			if err != nil {
				return errFail("file referenced in docs not found: " + c.file)
			}
			_ = info
			if !strings.Contains(content, c.pattern) {
				return errFail("dotfiles/README.md missing reference to: " + c.pattern)
			}
		}
		return nil
	})

	testutil.RunVerify(t, "Makefile .PHONY targets have matching recipes", func() error {
		makefile, err := os.ReadFile(tc.DotfilesPath("Makefile"))
		if err != nil {
			return err
		}

		makefileContent := string(makefile)
		phonyLine := extractPhonyLine(makefileContent)

		for _, target := range strings.Fields(phonyLine) {
			target = strings.TrimRight(target, "\\")
			if target == ".PHONY" || target == "" {
				continue
			}
			if !strings.Contains(makefileContent, target+":") {
				return errFail("Makefile .PHONY target missing recipe: " + target)
			}
		}
		return nil
	})
}

func TestBranchingModelConsistency(t *testing.T) {
	tc, err := testutil.NewTestContext()
	if err != nil {
		t.Fatal(err)
	}

	testutil.RunVerify(t, "current branch is a user branch (not master)", func() error {
		b, err := gitBranch()
		if err != nil {
			return err
		}
		if b == "main" || b == "master" {
			return errFail("On master branch — create a user branch: make init")
		}
		return nil
	})

	docFiles := []string{
		tc.RepoPath("AGENTS.md"),
		tc.RepoPath("README.md"),
		tc.DotfilesPath("README.md"),
	}

	for _, doc := range docFiles {
		doc := doc
		name := filepath.Base(filepath.Dir(doc)) + "/" + filepath.Base(doc)
		if doc == tc.RepoPath("README.md") {
			name = "root/README.md"
		}
		testutil.RunVerify(t, name+" mentions user branch model", func() error {
			data, err := os.ReadFile(doc)
			if err != nil {
				return err
			}
			content := string(data)
			if !strings.Contains(content, "user/") &&
				!strings.Contains(content, "user branch") &&
				!strings.Contains(content, "user/<name>") {
				return errFail(name + " missing user branch documentation")
			}
			return nil
		})
	}
}

func extractMakeTargets(content string) []string {
	var targets []string
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if strings.Contains(line, ":") && !strings.HasPrefix(line, "#") &&
			!strings.HasPrefix(line, "\t") && !strings.HasPrefix(line, " ") &&
			!strings.HasPrefix(line, ".") && !strings.HasPrefix(line, "export") {
			target := strings.Split(line, ":")[0]
			targets = append(targets, strings.TrimSpace(target))
		}
	}
	return targets
}

func extractMarkdownCodeRefs(content string) []string {
	var refs []string
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "`") && strings.Contains(line, "`") {
			inside := strings.Split(line, "`")[1]
			refs = append(refs, inside)
		}
	}
	return refs
}

func extractPhonyLine(content string) string {
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, ".PHONY:") {
			return strings.TrimPrefix(line, ".PHONY:")
		}
	}
	return ""
}

func contains(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}
