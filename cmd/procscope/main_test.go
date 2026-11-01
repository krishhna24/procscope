package main

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		wantCode   int
		wantStdout string //substring that must appear in stdout
		wantStderr string //substring that must appear in stderr
	}{
		{name: "version", args: []string{"version"}, wantCode: exitOK, wantStdout: "procscope "},
		{name: "version short", args: []string{"version", "--short"}, wantCode: exitOK, wantStdout: "dev"},
		{name: "help", args: []string{"help"}, wantCode: exitOK, wantStdout: "Usage:"},
		{name: "no args", args: nil, wantCode: exitUsage, wantStderr: "Usage:"},
		{name: "unknown command", args: []string{"bogus"}, wantCode: exitUsage, wantStderr: `unknown command "bogus"`},
		{name: "unknown flag", args: []string{"version", "--nope"}, wantCode: exitUsage, wantStderr: "unknown flag"},
		// pflag keeps parsing flags after a positional argument ("interspersed"),
		// so --short is understood and the complaint is about "x", not the flag.
		{name: "interspersed flags", args: []string{"version", "x", "--short"}, wantCode: exitUsage, wantStderr: `unexpected argument "x"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			code := run(context.Background(), tt.args, &stdout, &stderr)
			if code != tt.wantCode {
				t.Errorf("exit code %d , want %d", code, tt.wantCode)
			}
			checkOutput(t, "stdout", stdout.String(), tt.wantStdout)
			checkOutput(t, "stderr", stderr.String(), tt.wantStderr)
		})
	}
}

func checkOutput(t *testing.T, stream, got, want string) {
	t.Helper()
	if want == "" {
		if got != "" {
			t.Errorf("%s should be empty , got %q", stream, got)
		}
		return
	}
	if !strings.Contains(got, want) {
		t.Errorf("%s = %q, want it to contain %q", stream, got, want)
	}
}
