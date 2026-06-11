package tests

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"eande.uk/eande-os/tests/testutil"
)

func TestErchCommandsDiscovered(t *testing.T) {
	tc, err := testutil.NewTestContext()
	if err != nil {
		t.Fatal(err)
	}

	testutil.RunVerify(t, "erch CLI is available", func() error {
		return testutil.CommandExists("erch")
	})

	testutil.RunVerify(t, "erch commands --json succeeds", func() error {
		_, err := testutil.RunErch("commands", "--json")
		return err
	})

	testutil.RunVerify(t, "erch scripts exist in repo with metadata headers", func() error {
		return checkErchScripts(tc.LocalBinDir)
	})

	if _, err := os.Stat(tc.HomeLocalBinPath("erch-deploy")); err == nil {
		testutil.RunVerify(t, "erch commands are discoverable", func() error {
			expected := []string{
				"deploy",
				"status",
				"test",
				"commit",
				"pr",
				"layer-zero",
				"docs-verify",
			}
			for _, ec := range expected {
				exists, err := testutil.ErchCommandExists(ec)
				if err != nil {
					return err
				}
				if !exists {
					return errFail("erch " + ec + " not found in command list")
				}
			}
			return nil
		})
	} else {
		t.Log("erch commands not yet deployed to ~/.local/bin/ — skipping discovery test (run 'make deploy')")
	}
}

func checkErchScripts(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	var missing []string
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasPrefix(entry.Name(), "erch-") {
			continue
		}
		path := filepath.Join(dir, entry.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			missing = append(missing, entry.Name()+": read error")
			continue
		}
		if !strings.Contains(string(data), "erch:summary=") {
			missing = append(missing, entry.Name()+": missing erch:summary")
		}
	}
	if len(missing) > 0 {
		return errFail(strings.Join(missing, "; "))
	}
	return nil
}
