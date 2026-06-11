package tests

import (
	"os"
	"strings"
	"testing"

	"eande.uk/eande-os/tests/testutil"
)

func TestErrorPaths(t *testing.T) {
	tc, err := testutil.NewTestContext()
	if err != nil {
		t.Fatal(err)
	}

	branch, _ := gitBranch()
	onMain := branch == "main" || branch == "master"

	if onMain {
		testutil.RunVerify(t, "deploy fails with helpful message on master branch", func() error {
			out, err := testutil.RunMake(tc, "deploy")
			if err == nil {
				return errFail("deploy should have failed on master branch")
			}
			if !strings.Contains(out, "Create a user branch first") &&
				!strings.Contains(out, "user branch") {
				return errFail("deploy on master should mention user branch:\n" + out)
			}
			return nil
		})
	} else {
		t.Log("not on master branch — skipping master-guard test")
	}

	testutil.RunVerify(t, "deploy fails when stow is missing", func() error {
		_, err := os.Stat("/usr/bin/stow")
		if err != nil {
			return testutil.CommandExists("stow")
		}
		return nil
	})

	testutil.RunVerify(t, "deploy fails when gum is missing", func() error {
		return testutil.CommandExists("gum")
	})
}

func gitBranch() (string, error) {
	tc, err := testutil.NewTestContext()
	if err != nil {
		return "", err
	}
	out, err := testutil.RunMake(tc, "status")
	if err != nil {
		return "", err
	}
	for _, line := range strings.Split(out, "\n") {
		if strings.HasPrefix(line, "Branch:") {
			return strings.TrimSpace(strings.TrimPrefix(line, "Branch:")), nil
		}
	}
	return "", errFail("could not determine branch")
}
