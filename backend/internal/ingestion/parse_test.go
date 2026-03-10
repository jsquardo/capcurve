package ingestion

import "testing"

func TestParseStringFloat(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  float64
	}{
		{name: "decimal string", input: ".285", want: 0.285},
		{name: "whole string", input: "130.1", want: 130.1},
		{name: "placeholder", input: ".---", want: 0},
		{name: "invalid", input: "n/a", want: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := parseStringFloat(tt.input); got != tt.want {
				t.Fatalf("parseStringFloat(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestParseOptionalStringFloat(t *testing.T) {
	t.Parallel()

	if got := parseOptionalStringFloat(".---"); got != nil {
		t.Fatalf("parseOptionalStringFloat returned %v, want nil", *got)
	}

	got := parseOptionalStringFloat("1.066")
	if got == nil || *got != 1.066 {
		t.Fatalf("parseOptionalStringFloat returned %v, want 1.066", got)
	}
}
