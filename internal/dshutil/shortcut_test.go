package dshutil

import "testing"

func TestPSEscape(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{`C:\Program Files`, `C:\Program Files`},
		{`O'Reilly`, `O''Reilly`},
		{"", ""},
	}
	for _, tt := range tests {
		if got := psEscape(tt.in); got != tt.want {
			t.Fatalf("psEscape(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}
