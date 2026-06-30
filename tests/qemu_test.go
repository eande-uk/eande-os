package tests

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"eande.uk/eande-os/tests/testutil"
)

func hasQEMU() bool {
	_, err := exec.LookPath("qemu-system-x86_64")
	return err == nil
}

func hasKVM() bool {
	_, err := os.Stat("/dev/kvm")
	return err == nil
}

func hasPexpect() bool {
	cmd := exec.Command("python3", "-c", "import pexpect")
	return cmd.Run() == nil
}

func findISO(hubRoot, profile string) (string, error) {
	isoDir := filepath.Join(hubRoot, "iso", "out")
	entries, err := os.ReadDir(isoDir)
	if err != nil {
		return "", err
	}
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), profile) && strings.HasSuffix(e.Name(), ".iso") {
			return filepath.Join(isoDir, e.Name()), nil
		}
	}
	return "", os.ErrNotExist
}

func TestQEMUBoot(t *testing.T) {
	if !hasQEMU() {
		t.Skip("qemu-system-x86_64 not found — skipping QEMU boot test")
	}
	if !hasKVM() {
		t.Skip("/dev/kvm not available — skipping QEMU boot test")
	}
	if !hasPexpect() {
		t.Skip("pexpect not installed — skipping QEMU boot test (pip install pexpect)")
	}

	tc, err := testutil.NewTestContext()
	if err != nil {
		t.Fatal(err)
	}

	profiles := []string{"erch", "e-os-console"}

	for _, profile := range profiles {
		profile := profile
		t.Run(profile, func(t *testing.T) {
			isoPath, err := findISO(tc.RepoRoot, profile)
			if err != nil {
				t.Skipf("No ISO found for %s — build first: make iso/build/%s", profile, profile)
			}

			script := filepath.Join(tc.RepoRoot, "tests", "qemu", "boot_test.py")
			cmd := exec.Command("python3", script, "--iso", isoPath, "--timeout", "120")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				t.Errorf("Boot test failed for %s: %v", profile, err)
			}
		})
	}
}

func TestQEMUVerify(t *testing.T) {
	if !hasQEMU() {
		t.Skip("qemu-system-x86_64 not found — skipping QEMU verify test")
	}
	if !hasKVM() {
		t.Skip("/dev/kvm not available — skipping QEMU verify test")
	}

	tc, err := testutil.NewTestContext()
	if err != nil {
		t.Fatal(err)
	}

	diskPath := "/tmp/eos-verify-test.qcow2"

	// Create test disk
	cmd := exec.Command("qemu-img", "create", "-f", "qcow2", diskPath, "20G")
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to create disk: %v\n%s", err, output)
	}
	defer os.Remove(diskPath)

	// Find an ISO to install from
	profiles := []string{"erch", "e-os-console"}
	var isoPath string
	for _, profile := range profiles {
		p, err := findISO(tc.RepoRoot, profile)
		if err == nil {
			isoPath = p
			break
		}
	}

	if isoPath == "" {
		t.Skip("No ISO found — build first: make iso/build")
	}

	// Boot and install
	installScript := filepath.Join(tc.RepoRoot, "tests", "qemu", "install_test.py")
	installCmd := exec.Command("python3", installScript, "--iso", isoPath, "--disk", diskPath, "--timeout", "600")
	installCmd.Stdout = os.Stdout
	installCmd.Stderr = os.Stderr
	if err := installCmd.Run(); err != nil {
		t.Fatalf("Install test failed: %v", err)
	}

	// Verify installed system
	verifyScript := filepath.Join(tc.RepoRoot, "tests", "qemu", "verify_test.py")
	verifyCmd := exec.Command("python3", verifyScript, "--disk", diskPath, "--profile", "console")
	verifyCmd.Stdout = os.Stdout
	verifyCmd.Stderr = os.Stderr
	if err := verifyCmd.Run(); err != nil {
		t.Errorf("Verify test failed: %v", err)
	}
}

func TestQEMUScriptsExist(t *testing.T) {
	tc, err := testutil.NewTestContext()
	if err != nil {
		t.Fatal(err)
	}

	scripts := []string{
		"boot_test.py",
		"install_test.py",
		"verify_test.py",
		"utils.py",
	}

	qemuDir := filepath.Join(tc.RepoRoot, "tests", "qemu")

	for _, script := range scripts {
		script := script
		testutil.RunVerify(t, "qemu/"+script+" exists", func() error {
			return testutil.FileExists(filepath.Join(qemuDir, script))
		})
	}

	// Verify scripts have correct shebang
	for _, script := range scripts {
		script := script
		testutil.RunVerify(t, "qemu/"+script+" has python3 shebang", func() error {
			return testutil.FileContains(filepath.Join(qemuDir, script), "#!/usr/bin/env python3")
		})
	}
}
