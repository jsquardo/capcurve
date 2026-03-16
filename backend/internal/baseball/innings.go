package baseball

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

const inningsPrecisionEpsilon = 0.000001

// ParseBaseballInnings converts MLB baseball notation like "130.1" and "130.2"
// into outs so downstream math never treats partial innings as decimal tenths.
func ParseBaseballInnings(value string) (int, error) {
	value = strings.TrimSpace(value)
	if value == "" || value == ".---" || value == "--" {
		return 0, nil
	}

	parts := strings.Split(value, ".")
	if len(parts) > 2 {
		return 0, fmt.Errorf("invalid innings %q", value)
	}

	whole, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, fmt.Errorf("invalid innings whole value %q: %w", value, err)
	}

	return outsFromWholeAndPartial(value, whole, partialDigits(parts))
}

// OutsFromBaseballNotation validates a persisted MLB-style innings value and
// converts it back into outs for workload and rate math.
func OutsFromBaseballNotation(innings float64) (int, error) {
	if innings <= 0 {
		return 0, nil
	}

	scaled := innings * 10
	rounded := math.Round(scaled)
	if math.Abs(scaled-rounded) > inningsPrecisionEpsilon {
		return 0, fmt.Errorf("invalid innings %.3f: must use at most one decimal place", innings)
	}

	tenths := int(rounded)
	whole := tenths / 10
	partial := strconv.Itoa(tenths % 10)

	return outsFromWholeAndPartial(strconv.FormatFloat(innings, 'f', -1, 64), whole, partial)
}

func TrueInningsFromBaseballNotation(innings float64) (float64, error) {
	outs, err := OutsFromBaseballNotation(innings)
	if err != nil || outs <= 0 {
		return 0, err
	}

	return float64(outs) / 3, nil
}

func BaseballInningsFromOuts(outs int) float64 {
	if outs <= 0 {
		return 0
	}

	return float64(outs/3) + (float64(outs%3) / 10)
}

func partialDigits(parts []string) string {
	if len(parts) != 2 {
		return ""
	}
	return parts[1]
}

func outsFromWholeAndPartial(original string, whole int, partial string) (int, error) {
	partialOuts := 0
	switch partial {
	case "", "0":
		partialOuts = 0
	case "1":
		partialOuts = 1
	case "2":
		partialOuts = 2
	default:
		return 0, fmt.Errorf("invalid innings partial value %q", original)
	}

	return (whole * 3) + partialOuts, nil
}
