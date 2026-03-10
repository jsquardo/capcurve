package ingestion

import (
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
