package tests

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"eande.uk/eande-os/tests/testutil"
)

func TestDocsBehaviourCoherency(t *testing.T) {
	tc, err := testutil.NewTestContext()
	if err != nil {
		t.Fatal(err)
	}

	testutil.RunVerify(t, "AGENTS.md references make targets that exist in Makefile", func() error {
		makefile, err := os.ReadFile(tc.RepoPath("Makefile"))
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

	testutil.RunVerify(t, "Makefile .PHONY targets have matching recipes", func() error {
		makefile, err := os.ReadFile(tc.RepoPath("Makefile"))
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

	testutil.RunVerify(t, "AGENTS.md layer table matches erch layer system", func() error {
		agents, err := os.ReadFile(tc.RepoPath("AGENTS.md"))
		if err != nil {
			return err
		}
		content := string(agents)

		layers := []string{"L0", "L1", "L2", "L3", "L4"}
		for _, layer := range layers {
			if !strings.Contains(content, layer) {
				return errFail("AGENTS.md missing layer reference: " + layer)
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
	}

	for _, doc := range docFiles {
		doc := doc
		name := filepath.Base(doc)
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

func TestSubmoduleReferences(t *testing.T) {
	tc, err := testutil.NewTestContext()
	if err != nil {
		t.Fatal(err)
	}

	testutil.RunVerify(t, ".gitmodules references erch submodule", func() error {
		return testutil.FileContains(tc.RepoPath(".gitmodules"), "erch")
	})

	testutil.RunVerify(t, "AGENTS.md mentions erch as active distro", func() error {
		agents, err := os.ReadFile(tc.RepoPath("AGENTS.md"))
		if err != nil {
			return err
		}
		content := string(agents)
		if !strings.Contains(content, "erch") {
			return errFail("AGENTS.md missing erch")
		}
		if !strings.Contains(content, "Active") {
			return errFail("AGENTS.md missing Active status for erch")
		}
		return nil
	})

	testutil.RunVerify(t, "AGENTS.md mentions E-OS and E-OS-AI as planned", func() error {
		agents, err := os.ReadFile(tc.RepoPath("AGENTS.md"))
		if err != nil {
			return err
		}
		content := string(agents)
		if !strings.Contains(content, "E-OS") {
			return errFail("AGENTS.md missing E-OS")
		}
		if !strings.Contains(content, "E-OS-AI") {
			return errFail("AGENTS.md missing E-OS-AI")
		}
		if !strings.Contains(content, "Planned") {
			return errFail("AGENTS.md missing Planned status")
		}
		return nil
	})
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
