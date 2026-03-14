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

func TestParseBaseballInnings(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    int
		wantErr bool
	}{
		{name: "whole innings", input: "130.0", want: 390},
		{name: "one out", input: "130.1", want: 391},
		{name: "two outs", input: "130.2", want: 392},
		{name: "placeholder", input: ".---", want: 0},
		{name: "invalid partial", input: "130.3", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := parseBaseballInnings(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("parseBaseballInnings(%q) error = nil, want error", tt.input)
				}
				return
			}
			if err != nil {
				t.Fatalf("parseBaseballInnings(%q) error = %v", tt.input, err)
			}
			if got != tt.want {
				t.Fatalf("parseBaseballInnings(%q) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}
