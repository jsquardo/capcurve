package ingestion

import (
	"strconv"
	"strings"

	"github.com/jsquardo/capcurve/internal/baseball"
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
	return baseball.ParseBaseballInnings(value)
}

func baseballInningsFromOuts(outs int) float64 {
	return baseball.BaseballInningsFromOuts(outs)
}
