package ingestion

import (
	"fmt"
	"strconv"
	"strings"
)

// parseStringInt converts MLB string numerics and treats placeholders as empty.
func parseStringInt(value string) int {
	value = strings.TrimSpace(value)
	if value == "" || value == ".---" || value == "--" {
		return 0
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}

	return intValue
}

// parseStringFloat converts MLB rate strings like ".285" and "130.1".
func parseStringFloat(value string) float64 {
	value = strings.TrimSpace(value)
	if value == "" || value == ".---" || value == "--" {
		return 0
	}

	floatValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0
	}

	return floatValue
}

// parseOptionalStringFloat preserves nil for fields that are truly absent.
func parseOptionalStringFloat(value string) *float64 {
	value = strings.TrimSpace(value)
	if value == "" || value == ".---" || value == "--" {
		return nil
	}

	floatValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return nil
	}

	return &floatValue
}

// parseBaseballInnings converts MLB baseball notation like "130.1" and "130.2"
// into outs so merged-season math does not treat partial innings as tenths.
func parseBaseballInnings(value string) (int, error) {
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

	partialOuts := 0
	if len(parts) == 2 {
		switch parts[1] {
		case "0", "":
			partialOuts = 0
		case "1":
			partialOuts = 1
		case "2":
			partialOuts = 2
		default:
			return 0, fmt.Errorf("invalid innings partial value %q", value)
		}
	}

	return (whole * 3) + partialOuts, nil
}

func baseballInningsFromOuts(outs int) float64 {
	if outs <= 0 {
		return 0
	}

	return float64(outs/3) + (float64(outs%3) / 10)
}
