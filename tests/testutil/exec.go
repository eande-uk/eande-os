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
	makefile := filepath.Join(tc.RepoRoot, "Makefile")
	if _, err := os.Stat(makefile); err != nil {
		return "", fmt.Errorf("Makefile not found at %s: %w", makefile, err)
	}

	cmdArgs := append([]string{"-C", tc.RepoRoot, target}, args...)
	return runCommand("make", cmdArgs...)
}

func RunErch(args ...string) (string, error) {
	return runCommand("erch", args...)
}

func RunErchJSON(args ...string) (string, error) {
	args = append(args, "--json")
	return runCommand("erch-commands", args...)
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

func ErchCommandExists(route string) (bool, error) {
	out, err := RunErch("commands", "--json")
	if err != nil {
		return false, fmt.Errorf("erch commands --json failed: %w", err)
	}
	cmdRoute := fmt.Sprintf(`"route": "erch %s"`, route)
	return strings.Contains(out, cmdRoute), nil
}
