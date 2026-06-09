package testutil

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

type Verify func() error

func RunVerify(t *testing.T, name string, fn Verify) {
	t.Helper()
	t.Run(name, func(t *testing.T) {
		t.Helper()
		if err := fn(); err != nil {
			t.Errorf("%s: %v", name, err)
		}
	})
}

type TestContext struct {
	RepoRoot      string
	DotfilesRoot  string
	HomeDir       string
	BloatDir      string
	LayerZeroDir  string
	TestsDir      string
	ScriptsDir    string
	LocalBinDir   string
	Logger        *log.Logger
}

func NewTestContext() (*TestContext, error) {
	tc := &TestContext{
		Logger: log.New(os.Stderr, "", 0),
		HomeDir: os.Getenv("HOME"),
	}

	if v := os.Getenv("DOTFILES_ROOT"); v != "" {
		tc.DotfilesRoot = v
	} else {
		_, file, _, ok := runtime.Caller(0)
		if !ok {
			return nil, fmt.Errorf("cannot determine dotfiles root")
		}
		tc.DotfilesRoot = filepath.Dir(filepath.Dir(filepath.Dir(file)))
	}

	tc.RepoRoot = filepath.Dir(tc.DotfilesRoot)
	tc.TestsDir = filepath.Join(tc.DotfilesRoot, "tests")
	tc.ScriptsDir = filepath.Join(tc.DotfilesRoot, "scripts")
	tc.LocalBinDir = filepath.Join(tc.DotfilesRoot, "home", ".local", "bin")
	tc.BloatDir = filepath.Join(tc.DotfilesRoot, "layer-zero", "bloat")
	tc.LayerZeroDir = filepath.Join(tc.DotfilesRoot, "layer-zero")

	return tc, nil
}

func (tc *TestContext) HomePath(parts ...string) string {
	return filepath.Join(append([]string{tc.HomeDir}, parts...)...)
}

func (tc *TestContext) RepoPath(parts ...string) string {
	return filepath.Join(append([]string{tc.RepoRoot}, parts...)...)
}

func (tc *TestContext) DotfilesPath(parts ...string) string {
	return filepath.Join(append([]string{tc.DotfilesRoot}, parts...)...)
}

func (tc *TestContext) HomeConfigPath(app string, file string) string {
	return tc.HomePath(".config", app, file)
}

func (tc *TestContext) RepoConfigPath(app string, file string) string {
	return tc.DotfilesPath("home", ".config", app, file)
}

func (tc *TestContext) LocalBinPath(name string) string {
	return filepath.Join(tc.LocalBinDir, name)
}

func (tc *TestContext) HomeLocalBinPath(name string) string {
	return tc.HomePath(".local", "bin", name)
}

func (tc *TestContext) IsWSL() bool {
	data, err := os.ReadFile("/proc/version")
	return err == nil && (strings.Contains(string(data), "microsoft") || strings.Contains(string(data), "WSL"))
}
