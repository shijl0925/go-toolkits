package systemx

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSplitCommandLine(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []string
		wantErr string
	}{
		{
			name:  "multiple spaces",
			input: `printf   "%s"   "hello world"`,
			want:  []string{"printf", "%s", "hello world"},
		},
		{
			name:  "escaped quotes",
			input: `printf "%s" "hello \"quoted\" world"`,
			want:  []string{"printf", "%s", `hello "quoted" world`},
		},
		{
			name:  "mixed quotes",
			input: `printf "%s" 'hello world'`,
			want:  []string{"printf", "%s", "hello world"},
		},
		{
			name:    "unterminated quote",
			input:   `printf "hello`,
			wantErr: "unterminated quote",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := splitCommandLine(tt.input)
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("expected error containing %q, got %v", tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(got) != len(tt.want) {
				t.Fatalf("unexpected arg count: got %d want %d (%v)", len(got), len(tt.want), got)
			}

			for i := range got {
				if got[i] != tt.want[i] {
					t.Fatalf("arg %d mismatch: got %q want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestResolveExecutable(t *testing.T) {
	executablePath, err := os.Executable()
	if err != nil {
		t.Fatalf("os.Executable error: %v", err)
	}

	tests := []struct {
		name      string
		input     string
		assertion func(t *testing.T, resolved string, err error)
	}{
		{
			name:  "absolute path",
			input: executablePath,
			assertion: func(t *testing.T, resolved string, err error) {
				t.Helper()
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if resolved != filepath.Clean(executablePath) {
					t.Fatalf("unexpected resolved path: got %q want %q", resolved, filepath.Clean(executablePath))
				}
			},
		},
		{
			name:  "null byte",
			input: "go\x00test",
			assertion: func(t *testing.T, resolved string, err error) {
				t.Helper()
				if err == nil || !strings.Contains(err.Error(), "null byte") {
					t.Fatalf("expected null byte error, got path=%q err=%v", resolved, err)
				}
			},
		},
		{
			name:  "relative path",
			input: "." + string(os.PathSeparator) + "go",
			assertion: func(t *testing.T, resolved string, err error) {
				t.Helper()
				if err == nil || !strings.Contains(err.Error(), "relative command paths are not allowed") {
					t.Fatalf("expected relative path error, got path=%q err=%v", resolved, err)
				}
			},
		},
		{
			name:  "lookup path",
			input: "go",
			assertion: func(t *testing.T, resolved string, err error) {
				t.Helper()
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if !filepath.IsAbs(resolved) {
					t.Fatalf("expected absolute resolved path, got %q", resolved)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resolved, err := resolveExecutable(tt.input)
			tt.assertion(t, resolved, err)
		})
	}
}
