package tests

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"eande.uk/os-conf/tests/testutil"
)

func TestOmarchyCommandsDiscovered(t *testing.T) {
	tc, err := testutil.NewTestContext()
	if err != nil {
		t.Fatal(err)
	}

	testutil.RunVerify(t, "omarchy CLI is available", func() error {
		return testutil.CommandExists("omarchy")
	})

	testutil.RunVerify(t, "omarchy commands --json succeeds", func() error {
		_, err := testutil.RunOmarchy("commands", "--json")
		return err
	})

	testutil.RunVerify(t, "os-conf scripts exist in repo with metadata headers", func() error {
		return checkOmarchyScripts(tc.LocalBinDir)
	})

	if _, err := os.Stat(tc.HomeLocalBinPath("omarchy-os-conf-deploy")); err == nil {
		testutil.RunVerify(t, "os-conf commands are discoverable by omarchy", func() error {
			type expectedCmd struct {
				group string
				cmd   string
			}
			expected := []expectedCmd{
				{"os-conf", "deploy"},
				{"os-conf", "status"},
				{"os-conf", "test"},
				{"os-conf", "commit"},
				{"os-conf", "pr"},
				{"os-conf", "layer-zero"},
				{"os-conf", "docs-verify"},
				{"os-conf", "idle"},
				{"os-conf", "idle-resume"},
				{"os-conf", "scaling-cycle"},
				{"os-conf", "brightness-ddc"},
				{"os-conf", "source-ddc"},
				{"os-conf", "rename"},
				{"os-conf", "md-to-html"},
				{"os-conf", "mmd"},
			}
			for _, ec := range expected {
				exists, err := testutil.OmarchyCommandExists(ec.group, ec.cmd)
				if err != nil {
					return err
				}
				if !exists {
					return errFail("omarchy " + ec.group + "/" + ec.cmd + " not found in omarchy command list")
				}
			}
			return nil
		})
	} else {
		t.Log("os-conf commands not yet deployed to ~/.local/bin/ — skipping discovery test (run 'make deploy')")
	}
}

func checkOmarchyScripts(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	var missing []string
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasPrefix(entry.Name(), "omarchy-os-conf-") {
			continue
		}
		path := filepath.Join(dir, entry.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			missing = append(missing, entry.Name()+": read error")
			continue
		}
		if !strings.Contains(string(data), "omarchy:summary=") {
			missing = append(missing, entry.Name()+": missing omarchy:summary")
		}
		if !strings.Contains(string(data), "omarchy:group=os-conf") {
			missing = append(missing, entry.Name()+": missing omarchy:group=os-conf")
		}
	}
	if len(missing) > 0 {
		return errFail(strings.Join(missing, "; "))
	}
	return nil
}
