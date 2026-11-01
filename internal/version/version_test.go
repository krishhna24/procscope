package version

import "testing"

func TestInfoString(t *testing.T) {
	tests := []struct {
		name string
		info Info
		want string
	}{
		{
			name: "all fields set",
			info: Info{Version: "v0.1.0", Commit: "abc123", Date: "2026-11-01T10:00:00Z", GoVersion: "go1.26.0"},
			want: "procscope v0.1.0\n  commit: abc123\n  built:  2026-11-01T10:00:00Z\n  go:     go1.26.0",
		},
		{
			name: "missing commit and date show unknown",
			info: Info{Version: "dev", GoVersion: "go1.26.0"},
			want: "procscope dev\n  commit: unknown\n  built:  unknown\n  go:     go1.26.0",
		},
		{
			name: "dirty tree is marked",
			info: Info{Version: "dev", Commit: "abc123", Modified: true, GoVersion: "go1.26.0"},
			want: "procscope dev\n  commit: abc123 (dirty)\n  built:  unknown\n  go:     go1.26.0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.info.String(); got != tt.want {
				t.Errorf("String() mismatch\n got:%q\n want:%q", got, tt.want)
			}
		})
	}
}
