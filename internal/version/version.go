package version

import (
	"fmt"
	"runtime"
	"runtime/debug"
)

var (
	Version = "dev"
	Commit  = ""
	Date    = ""
)

type Info struct {
	Version   string
	Commit    string
	Date      string
	Modified  bool // True if built from a tree with uncommited changes
	GoVersion string
}

func Get() Info {
	info := Info{
		Version:   Version,
		Commit:    Commit,
		Date:      Date,
		Modified:  false,
		GoVersion: runtime.Version(),
	}
	bi, ok := debug.ReadBuildInfo()

	if !ok {
		return info
	}

	if info.Version == "dev" && bi.Main.Version != "" && bi.Main.Version != "(devel)" {
		info.Version = bi.Main.Version //Set by: go install @v0.1.0
	}

	for _, s := range bi.Settings {
		switch s.Key {
		case "vcs.revision":
			if info.Commit == "" {
				info.Commit = s.Value
			}
		case "vcs.time":
			if info.Date == "" {
				info.Date = s.Value
			}
		case "vcs.modified":
			info.Modified = s.Value == "true"
		}
	}
	return info
}

// String formats the information for `procscope version`.

func (i Info) String() string {
	commit := orUnknown(i.Commit)

	if i.Modified {
		commit += " (dirty)"
	}

	return fmt.Sprintf("procscope %s\n  commit: %s\n  built:  %s\n  go:     %s", i.Version, commit, orUnknown(i.Date), i.GoVersion)
}

func orUnknown(s string) string {
	if s == "" {
		return "unknown"
	}
	return s
}
