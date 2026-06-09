package testutil

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

func FileExists(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file not found: %s", path)
		}
		return fmt.Errorf("stat %s: %w", path, err)
	}
	if info.IsDir() {
		return fmt.Errorf("path is a directory, not a file: %s", path)
	}
	return nil
}

func DirExists(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("directory not found: %s", path)
		}
		return fmt.Errorf("stat %s: %w", path, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("path is a file, not a directory: %s", path)
	}
	return nil
}

func FileNotExists(path string) error {
	_, err := os.Stat(path)
	if err == nil {
		return fmt.Errorf("file still exists (should be removed): %s", path)
	}
	if os.IsNotExist(err) {
		return nil
	}
	return fmt.Errorf("stat %s: %w", path, err)
}

func IsSymlink(path string) error {
	info, err := os.Lstat(path)
	if err != nil {
		return fmt.Errorf("lstat %s: %w", path, err)
	}
	if info.Mode()&os.ModeSymlink == 0 {
		return fmt.Errorf("not a symlink: %s", path)
	}
	return nil
}

func IsNotSymlink(path string) error {
	info, err := os.Lstat(path)
	if err != nil {
		return fmt.Errorf("lstat %s: %w", path, err)
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("is a symlink (should be real file): %s", path)
	}
	return nil
}

func FileContains(path, pattern string) error {
	if err := FileExists(path); err != nil {
		return fmt.Errorf("cannot check content: %w", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}
	if !strings.Contains(string(data), pattern) {
		return fmt.Errorf("pattern %q not found in %s", pattern, path)
	}
	return nil
}

func FileMatches(path, pattern string) error {
	if err := FileExists(path); err != nil {
		return fmt.Errorf("cannot check content: %w", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return fmt.Errorf("invalid regex %q: %w", pattern, err)
	}
	if !re.Match(data) {
		return fmt.Errorf("pattern %q not matched in %s", pattern, path)
	}
	return nil
}

func FileLines(path string) (int, error) {
	if err := FileExists(path); err != nil {
		return 0, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, fmt.Errorf("read %s: %w", path, err)
	}
	if len(data) == 0 {
		return 0, nil
	}
	return len(strings.Split(strings.TrimRight(string(data), "\n"), "\n")), nil
}

func FileMinLines(path string, min int) error {
	lines, err := FileLines(path)
	if err != nil {
		return err
	}
	if lines < min {
		return fmt.Errorf("%s has %d lines, expected >= %d", path, lines, min)
	}
	return nil
}

func CommandExists(name string) error {
	if _, err := exec.LookPath(name); err != nil {
		return fmt.Errorf("command not found: %s", name)
	}
	return nil
}

func FilesIdentical(a, b string) error {
	dataA, err := os.ReadFile(a)
	if err != nil {
		return fmt.Errorf("read %s: %w", a, err)
	}
	dataB, err := os.ReadFile(b)
	if err != nil {
		return fmt.Errorf("read %s: %w", b, err)
	}
	if !bytesEqual(dataA, dataB) {
		return fmt.Errorf("files differ:\n  %s\n  %s", a, b)
	}
	return nil
}

func bytesEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func StubFree(root string) error {
	var found []string
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		ext := filepath.Ext(path)
		if ext != ".conf" && ext != ".sh" && ext != ".toml" && filepath.Base(path) != "config" {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		if strings.Contains(string(data), "Full config goes here") {
			found = append(found, path)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("walk %s: %w", root, err)
	}
	if len(found) > 0 {
		return fmt.Errorf("stub files found: %s", strings.Join(found, ", "))
	}
	return nil
}
