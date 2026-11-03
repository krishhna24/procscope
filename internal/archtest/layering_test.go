package archtest

import (
	"errors"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

const module = "github.com/krishhna24/procscope"

//allowed lists, for every package in this module, the other procscope
//packages it may import directly. It encodes ROADMAP.md section 3.3:
// every arrow points toward internal/sample, so no import cycle can form.
// Standard-library and third party imports are not restricted here.

// Package that don't exist yet are listed so the rules are decided before
// the code it written. Changing a rule is a design decision: do it in its  own commit and say why
var allowed = map[string][]string{
	"cmd/procscope": {
		"internal/version", "internal/collect", "internal/engine",
		"internal/view", "internal/output", "internal/record", "internal/tui",
	},
	"internal/archtest": {},
	"internal/version":  {},
	"internal/sample":   {},
	"internal/procfs":   {},
	"internal/collect":  {"internal/procfs", "internal/sample"},
	"internal/engine":   {"internal/sample"},
	"internal/view":     {"internal/engine", "internal/sample"},
	"internal/output":   {"internal/view", "internal/engine", "internal/sample"},
	"internal/record":   {"internal/sample"},
	"internal/tui":      {"internal/view", "internal/engine", "internal/sample"},
}

func TestLayering(t *testing.T) {
	// One line per package: its import path, then its direct imports
	// Imports excludes _test.go files, so tests may import more freely.

	cmd := exec.Command("go", "list", "-f", "{{.ImportPath}}{{range .Imports}} {{.}}{{end}}", module+"/...")

	out, err := cmd.Output()
	watchSources(t, "../..")

	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			t.Fatalf("go list failed: %v\n%s", err, exitErr.Stderr)
		}
		t.Fatalf("go list failed: %v", err)
	}

	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		fields := strings.Fields(line)
		pkg, imports := rel(fields[0]), fields[1:]

		rules, known := allowed[pkg]
		if !known {
			t.Errorf("%s has no layering rule: add it to allowed in layering_test.go", pkg)
			continue
		}
		for _, imp := range imports {
			if !strings.HasPrefix(imp, module+"/") {
				continue // standard library or third-party
			}
			if !slices.Contains(rules, rel(imp)) {
				t.Errorf("forbidden import %s -> %s (allowed: %v)", pkg, rel(imp), rules)
			}
		}
	}
}

//rel turns "github.com/krishhna24/procscope/internal/sample" into "internal/sample"

func rel(importPath string) string {
	return strings.TrimPrefix(importPath, module+"/")
}

// watchSources stats every .go file in the module. go list runs in a
// subprocess, so go test can't see which files the result depends on and
// would replay a cached PASS after a forbidden import is added. Files and
// directories the test itself opens or stats become part of the cache key.
func watchSources(t *testing.T, root string) {
	t.Helper()
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if path != root && (d.Name() == "scratch" || d.Name() == "learning" || strings.HasPrefix(d.Name(), ".")) {
				return filepath.SkipDir
			}
			return nil
		}
		if strings.HasSuffix(path, ".go") {
			_, err = os.Stat(path)
		}
		return err
	})
	if err != nil {
		t.Fatalf("walking module sources: %v", err)
	}
}
