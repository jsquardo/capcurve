package baseball

import (
	"testing"

	"github.com/stretchr/testify/require"
)

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

			got, err := ParseBaseballInnings(tt.input)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestOutsFromBaseballNotation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   float64
		want    int
		wantErr bool
	}{
		{name: "whole innings", input: 29.0, want: 87},
		{name: "one out", input: 29.1, want: 88},
		{name: "two outs", input: 29.2, want: 89},
		{name: "zero innings", input: 0, want: 0},
		{name: "invalid partial", input: 10.4, wantErr: true},
		{name: "too many decimals", input: 10.25, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := OutsFromBaseballNotation(tt.input)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestTrueInningsFromBaseballNotation(t *testing.T) {
	t.Parallel()

	innings, err := TrueInningsFromBaseballNotation(29.2)
	require.NoError(t, err)
	require.InDelta(t, 29.6666667, innings, 0.000001)

	innings, err = TrueInningsFromBaseballNotation(10.1)
	require.NoError(t, err)
	require.InDelta(t, 10.3333333, innings, 0.000001)

	innings, err = TrueInningsFromBaseballNotation(0)
	require.NoError(t, err)
	require.Equal(t, 0.0, innings)

	innings, err = TrueInningsFromBaseballNotation(10.4)
	require.Error(t, err)
	require.Equal(t, 0.0, innings)
}
