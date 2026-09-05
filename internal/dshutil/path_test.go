package dshutil

import "testing"

func TestMergePath(t *testing.T) {
	tests := []struct {
		name string
		in   []string
		want string
	}{
		{"empty", []string{}, ""},
		{"single", []string{`C:\a;C:\b`}, `C:\a;C:\b`},
		{"dedup case-insensitive", []string{`C:\A;C:\a`}, `C:\A`},
		{"trim trailing backslash", []string{`C:\a\`}, `C:\a`},
		{"drive root keeps slash", []string{`C:\`}, `C:\`},
		{"trim whitespace", []string{` C:\a ;`}, `C:\a`},
		{"skip empty parts", []string{`C:\a;;C:\b`}, `C:\a;C:\b`},
		{"multiple sources keep order", []string{`C:\a`, `C:\b;C:\c`, `C:\a`}, `C:\a;C:\b;C:\c`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MergePath(tt.in...)
			if got != tt.want {
				t.Fatalf("MergePath(%v) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}
