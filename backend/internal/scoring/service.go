package scoring

import (
	"context"
	"sort"

	"gorm.io/gorm"

	"github.com/jsquardo/capcurve/internal/models"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

// RecalculateYears refreshes value_score for every row in the affected MLB
// seasons so percentiles remain scoped to the same season instead of all-time.
func (s *Service) RecalculateYears(ctx context.Context, db *gorm.DB, years []int) error {
	for _, year := range uniqueYears(years) {
		var stats []models.SeasonStat
		if err := db.WithContext(ctx).Where("year = ?", year).Find(&stats).Error; err != nil {
			return err
		}

		scores := ScoreSeasonStats(stats)
		for _, stat := range stats {
			if err := db.WithContext(ctx).
				Model(&models.SeasonStat{}).
				Where("id = ?", stat.ID).
				Update("value_score", scores[stat.ID].FinalScore).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

func uniqueYears(years []int) []int {
	set := make(map[int]struct{}, len(years))
	for _, year := range years {
		if year == 0 {
			continue
		}
		set[year] = struct{}{}
	}

	unique := make([]int, 0, len(set))
	for year := range set {
		unique = append(unique, year)
	}
	sort.Ints(unique)
	return unique
}
