package models

import (
	"sync"
	"testing"

	"gorm.io/gorm/schema"
)

func TestSeasonStatColumnNames(t *testing.T) {
	t.Parallel()

	s, err := schema.Parse(&SeasonStat{}, &sync.Map{}, schema.NamingStrategy{})
	if err != nil {
		t.Fatalf("schema.Parse returned error: %v", err)
	}

	expected := map[string]string{
		"BABIP":          "babip",
		"WHIP":           "whip",
		"StrikeoutsPer9": "strikeouts_per_9",
		"WalksPer9":      "walks_per_9",
		"HitsPer9":       "hits_per_9",
		"HomeRunsPer9":   "home_runs_per_9",
	}

	for fieldName, columnName := range expected {
		field := s.LookUpField(fieldName)
		if field == nil {
			t.Fatalf("field %s not found", fieldName)
		}
		if field.DBName != columnName {
			t.Fatalf("%s DBName = %s, want %s", fieldName, field.DBName, columnName)
		}
	}
}
