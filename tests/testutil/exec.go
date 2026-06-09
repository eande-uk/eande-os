package testutil

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func RunMake(tc *TestContext, target string, args ...string) (string, error) {
	makefile := filepath.Join(tc.DotfilesRoot, "Makefile")
	if _, err := os.Stat(makefile); err != nil {
		return "", fmt.Errorf("Makefile not found at %s: %w", makefile, err)
	}

	cmdArgs := append([]string{"-C", tc.DotfilesRoot, target}, args...)
	return runCommand("make", cmdArgs...)
}

func RunOmarchy(args ...string) (string, error) {
	return runCommand("omarchy", args...)
}

func RunOmarchyJSON(args ...string) (string, error) {
	args = append(args, "--json")
	return runCommand("omarchy-commands", args...)
}

func runCommand(name string, args ...string) (string, error) {
	var stdout, stderr bytes.Buffer
	cmd := exec.Command(name, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	out := strings.TrimSpace(stdout.String())
	if err != nil {
		return out, fmt.Errorf("%s %s failed: %w\nstderr: %s", name, strings.Join(args, " "), err, strings.TrimSpace(stderr.String()))
	}
	return out, nil
}

func OmarchyCommandExists(group, cmd string) (bool, error) {
	out, err := RunOmarchy("commands", "--json")
	if err != nil {
		return false, fmt.Errorf("omarchy commands --json failed: %w", err)
	}
	groupPrefix := fmt.Sprintf(`"group": "%s"`, group)
	cmdRoute := fmt.Sprintf(`"route": "omarchy %s %s"`, group, cmd)
	return strings.Contains(out, cmdRoute) && strings.Contains(out, groupPrefix), nil
}
